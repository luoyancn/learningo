package rpc

import (
	"fmt"
	"net"
	"oceanstack/conf"
	"oceanstack/logging"
	"runtime"
	"sync"

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

/****************************************************************************/
var once sync.Once
var GRPC *grpc.Server

func StartServer() {
	once.Do(func() {
		runtime.GOMAXPROCS(conf.RPC_WORKERS)
		listener, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%d", "0.0.0.0", conf.RPC_PORT))
		if nil != err {
			logging.LOG.Fatalf(
				"Failed to start grpc server on 0.0.0.0:%d, %v\n",
				conf.RPC_PORT, err)
		}
		GRPC = grpc.NewServer()
		RegisterReQRePServer(GRPC, &Server{})
		logging.LOG.Infof("Grpc Server start on 0.0.0.0:%d\n", conf.RPC_PORT)
		GRPC.Serve(listener)
	})
}

func StopServer() {
	logging.LOG.Infof("Terminate the Grpc Server ...\n")
	GRPC.Stop()
}

func GrpcClient() *Response {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		conf.RPC_SERVER, conf.RPC_PORT), grpc.WithInsecure())
	if nil != err {
		logging.LOG.Errorf("Failed to connect grpc server on %s:%d, %v\n",
			conf.RPC_SERVER, conf.RPC_PORT, err)
		return nil
	}
	defer conn.Close()
	client := NewReQRePClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), conf.RPC_TIMEOUT)
	defer cancle()
	resp, err := client.GetResp(ctx, &Request{Req: "luoyan"})
	if nil != err {
		logging.LOG.Errorf("RPC ERROR:%v\n", err)
		return nil
	}
	return resp
}
