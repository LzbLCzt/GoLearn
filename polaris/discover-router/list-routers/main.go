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

// Package main 提供 discover-router ListRouters 接口的最小化 trpc 调用示例。
// 服务端对应 polaris-service cmd/discover-router/trpcserver/access.go 中的 ListRouters。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/client"
	api "git.woa.com/polaris/polaris-server-api/api/router"
)

func main() {
	var (
		serverAddr          string
		namespace           string
		serviceName         string
		serviceClusterType  int
		resourceClusterType int
		resourceType        int
		resourceID          string
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8098",
		"discover-router trpc 服务地址，格式 host:port")
	flag.StringVar(&namespace, "namespace", "Polaris", "北极星命名空间")
	flag.StringVar(&serviceName, "service", "", "北极星服务名（必填）")
	flag.IntVar(&serviceClusterType, "service_cluster_type", int(api.ClusterType_ClusterTypeDiscover),
		"ServiceRouterReq.ClusterType，可选 0-6，1=Discover")
	flag.IntVar(&resourceClusterType, "resource_cluster_type", int(api.ClusterType_ClusterTypeHeartbeat),
		"ResourceRouterReq.ClusterType，可选 0-6，2=Heartbeat")
	flag.IntVar(&resourceType, "resource_type", int(api.ResourceType_ResourceInstanceId),
		"ResourceRouterReq.ResourceType，可选 0-3，1=InstanceId")
	flag.StringVar(&resourceID, "resource_id", "", "ResourceRouterReq.ResourceId，资源ID（查询资源路由时使用）")
	flag.Parse()

	if serviceName == "" {
		flag.Usage()
		log.Fatalf("service 参数必填")
	}

	target := fmt.Sprintf("ip://%s", serverAddr)
	fmt.Printf("target: %s, ns=%s, svc=%s, service_cluster_type=%d, resource_cluster_type=%d, resource_type=%d, resource_id=%s\n",
		target, namespace, serviceName, serviceClusterType, resourceClusterType, resourceType, resourceID)

	proxy := api.NewRouterApiClientProxy(
		client.WithTarget(target),
		client.WithTimeout(3*time.Second),
	)

	req := &api.ListRoutersReq{
		Services: []*api.ServiceRouterReq{
			{
				ClusterType: api.ClusterType(serviceClusterType),
				Namespace:   namespace,
				Service:     serviceName,
			},
		},
		Resources: []*api.ResourceRouterReq{
			{
				ClusterType:  api.ClusterType(resourceClusterType),
				ResourceType: api.ResourceType(resourceType),
				ResourceId:   resourceID,
			},
		},
	}

	rsp, err := proxy.ListRouters(context.Background(), req)
	if err != nil {
		log.Fatalf("ListRouters 请求失败: %v", err)
	}

	fmt.Printf("ListRouters 响应: code=%d, msg=%s\n", rsp.GetCode(), rsp.GetMsg())
	for i, r := range rsp.GetServiceRoutes() {
		fmt.Printf("  [service#%d] code=%d msg=%s router=%+v\n", i, r.GetCode(), r.GetMsg(), r.GetRouter())
	}
	for i, r := range rsp.GetResourceRoutes() {
		fmt.Printf("  [resource#%d] code=%d msg=%s router=%+v\n", i, r.GetCode(), r.GetMsg(), r.GetRouter())
	}
}
