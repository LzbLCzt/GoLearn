package discover_proxy

import (
	"context"

	grpcpb "git.woa.com/polaris/polaris-server-api/api/v1/grpc"
	grpcnet "google.golang.org/grpc"
)

type GRPCClient struct {
	client grpcpb.PolarisGRPC_DiscoverClient
	conn   *grpcnet.ClientConn
}

// 连接discover proxy
func connectToProxy(address string) (*grpcnet.ClientConn, grpcpb.PolarisGRPC_DiscoverClient, error) {
	clientConn, err := grpcnet.Dial(address, grpcnet.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := grpcpb.NewPolarisGRPCClient(clientConn)
	discoverClient, err := client.Discover(context.Background())
	if err != nil {
		clientConn.Close()
		return nil, nil, err
	}
	return clientConn, discoverClient, nil
}
