package sdk

import (
	"strings"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.woa.com/polaris/polaris-go/v2/api"
)

var Consumer api.ConsumerAPI
var Provider api.ProviderAPI

func InitPolarisByYaml(path string) error {
	if Consumer != nil {
		return nil
	}
	var err error
	Consumer, err = api.NewConsumerAPIByFile(path)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}

func InitPolaris(address string) error {
	if Consumer != nil {
		return nil
	}

	cfg := api.NewConfiguration()
	if len(address) > 0 {
		cfg.GetGlobal().GetServerConnector().SetAddresses(strings.Split(address, ","))
	}

	var err error
	Consumer, err = api.NewConsumerAPIByConfig(cfg)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}

func InitPolarisProviderByYaml(path string) error {
	if Provider != nil {
		return nil
	}
	var err error
	Provider, err = api.NewProviderAPIByFile(path)
	if err != nil {
		log.Errorf("fail to create ConsumerAPI by default configuration, err is %v", err)
		return err
	}
	return nil
}
