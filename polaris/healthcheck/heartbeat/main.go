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

// Package main 提供直连 healthcheck-server 的 gRPC 心跳上报脚本。
// 服务端对应 HcGrpcHandler.Heartbeat（unary），
// 参考：polaris-service/cmd/healthcheck/handler/grpchandler.go。
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	grpcpb "git.woa.com/polaris/polaris-server-api/api/v1/grpc"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	serverAddr := flag.String("server", "9.134.117.127:8101", "healthcheck server grpc address")
	namespace := flag.String("namespace", "Test", "polaris namespace")
	service := flag.String("service", "lzb_healthcheck", "polaris service name")
	token := flag.String("token", "5280e43390604b0d85b8e2f6c592f679", "polaris service token")
	hostPort := flag.String("addr", "1.1.1.1:8080", "instance ip:port to heartbeat")
	count := flag.Int("count", 5, "heartbeat times")
	interval := flag.Duration("interval", time.Second, "heartbeat interval")
	timeout := flag.Duration("timeout", 3*time.Second, "per heartbeat rpc timeout")
	flag.Parse()

	host, portStr, err := net.SplitHostPort(*hostPort)
	if err != nil {
		log.Fatalf("invalid -addr %q: %v", *hostPort, err)
	}
	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		log.Fatalf("invalid port in -addr %q: %v", *hostPort, err)
	}

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial %s err: %v", *serverAddr, err)
	}
	defer conn.Close()

	client := grpcpb.NewPolarisGRPCClient(conn)

	instance := &apiV1Model.Instance{
		Namespace:    &wrapperspb.StringValue{Value: *namespace},
		Service:      &wrapperspb.StringValue{Value: *service},
		Host:         &wrapperspb.StringValue{Value: host},
		Port:         &wrapperspb.UInt32Value{Value: uint32(port)},
		ServiceToken: &wrapperspb.StringValue{Value: *token},
	}

	log.Printf("start heartbeat: server=%s namespace=%s service=%s addr=%s count=%d interval=%s",
		*serverAddr, *namespace, *service, *hostPort, *count, *interval)

	for i := 1; i <= *count; i++ {
		doHeartbeat(client, instance, i, *count, *timeout)
		if i < *count {
			time.Sleep(*interval)
		}
	}
}

// doHeartbeat 执行一次心跳上报并打印日志。
func doHeartbeat(client grpcpb.PolarisGRPCClient, in *apiV1Model.Instance, idx, total int, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := client.Heartbeat(ctx, in)
	if err != nil {
		log.Printf("[%d/%d] heartbeat failed, err=%v", idx, total, err)
		return
	}
	log.Printf("[%d/%d] heartbeat ok, code=%d info=%q",
		idx, total, resp.GetCode().GetValue(), resp.GetInfo().GetValue())
}
