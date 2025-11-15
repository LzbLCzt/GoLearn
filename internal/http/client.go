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

package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	api "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"github.com/golang/protobuf/jsonpb"
)

// NewClient 创建HTTP客户端
func NewClient(address, version string) *Client {
	return &Client{
		Address: address,
		Version: version,
		Worker:  &http.Client{},
	}
}

// NewClientWithPlatform 创建HTTP客户端
func NewClientWithPlatform(address, version, platformId, platformToken string) *Client {
	return &Client{
		Address: address,
		Version: version,
		Worker: &http.Client{
			Timeout: 10 * time.Second,
		},
		PlatformId:    platformId,
		PlatformToken: platformToken,
	}
}

func NewClientWithPlatformAndProxyToken(address, version, platformId, platformToken string, proxyToken string) *Client {
	return &Client{
		Address: address,
		Version: version,
		Worker: &http.Client{
			Timeout: 10 * time.Second,
		},
		PlatformId:    platformId,
		PlatformToken: platformToken,
		ProxyToken:    proxyToken,
	}
}

func NewClientWithPolarisToken(address, version, polarisToken string, timeout int) *Client {
	c := &Client{
		Address:      address,
		Version:      version,
		Worker:       &http.Client{},
		PolarisToken: polarisToken,
	}
	// if timeout > 0 {
	// 	c.Worker.Timeout = time.Duration(timeout) * time.Second
	// }
	return c
}

// Client HTTP客户端
type Client struct {
	Address          string
	Version          string
	Worker           *http.Client
	PlatformId       string
	PlatformToken    string
	ProxyToken       string
	AliasUpdateToken string
	PolarisToken     string
}

// SendRequest 发送 HTTP Post/Put
func (c *Client) SendRequest(method string, url string, body *bytes.Buffer) (*http.Response, error) {
	var request *http.Request
	var err error

	var bodyStr string
	if body == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		bodyStr = body.String()
		request, err = http.NewRequest(method, url, body)
	}

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	// request.Header.Add("Request-Id", "test")
	request.Header.Add("Platform-Id", c.PlatformId)
	request.Header.Add("Platform-Token", c.PlatformToken)

	// Keep-Alive
	request.Header.Add("Connection", "keep-alive")

	if c.ProxyToken != "" {
		request.Header.Add("Proxy-Token", c.ProxyToken)
	}

	if c.PolarisToken != "" {
		request.Header.Add("Polaris-Token", c.PolarisToken)
	}

	response, err := c.Worker.Do(request)
	if err != nil {
		// 增加详细的错误信息，帮助诊断EOF错误
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("HTTP请求EOF错误 - 可能原因: 1.服务器未启动 2.网络连接中断 3.服务器提前关闭连接 4.超时(当前超时: %v)", c.Worker.Timeout)
		}
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Printf("request url:%s body:%s, response:%+v\n", url, bodyStr, response)
	}

	return response, nil
}

// SendRequestWithGlobalToken 发送带全局token的请求
func (c *Client) SendRequestWithGlobalToken(method string, url string, body *bytes.Buffer) (
	*http.Response, error) {
	var request *http.Request
	var err error

	if body == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, body)
	}

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Request-Id", "test")
	request.Header.Add("Polaris-token", "polaris@12345678")

	response, err := c.Worker.Do(request)
	if err != nil {
		// 增加详细的错误信息，帮助诊断EOF错误
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("HTTP请求EOF错误 - 可能原因: 1.服务器未启动 2.网络连接中断 3.服务器提前关闭连接 4.超时(当前超时: %v)", c.Worker.Timeout)
		}
		return nil, err
	}

	return response, nil
}

// SendRequestWithHeader 带指定header 发送 HTTP Post/Put
func (c *Client) SendRequestWithHeader(method string, url string,
	body *bytes.Buffer, header map[string]string) (*http.Response, error) {
	var request *http.Request
	var err error

	var bodyStr string
	if body == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		bodyStr = body.String()
		request, err = http.NewRequest(method, url, body)
	}

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Request-Id", "test")
	for key, value := range header {
		request.Header.Set(key, value)
	}

	response, err := c.Worker.Do(request)
	if err != nil {
		// 增加详细的错误信息，帮助诊断EOF错误
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("HTTP请求EOF错误 - 可能原因: 1.服务器未启动 2.网络连接中断 3.服务器提前关闭连接 4.超时(当前超时: %v)", c.Worker.Timeout)
		}
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Printf("request url:%s body:%s, response:%+v\n", url, bodyStr, response)
	}

	return response, nil
}

// CompleteURL 生成GET请求的完整URL
func (c *Client) CompleteURL(url string, params map[string][]interface{}) string {
	count := 1
	url += "?"

	num := 0
	for _, param := range params {
		num += len(param)
	}

	for index, param := range params {
		for _, item := range param {
			url += fmt.Sprintf("%v=%v", index, item)
			if count != num {
				url += "&"
			}
			count++
		}
	}
	return url
}

// GetBatchWriteResponse 获取BatchWriteResponse
func GetBatchWriteResponse(response *http.Response) (*api.BatchWriteResponse, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.BatchWriteResponse{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("unmarshal resp err:%+v, %v\n", response, checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid http code, code: %d, ret: %+v", response.StatusCode, ret)
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}

// GetBatchQueryResponse 获取BatchQueryResponse
func GetBatchQueryResponse(response *http.Response) (*api.BatchQueryResponse, error) {
	// 打印和检查回复
	fmt.Printf("http code: %v\n", response.StatusCode)
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("invalid http code:%d", response.StatusCode))
	}

	ret := &api.BatchQueryResponse{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("ret: %+v\n", ret)
	} else {
		fmt.Printf("err: %v\n", checkErr)
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}

// GetSimpleResponse 获取SimpleResponse
func GetSimpleResponse(response *http.Response) (*api.Response, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.Response{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("%v\n", checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid http code, code: %d, ret: %+v", response.StatusCode, ret)
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}

// GetDiscoverResponse 解码相应
func GetDiscoverResponse(response *http.Response) (*api.DiscoverResponse, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.DiscoverResponse{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("%v\n", checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, errors.New("invalid http code")
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}
