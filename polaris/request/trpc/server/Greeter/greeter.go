package Greeter

import (
	"context"
	"fmt"
	"github.com/LzbLCzt/GoLearn/polaris/request/trpc/Greeter/helloworld.pb"
)

// 服务实现
type greeterServer struct {
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    message := fmt.Sprintf("Hello, %s! Welcome to tRPC!", req.Name)
    return &pb.HelloReply{Message: message}, nil
})
