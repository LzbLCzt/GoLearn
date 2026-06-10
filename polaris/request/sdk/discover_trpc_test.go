package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"testing"

	grpcpb "git.woa.com/polaris/polaris-server-api/api/v1/grpc"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDiscoverTrpc(t *testing.T) {
	req := &apiV1Model.DiscoverRequest{
		Type: apiV1Model.DiscoverRequest_INSTANCE,
		Service: &apiV1Model.Service{
			Namespace: &wrapperspb.StringValue{Value: "Development"},
			Name:      &wrapperspb.StringValue{Value: "personal_test_polaris"},
		},
	}

	// 使用 gRPC 客户端连接 8091 端口
	conn, err := grpc.Dial("30.163.16.103:8091", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial err: %v", err)
	}
	defer conn.Close()

	client := grpcpb.NewPolarisGRPCClient(conn)
	discoverClient, err := client.Discover(context.Background())
	if err != nil {
		t.Fatalf("create discover client err: %v", err)
	}

	err = discoverClient.Send(req)
	if err != nil {
		t.Fatalf("send err: %v", err)
	}

	rsp, err := discoverClient.Recv()
	if err != nil {
		t.Fatalf("recv err: %v", err)
	}

	fmt.Printf("rsp: %+v", rsp)
}


func TestAA(t *testing.T) {

}

func IPV4LegalCheck(ip string) error {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return err
	}
	if addr.Is6() {
		return errors.New("not support ipv6")
	}
	return nil
}