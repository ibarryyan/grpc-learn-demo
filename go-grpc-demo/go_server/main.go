package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"go-grpc-demo/proto/hello_proto"
	
	"google.golang.org/grpc"
)

type HelloServer struct{}

func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

func (s *HelloServer) SayHello(ctx context.Context, req *hello_proto.HelloRequest) (*hello_proto.HelloResponse, error) {
	resp := &hello_proto.HelloResponse{}

	defer func() {
		log.Printf("HelloServer SayHello req = %v ,resp = %v", req, resp)
	}()

	if req.GetName() == "" {
		return resp, errors.New("request name is not nil")
	} else {
		resp.Message = fmt.Sprintf("Hello %s", req.Name)
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatalf("端口监听错误 : %v\n", err)
	}
	srv := grpc.NewServer()
	hello_proto.RegisterHelloServiceServer(srv, NewHelloServer())
	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("端口监听错误 : %v\n", err)
	}
}
