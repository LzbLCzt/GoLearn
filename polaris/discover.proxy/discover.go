package discover_proxy

import (
	"context"
	"fmt"
	"git.code.oa.com/polaris/polaris-go/api"
	"git.code.oa.com/trpc-go/trpc-go/client"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"git.woa.com/polaris/polaris-server-api/api/v1/trpc"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"time"
)

func init() {
	err := DefaultBackend.InitPolarisByYaml("polaris.yaml")
	if err != nil {
		panic(err)
	}
}

var DefaultBackend = &backend{
	polarisTRPCClientProxy: trpc.NewPolarisTRPCClientProxy(),
}

type backend struct {
	polarisTRPCClientProxy trpc.PolarisTRPCClientProxy
	consumer               api.ConsumerAPI
}

func (b *backend) discover(ctx context.Context, backendService, serviceName, namespace string,
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
