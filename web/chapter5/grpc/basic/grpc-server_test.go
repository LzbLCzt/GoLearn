package basic

import (
	pb "GoLearn/web/chapter5/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	fmt.Printf("listen %s\n", "8088")
	s := grpc.NewServer()
	pb.RegisterProgrammerServiceServer(s, &ProgrammerServiceServer{})
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("grpc service serve error: %v\n", err)
	}
}

type ProgrammerServiceServer struct{}

func (p *ProgrammerServiceServer) GetProgrammerInfo(ctx context.Context, req *pb.Request) (resp *pb.Response, err error) {
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