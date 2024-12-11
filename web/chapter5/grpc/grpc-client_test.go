package grpc

import (
	pb "GoLearn/web/chapter5/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestGrpcClient(t *testing.T) {
	//conn, err := grpc.Dial(":8088", grpc.WithInsecure())
	conn, err := grpc.Dial("9.134.117.127:8088", grpc.WithInsecure())	//远程调用个人开发机上部署的服务
	if err != nil {
		log.Fatalf("dial error: %v\n", err)
	}

	defer conn.Close()

	client := pb.NewProgrammerServiceClient(conn)
	req := new(pb.Request)
	req.Name = "zhengbangli"
	//req.Name = "aaa"
	resp, err := client.GetProgrammerInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("get programmer info error: %v\n", err)
	}

	fmt.Printf("Recevied: %v\n", resp)
}