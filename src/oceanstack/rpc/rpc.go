package rpc

import (
	"fmt"
	"net"
	"oceanstack/conf"
	"oceanstack/logging"
	"sync"

	empty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Server struct {
}

func (this *Server) GetResp(
	ctx context.Context, req *Request) (*Response, error) {
	return &Response{
		Resp: "lucifer",
	}, nil
}

func (this *Server) Cast(ctx context.Context,
	req *Request) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

/****************************************************************************/
var once sync.Once
var GRPC *grpc.Server

func StartServer() {
	once.Do(func() {
		listener, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%d", "0.0.0.0", conf.GRPC_PORT))
		if nil != err {
			logging.LOG.Fatalf(
				"Failed to start grpc server on 0.0.0.0:%d, %v\n",
				conf.GRPC_PORT, err)
		}
		GRPC = grpc.NewServer()
		RegisterReQRePServer(GRPC, &Server{})
		GRPC.Serve(listener)
		logging.LOG.Infof("Grpc Server started on 0.0.0.0:%d\n",
			conf.GRPC_PORT)
	})
}

func StopServer() {
	logging.LOG.Infof("Terminate the Grpc Server ...\n")
	GRPC.Stop()
}

/****************************************************************************/
type grpcPool struct {
	conn chan *grpc.ClientConn
	addr string
}

var gonce sync.Once
var pool *grpcPool

func InitGrpcClientPool() {
	gonce.Do(func() {
		pool = new(grpcPool)
		pool.addr = fmt.Sprintf("%s:%d", conf.GRPC_SERVER, conf.GRPC_PORT)
		pool.conn = make(chan *grpc.ClientConn, conf.GRPC_POOL_SIZE)
		conn := pool.dialNew()
		pool.conn <- conn
	})
}

func (this *grpcPool) dialNew() *grpc.ClientConn {
	conn, err := grpc.Dial(this.addr, grpc.WithInsecure(), grpc.WithInsecure())
	if nil != err {
		logging.LOG.Errorf("RPC ERROR:%v\n", err)
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

func GrpcClient() *Response {
	ctx, cancle := context.WithTimeout(context.Background(), conf.GRPC_TIMEOUT)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := NewReQRePClient(conn)
	resp, err := client.GetResp(ctx, &Request{Req: "luoyan"})
	if nil != err {
		logging.LOG.Errorf("RPC ERROR:%v\n", err)
		return nil
	}
	return resp
}
