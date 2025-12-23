package main

import (
	"flag"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	plog "git.woa.com/polaris/polaris-go/v2/pkg/log"
)

/*
go-sdk 上报 服务调用统计 数据
*/

const logLevel = plog.InfoLog

var (
	addresses string
	joinPoint string
	namespace string
	service   string
	seconds   int
	slice     string
)

func initArgs() {
	// [必填参数]
	flag.StringVar(&namespace, "namespace", "Test", "namespace") // 服务命名空间，格式 -namespace=xxx
	flag.StringVar(&service, "service", "lzb_test", "service")   // 服务名，格式 -service=xxx
	// [选填参数]
	flag.IntVar(&seconds, "round", 20, "loop seconds")                       // 循环次数
	flag.StringVar(&addresses, "addresses", "", "polaris default addresses") // 埋点地址，默认接入公共集群
	flag.StringVar(&joinPoint, "joinPoint", "", "joinPoint")                 // 接入点，默认接入公共集群
	flag.StringVar(&slice, "slice", "", "slice name of service")             // 分片名
}

func GetConsumerProd() (api.ConsumerAPI, error) {
	//创建consumerAPI实例
	//注意该实例所有方法都是协程安全，一般用户进程只需要创建一个consumerAPI,重复使用即可
	//切勿每次调用之前都创建一个consumerAPI
	//设置直接埋点地址与设置jointPoint不能同时生效
	cfg := api.NewConfiguration()
	if len(addresses) > 0 {
		cfg.GetGlobal().GetServerConnector().SetAddresses(strings.Split(addresses, ","))
	} else if len(joinPoint) > 0 {
		cfg.GetGlobal().GetServerConnector().SetJoinPoint(joinPoint)
	}
	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if nil != err {
		return nil, err
	}
	return consumer, nil
}

func GetConsumerDev(path string) (api.ConsumerAPI, error) {
	consumer, err := api.NewConsumerAPIByFile(path)
	if nil != err {
		return nil, err
	}
	return consumer, nil
}

// 主入口函数
func main() {
	initArgs()
	flag.Parse()
	var err error
	err = api.ConfigLoggers("polaris", logLevel)
	if nil != err {
		log.Fatalf("fail to SetLogLevel, err is %v", err)
	}

	if nil != err {
		log.Fatalf("fail to create ConsumerAPI by default configuration, err is %v", err)
	}

	consumer, err := GetConsumerDev("polaris.yaml")
	if nil != err {
		log.Fatalf("fail to create ConsumerAPI by default configuration, err is %v", err)
	}

	deadline := time.Now().Add(time.Duration(seconds) * time.Second)
	for {
		if time.Now().After(deadline) {
			break
		}
		var flowId uint64
		getOneInstanceReq := &api.GetOneInstanceRequest{}
		getOneInstanceReq.FlowID = atomic.AddUint64(&flowId, 1)
		getOneInstanceReq.Namespace = namespace
		getOneInstanceReq.Service = service
		//设置负载均衡算法为权重随机
		getOneInstanceReq.LbPolicy = api.LBPolicyWeightedRandom
		//设置服务的分片名
		getOneInstanceReq.SetSlice(slice)

		startTime := time.Now()
		//进行服务发现，获取单一服务实例
		getInstResp, err := consumer.GetOneInstance(getOneInstanceReq)
		if nil != err {
			log.Fatalf("fail to sync GetOneInstance, err is %v", err)
		}
		consumeDuration := time.Since(startTime)
		log.Printf("success to sync GetOneInstance, count is %d, consume is %v\n",
			len(getInstResp.Instances), consumeDuration)
		targetInstance := getInstResp.Instances[0]
		log.Printf("sync instance is id=%s, address=%s:%d\n",
			targetInstance.GetId(), targetInstance.GetHost(), targetInstance.GetPort())
		//构造请求，进行服务调用结果上报
		svcCallResult := &api.ServiceCallResult{}
		//设置被调的实例信息
		svcCallResult.SetCalledInstance(targetInstance)
		//设置服务调用结果，枚举，成功或者失败
		svcCallResult.SetRetStatus(api.RetSuccess)
		//设置服务调用返回码
		svcCallResult.SetRetCode(0)
		//设置服务调用时延信息
		svcCallResult.SetDelay(consumeDuration)
		//执行调用结果上报
		err = consumer.UpdateServiceCallResult(svcCallResult)
		if nil != err {
			log.Fatalf("fail to UpdateServiceCallResult, err is %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	log.Printf("success to sync get one instance")

	// 一般用户进程只需要创建一个consumerAPI，重复使用即可。在退出main前释放consumerAPI。
	consumer.Destroy()
}
