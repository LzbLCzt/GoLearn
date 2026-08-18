/**
 * Tencent is pleased to support the open source community by making CL5 available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"git.code.oa.com/polaris/polaris-go/api"
	"git.code.oa.com/polaris/polaris-go/pkg/config"
	"git.code.oa.com/polaris/polaris-go/pkg/model"
)

const logLevel = api.InfoLog

// GetInstanceEvent 从事件channel中非阻塞获取一个订阅事件
func GetInstanceEvent(ch <-chan model.SubScribeEvent) (model.SubScribeEvent, error) {
	select {
	case e := <-ch:
		return e, nil
	default:
		return nil, nil
	}
}

// 主入口函数：订阅指定服务的实例变更事件并打印
func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		log.Fatalf("using %s <namespace> <service>", os.Args[0])
	}
	namespace := argsWithoutProg[0]
	service := argsWithoutProg[1]
	err := api.ConfigLoggers("polaris", logLevel)
	if nil != err {
		log.Fatalf("fail to SetLogLevel, err is %v", err)
	}
	cfg := api.NewConfiguration()
	//设置负载均衡算法为hashRing
	cfg.GetConsumer().GetLoadbalancer().SetType(api.LBPolicyRingHash)
	cfg.GetConsumer().GetSubScribe().SetType(config.SubscribeLocalChannel)
	//创建consumerAPI实例
	//注意该实例所有方法都是协程安全，一般用户进程只需要创建一个consumerAPI,重复使用即可
	//切勿每次调用之前都创建一个consumerAPI
	consumer, err := api.NewConsumerAPIByConfig(cfg)
	if nil != err {
		log.Fatalf("fail to create ConsumerAPI by default configuration, err is %v", err)
	}
	defer consumer.Destroy() // 在进程退出时销毁consumerAPI，不要每次调用都执行API创建和销毁

	key := model.ServiceKey{
		Namespace: namespace,
		Service:   service,
	}
	watchReq := api.WatchServiceRequest{}
	watchReq.Key = key
	watchRsp, err := consumer.WatchService(&watchReq)
	if err != nil {
		fmt.Println("WatchService err: ", err)
		return
	}
	ch := watchRsp.EventChannel
	for {
		event, err := GetInstanceEvent(ch)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		}
		if event == nil {
			fmt.Println("event is nil")
			time.Sleep(time.Second * 3)
			continue
		}
		eType := event.GetSubScribeEventType()
		fmt.Println(eType)
		if eType == api.EventInstance {
			insEvent := event.(*model.InstanceEvent)
			if insEvent.AddEvent != nil {
				fmt.Println("==========add instances: ")
				for _, v := range insEvent.AddEvent.Instances {
					fmt.Println(v.GetId(), v.GetHost(), v.GetPort())
				}
			}
			if insEvent.UpdateEvent != nil {
				fmt.Println("==========update instances: ")
				for _, v := range insEvent.UpdateEvent.UpdateList {
					fmt.Printf("host:%s before: %s after:%s", v.After.GetHost(), v.Before.GetRevision(), v.After.GetRevision())
				}
			}
			if insEvent.DeleteEvent != nil {
				fmt.Println("==========delete instances: ")
				for _, v := range insEvent.DeleteEvent.Instances {
					fmt.Println(v.GetId(), v.GetHost(), v.GetPort())
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}
