package rpc

import (
	"fmt"
	"io"
	"net"
	"oceanstack/conf"
	"oceanstack/logging"
	"sync"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	netutil "golang.org/x/net/netutil"
	"golang.org/x/time/rate"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
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
	ticker := time.After(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker:
				logging.LOG.Noticef("Time out after 30 second\n")
				return
			default:
				logging.LOG.Infof("waiting for timeout\n")
				time.Sleep(time.Second * 1)
			}
		}
	}()
	return &empty.Empty{}, nil
}

func (this *Server) StreamResponse(req *Request,
	s_server ReQReP_StreamResponseServer) error {
	s_server.Send(&Response{Resp: "zhangzz"})
	return nil
}

func (this *Server) StreamRequest(s_server ReQReP_StreamRequestServer) error {
	for {
		_, err := s_server.Recv()
		if io.EOF == err {
			s_server.SendAndClose(&Response{Resp: "lalala"})
			return nil
		}
		if nil != err {
			return err
		}
	}
	return nil
}

func (this *Server) StreamReqRep(s_server ReQReP_StreamReqRepServer) error {
	for {
		_, err := s_server.Recv()
		if io.EOF == err {
			return nil
		}
		if nil != err {
			logging.LOG.Errorf("Failed to recive the message :%v\n", err)
			return err
		}
		s_server.Send(&Response{Resp: "good night"})
	}
	return nil
}

/****************************************************************************/
var once sync.Once
var GRPC *grpc.Server

func withServerInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	logging.LOG.Infof("invoke server method=%s duration=%s error=%v",
		info.FullMethod, time.Since(start), err)
	return resp, err
}

func streamServerInterceptor(srv interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, ss)
	logging.LOG.Infof("invoke server stream method=%s duration=%s error=%v",
		info.FullMethod, time.Since(start), err)
	return err
}

type tapp struct {
	lim *rate.Limiter
}

// 限制访问频率。默认1024次/秒
func newTap() *tapp {
	return &tapp{rate.NewLimiter(rate.Limit(conf.GRPC_SERVER_REQ_MAX_FREQUENCY),
		conf.GRPC_SERVER_REQ_BURST_FREQUENCY)}
}
func (t *tapp) handler(ctx context.Context,
	info *tap.Info) (context.Context, error) {
	if !t.lim.Allow() {
		return nil, status.Errorf(
			codes.ResourceExhausted, "Service is over rate limit")
	}
	return ctx, nil
}

func StartServer() {
	once.Do(func() {
		listener, err := net.Listen(
			"tcp", fmt.Sprintf("%s:%d", "0.0.0.0", conf.GRPC_PORT))
		if nil != err {
			logging.LOG.Fatalf(
				"Failed to start grpc server on 0.0.0.0:%d, %v\n",
				conf.GRPC_PORT, err)
		}
		// grpc.MaxConcurrentStreams限定每个grpc连接可以有多少个并发
		GRPC = grpc.NewServer(withServerInterceptor(),
			grpc.StreamInterceptor(streamServerInterceptor),
			grpc.MaxConcurrentStreams(uint32(conf.GRPC_CONCURRENCY)),
			grpc.InTapHandle(newTap().handler))
		RegisterReQRePServer(GRPC, &Server{})
		// netutil.LimitListener限定总共可以对外提供多少连接
		limit_lister := netutil.LimitListener(
			listener, conf.GRPC_SERVER_CONN_LIMITS)
		GRPC.Serve(limit_lister)
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

func withClientInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

func clientInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	logging.LOG.Infof("invoke remote method=%s duration=%s error=%v",
		method, time.Since(start), err)
	return err
}

func clientstreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	opts ...grpc.CallOption) (grpc.ClientStream, error) {
	start := time.Now()
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	logging.LOG.Infof("invoke remote stream method=%s duration=%s error=%v",
		method, time.Since(start), err)
	return clientStream, err
}

func (this *grpcPool) dialNew() *grpc.ClientConn {
	// grpc.MaxCallSendMsgSize 设置客户端最大可以发送的消息体大小。默认为1M.
	conn, err := grpc.Dial(this.addr, grpc.WithInsecure(),
		grpc.WithInsecure(), withClientInterceptor(),
		grpc.WithStreamInterceptor(clientstreamClientInterceptor),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(conf.GRPC_REQ_MSG_SIZE)))
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

func GrpcCast() {
	ctx, cancle := context.WithTimeout(context.Background(), conf.GRPC_TIMEOUT)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := NewReQRePClient(conn)
	_, err := client.Cast(ctx, &Request{Req: "zhangjl"})
	if nil != err {
		logging.LOG.Errorf("Grpc Cast Error: %v\n", err)
	}
}

func GrpcStreamResponse() {
	ctx, cancle := context.WithTimeout(context.Background(), conf.GRPC_TIMEOUT)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := NewReQRePClient(conn)

	stream, err := client.StreamResponse(ctx, &Request{Req: "chidoli"})
	if nil != err {
		logging.LOG.Errorf("Request the stream request faild :%v\n", err)
		return
	}

	for {
		reply, err := stream.Recv()
		if io.EOF == err {
			break
		}
		if nil != err {
			logging.LOG.Errorf("Failed to recive the message :%v\n", err)
			break
		}
		logging.LOG.Infof("Recived the message:%v\n", reply.Resp)
	}
	return
}

func GrpcStreamRequest() {
	ctx, cancle := context.WithTimeout(context.Background(), conf.GRPC_TIMEOUT)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := NewReQRePClient(conn)

	stream, err := client.StreamRequest(ctx)
	if nil != err {
		logging.LOG.Errorf("Failed to connect the stream server:%v\n", err)
		return
	}
	stream.Send(&Request{Req: "hahahha"})
	reply, err := stream.CloseAndRecv()
	if nil != err {
		logging.LOG.Errorf("Failed to recive the message :%v\n", err)
		return
	}

	logging.LOG.Infof("Recived the message:%v\n", reply.Resp)
	return
}

func GrpcStreamRequestResponse() {
	ctx, cancle := context.WithTimeout(context.Background(), conf.GRPC_TIMEOUT)
	defer cancle()
	conn := pool.get()
	defer pool.put(conn)
	client := NewReQRePClient(conn)

	stream, err := client.StreamReqRep(ctx)
	if nil != err {
		logging.LOG.Errorf("Failed to connect the stream server:%v\n", err)
		return
	}
	stream.Send(&Request{Req: "good evening"})
	reply, err := stream.Recv()
	if nil != err {
		logging.LOG.Errorf("Failed to recive the message :%v\n", err)
		return
	}

	logging.LOG.Infof("Recived the message:%v\n", reply.Resp)
	return
}
