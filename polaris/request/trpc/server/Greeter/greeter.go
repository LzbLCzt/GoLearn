package Greeter

import (
	greeter "GoLearn/polaris/request/trpc/api/Greeter"
	"context"
	"fmt"
)

// 服务实现
type GreeterServer struct {
}

func (s *GreeterServer) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	message := fmt.Sprintf("Hello, %s! Welcome to tRPC!", req.Name)
	fmt.Println(message)
	return &greeter.HelloReply{Message: message}, nil
}
