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

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"git.woa.com/polaris/polaris-server-api/api/monitor/polaris/v1"

	"GoLearn/polaris/monitor/grpc/global-scheduler"
)

func main() {
	var (
		addr        string
		namespace   string
		service     string
		lbType      string
		count       int
		dialTimeout time.Duration
		callTimeout time.Duration
		id          string
	)
	flag.StringVar(&addr, "addr", "9.134.117.127:8086", "CollectScheduleRequestStat gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "Test", "北极星命名空间")
	flag.StringVar(&service, "service", "lzb_test", "业务服务名")
	flag.StringVar(&lbType, "lb", "GLOBAL_P2C", "负载均衡策略")
	flag.IntVar(&count, "count", 3, "上报请求条数")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second, "gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 10*time.Second, "流式调用超时")
	flag.StringVar(&id, "id", "", "记录唯一 ID 前缀（为空则自动生成）")
	flag.Parse()

	log.Printf("start CollectScheduleRequestStat, addr=%s namespace=%s service=%s count=%d",
		addr, namespace, service, count)

	option := &global_scheduler.CollectOption{
		Addr:            addr,
		Namespace:       namespace,
		Service:         service,
		LoadBalanceType: lbType,
		Count:           count,
		DialTimeout:     dialTimeout,
		CallTimeout:     callTimeout,
		ID:              id,
	}

	handle := func(index int, resp *v1.StatResponse) {
		fmt.Printf("[resp %d] code=%d info=%s\n", index, resp.GetCode().GetValue(), resp.GetInfo().GetValue())
	}

	if err := global_scheduler.CollectScheduleRequestStat(option, handle); err != nil {
		log.Fatalf("CollectScheduleRequestStat failed: %v", err)
	}
	log.Printf("CollectScheduleRequestStat done, total response=%d", count)
}
