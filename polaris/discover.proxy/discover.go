package discover_proxy

import (
	"context"
	"fmt"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/client"
	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/model"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"git.woa.com/polaris/polaris-server-api/api/v1/trpc"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

func init() {
	err := DefaultBackend.InitPolarisByYaml("polaris.yaml")
	if err != nil {
		panic(err)
	}
	proxyConnection, proxyDiscoverClient, err := connectToProxy("9.141.112.135:8081")
	if err != nil {
		panic(err)
	}

	DefaultBackend.gClient = &GRPCClient{
		conn:   proxyConnection,
		client: proxyDiscoverClient,
	}

}

var DefaultBackend = &backend{
	polarisTRPCClientProxy: trpc.NewPolarisTRPCClientProxy(),
}

type backend struct {
	polarisTRPCClientProxy IDiscover
	consumer               api.ConsumerAPI
	gClient                *GRPCClient
}

func (b *backend) getOneInstance(req *api.GetOneInstanceRequest) (*model.InstancesResponse, error) {
	return b.consumer.GetOneInstance(req)
}

// trpc服务对应的端口非8080或8081， 开发环境网络不通
func (b *backend) discover(ctx context.Context, backendService, serviceName, namespace string,
	reqTyp apiV1Model.DiscoverRequest_DiscoverRequestType) (*apiV1Model.DiscoverResponse, error) {

	req := api.GetOneInstanceRequest{}
	req.Namespace = "Polaris"
	//req.Service = backendService
	req.Service = "polaris.discover"
	rsp, err := b.consumer.GetOneInstance(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get instance for service %s", backendService)
	}

	if len(rsp.Instances) == 0 {
		return nil, errors.Errorf("instance not found for service %s", backendService)
	}

	instance := rsp.Instances[0]

	fmt.Printf("instances: +%v\n", rsp.Instances)

	discoverReq := &apiV1Model.DiscoverRequest{
		Type: reqTyp,
		Service: &apiV1Model.Service{
			Namespace: &wrappers.StringValue{Value: namespace},
			Name:      &wrappers.StringValue{Value: serviceName},
		},
	}

	address := instance.GetHost() + ":" + fmt.Sprintf("%d", instance.GetPort())
	dialCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	target := fmt.Sprintf("ip://%s", address)
	fmt.Printf("target: %s\n", target)
	trpcClient := client.WithTarget(target)

	discoverRsp, err := b.polarisTRPCClientProxy.Discover(dialCtx, discoverReq, trpcClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to discover service %s", serviceName)
	}

	if discoverRsp.GetCode().GetValue() != apiV1Model.ExecuteSuccess {
		return nil, errors.Errorf("failed to discover service %s, rsp code: %d", serviceName, discoverRsp.GetCode().GetValue())
	}

	return discoverRsp, nil
}

func (b *backend) InitPolarisByYaml(path string) error {
	if b.consumer != nil {
		return nil
	}
	var err error
	b.consumer, err = api.NewConsumerAPIByFile(path)
	if err != nil {
		fmt.Printf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}

// 改用grpc
func (b *backend) discoverWithGrpc(ctx context.Context, backendService, serviceName, namespace string,
	reqTyp apiV1Model.DiscoverRequest_DiscoverRequestType) (*apiV1Model.DiscoverResponse, error) {
	req := api.GetOneInstanceRequest{}
	req.Namespace = "Polaris"
	req.Service = backendService
	rsp, err := b.consumer.GetOneInstance(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get instance for service %s", backendService)
	}

	if len(rsp.Instances) == 0 {
		return nil, errors.Errorf("instance not found for service %s", backendService)
	}

	instance := rsp.Instances[0]

	fmt.Printf("instances: +%v\n", instance)

	request := &apiV1Model.DiscoverRequest{
		Service: &apiV1Model.Service{
			Namespace: &wrappers.StringValue{Value: namespace},
			Name:      &wrappers.StringValue{Value: serviceName},
		},
		Type: reqTyp,
	}
	err = b.gClient.client.Send(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to discover service %s", serviceName)
	}

	recv, err := b.gClient.client.Recv()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to discover service %s", serviceName)
	}

	if recv.GetCode().GetValue() != apiV1Model.ExecuteSuccess {
		return nil, errors.Errorf("failed to discover service %s, rsp code: %d", serviceName, recv.GetCode().GetValue())
	}

	return recv, nil
}
