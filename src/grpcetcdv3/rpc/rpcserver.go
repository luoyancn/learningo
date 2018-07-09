package rpc

import (
	"fmt"
	"net"
	"sync"

	"grpcetcdv3/rpc/messages"

	"google.golang.org/grpc"
)

var once sync.Once
var _grpc *grpc.Server

func StartServer(port int) {
	once.Do(func() {
		listener, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%d", "0.0.0.0", port))
		if nil != err {
			panic(err)
		}
		_grpc = grpc.NewServer()
		messages.RegisterReQRePServer(_grpc, &messages.Service{})
		_grpc.Serve(listener)
	})
}

func StopServer() {
	if nil != _grpc {
		_grpc.Stop()
	}
}
