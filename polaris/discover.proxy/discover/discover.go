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

package discover

import (
	"errors"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// 默认命名空间，兼容老版本 SDK 不传命名空间的情况（对齐服务端 checkParam 逻辑）。
const defaultNamespace = "Polaris"

// Request 是 Discover 请求的业务入参，屏蔽底层 proto 字段的构造细节。
type Request struct {
	// Type 发现类型，例如 apiV1Model.DiscoverRequest_INSTANCE。
	Type apiV1Model.DiscoverRequest_DiscoverRequestType
	// Namespace 命名空间，为空时自动使用 defaultNamespace。
	Namespace string
	// ServiceName 服务名，不可为空。
	ServiceName string
	// Revision 服务版本号，用于增量发现；首次发现可传空。
	Revision string
	// Cluster 客户端期望的集群，可选。
	Cluster string
}

// Discover 向 discover-proxy 发起一次服务发现请求并返回响应。
// 该方法对齐服务端 GrpcProxyHandler.Discover：服务名为空时返回错误；
// 命名空间为空时回退到默认命名空间。
func (c *DiscoverClient) Discover(req *Request) (*apiV1Model.DiscoverResponse, error) {
	if req == nil {
		return nil, errors.New("discover request is nil")
	}
	if req.ServiceName == "" {
		return nil, errors.New("service name is required")
	}
	namespace := req.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	request := &apiV1Model.DiscoverRequest{
		Type: req.Type,
		Service: &apiV1Model.Service{
			Namespace: &wrapperspb.StringValue{Value: namespace},
			Name:      &wrapperspb.StringValue{Value: req.ServiceName},
			Revision:  &wrapperspb.StringValue{Value: req.Revision},
		},
	}
	if req.Cluster != "" {
		request.Service.Extend = &apiV1Model.ServiceExtend{
			Cluster: req.Cluster,
		}
	}

	if err := c.stream.Send(request); err != nil {
		return nil, err
	}
	return c.stream.Recv()
}
