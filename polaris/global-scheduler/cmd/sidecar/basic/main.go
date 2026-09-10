// Tencent is pleased to support the open source community by making polaris available.
//
// Copyright (C) 2024 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the BSD 3-Clause License (the "License"); you may not use
// this file except in compliance with the License. You may obtain a copy of
// the License at http://opensource.org/licenses/BSD-3-Clause
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// 该脚本演示如何以 gRPC 方式调用 polaris-sidecar 提供的 ConsumerService，
// 依次调用 GetOneInstance、Discover(INSTANCE/ROUTING)、UpdateCallResult、ReportClient。
// 服务端实现参考: polaris-sidecar/grpcserver/consumer_access.go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	apiV2 "git.woa.com/polaris/polaris-server-api/sidecar/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		serverAddr  string
		namespace   string
		service     string
		dialTimeout time.Duration
		callTimeout time.Duration
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8092", "polaris-sidecar gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "Test", "北极星命名空间")
	flag.StringVar(&service, "service", "lzb_test", "北极星服务名")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second, "gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 5*time.Second, "gRPC 调用超时")
	flag.Parse()

	// 建立 gRPC 连接
	dialCtx, dialCancel := context.WithTimeout(context.Background(), dialTimeout)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("连接服务器 %s 失败: %v", serverAddr, err)
	}
	defer conn.Close()

	client := apiV2.NewConsumerServiceClient(conn)

	fmt.Printf("========== 目标: %s | namespace=%s service=%s ==========\n",
		serverAddr, namespace, service)

	// 1) GetOneInstance
	callGetOneInstance(client, namespace, service, callTimeout)

	// 2) Discover: 拉取实例
	//callDiscover(client, namespace, service, apiV2.ResourceType_INSTANCE, callTimeout)

	// 3) Discover: 拉取路由
	//callDiscover(client, namespace, service, apiV2.ResourceType_ROUTING, callTimeout)

	// 4) ReportClient
	//callReportClient(client, callTimeout)

	// 5) UpdateCallResult
	//callUpdateCallResult(client, namespace, service, callTimeout)
}

// callGetOneInstance 调用单次服务发现
func callGetOneInstance(client apiV2.ConsumerServiceClient, namespace, service string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &apiV2.InstanceRequest{
		Namespace: namespace,
		Service:   service,
		LbType:    apiV2.InstanceRequest_LB_DEFAULT,
	}
	resp, err := client.GetOneInstance(ctx, req)
	if err != nil {
		log.Printf("[GetOneInstance] 调用失败: %v", err)
		return
	}
	fmt.Println("---------- GetOneInstance ----------")
	fmt.Printf("Code=%d Info=%s\n", resp.GetCode(), resp.GetInfo())
	for i, ins := range resp.GetInstances() {
		fmt.Printf("  [%d] id=%s host=%s port=%d weight=%d healthy=%v\n",
			i, ins.GetId(), ins.GetHost(), ins.GetPort(), ins.GetWeight(), ins.GetHealthy())
	}
	if loc := resp.GetLocation(); loc != nil {
		fmt.Printf("  Location: region=%s zone=%s campus=%s\n",
			loc.GetRegion(), loc.GetZone(), loc.GetCampus())
	}
}

// callDiscover 调用批量服务发现（双向流，只发送 1 条请求后读取应答）
func callDiscover(client apiV2.ConsumerServiceClient, namespace, service string,
	rt apiV2.ResourceType, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := client.Discover(ctx)
	if err != nil {
		log.Printf("[Discover] 建流失败: %v", err)
		return
	}

	req := &apiV2.DiscoverRequest{
		Type:      rt,
		Namespace: namespace,
		Service:   service,
	}
	if err := stream.Send(req); err != nil {
		log.Printf("[Discover] 发送失败: %v", err)
		return
	}
	// 关闭发送端，触发服务端 io.EOF 结束
	if err := stream.CloseSend(); err != nil {
		log.Printf("[Discover] CloseSend 失败: %v", err)
	}

	fmt.Printf("---------- Discover(%s) ----------\n", rt.String())
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("[Discover] 接收失败: %v", err)
			return
		}
		fmt.Printf("Code=%d Info=%s Type=%s Namespace=%s Service=%s Revision=%s\n",
			resp.GetCode(), resp.GetInfo(), resp.GetType().String(),
			resp.GetNamespace(), resp.GetService(), resp.GetRevision())
		if rt == apiV2.ResourceType_INSTANCE {
			fmt.Printf("  Instances(%d):\n", len(resp.GetInstances()))
			for i, ins := range resp.GetInstances() {
				fmt.Printf("    [%d] id=%s host=%s port=%d weight=%d healthy=%v\n",
					i, ins.GetId(), ins.GetHost(), ins.GetPort(), ins.GetWeight(), ins.GetHealthy())
			}
		} else if r := resp.GetRouting(); r != nil {
			fmt.Printf("  Routing: inbounds=%d outbounds=%d\n",
				len(r.GetInbounds()), len(r.GetOutbounds()))
		}
	}
}

// callReportClient 上报客户端 SDK 信息，获取位置
func callReportClient(client apiV2.ConsumerServiceClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.ReportClient(ctx, &apiV2.Client{
		SdkType: "sidecar-basic-demo",
		Version: "1.0.0",
	})
	if err != nil {
		log.Printf("[ReportClient] 调用失败: %v", err)
		return
	}
	fmt.Println("---------- ReportClient ----------")
	fmt.Printf("Code=%d Info=%s\n", resp.GetCode(), resp.GetInfo())
	if loc := resp.GetLocation(); loc != nil {
		fmt.Printf("  Location: region=%s zone=%s campus=%s\n",
			loc.GetRegion(), loc.GetZone(), loc.GetCampus())
	}
}

// callUpdateCallResult 上报一次调用结果
func callUpdateCallResult(client apiV2.ConsumerServiceClient, namespace, service string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := client.UpdateCallResult(ctx)
	if err != nil {
		log.Printf("[UpdateCallResult] 建流失败: %v", err)
		return
	}

	stat := &apiV2.CallResultStat{
		Namespace: namespace,
		Service:   service,
		CallResults: []*apiV2.CallResultStat_CallResult{
			{
				Host:      "127.0.0.1",
				Port:      8080,
				RetStatus: apiV2.CallResultStat_CallResult_RetOk,
				RetCode:   0,
				Count:     1,
				CallDelay: int64(10 * time.Millisecond),
			},
		},
	}
	if err := stream.Send(stat); err != nil {
		log.Printf("[UpdateCallResult] 发送失败: %v", err)
		return
	}
	if err := stream.CloseSend(); err != nil {
		log.Printf("[UpdateCallResult] CloseSend 失败: %v", err)
	}

	fmt.Println("---------- UpdateCallResult ----------")
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("[UpdateCallResult] 接收失败: %v", err)
			return
		}
		fmt.Printf("Code=%d Info=%s\n", resp.GetCode(), resp.GetInfo())
	}
}
