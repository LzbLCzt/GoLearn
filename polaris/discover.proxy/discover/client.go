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

// Package discover 提供 PolarisGRPC Discover 接口的客户端实现。
// 服务端对应 GrpcProxyHandler.Discover（服务端流：循环接收 DiscoverRequest 并回包 DiscoverResponse）。
package discover

import (
	"context"

	grpcpb "git.woa.com/polaris/polaris-server-api/api/v1/grpc"
	"google.golang.org/grpc"
)

// DiscoverClient 封装与 discover-proxy 之间的一条 Discover 双向流连接。
// Discover 接口为服务端流，客户端可多次 Send 请求并对应 Recv 响应。
type DiscoverClient struct {
	conn   *grpc.ClientConn
	stream grpcpb.PolarisGRPC_DiscoverClient
}

// NewClient 拨号指定地址并建立 Discover 流连接。
// address 格式为 host:port，例如 "9.141.112.135:8081"。
func NewClient(address string) (*DiscoverClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	stream, err := grpcpb.NewPolarisGRPCClient(conn).Discover(context.Background())
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &DiscoverClient{conn: conn, stream: stream}, nil
}

// Stream 返回底层的 Discover 流，便于需要自定义循环 Send/Recv 的调用方使用。
func (c *DiscoverClient) Stream() grpcpb.PolarisGRPC_DiscoverClient {
	return c.stream
}

// Close 关闭底层连接，释放资源。
func (c *DiscoverClient) Close() error {
	return c.conn.Close()
}
