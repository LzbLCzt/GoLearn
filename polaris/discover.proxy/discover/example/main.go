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

// Package main 提供 discover.proxy Discover 接口的最小化调用示例。
// 通过命令行参数 -addr 指定 discover-proxy 服务地址，发起一次 INSTANCE 类型的服务发现请求并打印响应。
package main

import (
	"flag"
	"fmt"
	"log"

	"GoLearn/polaris/discover.proxy/discover"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
)

func main() {
	var (
		serverAddr  string
		namespace   string
		serviceName string
	)
	flag.StringVar(&serverAddr, "addr", "",
		"discover-proxy gRPC 服务地址，格式 host:port（必填）")
	flag.StringVar(&namespace, "namespace", "Polaris", "北极星命名空间")
	flag.StringVar(&serviceName, "service", "", "北极星服务名（必填）")
	flag.Parse()

	if serverAddr == "" || serviceName == "" {
		flag.Usage()
		log.Fatalf("addr 和 service 参数必填")
	}

	client, err := discover.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("连接 %s 失败: %v", serverAddr, err)
	}
	defer client.Close()

	fmt.Printf("已连接 %s, 目标: %s/%s, 类型: INSTANCE\n",
		serverAddr, namespace, serviceName)

	resp, err := client.Discover(&discover.Request{
		Type: apiV1Model.DiscoverRequest_INSTANCE,
		//Type:        apiV1Model.DiscoverRequest_RATE_LIMIT,
		Namespace:   namespace,
		ServiceName: serviceName,
	})
	if err != nil {
		log.Fatalf("Discover 请求失败: %v", err)
	}

	fmt.Printf("Discover 响应: code=%d, info=%s, 实例数=%d\n",
		resp.GetCode().GetValue(),
		resp.GetInfo().GetValue(),
		len(resp.GetInstances()))
	for i, ins := range resp.GetInstances() {
		fmt.Printf("  [%d] host=%s port=%d weight=%d healthy=%v\n",
			i,
			ins.GetHost().GetValue(),
			ins.GetPort().GetValue(),
			ins.GetWeight().GetValue(),
			ins.GetHealthy().GetValue())
	}
}
