package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/gunsluo/go-example/multiplex_port/pb"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
)

type exampleHTTPHandler struct{}

func (h *exampleHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example http response")
}

func serveHTTP(l net.Listener) {
	s := &http.Server{
		Handler: &exampleHTTPHandler{},
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func EchoServer(ws *websocket.Conn) {
	if _, err := io.Copy(ws, ws); err != nil {
		panic(err)
	}
}

func serveWS(l net.Listener) {
	s := &http.Server{
		Handler: websocket.Handler(EchoServer),
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

type ExampleRPCRcvr struct{}

func (r *ExampleRPCRcvr) Cube(i int, j *int) error {
	*j = i * i
	return nil
}

func serveRPC(l net.Listener) {
	s := rpc.NewServer()
	if err := s.Register(&ExampleRPCRcvr{}); err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			if err != cmux.ErrListenerClosed {
				panic(err)
			}
			return
		}
		go s.ServeConn(conn)
	}
}

type grpcServer struct{}

func (s *grpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (
	*pb.HelloReply, error) {

	return &pb.HelloReply{Message: "Hello " + in.Name + " from cmux"}, nil
}

func serveGRPC(l net.Listener) {
	grpcs := grpc.NewServer()
	pb.RegisterGreeterServer(grpcs, &grpcServer{})
	if err := grpcs.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Panic(err)
	}

	m := cmux.New(l)

	// We first match the connection against HTTP2 fields. If matched, the
	// connection will be sent through the "grpcl" listener.
	grpcl := m.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	//Otherwise, we match it againts a websocket upgrade request.
	wsl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// Otherwise, we match it againts HTTP1 methods. If matched,
	// it is sent through the "httpl" listener.
	httpl := m.Match(cmux.HTTP1Fast())
	// If not matched by HTTP, we assume it is an RPC connection.
	rpcl := m.Match(cmux.Any())

	// Then we used the muxed listeners.
	go serveGRPC(grpcl)
	go serveWS(wsl)
	go serveHTTP(httpl)
	go serveRPC(rpcl)

	logrus.Info("Listen Port :50051")
	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
