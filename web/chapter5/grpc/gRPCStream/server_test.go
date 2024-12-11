package gRPCStream

import (
	"GoLearn/web/chapter5/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	listen, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	proto.RegisterHelloServiceServer(s, &HelloServiceImpl{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("server start err:%v", err)
	}
}