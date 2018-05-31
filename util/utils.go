/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package util

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

const (
	MSG_FORMAT_RPC            = "Rpc"
	MSG_FORMAT_BROKER_ID 	  = "BrokerId"
	MSG_FORMAT_CLIENT_ID      = "ClientId"
	MSG_FORMAT_TRANSACTION_ID = "TransactionId"
	MSG_FORMAT_RPC_ID         = "RpcId"
	MSG_FORMAT_IPADDRESS      = "IpAddress"
	MSG_FORMAT_PORT           = "Port"
	MSG_FORMAT_VALUE          = "Value"
	MSG_FORMAT_METADATA       = "Metadata"
)

/*
 * Message format for data sent between client
 * and server plugins
 */
type MsgFormat struct {
	Rpc           string // RPC that needs to be processed
	BrokerId	  string // String to uniquely identify the broker that should process the message
	ClientId      string // String to uniquely identify each client
	TransactionId string // String to uniquely identify each transaction
	RpcId         string // String to uniquely identify an RPC in this transaction
	IpAddress     string // Ip address to connect, if needed
	Port          string
	Value         interface{}       // The payload on which processing needs to be done
	Metadata      map[string]string // Any additional metadata as key value pairs
}

const (
	ROUTINE_NORMAL_EXIT = 0
	ROUTINE_ERROR_EXIT  = 1
)

type TermStatus struct {
	Status int
	ErrStr string // Will be set on error
}

/*
 * Types of codec list
 */
const (
	CODEC_JSON_GRPC = 1
)

/*
 * List of RPC definitions
 */
const (
	CALL_CANCELLATION_RPC = "/wedge.Grpc/cancelCall"
	CLIENT_TERM_RPC       = "/wedge.Grpc/terminateClient"
)

/*
 * Type of gRPC call
 */
const (
	RPC_TYPE_UNARY            = 0
	RPC_TYPE_SERVER_STREAMING = 1
	RPC_TYPE_CLIENT_STREAMING = 2
	RPC_TYPE_BIDISTREAMING    = 3
)

/*
 * Generic codec interface superseding grpc codec
 */
type GenCodec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// String returns the name of the Codec implementation. The returned
	// string will be used as part of content type in transmission.
	String() string

	// Method to encode the input to the format needed by codec
	EncodeInput(rpc string, input interface{}) (interface{}, error)

	// Method to decode the ouput from the codec format
	DecodeOutput(codedOutput interface{}) interface{}

	// Method to create the output with RPC related parameters
	// copied
	CreateRPCOutputObj(input interface{}) interface{}

	// Method to indicate if the payload is empty as per the
	// codec format
	IsEmpty(codedVal interface{}) bool
}

/*
 * Function type for the function to be invoked
 * when data is received on client side plugin
 */
type PostReadFunc func(input interface{}) ([]MsgFormat, error)

/*
 * Function type for the function to be invoked
 * when data is to be written on server plugin
 */
type PreWriteFunc func(msg MsgFormat) ([]interface{}, error)

/*
 * Function to report processing errors
 */
type ErrorFunc func(msg *MsgFormat, input interface{},
	errStr string) (interface{}, error)

/*
 * Strut to register a plugin and be invoked from the main
 * function. The plugin needs to implement the functions and register
 * with the mail module
 */
type Plugin struct {
	ClientPlugin            ClientSidePlugin         // Plugin function/method to run when client facing
	ClientSideCodec         int                      // Codec to be used for client side
	ClientSidePluginSupport bool                     // Indicates if a client side plugin is available
	EmitClientSideRefConfig func() string            // Function to emit reference config for client side
	EmitServerSideRefConfig func() string            // Function to emit reference config for server side
	FetchErrorFunc          func(interface{}) string // Error function to be used by the plugin
	ParseClientSideConfig   func([]byte) interface{} // Function to parse the yaml input if plugin is client facing
	ParseServerSideConfig   func([]byte) interface{} // Function to parse the yaml input if plugin is client facing
	ServerPlugin            ServerSidePlugin         // Plugin function/method to run when server facing
	ServerSideCodec         int                      // Codec to be used for server side
	ServerSidePluginSupport bool                     // Indicates if a server side plugin is available
}

var PluginMap = make(map[string]Plugin)

/*
 * Interface defining methods to be implemented
 * when a plugin is used to connect to the client.
 * In this case, the broker plugin will act as consumer
 * of the client (i.e server to the client).
 */
type ClientSidePlugin interface {
	/*
	 * Main Method to beCALL_CANCELLATION_RPC run for client side plugin.
	 * codecType - Type of grpc codec to use if needed
	 * params - User defined input
	 * sendChan - To send data to server side plugin.
	 * recvChan - To read data from server side plugin.
	 * termChan - To terminate the client side Go routine
	 * statusChan   - To pass the execution result to the caller of the routine.
	 *                0 would indicate success and 1 would
	 */
	RunClientPlugin(codecType int, params interface{}, readFunc PostReadFunc,
		writeFunc PreWriteFunc, errorFunc ErrorFunc, sendChan chan<- MsgFormat,
		recvChan <-chan MsgFormat, termChan <-chan struct{},
		statusChan chan<- TermStatus)
}

/*
 * Interface defining methods to be implemented
 * when a plugin is used to connect to the server.
 * In thies case, the plugin will be a producer to the
 * server (i.e a client to the server).
 */
type ServerSidePlugin interface {
	/*
	 * Main Method to be run for server side plugin.
	 * codec - Type of grpc codec to use if needed
	 * params - User defined input
	 * sendChan - To send data to client side plugin.
	 * recvChan - To read data from client side plugin.
	 * termChan - To terminate the side side Go routine
	 */
	RunServerPlugin(codecType int, params interface{}, readFunc PostReadFunc,
		writeFunc PreWriteFunc, errorFunc ErrorFunc, sendChan chan<- MsgFormat,
		recvChan <-chan MsgFormat, termChan <-chan struct{},
		statusChan chan<- TermStatus)
}

/*
 * Interface defining methods to be implemented for an
 * implementation. The channels for spawning the channels
 * must be implemented as members of the struct implementing
 * this interface. Naming convention for the struct implementing
 * the interface: <client-side><server-side>Impl; Eg: KafkaGrpcImpl
 */
type Implementation interface {
	// Main method for the implementation
	Run(waitGroup sync.WaitGroup, clientParams interface{},
		serverParams interface{}, clientSide ClientSidePlugin,
		clientReadFunc PostReadFunc, clientWriteFunc PreWriteFunc,
		clientErrorFunc ErrorFunc, serverSide ServerSidePlugin,
		serverReadFunc PostReadFunc, serverWriteFunc PreWriteFunc,
		serverErrorFunc ErrorFunc)

	/*
	 *  Method to stop running the implementation. The method
	 *  should return the client and server side plugin termination
	 *  status
	 */
	Stop() (clientStatus, serverStatus TermStatus)
}

/*
 * Function to get the string equivalent of an RPC type
 */
func GetRpcTypeString(rpcType int) (string, error) {
	switch rpcType {
	case RPC_TYPE_UNARY:
		return "RPC_TYPE_UNARY", nil
	case RPC_TYPE_SERVER_STREAMING:
		return "RPC_TYPE_SERVER_STREAMING", nil
	case RPC_TYPE_CLIENT_STREAMING:
		return "RPC_TYPE_CLIENT_STREAMING", nil
	case RPC_TYPE_BIDISTREAMING:
		return "RPC_TYPE_BIDISTREAMING", nil
	default:
		errorStr := fmt.Sprintf("Unknown RPC type")
		return "", errors.New(errorStr)
	}
}

/**
 * Function to get the current function name for
 * Error reporting
 */
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
