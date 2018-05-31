/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package implementations

import (
	"fmt"
	"log"
	"strings"
	"sync"

	wu "github.com/Juniper/wedge/util"
)

//"github.com/confluentinc/confluent-kafka-go/kafka"

type WedgeImpl struct {
	servertoclientPluginChan chan wu.MsgFormat
	clienttoserverPluginChan chan wu.MsgFormat
	termChan                 chan struct{}
	clientStatusChan         chan wu.TermStatus
	serverStatusChan         chan wu.TermStatus
	implTermChan             chan struct{}
	implTermStatusChan       chan wu.TermStatus
	ClientCodec              int
	ServerCodec              int
}

func (impl *WedgeImpl) Run(waitGroup sync.WaitGroup,
	clientParams interface{}, serverParams interface{},
	clientSide wu.ClientSidePlugin, clientReadFunc wu.PostReadFunc,
	clientWriteFunc wu.PreWriteFunc, clientErrorFunc wu.ErrorFunc,
	serverSide wu.ServerSidePlugin, serverReadFunc wu.PostReadFunc,
	serverWriteFunc wu.PreWriteFunc, serverErrorFunc wu.ErrorFunc) {

	var status, clientStatus, serverStatus wu.TermStatus
	var clientTerminated, serverTerminated bool

	// Increment the WaitGroup counter.
	defer waitGroup.Done()

	impl.servertoclientPluginChan = make(chan wu.MsgFormat, 1000)
	impl.clienttoserverPluginChan = make(chan wu.MsgFormat, 1000)
	impl.termChan = make(chan struct{}, 2)
	impl.clientStatusChan = make(chan wu.TermStatus, 2)
	impl.serverStatusChan = make(chan wu.TermStatus, 2)

	impl.implTermChan = make(chan struct{}, 2)
	impl.implTermStatusChan = make(chan wu.TermStatus, 2)

	if clientSide == nil {
		status.ErrStr = fmt.Sprintf("%s: Client side plugin is nil",
			wu.FuncName())
		status.Status = wu.ROUTINE_ERROR_EXIT
		impl.implTermStatusChan <- status
	}

	if serverSide == nil {
		status.ErrStr = fmt.Sprintf("%s: Server side plugin is nil",
			wu.FuncName())
		status.Status = wu.ROUTINE_ERROR_EXIT
		impl.implTermStatusChan <- status
		return
	}

	if clientParams == nil {
		status.ErrStr = fmt.Sprintf("%s: Client side parameters are not "+
			"provided", wu.FuncName())
		status.Status = wu.ROUTINE_ERROR_EXIT
		impl.implTermStatusChan <- status
		return
	}

	go clientSide.RunClientPlugin(impl.ClientCodec, clientParams,
		clientReadFunc, clientWriteFunc, clientErrorFunc, impl.clienttoserverPluginChan,
		impl.servertoclientPluginChan, impl.termChan, impl.clientStatusChan)

	if serverParams == nil {
		status.ErrStr = fmt.Sprintf("%s: Server side parameters are not "+
			"provided", wu.FuncName())
		status.Status = wu.ROUTINE_ERROR_EXIT
		impl.implTermStatusChan <- status
		return
	}

	go serverSide.RunServerPlugin(impl.ServerCodec, serverParams,
		serverReadFunc, serverWriteFunc, serverErrorFunc,
		impl.servertoclientPluginChan, impl.clienttoserverPluginChan,
		impl.termChan, impl.serverStatusChan)

	for {
		select {
		case <-impl.implTermChan:
			close(impl.termChan)
		case clientStatus, _ = <-impl.clientStatusChan:
			log.Printf("%s: Client terminated with status %v", wu.FuncName(),
				clientStatus)
			clientTerminated = true
			status.ErrStr = "Client Termination Status Info: " + status.ErrStr
			impl.implTermStatusChan <- status
		case serverStatus, _ = <-impl.serverStatusChan:
			log.Printf("%s: Server terminated with status %v", wu.FuncName(),
				serverStatus)
			serverTerminated = true
			status.ErrStr = "Client Termination Status Info: " + status.ErrStr
			impl.implTermStatusChan <- status
		default:
			if clientTerminated == true && serverTerminated == true {
				close(impl.implTermStatusChan)
			}

		}
	}
}

func (impl *WedgeImpl) Stop() (clientStatus, serverStatus wu.TermStatus) {
	close(impl.implTermChan)

	for {
		select {
		case status, ok := <-impl.implTermStatusChan:
			if !ok {
				log.Printf("%s: Error Fetching termination status",
					wu.FuncName())
			} else {
				if strings.Contains(status.ErrStr, "Client") {
					clientStatus = status
				} else {
					serverStatus = status
				}
			}

		}
	}
}
