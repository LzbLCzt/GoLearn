package basic

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	pb "GoLearn/web/chapter5/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// todo UnaryInterceptor 用于拦截每次的请求调用，可以打印请求的开始时间、结束时间、请求参数、响应参数、错误信息等信息
func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// 打印请求的基本信息
	start := time.Now()
	log.Printf("Unary Request - Method: %s, StartTime: %s", info.FullMethod, start.Format(time.RFC3339))

	// 如果需要，可以从上下文中提取 metadata（如 headers）
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("Metadata: %v", md)
	}

	// 处理请求
	resp, err := handler(ctx, req)

	// 打印响应时间及错误信息
	duration := time.Since(start)
	if err != nil {
		log.Printf("Unary Request - Method: %s, Duration: %s, Error: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("Unary Request - Method: %s, Duration: %s, Success", info.FullMethod, duration)
	}

	return resp, err
}

func TestGrpcServerWithInterceptor(t *testing.T) {
	// 创建监听器
	listen, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf("Listen error: %v\n", err)
	}
	fmt.Printf("Listening on %s\n", ":8088")

	// todo 创建 gRPC 服务并添加拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor), // 添加 Unary 拦截器
	)

	// 注册服务
	pb.RegisterProgrammerServiceServer(s, &ProgrammerServiceServer2{})

	// 启动服务
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("gRPC service serve error: %v\n", err)
	}
}

type ProgrammerServiceServer2 struct{}

func (p *ProgrammerServiceServer2) GetProgrammerInfo(ctx context.Context, req *pb.Request) (resp *pb.Response, err error) {
	name := req.Name
	if name == "zhengbangli" {
		resp = &pb.Response{
			Uid:      6,
			Username: name,
			Job:      "CTO",
			GoodAt:   []string{"Go", "Java", "PHP", "Python"},
		}
	} else {
		resp = &pb.Response{}
	}
	err = nil
	return
}