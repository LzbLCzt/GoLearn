/**
 * Tencent is pleased to support the open source community by making Polaris available.
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

// Package main 通过固定主调 metadata 使每次 select 命中同一实例，
// 并连续上报失败以触发被调实例的熔断。
package main

import (
	"flag"
	"log"
	"sync/atomic"
	"time"

	"git.woa.com/polaris/polaris-go/v2/api"
	plog "git.woa.com/polaris/polaris-go/v2/pkg/log"
	"git.woa.com/polaris/polaris-go/v2/pkg/model"
)

const logLevel = plog.InfoLog

var (
	namespace    string
	service      string
	srcNamespace string
	srcService   string
	count        int
	interval     time.Duration
)

func initArgs() {
	// 被调服务
	flag.StringVar(&namespace, "namespace", "Test", "被调服务命名空间")
	flag.StringVar(&service, "service", "lzb_test", "被调服务名")
	// 主调服务（用于携带 metadata，保证每次 select 到同一实例）
	flag.StringVar(&srcNamespace, "src-namespace", "Test", "主调服务命名空间")
	flag.StringVar(&srcService, "src-service", "lzb_test2", "主调服务名")
	// 循环次数
	flag.IntVar(&count, "count", 20, "select 次数")
	// 每次 select 后的间隔
	flag.DurationVar(&interval, "interval", time.Second, "每次 select 后 sleep 时长")
}

func main() {
	initArgs()
	flag.Parse()

	if err := api.ConfigLoggers("polaris", logLevel); err != nil {
		log.Fatalf("fail to SetLogLevel, err is %v", err)
	}

	consumer, err := api.NewConsumerAPI()
	if err != nil {
		log.Fatalf("fail to create ConsumerAPI, err is %v", err)
	}
	defer consumer.Destroy()

	var flowID uint64
	log.Printf("start circuitbreaker basic test, namespace=%s service=%s count=%d interval=%s",
		namespace, service, count, interval)

	for i := 1; i <= count; i++ {
		req := &api.GetOneInstanceRequest{}
		req.FlowID = atomic.AddUint64(&flowID, 1)
		req.Namespace = namespace
		req.Service = service
		// 通过主调 metadata k4:v4 让每次 select 命中同一实例
		req.SourceService = &model.ServiceInfo{
			Namespace: srcNamespace,
			Service:   srcService,
			Metadata: map[string]string{
				"k4": "v4",
			},
		}

		startTime := time.Now()
		resp, err := consumer.GetOneInstance(req)
		if err != nil {
			log.Printf("[%d] fail to GetOneInstance, err is %v", i, err)
			time.Sleep(interval)
			continue
		}
		consume := time.Since(startTime)
		if len(resp.Instances) == 0 {
			log.Printf("[%d] GetOneInstance empty instances", i)
			time.Sleep(interval)
			continue
		}
		target := resp.Instances[0]
		var status model.Status
		if target.GetCircuitBreakerStatus() != nil {
			status = target.GetCircuitBreakerStatus().GetStatus()
		}
		log.Printf("[%d] select instance id=%s addr=%s:%d weight=%d circuit breaker status=%d consume=%v",
			i, target.GetId(), target.GetHost(), target.GetPort(), target.GetWeight(), status, consume)
		// 上报失败结果，用于触发熔断
		svcCallResult := &api.ServiceCallResult{}
		svcCallResult.SetCalledInstance(target)
		svcCallResult.SetRetStatus(api.RetFail)
		svcCallResult.SetRetCode(-1)
		svcCallResult.SetDelay(consume)
		if err := consumer.UpdateServiceCallResult(svcCallResult); err != nil {
			log.Printf("[%d] fail to UpdateServiceCallResult, err is %v", i, err)
		} else {
			log.Printf("[%d] report fail success, instance=%s:%d", i, target.GetHost(), target.GetPort())
		}

		time.Sleep(interval)
	}
	log.Printf("circuitbreaker basic test done, total=%d", count)
}
