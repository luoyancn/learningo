package rpclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	registry "grpcetcdv3/rpc/etcdv3"
	"grpcetcdv3/rpc/messages"

	etcd "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
)

type grpcPool struct {
	conn chan *grpc.ClientConn
	addr string
}

var gonce sync.Once
var pool *grpcPool

func InitGrpcClientPool() {
	gonce.Do(func() {
		pool = new(grpcPool)
		pool.conn = make(chan *grpc.ClientConn, 1024)
		conn := pool.dialNew()
		pool.conn <- conn
	})
}

func (this *grpcPool) dialNew() *grpc.ClientConn {
	var err error
	var conn *grpc.ClientConn
	config := etcd.Config{
		Endpoints: []string{"http://localhost:2379"},
	}
	resolver := registry.NewResolver("grpc-lb", "zhangjl", config)
	balancer_etcd := grpc.RoundRobin(resolver)
	conn, err = grpc.Dial("", grpc.WithInsecure(),
		grpc.WithBalancer(balancer_etcd), grpc.WithBlock())
	if nil != err {
		return nil
	}
	return conn
}

func (this *grpcPool) get() *grpc.ClientConn {
	select {
	case conn := <-this.conn:
		return conn
	default:
		return this.dialNew()
	}
}

func (this *grpcPool) put(conn *grpc.ClientConn) error {
	select {
	case this.conn <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func Call() string {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := messages.NewReQRePClient(conn)
	resp, err := client.Call(ctx, &messages.Request{Req: "luoyan"})
	if nil != err {
		fmt.Printf("Failed to get response from grpc server:%v\n", err)
		return ""
	}
	return resp.GetResp()
}
