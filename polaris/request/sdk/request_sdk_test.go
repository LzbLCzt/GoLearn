package sdk

import (
	"git.code.oa.com/polaris/polaris-go/api"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"strings"
	"testing"
)

var consumer api.ConsumerAPI

func TestSdkRequest_GetOneInstance(t *testing.T) {
	err := InitPolaris("")
	if err != nil {
		t.Errorf("fail to init polaris, err is %v", err)
		return
	}

	req := api.GetOneInstanceRequest{}
	req.Namespace = "Test"
	req.Service = "shennong-backend-risk.test"
	rsp, err := consumer.GetOneInstance(&req)
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

func InitPolaris(address string) error {
	if consumer != nil {
		return nil
	}

	cfg := api.NewConfiguration()
	if len(address) > 0 {
		cfg.GetGlobal().GetServerConnector().SetAddresses(strings.Split(address, ","))
	}

	var err error
	consumer, err = api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}
