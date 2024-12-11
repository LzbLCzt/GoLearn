package gRPCStream

import (
	"GoLearn/web/chapter5/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
	client := proto.NewHelloServiceClient(conn)
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	//发送消息给server
	go func() {
		for {
			if err := stream.Send(&proto.String{Value: "this is zhengbangli"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	//接受server消息
	for {
		reply, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("response: %s\n", reply.GetValue())
	}
}