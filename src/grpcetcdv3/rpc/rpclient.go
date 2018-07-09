package rpc

import (
	"context"
	"sync"
	"time"

	"grpcetcdv3/rpc/messages"

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
	conn, err = grpc.Dial("localhost:8080", grpc.WithInsecure())
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
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := messages.NewReQRePClient(conn)
	resp, err := client.Call(ctx, &messages.Request{Req: "luoyan"})
	if nil != err {
		panic(err)
	}
	return resp.GetResp()
}
