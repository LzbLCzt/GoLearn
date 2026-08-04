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

package global_scheduler

import (
	"context"
	"fmt"
	"io"
	"time"

	api "git.woa.com/polaris/polaris-server-api/api/monitor/polaris/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CollectOption 调用 CollectScheduleRequestStat 接口的可选参数
type CollectOption struct {
	Addr            string        // gRPC 服务地址，如 9.134.117.127:8086
	Namespace       string        // 北极星命名空间
	Service         string        // 业务服务名
	LoadBalanceType string        // 负载均衡策略
	Count           int           // 上报请求条数
	DialTimeout     time.Duration // 拨号超时
	CallTimeout     time.Duration // 流式调用超时
	ID              string        // 记录唯一 ID，为空时自动生成
}

// ResponseHandler 每收到一条服务端响应时的回调，用于打印或断言
type ResponseHandler func(index int, resp *api.StatResponse)

// CollectScheduleRequestStat 调用 CollectScheduleRequestStat 流式接口。
// 循环上报 Count 条全局调度请求统计，并在 CloseSend 后读取服务端回包直到 EOF。
func CollectScheduleRequestStat(option *CollectOption, handle ResponseHandler) error {
	if option == nil {
		return fmt.Errorf("option is nil")
	}
	if option.Count <= 0 {
		return fmt.Errorf("count must be positive, got %d", option.Count)
	}

	dialCtx, dialCancel := context.WithTimeout(context.Background(), option.DialTimeout)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, option.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("dial %s failed: %w", option.Addr, err)
	}
	defer conn.Close()

	client := api.NewGrpcAPIClient(conn)
	callCtx, callCancel := context.WithTimeout(context.Background(), option.CallTimeout)
	defer callCancel()

	stream, err := client.CollectScheduleRequestStat(callCtx)
	if err != nil {
		return fmt.Errorf("open CollectScheduleRequestStat stream failed: %w", err)
	}

	// 循环上报
	for i := 0; i < option.Count; i++ {
		req := buildRequest(option, i)
		if err := stream.Send(req); err != nil {
			return fmt.Errorf("send request[%d] failed: %w", i, err)
		}
	}

	// 关闭发送端，通知服务端不再有请求
	if err := stream.CloseSend(); err != nil {
		return fmt.Errorf("close send failed: %w", err)
	}

	// 读取服务端响应直到 EOF
	index := 0
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv response[%d] failed: %w", index, err)
		}
		if handle != nil {
			handle(index, resp)
		}
		index++
	}
	return nil
}

// buildRequest 构造第 index 条 GlobalScheduleRequestStat 上报请求
func buildRequest(option *CollectOption, index int) *api.GlobalScheduleRequestStat {
	req := &api.GlobalScheduleRequestStat{
		Id:              genRequestID(option, index),
		Service:         option.Service,
		Namespace:       option.Namespace,
		LoadBalanceType: option.LoadBalanceType,
		SdkToken: &api.SDKToken{
			Ip:      "127.0.0.1",
			Client:  "polaris-go",
			Version: "2.6.0",
		},
		Time: timestamppb.Now(),
		RequestStat: []*api.RequestStat{
			{StatType: api.RequestStatType_RequestSuccess, PeriodTimes: 1, Reason: "ok"},
			{StatType: api.RequestStatType_RequestDegrade, PeriodTimes: 1, Reason: "degrade"},
			{StatType: api.RequestStatType_RequestFailed, PeriodTimes: 1, Reason: "failed"},
		},
	}
	return req
}

// genRequestID 生成记录唯一 ID，option.ID 非空时复用并追加序号
func genRequestID(option *CollectOption, index int) string {
	if option.ID == "" {
		return fmt.Sprintf("test-schedule-stat-%d-%d", time.Now().UnixNano(), index)
	}
	return fmt.Sprintf("%s-%d", option.ID, index)
}
