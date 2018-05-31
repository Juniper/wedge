/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package gRPC

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	wc "github.com/Juniper/wedge/codecs"
	wu "github.com/Juniper/wedge/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	yaml "gopkg.in/yaml.v2"
)

const (
	RESPONSE_NORMAL = 0
	RESPONSE_ERROR
)

const GRPC_SERVER_MAX_CLIENTS = 10
const GRPC_SERVER_MAX_CALLS_PER_CLIENT = 5

/*
 * GrpcCall structure From Client Handler perspective
 */
type GrpcClientCall struct {
	rpcType       int
	termChan      chan struct{}
	callInputChan chan interface{} // Needed for input streaming
	format        wu.MsgFormat
	respMutex     sync.Mutex // Used for synchronization in reponses by call routines
	doneMutex     sync.Mutex // Used for sychronization indicating end of call
	rpcInput      interface{}
	rpcOutput     interface{}
	CallTimeout   uint
}

func createFormatCopy(old, resp *wu.MsgFormat) {
	resp.ClientId = old.ClientId
	resp.IpAddress = old.IpAddress
	resp.TransactionId = old.TransactionId
	resp.Port = old.Port
	resp.Metadata = old.Metadata
	resp.Rpc = old.Rpc
	resp.Value = ""
}

func (g *GrpcClientCall) writeResponse(ret wu.MsgFormat,
	respChan chan<- wu.MsgFormat) {
	g.respMutex.Lock()

	respChan <- ret

	g.respMutex.Unlock()
}

func (g *GrpcClientCall) writeDone(callDoneChan chan<- string) {
	var rpcId string

	rpcId = g.format.TransactionId + "_" + g.format.RpcId + "_" + g.format.Rpc

	g.doneMutex.Lock()

	log.Printf("%s: Execute call Done for RPC %s", wu.FuncName(), g.format.Rpc)

	callDoneChan <- rpcId

	g.doneMutex.Unlock()
}

/*
 * Method execute a unary RPC
 */
func (g *GrpcClientCall) ExecuteUnaryRPC(conn *grpc.ClientConn,
	codec wu.GenCodec, rpcInput, rpcOutput interface{},
	termChan <-chan struct{}, respChan chan<- wu.MsgFormat,
	callDoneChan chan<- string) {
	var err error
	var ret wu.MsgFormat

	createFormatCopy(&g.format, &ret)

	for {
		select {
		case <-termChan:
			ret.Value = "Execution was cancelled for RPC " + g.format.Rpc
			ret.Rpc = wu.CALL_CANCELLATION_RPC
			g.writeResponse(ret, respChan)
			g.writeDone(callDoneChan)
			log.Printf("%s: ****Got Term, Call Done******", wu.FuncName())
			return
		default:
			err = grpc.Invoke(context.Background(), g.format.Rpc, g.rpcInput,
				g.rpcOutput, conn)
			if err != nil {
				errorStr := fmt.Sprintf("%s: Execution of Unary RPC %s failed "+
					"with error %v", wu.FuncName(), g.format.Rpc, err)
				ret.Value = errorStr
				ret.Metadata = nil
				g.writeResponse(ret, respChan)
				g.writeDone(callDoneChan)
				return
			}

			ret.Value = codec.DecodeOutput(g.rpcOutput)
			ret.Metadata = nil
			g.writeResponse(ret, respChan)
			g.writeDone(callDoneChan)
			log.Printf("%s: Call Done******", wu.FuncName())
			return
		}
	}
}

/*
 * Method to execute a server side streaming call
 */
func (g *GrpcClientCall) ExecuteServerStreamRPC(conn *grpc.ClientConn,
	codec wu.GenCodec, rpcInput, rpcOutput interface{},
	termChan <-chan struct{}, respChan chan<- wu.MsgFormat,
	callDoneChan chan<- string) {
	var err error
	var errorStr string
	var ret wu.MsgFormat

	createFormatCopy(&g.format, &ret)

	log.Printf("%s: format is %s", wu.FuncName(), g.format)

	log.Printf("%s: ret is %s", wu.FuncName(), ret)

	name := strings.Split(g.format.Rpc, "/")
	serverStreamDesc := grpc.StreamDesc{
		StreamName:    name[len(name)-1],
		Handler:       nil,
		ServerStreams: true,
	}

	//ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	stream, err := grpc.NewClientStream(context.Background(),
		&serverStreamDesc, conn, g.format.Rpc)

	if err != nil {
		errorStr = fmt.Sprintf("%s: Stream creation for RPC %s failed with "+
			"error: %v", wu.FuncName(), g.format.Rpc, err)
		ret.Value = errorStr
		ret.Metadata = nil
		g.writeDone(callDoneChan)
		return
	}

	if err = stream.SendMsg(rpcInput); err != nil {
		errorStr = fmt.Sprintf("%s: SendMsg() for RPC %s failed with error: %v",
			wu.FuncName(), g.format.Rpc, err)
		ret.Value = errorStr
		ret.Metadata = nil
		g.writeResponse(ret, respChan)
		g.writeDone(callDoneChan)
		return
	}

	if err = stream.CloseSend(); err != nil {
		errorStr = fmt.Sprintf("%s: CloseSend() for RPC %s failed with "+
			"error: %v", wu.FuncName(), g.format.Rpc, err)
		ret.Value = errorStr
		ret.Metadata = nil
		g.writeResponse(ret, respChan)
		g.writeDone(callDoneChan)
		return
	}

	for {
		select {
		case <-termChan:
			ret.Value = "Execution was cancelled for RPC " + g.format.Rpc
			ret.Rpc = wu.CALL_CANCELLATION_RPC
			g.writeResponse(ret, respChan)
			g.writeDone(callDoneChan)
			return
		default:
			if err = stream.RecvMsg(rpcOutput); err != nil {
				if grpc.Code(err) == codes.DeadlineExceeded {
					break
				} else if err != io.EOF {
					errorStr = fmt.Sprintf("%s: RecvMsg() for RPC %s failed "+
						"with error: %v", wu.FuncName(), g.format.Rpc, err)
					ret.Value = errorStr
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
					return
				} else {
					ret.Value = codec.DecodeOutput(rpcOutput)
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
					return
				}
			} else {
				ret.Value = codec.DecodeOutput(rpcOutput)
				ret.Metadata = nil
				g.writeResponse(ret, respChan)
			}
		}
	}
}

/*
 * Method to execute a client side streaming call
 */
func (g *GrpcClientCall) ExecuteClientStreamRPC(conn *grpc.ClientConn,
	codec wu.GenCodec, rpcOutput interface{},
	termChan <-chan struct{}, inputChan <-chan interface{},
	respChan chan<- wu.MsgFormat, callDoneChan chan<- string) {
	var err error
	var errorStr string
	var rpcInput interface{}
	var ok bool
	var ret wu.MsgFormat

	createFormatCopy(&g.format, &ret)

	name := strings.Split(g.format.Rpc, "/")
	clientStreamDesc := grpc.StreamDesc{
		StreamName:    name[len(name)-1],
		Handler:       nil,
		ClientStreams: true,
	}

	stream, err := grpc.NewClientStream(context.Background(), &clientStreamDesc,
		conn, g.format.Rpc)

	if err != nil {
		errorStr = fmt.Sprintf("%s: Stream creation for RPC %s failed with "+
			"error: %v", wu.FuncName(), g.format.Rpc, err)
		ret.Value = errorStr
		ret.Metadata = nil
		g.writeResponse(ret, respChan)
		g.writeDone(callDoneChan)
		return
	}

	for {
		select {
		case <-termChan:
			ret.Value = "Execution was cancelled for RPC " + g.format.Rpc
			ret.Rpc = wu.CALL_CANCELLATION_RPC
			g.writeResponse(ret, respChan)
			g.writeDone(callDoneChan)
			return
		case rpcInput, ok = <-inputChan:
			if !ok {
				ret.Value = "Client side stream receive channel closed"
				ret.Metadata = nil
				g.writeResponse(ret, respChan)
				g.writeDone(callDoneChan)
				return
			} else if codec.IsEmpty(rpcInput) {
				if err = stream.SendMsg(rpcInput); err != nil {
					errorStr = fmt.Sprintf("%s: SendMsg() for RPC %s failed "+
						"with error: %v", wu.FuncName(), g.format.Rpc, err)
					ret.Value = errorStr
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
				}
			} else {
				if err = stream.CloseSend(); err != nil {
					errorStr = fmt.Sprintf("%s: CloseSend() for RPC %s failed "+
						"with error: %v", wu.FuncName(), g.format.Rpc, err)
					ret.Value = errorStr
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
					return
				}

				if err = stream.RecvMsg(rpcOutput); err != nil {
					errorStr = fmt.Sprintf("%s: RecvMsg() for RPC %s failed "+
						"with error: %v", wu.FuncName(), g.format.Rpc, err)
					ret.Value = errorStr
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
					return
				} else {
					ret.Value = codec.DecodeOutput(rpcOutput)
					g.writeResponse(ret, respChan)
					g.writeDone(callDoneChan)
					return
				}
			}
		}
	}
}

/*
 * Method to execute a Bidirectional streaming call
 */
func (g *GrpcClientCall) ExecuteBidiStreamRPC(conn *grpc.ClientConn,
	codec wu.GenCodec, rpcOutput interface{},
	termChan <-chan struct{}, inputChan <-chan interface{},
	respChan chan<- wu.MsgFormat, callDoneChan chan<- string) {
	var err error
	var inErr, outErr string
	var rpcInput interface{}
	var ok bool
	var ret wu.MsgFormat

	createFormatCopy(&g.format, &ret)

	name := strings.Split(g.format.Rpc, "/")
	clientStreamDesc := grpc.StreamDesc{
		StreamName:    name[len(name)-1],
		Handler:       nil,
		ServerStreams: true,
		ClientStreams: true,
	}

	stream, err := grpc.NewClientStream(context.Background(), &clientStreamDesc,
		conn, g.format.Rpc)

	if err != nil {
		inErr = fmt.Sprintf("%s: Stream creation for RPC %s failed with "+
			"error: %v", wu.FuncName(), g.format.Rpc, err)
		ret.Value = inErr
		ret.Metadata = nil
		g.writeResponse(ret, respChan)
		g.writeDone(callDoneChan)
	}

	/*
	 * Start a Go routine for receiving streaming data
	 * from the server i.e output/server side streaming
	 */
	waitc := make(chan struct{})
	go func() {
		for {
			select {
			case <-termChan:
				close(waitc)
				return
			default:
				if err = stream.RecvMsg(rpcOutput); err != nil {
					if err != io.EOF {
						outErr = fmt.Sprintf("%s: RecvMsg() for RPC %s "+
							"failed with error: %v", wu.FuncName(),
							g.format.Rpc, err)
						close(waitc)
						break
					} else {
						ret.Value = codec.DecodeOutput(rpcOutput)
						break
					}
				} else {
					ret.Value = codec.DecodeOutput(rpcOutput)
					ret.Metadata = nil
					g.writeResponse(ret, respChan)
				}
			}
		}
	}()

	/*
	 * Loop to send streaming data to the server ie.
	 * input/client side streaming
	 */
	for {
		select {
		case <-termChan:
			ret.Value = "Execution was cancelled for RPC " + g.format.Rpc
			ret.Rpc = wu.CALL_CANCELLATION_RPC
			g.writeResponse(ret, respChan)
			goto INPUT_STREAM_DONE
		case rpcInput, ok = <-inputChan:
			if !ok {
				inErr = "Client side stream receive channel closed"
				goto INPUT_STREAM_DONE
			} else {
				if err = stream.SendMsg(rpcInput); err != nil {
					inErr = fmt.Sprintf("%s: SendMsg() for RPC %s failed with "+
						"error: %v", wu.FuncName(), g.format.Rpc, err)
				}
				if rpcInput == nil {
					goto INPUT_STREAM_DONE
				}
			}
		}
	}

INPUT_STREAM_DONE:
	<-waitc

	/*
	 * Check for any error and respond accordingly
	 */
	inErr += outErr

	if inErr != "" {
		ret.Value = inErr
	}

	g.writeDone(callDoneChan)

}

/*
 * GrpcClient structure from grpc plugin perspective
 */
type GrpcClient struct {
	conn         *grpc.ClientConn
	codec        wu.GenCodec
	maxCalls     uint             // Indicates max number of parallel calls for the client
	termChan     chan struct{}    //send termination to this client
	inputChan    chan interface{} // To send calls to client
	clientMutex  sync.Mutex       // Used between Client's request and response handlers
	respMutex    sync.Mutex       // Used for synchronization in reponses by client routines
	doneMutex    sync.Mutex       // Used for sychronization indicating client termination
	curCallCount uint             // Indicates the call capacity
	callMap      map[string]*GrpcClientCall
	CallTimeout  uint
}

/*
 * Options to be be specified for a
 * gRPC connection
 */
type GrpcConnOptions struct {
	InitWindowSize     int32                            `yaml:"init-window-size"`
	InitConnWindowSize int32                            `yaml:"init-conn-window-size"`
	Compressor         string                           `yaml:"compressor"`   // gRPC supports gzip
	Decompressor       string                           `yaml:"decompressor"` // gRPC supports gzip
	MaxBackoffDelay    time.Duration                    `yaml:"max-backoff-delay"`
	Insecure           bool                             `yaml:"insecure"` //This bool val or Creds/PerCallCreds needs to set
	Creds              credentials.TransportCredentials `yaml:"credentials"`
	PerCallCreds       credentials.PerRPCCredentials    `yaml:"per-call-credentials"`
	MaxCallSendSize    int                              `yaml:"max-call-sendsize"`
	MaxCallRecvSize    int                              `yaml:"max-call-recvsize"`
	Block              bool                             `yaml:"block"`
	Timeout            time.Duration                    `yaml:"timeout"`
}

/*
 * Function to create a gRPC client
 */
func CreateGrpcClient(address string, maxCalls uint,
	gOpts GrpcConnOptions, codecType int) (*GrpcClient, error) {
	var err error
	var g *GrpcClient = new(GrpcClient)
	var opts []grpc.DialOption
	g.codec = wc.GetCodecObject(codecType)

	if g.codec == nil {
		errorStr := fmt.Sprintf("%s: Could not create a codec with type %d",
			wu.FuncName(), codecType)
		return nil, errors.New(errorStr)
	}

	/*
	 * Parse the structure and populate the options
	 */
	if gOpts.InitWindowSize > 0 {
		if gOpts.InitWindowSize < 65536 {
			errorStr := fmt.Sprintf("%s: InitWindowSize has a value %d that "+
				"is less that the min value of 64K", wu.FuncName(),
				gOpts.InitWindowSize)
			return nil, errors.New(errorStr)
		}

		opts = append(opts, grpc.WithInitialWindowSize(gOpts.InitWindowSize))
	}

	if gOpts.InitConnWindowSize > 0 {
		if gOpts.InitWindowSize < 65536 {
			errorStr := fmt.Sprintf("%s: InitWindowSize has a value %d that "+
				"is less that the min value of 64K", wu.FuncName(),
				gOpts.InitWindowSize)
			return nil, errors.New(errorStr)
		}

		opts = append(opts,
			grpc.WithInitialConnWindowSize(gOpts.InitConnWindowSize))
	}

	opts = append(opts, grpc.WithCodec(g.codec))

	if gOpts.Compressor != "" {
		if gOpts.Compressor == "gzip" {
			opts = append(opts, grpc.WithCompressor(grpc.NewGZIPCompressor()))
		}
	}

	if gOpts.Decompressor != "" {
		if gOpts.Decompressor == "gzip" {
			opts = append(opts,
				grpc.WithDecompressor(grpc.NewGZIPDecompressor()))
		}
	}

	if gOpts.MaxBackoffDelay > 0 {
		opts = append(opts, grpc.WithBackoffMaxDelay(gOpts.MaxBackoffDelay))
	}

	/*
	 * Setting Backoffconfig parameters
	 */
	if gOpts.MaxBackoffDelay > 0 {
		opts = append(opts, grpc.WithBackoffMaxDelay(gOpts.MaxBackoffDelay))
	}

	if gOpts.Insecure == true {
		opts = append(opts, grpc.WithInsecure())
	} else if gOpts.Creds != nil {
		opts = append(opts, grpc.WithTransportCredentials(gOpts.Creds))
	} else if gOpts.PerCallCreds != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(gOpts.PerCallCreds))
	} else {
		errorStr := fmt.Sprintf("%s: Either Insecure should be true or Creds/"+
			"/PerCallCreds should be set for gRPC connection", wu.FuncName())
		return nil, errors.New(errorStr)
	}

	if gOpts.MaxCallRecvSize > 0 {
		opts = append(opts,
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(
				gOpts.MaxCallRecvSize)))
	}

	if gOpts.MaxCallSendSize > 0 {
		opts = append(opts,
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(
				gOpts.MaxCallSendSize)))
	}

	if gOpts.Block == true && gOpts.Timeout > 0 {
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(gOpts.Timeout))
	}

	g.conn, err = grpc.Dial(address, opts...)

	if err != nil {
		errorStr := fmt.Sprintf("%s: Creating a gRPC client dial option failed "+
			"with error: %v", wu.FuncName(), err)
		return nil, errors.New(errorStr)
	}

	/*
	 * Create the required client parameters
	 */
	g.maxCalls = maxCalls

	g.callMap = make(map[string]*GrpcClientCall)

	return g, nil
}

/*
 * Method to destroy a gRPC client.
 */
func (g *GrpcClient) DestroyGrpcClient() {
	g.conn.Close()
}

func (gc *GrpcClient) callCleanup(rpcId string, call *GrpcClientCall) {
	var ok bool

	if call == nil {
		call, ok = gc.callMap[rpcId]
	}

	if ok {
		close(call.termChan)
		if call.callInputChan != nil {
			close(call.callInputChan)
		}

		delete(gc.callMap, rpcId)
	}
}

func (gc *GrpcClient) startClientCall(call *GrpcClientCall,
	conn *grpc.ClientConn, codec wu.GenCodec, rpcInput, rpcOutput interface{},
	termChan <-chan struct{}, inputChan <-chan interface{},
	respChan chan<- wu.MsgFormat, doneChan chan<- string) {

	/*
	 * Start a Go routine for the call depending
	 * on the call type
	 */
	switch call.rpcType {
	case wu.RPC_TYPE_UNARY:
		go call.ExecuteUnaryRPC(conn, codec, rpcInput, rpcOutput, termChan,
			respChan, doneChan)
	case wu.RPC_TYPE_SERVER_STREAMING:
		go call.ExecuteServerStreamRPC(conn, codec, rpcInput, rpcOutput,
			termChan, respChan, doneChan)
	case wu.RPC_TYPE_CLIENT_STREAMING:
		go call.ExecuteClientStreamRPC(conn, codec, rpcOutput, termChan,
			inputChan, respChan, doneChan)
	case wu.RPC_TYPE_BIDISTREAMING:
		go call.ExecuteBidiStreamRPC(conn, codec, rpcOutput, termChan,
			inputChan, respChan, doneChan)
	}
}

/*
 * Method to handle requests for a client
 */
func (gc *GrpcClient) clientRequestHandler(codecType int, clientId string,
	inputChan <-chan interface{}, outputChan chan<- wu.MsgFormat,
	termChan <-chan struct{}, clientDoneChan chan<- string) {
	var rpc, errorStr string
	var err error
	var input interface{}
	var format wu.MsgFormat
	var ok bool
	var rpcType int
	var call *GrpcClientCall
	var pendingCallQueue []*GrpcClientCall

	respChan := make(chan wu.MsgFormat, 1000)
	doneChan := make(chan string, gc.maxCalls)
	gc.callMap = make(map[string]*GrpcClientCall)

	/*
	 * Start the response handler go routine
	 */
	go gc.clientResponseHandler(clientId, outputChan, clientDoneChan,
		respChan, doneChan, termChan)

	for {
		select {
		case <-termChan:
			/*
			 * Walk the call map for this client
			 * and terminate them
			 */
			log.Printf("%s: Received termination request for client %s",
				wu.FuncName(), clientId)
			gc.clientMutex.Lock()
			for rpc, call = range gc.callMap {
				gc.callCleanup(rpc, call)
			}
			gc.clientMutex.Unlock()
			return
		case input, ok = <-inputChan:
			if !ok {
				log.Printf("%s: Input Channel closed for client %s",
					wu.FuncName(), clientId)
				gc.clientMutex.Lock()
				for rpc, call = range gc.callMap {
					gc.callCleanup(rpc, call)
				}
				gc.clientMutex.Unlock()
				return
			}

			format = input.(wu.MsgFormat)
			log.Printf("%s: Got input with Id: %s, Ipaddress: %s, "+
				"Metadata :%s, Port: %s, Rpc: %s, TransactionId: %s, RpcId: %s",
				wu.FuncName(), format.ClientId, format.IpAddress,
				format.Metadata, format.Port, format.Rpc, format.TransactionId,
				format.RpcId)

			if format.TransactionId == "" || format.RpcId == "" {
				errorStr = fmt.Sprintf("Mandatory parameter(s) TransactionId "+
					"and/or RpcId not specified for RPC %s", format.Rpc)
				format.Value = errorStr
				gc.writeResponse(format, outputChan)
				break
			}

			/*
			 * Check if the RPC is already under execution, this will
			 * be mostly the case for client streaming or bidirectional
			 * streaming RPC or a call cancellation is performed
			 */
			var callId string
			cmp := strings.Compare(format.Rpc, wu.CALL_CANCELLATION_RPC)

			callId = format.TransactionId + "_" + format.RpcId + "_" +
				format.Rpc

			if cmp != 0 {
				log.Printf("%s: %s is not a cancellation RPC", wu.FuncName(),
					format.Rpc)
				call, ok = gc.callMap[callId]
				rpcType, err = wu.GetRpcType(format.Rpc, nil)
			} else {
				log.Printf("%s: Received cancellation RPC for call %s with "+
					"transaction id %s and rpc id %s",
					wu.FuncName(), format.TransactionId, format.RpcId,
					format.Value.(string))

				callId = format.TransactionId + "_" + format.RpcId + "_" +
					format.Value.(string)
				call, ok = gc.callMap[callId]
				rpcType, err = wu.GetRpcType(format.Value.(string), nil)
			}

			if ok {
				if cmp == 0 {
					/*
					 * Cancel the call
					 */
					log.Printf("%s: Cancelling call %s for client %s",
						wu.FuncName(), format.Value.(string),
						call.format.ClientId)

					gc.clientMutex.Lock()
					gc.callCleanup(callId, nil)
					gc.clientMutex.Unlock()
				} else {
					if rpcType == wu.RPC_TYPE_CLIENT_STREAMING ||
						rpcType == wu.RPC_TYPE_BIDISTREAMING {
						/*
						 * Parse the input using the codec and encode
						 * it into proto
						 */
						call.rpcInput, err = gc.codec.EncodeInput(format.Rpc,
							format.Value)
						if err != nil {
							errorStr = fmt.Sprintf("Encoding input to proto "+
								"failied for RPC %s with error %v", format.Rpc,
								err)
							format.Value = errorStr
							gc.writeResponse(format, outputChan)
							break
						}

						/*
						 * Write the data to the correspoding input
						 * channel of the call
						 */
						call.callInputChan <- call.rpcInput
					} else {
						rtype, _ := wu.GetRpcTypeString(rpcType)
						errorStr = fmt.Sprintf("Got an unexpected data for RPC "+
							"%s of type %s with transaction id %s and rpc id %s"+
							" currently under execution", format.Rpc,
							format.TransactionId, format.RpcId, rtype)
						format.Value = errorStr
						gc.writeResponse(format, outputChan)
						break
					}
				}
			} else if cmp == 0 {
				log.Printf("%s: inside else-if for Rpc %s", wu.FuncName(), format.Rpc)
				log.Printf("%s: Termination request for a non-existent call %s",
					wu.FuncName(), format.Value.(string))

				errorStr = fmt.Sprintf("Termination request for a "+
					"non-existent call %s", format.Value.(string))
				format.Value = errorStr
				gc.writeResponse(format, outputChan)
				break
			} else {
				/*
				 * Create the GrpcCall structure for this call
				 */
				call = new(GrpcClientCall)
				gc.callMap[callId] = call

				call.rpcType = rpcType

				var callInputChan interface{}
				if rpcType == wu.RPC_TYPE_CLIENT_STREAMING ||
					rpcType == wu.RPC_TYPE_BIDISTREAMING {
					callInputChan = make(chan interface{}, 1000)
					call.callInputChan = callInputChan.(chan interface{})
				}

				callTermChan := make(chan struct{}, 2)
				call.termChan = callTermChan
				call.format = format
				call.CallTimeout = gc.CallTimeout

				call.rpcInput, err = gc.codec.EncodeInput(format.Rpc,
					format.Value)
				if err != nil {
					errorStr = fmt.Sprintf("Encoding input to proto "+
						"failied for RPC %s with error %v", format.Rpc,
						err)
					format.Value = errorStr
					gc.writeResponse(format, outputChan)
					break
				}

				call.rpcOutput = gc.codec.CreateRPCOutputObj(call.rpcInput)

				/*
				 * For a new call queue, first check if the curCallCount
				 * doesn't exceed max call count and also check if any calls
				 * are in pending call queue
				 */
				log.Printf("%s: before lock for Rpc %s",
					wu.FuncName(), format.Rpc)

				gc.clientMutex.Lock()

				log.Printf("%s: Adding RPC %s to call map for client %s",
					wu.FuncName(), format.Rpc, format.ClientId)

				if gc.curCallCount >= gc.maxCalls {
					pendingCallQueue = append(pendingCallQueue, call)
				} else if len(pendingCallQueue) > 0 {
					pendingCallQueue = append(pendingCallQueue, call)

					/*
					 * Walk the pending call queue and schedule calls
					 * till curCallCount becomes equal to maxCalls
					 */
					for {
						if gc.curCallCount < gc.maxCalls {
							call = pendingCallQueue[0]
							pendingCallQueue = pendingCallQueue[1:]
							gc.startClientCall(call, gc.conn, gc.codec,
								call.rpcInput, call.rpcOutput,
								call.termChan, call.callInputChan,
								respChan, doneChan)
							gc.curCallCount++
						} else {
							break
						}
					}
				} else {
					log.Printf("%s: Executing RPC %s to target %s",
						wu.FuncName(), call.format.Rpc, call.format.IpAddress)
					gc.startClientCall(call, gc.conn, gc.codec,
						call.rpcInput, call.rpcOutput, call.termChan,
						call.callInputChan, respChan, doneChan)
					gc.curCallCount++
				}

				gc.clientMutex.Unlock()

			}
		}
	}
}

/*
 * Method to process responses from call routines for a client
 */
func (gc *GrpcClient) clientResponseHandler(clientId string,
	outputChan chan<- wu.MsgFormat, clientDoneChan chan<- string,
	respChan chan wu.MsgFormat, doneChan chan string,
	termChan <-chan struct{}) {

	var term, ok bool
	var resp wu.MsgFormat
	var errorStr string
	var rpcId string
	var call *GrpcClientCall

	defer close(respChan)
	defer close(doneChan)

	for {
		select {
		case <-termChan:
			term = true
			gc.DestroyGrpcClient()
			gc.clientMutex.Lock()
			if gc.curCallCount <= 0 {
				log.Printf("%s: Call count is 0, returning...", wu.FuncName())
				gc.writeClientDone(clientId, clientDoneChan)
				gc.clientMutex.Unlock()
				return
			} else {
				log.Printf("%s: Call count is %d, will wait...", wu.FuncName(),
					gc.curCallCount)
				gc.clientMutex.Unlock()
			}
		case resp, ok = <-respChan:
			if !ok {
				errorStr = fmt.Sprintf("Call response channel closed for "+
					" client %s", clientId)
				resp.Value = errorStr
				resp.Metadata = nil
				gc.writeResponse(resp, outputChan)
				gc.clientMutex.Lock()
				for rpcId, call = range gc.callMap {
					gc.callCleanup(rpcId, call)
				}
				gc.clientMutex.Unlock()
				term = true
			} else {
				gc.writeResponse(resp, outputChan)
			}
		case rpcId, ok = <-doneChan:
			if !ok {
				log.Printf("%s: done Chan is closed", wu.FuncName())
				gc.DestroyGrpcClient()
				return
			}
			log.Printf("%s: Received call done for RPC %s", wu.FuncName(),
				rpcId)

			gc.clientMutex.Lock()

			gc.callCleanup(rpcId, nil)
			if gc.curCallCount > 0 {
				gc.curCallCount--
			}

			if gc.curCallCount <= 0 && term == true {
				log.Printf("%s: Terminating response handler for client %s",
					wu.FuncName(), clientId)
				defer gc.clientMutex.Unlock()
				gc.DestroyGrpcClient()
				/*
				 * Destroy the gRPC client connection to
				 * the router
				 */
				gc.writeClientDone(clientId, clientDoneChan)
				return
			}

			gc.clientMutex.Unlock()

		}
	}
}

/*
 * gRPC server side plugin
 */
type GrpcServerPlugin struct {
	clientMap      map[string]*GrpcClient
	pluginMutex    sync.Mutex // mutex between server plugiin routines
	curClientCount uint
	maxClients     uint
	maxCalls       uint
}

type GrpcServerPluginConfig struct {
	MaxClients    uint `yaml:"max-clients"`
	MaxCalls      uint `yaml:"max-calls"`
	CallTimeout   uint `yaml:"call-timeout"`
	Options       GrpcConnOptions
	ProtoTableLoc string `yaml:"proto-desc-table-loc"`
}

func (gs *GrpcServerPlugin) clientCleanup(clientId string, client *GrpcClient) {
	var ok bool

	if client == nil {
		client, ok = gs.clientMap[clientId]
	}

	if ok {
		close(client.termChan)
		close(client.inputChan)
		delete(gs.clientMap, clientId)
	}
}

func (gc *GrpcClient) writeResponse(resp wu.MsgFormat,
	outputChan chan<- wu.MsgFormat) {
	gc.respMutex.Lock()

	outputChan <- resp

	gc.respMutex.Unlock()
}

func (gc *GrpcClient) writeClientDone(clientId string,
	clientDoneChan chan<- string) {
	gc.doneMutex.Lock()
	log.Printf("%s: Execute client Done for clientId %s", wu.FuncName(),
		clientId)

	clientDoneChan <- clientId

	gc.doneMutex.Unlock()
}

func (gs *GrpcServerPlugin) runPluginClientRespAggregator(
	sendChan chan<- wu.MsgFormat, respChan chan wu.MsgFormat,
	doneChan chan string, termChan <-chan struct{},
	statusChan chan<- wu.TermStatus) {

	var term, ok bool
	var resp wu.MsgFormat
	var errorStr string
	var clientId string
	var client *GrpcClient
	var status wu.TermStatus

	defer close(respChan)
	defer close(doneChan)
	defer close(sendChan)

	for {
		select {
		case <-termChan:
			term = true
			gs.pluginMutex.Lock()
			if gs.curClientCount <= 0 {
				gs.pluginMutex.Unlock()
				status.Status = wu.ROUTINE_NORMAL_EXIT
				status.ErrStr = ""
				statusChan <- status
			} else {
				gs.pluginMutex.Unlock()
			}

		case resp, ok = <-respChan:
			if !ok {
				errorStr = fmt.Sprintf("Client response channel closed")
				resp.Value = errorStr
				resp.Metadata = nil
				sendChan <- resp

				gs.pluginMutex.Lock()
				for clientId, client = range gs.clientMap {
					gs.clientCleanup(clientId, client)
				}
				gs.pluginMutex.Unlock()
				term = true
			} else {
				sendChan <- resp
			}
		case clientId = <-doneChan:
			gs.pluginMutex.Lock()

			gs.clientCleanup(clientId, nil)
			if gs.curClientCount > 0 {
				gs.curClientCount--
			}

			log.Printf("%s: Completed termination for client %s", wu.FuncName(),
				clientId)

			if gs.curClientCount <= 0 && term == true {
				gs.pluginMutex.Unlock()
				status.Status = wu.ROUTINE_NORMAL_EXIT
				status.ErrStr = ""
				statusChan <- status
			} else {

				gs.pluginMutex.Unlock()
			}

		}
	}
}

func (gs GrpcServerPlugin) RunServerPlugin(codecType int, params interface{},
	readFunc wu.PostReadFunc, writeFunc wu.PreWriteFunc, errorFunc wu.ErrorFunc,
	sendChan chan<- wu.MsgFormat, recvChan <-chan wu.MsgFormat,
	termChan <-chan struct{}, statusChan chan<- wu.TermStatus) {

	var msg wu.MsgFormat
	var ok bool
	var err error
	var clientId, key, errorStr string
	var client *GrpcClient
	var nwAddress string
	var pendingConqueue []*GrpcClient
	var pluginConfig *GrpcServerPluginConfig

	log.Printf("%s: Starting gRPC server plugin", wu.FuncName())

	if params != nil {
		pluginConfig = params.(*GrpcServerPluginConfig)

		if pluginConfig.ProtoTableLoc == "" {
			status := new(wu.TermStatus)
			status.Status = wu.ROUTINE_ERROR_EXIT
			status.ErrStr = "ProtoTableLoc cannot be empty"
			statusChan <- *status
			return
		}

		if pluginConfig.MaxClients == 0 {
			pluginConfig.MaxClients = GRPC_SERVER_MAX_CLIENTS
		}
		gs.maxClients = pluginConfig.MaxClients

		if pluginConfig.MaxCalls == 0 {
			pluginConfig.MaxCalls = GRPC_SERVER_MAX_CALLS_PER_CLIENT
		} else {

		}
	} else {
		/*
		 * ProtoTableLoc is a mandatory config parameter that
		 * needs to be passed for the plugin to perform codec
		 * transformation. Return an error
		 */
		status := new(wu.TermStatus)
		status.Status = wu.ROUTINE_ERROR_EXIT
		status.ErrStr = "params cannot be empty. Please pass an object of " +
			"type GrpcServerPluginConfig with ProtoTableLoc set"
		statusChan <- *status
		return
	}

	wu.InitProtoParser(pluginConfig.ProtoTableLoc)

	clientRespChan := make(chan wu.MsgFormat, 1000)
	clientDoneChan := make(chan string, pluginConfig.MaxClients)
	gs.clientMap = make(map[string]*GrpcClient)

	/*
	 * Create the response aggregator go routine
	 */
	go gs.runPluginClientRespAggregator(sendChan, clientRespChan,
		clientDoneChan, termChan, statusChan)

	for {
		select {
		case <-termChan:
			/*
			 * Walk the client map and perform cleanup
			 */
			gs.pluginMutex.Lock()

			for clientId, client = range gs.clientMap {
				gs.clientCleanup(clientId, client)
			}

			gs.pluginMutex.Unlock()
			return
		case msg, ok = <-recvChan:
			if !ok {
				log.Printf("%s: Recv Channel closed for Server side plugin",
					wu.FuncName())
				gs.pluginMutex.Lock()

				for clientId, client = range gs.clientMap {
					gs.clientCleanup(clientId, client)
				}

				gs.pluginMutex.Unlock()
				return
			}

			log.Printf("%s: Received a request", wu.FuncName())

			/*
			 * Check if the client already exists and if so, send the call
			 * to the respective client(s). For this, the IpAddresses parameter
			 * within the key. A combination of the clientId and IpAddress
			 * will be used as a key since the same client can connect to
			 * different routers
			 */
			cmp := strings.Compare(msg.Rpc, wu.CLIENT_TERM_RPC)
			key = msg.ClientId + msg.IpAddress
			client, ok = gs.clientMap[key]

			/*
			 * If the client, ip address combination already exits,
			 * then send the call through the relevant channel
			 */
			if ok {
				/*
				 * Check if the RPC is for client terminate and if so,
				 * close the termChan for the client.
				 */
				if cmp == 0 {
					gs.pluginMutex.Lock()
					log.Printf("%s: Terminating client %s", wu.FuncName(),
						key)
					gs.clientCleanup(key, nil)
					gs.pluginMutex.Unlock()
					msg.Value = "Connection to router " + msg.IpAddress +
						" terminated for client " + msg.ClientId
					sendChan <- msg
				} else {
					client.inputChan <- msg
				}
			} else if cmp == 0 {
				errorStr = fmt.Sprintf("Grpc Client termination requrest for "+
					"a non-existent client connection %s", key)
				msg.Value = errorStr
				msg.Metadata = nil
				sendChan <- msg
			} else {
				/*
				 * based on if the ipaddress is ipv4 or v6, configure
				 * the network address for gRPC client connection
				 */
				if strings.Contains(msg.IpAddress, ":") {
					nwAddress = "[" + msg.IpAddress + "]:" + msg.Port
				} else {
					nwAddress = msg.IpAddress + ":" + msg.Port
				}

				/*
				 * Create a connection to the respective router ip
				 */
				client, err = CreateGrpcClient(nwAddress, pluginConfig.MaxCalls,
					pluginConfig.Options, codecType)
				if err != nil {
					errorStr = fmt.Sprintf("Grpc Client Creation for address "+
						"%s, port %s failed for client %s with error %v",
						msg.IpAddress, msg.Port, msg.ClientId, err)
					msg.Value = errorStr
					msg.Metadata = nil
					sendChan <- msg
					break
				}

				log.Printf("%s: Client Created successfully address %s, port %s",
					wu.FuncName(), msg.IpAddress, msg.Port)

				/*
				 * Add the client entry to the map
				 */
				key = msg.ClientId + msg.IpAddress
				gs.clientMap[key] = client
				client.maxCalls = pluginConfig.MaxCalls
				client.CallTimeout = pluginConfig.CallTimeout

				clientInputChan := make(chan interface{}, 1000)
				clientTermChan := make(chan struct{}, 2)

				client.inputChan = clientInputChan
				client.termChan = clientTermChan

				/*
				 * Check for pending connections
				 */
				gs.pluginMutex.Lock()

				if gs.curClientCount > gs.maxClients {
					log.Printf("%s: Inside if with gs.curClientCount %d and gs.maxClients %d",
						wu.FuncName(), gs.curClientCount, gs.maxClients)

					pendingConqueue = append(pendingConqueue, client)
				} else if len(pendingConqueue) > 0 {

					log.Printf("%s: Inside else if", wu.FuncName())

					pendingConqueue = append(pendingConqueue, client)

					/*
					 * Walk the pending client connection queue and
					 * start the client routines
					 */
					for {
						if gs.curClientCount < gs.maxClients {
							client = pendingConqueue[0]
							pendingConqueue = pendingConqueue[1:]
							/*
							 * Start the client request handler go routine
							 */
							go client.clientRequestHandler(codecType, key,
								client.inputChan, clientRespChan,
								client.termChan, clientDoneChan)
							client.inputChan <- msg
							gs.curClientCount++
						} else {
							break
						}
					}
				} else {
					log.Printf("%s: Inside else", wu.FuncName())

					/*
					 * Start the client request handler go routine
					 */
					log.Printf("%s: Starting the client request handler",
						wu.FuncName())
					go client.clientRequestHandler(codecType, key,
						client.inputChan, clientRespChan,
						client.termChan, clientDoneChan)
					client.inputChan <- msg
					gs.curClientCount++
				}

				gs.pluginMutex.Unlock()

			}
		}
	}
}

/*
 * Plugin initialization functions for gRPC
 */
var GrpcPlugin wu.Plugin

func grpcParseServerSideConfig(data []byte) interface{} {
	var gconf GrpcServerPluginConfig

	err := yaml.Unmarshal(data, &gconf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &gconf
}

func grpcFetchErrorFunc(params interface{}) string {
	return ""
}

func grpcEmitRefConfig() string {
	return `
            name: "grpc"
            client-post-read-cb: ""
            client-pre-write-cb: ""
            server-post-read-cb: ""
            server-pre-write-cb: ""
            config:
                max-clients: 0 # Maximum number of client connections the broker should support, Default: 10
                max-calls: 0 # Maximum number of calls per client # Default: 5
                call-timeout: 0 # gRPC call timeout, Default: Blocking
                options:
                    init-window-size: 0 # Initial window size on a stream, the lower bound is 64K 
                    init-conn-window-size: 0 # Initial window size on a connection, the lower bound is 64K
                    compressor: "" #  Compressor to use for message compression 
                    decompressor: "" # Decompressor to use for incoming message decompression
                    max-backoff-delay: 0s # Maximum delay when backing off after failed connection attempts
                    insecure: false  # MANDATORY,  Setting to true disables transport security for this Client Connection
                    credentials: null # Connection level security credentials (e.g., TLS/SSL), Datatype: credentials.TransportCredentials 
                    per-call-credentials: null # Credentials for each outbound RPC, Datatype: credentials.PerRPCCredentials
                    max-call-sendsize: 0  # Maximum message size the client can send
                    max-call-recvsize: 0  # Maximum message size the client can receive
                    block: false # Dial blocks until the underlying connection is up.
                                 # Without this, Dial returns immediately and connecting the server
                                 # happens in background
                    timeout: 0s # configures a timeout for dialing a ClientConn initially.
                                # This is valid if and only if block is true
                proto-desc-table-loc: "" # MANDATORY, Location of the proto descriptor table
                                    # generated from proto using protoc-wedge compiler`
}

var plugin GrpcServerPlugin

func GrpcInit(pluginMap map[string]wu.Plugin) {
	GrpcPlugin.EmitServerSideRefConfig = grpcEmitRefConfig
	GrpcPlugin.ParseServerSideConfig = grpcParseServerSideConfig
	GrpcPlugin.ClientSidePluginSupport = false
	GrpcPlugin.ServerSidePluginSupport = true
	GrpcPlugin.ClientSideCodec = wu.CODEC_JSON_GRPC
	GrpcPlugin.ServerSideCodec = wu.CODEC_JSON_GRPC
	GrpcPlugin.ServerPlugin = plugin
	GrpcPlugin.FetchErrorFunc = grpcFetchErrorFunc
	pluginMap["grpc"] = GrpcPlugin
}
