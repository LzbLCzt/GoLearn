package sdk

import (
	"testing"

	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
	"git.woa.com/polaris/polaris-go/v2/pkg/model"
)

func TestRuleRouteSdk(t *testing.T) {
	consumer, err := api.NewConsumerAPI()
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}

	defer consumer.Destroy()

	var req *api.GetOneInstanceRequest
	req = &api.GetOneInstanceRequest{}

	// 设置负载均衡算法
	req.LbPolicy = config.DefaultLoadBalancerMaglev

	// 设置被调服务信息
	req.Namespace = "Test"
	req.Service = "lzb_test"

	// 设置主调服务信息
	req.SourceService = &model.ServiceInfo{
		Namespace: "Test",
		Service:   "lzb_test2",
		Metadata: map[string]string{
			"k3": "v3",
		},
	}

	resp, err := consumer.GetOneInstance(req)
	if err != nil {
		t.Errorf("fail to get instance, err is %v", err)
		return
	}

	instance := resp.GetInstances()[0]
	t.Logf("instance: %v", instance.GetInstanceKey())
}
