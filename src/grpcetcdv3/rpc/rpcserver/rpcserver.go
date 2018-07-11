package rpcserver

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	registry "grpcetcdv3/rpc/etcdv3"
	"grpcetcdv3/rpc/messages"

	etcd "github.com/coreos/etcd/clientv3"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

var once sync.Once
var _grpc *grpc.Server

func StartServer(port int) {
	once.Do(func() {
		fmt.Printf("Start the grpc server\n")
		listener, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%d", "0.0.0.0", port))
		if nil != err {
			panic(err)
		}

		config := etcd.Config{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		_uuid, _ := uuid.NewV4()
		opt := registry.Option{
			NData:       fmt.Sprintf("%s:%d", "192.168.137.30", port),
			RegistryDir: "grpc-lb",
			ServiceName: "zhangjl",
			NodeID:      _uuid.String(),
		}
		registry.Register(config, opt)

		_grpc = grpc.NewServer()
		host_name, _ := os.Hostname()
		messages.RegisterReQRePServer(_grpc,
			&messages.Service{HostName: host_name, ListenPort: port})
		_grpc.Serve(listener)
	})
}

func StopServer() {
	registry.UnRegister()
	if nil != _grpc {
		_grpc.Stop()
	}
}
