package discover_proxy

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
)

func TestDiscover(t *testing.T) {

	rsp, err := DefaultBackend.discoverWithGrpc(context.Background(), "polaris.discover",
		"polaris.report", "Polaris", apiV1Model.DiscoverRequest_INSTANCE)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("instances: %+v", rsp.GetInstances())
}

func TestGetOneInstance(t *testing.T) {
	req := &api.GetOneInstanceRequest{}
	req.Namespace = "Polaris"
	req.Service = "polaris.report"

	resp, err := DefaultBackend.getOneInstance(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ins := resp.Instances[0]
	fmt.Printf("ins: %v", ins)
}

func TestGetInstance(t *testing.T) {
	req := &api.GetInstancesRequest{}
	req.Namespace = "Polaris"
	req.Service = "polaris.report"

	for {
		resp, err := DefaultBackend.getInstance(req)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		instances := resp.Instances
		fmt.Printf("len of ins: %d\n", len(instances))
		//for _, ins := range instances {
		//	fmt.Printf("ins: %v", ins.GetCircuitBreakerStatus())
		//}
		time.Sleep(1 * time.Second)
	}

}

func TestGrpcDial(t *testing.T) {

}
