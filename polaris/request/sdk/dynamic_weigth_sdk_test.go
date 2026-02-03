package sdk

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/config"
)

var ConsumerV2 api.ConsumerAPI
var ProviderV2 api.ProviderAPI

func InitPolarisProviderByYamlV2(path string) error {
	if Provider != nil {
		return nil
	}
	var err error
	ProviderV2, err = api.NewProviderAPIByFile(path)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}

func InitPolarisByYamlV2(path string) error {
	if Consumer != nil {
		return nil
	}
	var err error
	ConsumerV2, err = api.NewConsumerAPIByFile(path)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}

func TestDynamic_weight_Consumer(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err := InitPolarisByYamlV2("polaris_dynamic_weight.yaml")
	if err != nil {
		t.Fatalf("init polaris consumer error")
	}

	req := api.GetOneInstanceRequest{}
	req.Service = "lzb_test2"
	req.Namespace = "Test"
	req.LbPolicy = config.DefaultLoadBalancerDWR

	err = InitPolarisProviderByYamlV2("polaris_dynamic_weight.yaml")
	if err != nil {
		t.Fatalf("init polaris provider error")
	}
	go reportDynamicWeight(t, ctx)
	time.Sleep(5 * time.Second)

	maps := make(map[string]int)
	for i := 0; i < 1000; i++ {
		rsp, err := ConsumerV2.GetOneInstance(&req)
		if err != nil {
			t.Fatalf("getOneInstance fail")
		}
		for _, instance := range rsp.Instances {
			key := fmt.Sprintf("%s:%d", instance.GetHost(), instance.GetPort())
			maps[key] += 1
		}
		time.Sleep(10 * time.Millisecond)
	}

	for k, v := range maps {
		t.Logf("key: %s, value: %d", k, v)
	}

}

func reportDynamicWeight(t *testing.T, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			req := api.DynamicWeightReportRequest{}
			req.Service = "lzb_test2"
			req.Namespace = "Test"
			req.Host = "127.0.0.4"
			req.Port = 8080
			req.Token = "@bafd38f07a"

			metrics := make(map[string]string)
			metrics["capacity"] = "100"
			metrics["used"] = "80"
			req.SetMetrics(metrics)

			err := ProviderV2.ReportDynamicWeight(&req)
			if err != nil {
				t.Logf("report dynamic weight err: %+v", err)
			}

			req = api.DynamicWeightReportRequest{}
			req.Service = "lzb_test2"
			req.Namespace = "Test"
			req.Host = "127.0.0.1"
			req.Port = 8080
			req.Token = "@bafd38f07a"

			metrics = make(map[string]string)
			metrics["capacity"] = "100"
			metrics["used"] = "10"
			req.SetMetrics(metrics)

			err = ProviderV2.ReportDynamicWeight(&req)
			if err != nil {
				t.Logf("report dynamic weight err: %+v", err)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}
