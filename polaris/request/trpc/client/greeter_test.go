package main

import (
	pb "GoLearn/polaris/request/trpc/api/Greeter"
	"context"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go/client"
	"log"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	// 配置客户端
	opts := []client.Option{
		client.WithTarget("ip://127.0.0.1:8000"),
		client.WithTimeout(3 * time.Second),
	}

	// 创建客户端代理
	greeterClient := pb.NewGreeterClientProxy(opts...)

	// 准备请求
	req := &pb.HelloRequest{
		Name: "World",
	}

	// 发送请求
	ctx := context.Background()
	rsp, err := greeterClient.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	fmt.Printf("Server response: %s\n", rsp.Message)

	// 测试多个调用
	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		req.Name = name
		rsp, err := greeterClient.SayHello(ctx, req)
		if err != nil {
			log.Printf("Error calling for %s: %v", name, err)
			continue
		}
		fmt.Printf("Response for %s: %s\n", name, rsp.Message)
	}
}
