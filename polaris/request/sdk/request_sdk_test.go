package sdk

import (
	"testing"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/polaris/polaris-go/v2/api"
	"git.woa.com/polaris/polaris-go/v2/pkg/model"

	"GoLearn/polaris/util"
)

func TestSdkRequest_CreateInstance(t *testing.T) {
	err := InitPolarisProviderByYaml("polaris.yaml")
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}
	hosts := util.MockIPs(100)

	for _, host := range hosts {
		err := RegisterInstance(t, host)
		if err != nil {
			t.Errorf("fail to register instance, err is %v", err)
		}
		time.Sleep(time.Millisecond * 20)
	}
}

func RegisterInstance(t *testing.T, host string) error {
	req := api.InstanceRegisterRequest{}
	req.Namespace = "Polaris"
	req.Service = "polaris.report.maglev"
	req.Host = host
	req.Port = 8089
	protocol := "http"
	req.Protocol = &protocol
	weight := 100
	req.Weight = &weight
	req.Location = &model.InstanceLocation{
		Region: "华南",
		Zone:   "深圳",
		Campus: "深圳三区",
	}
	req.ServiceToken = "22741879d8c64abbb8518a10e4a6cd98"

	rsp, err := Provider.Register(&req)
	if err != nil {
		return err
	}
	log.Infof("success to create instance, instance is %v", rsp)
	return nil
}

// 通过配置文件启动polaris sdk(配置文件用的是开发环境的配置)
func TestSdkRequest_GetOneInstanceV2(t *testing.T) {
	err := InitPolarisByYaml("polaris.yaml")
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}
	req := api.GetOneInstanceRequest{}
	req.Service = "lzb_test"
	req.Namespace = "Test"
	//req.LbPolicy = config.DefaultLoadBalancerMaglev
	//req.HashKey = []byte("aaa")
	rsp, err := Consumer.GetOneInstance(&req)
	//rsp, err := Consumer.GetInstances(&req)
	if err != nil {
		t.Errorf("fail to get instance, err is %v", err)
		return
	}
	for _, instance := range rsp.Instances {
		log.Infof("success to get instance, instance is %v", instance)
	}
}

func TestSdkRequest_GetInstancesV2(t *testing.T) {
	err := InitPolarisByYaml("polaris.yaml")
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}
	req := api.GetInstancesRequest{}
	req.Service = "polaris.report.test"
	req.Namespace = "Polaris"
	rsp, err := Consumer.GetInstances(&req)
	if err != nil {
		t.Errorf("fail to get instance, err is %v", err)
		return
	}
	for _, instance := range rsp.Instances {
		log.Infof("success to get instance, instance is %v", instance)
	}
}

func TestSdkRequest_GetOneInstance(t *testing.T) {
	err := InitPolaris("")
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}

	req := api.GetOneInstanceRequest{}
	req.Namespace = "Test"
	req.Service = "shennong-backend-risk.test"
	rsp, err := Consumer.GetOneInstance(&req)
	if err != nil {
		t.Errorf("fail to get instance, err is %v", err)
		return
	}

	instance := rsp.Instances[0]
	log.Infof("success to get instance, instance is %v", instance)

	metadata := rsp.GetMetadata()
	for key, value := range metadata {
		log.Infof("metadata: key is %s, value is %s", key, value)
	}

}
