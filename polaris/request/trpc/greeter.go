// server/main.go
package main

import (
	"context"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go"
	"log"
)

func main() {
	// 创建服务器实例
	s := trpc.NewServer()

	// 注册服务实现
	pb.RegisterGreeterService(s, &greeterServer{})

	// 启动服务
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

// 服务实现
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	message := fmt.Sprintf("Hello, %s! Welcome to tRPC!", req.Name)
	fmt.Printf("Received request for: %s\n", req.Name)

	return &pb.HelloReply{
		Message: message,
	}, nil
}
