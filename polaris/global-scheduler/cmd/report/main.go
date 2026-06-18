// Tencent is pleased to support the open source community by making Polaris available.
//
// Copyright (C) 2026 THL A29 Limited, a Tencent company. All rights reserved.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// instanceFlags 支持 -instance 参数重复出现，每次追加一个实例
type instanceFlags []string

// String flag.Value 接口实现
func (f *instanceFlags) String() string {
	return strings.Join(*f, " | ")
}

// Set flag.Value 接口实现，每次 -instance 调用追加一项
func (f *instanceFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var (
		serverAddr   string
		namespace    string
		service      string
		defaultPort  uint
		metricTimeMs int64
		dialTimeout  time.Duration
		callTimeout  time.Duration
		instances    instanceFlags
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8085", "global-scheduler gRPC 服务地址")
	//flag.StringVar(&serverAddr, "addr", "30.163.16.111:8082", "global-scheduler gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "", "北极星命名空间（必填）")
	flag.StringVar(&service, "service", "", "北极星服务名（必填）")
	flag.UintVar(&defaultPort, "default-port", 8080, "实例未显式指定 port 时使用的默认端口")
	flag.Int64Var(&metricTimeMs, "metric-time-ms", 0, "指标时间戳(毫秒)，0 表示使用当前时间")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second, "gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 3*time.Second, "report 调用超时")
	flag.Var(&instances, "instance",
		"上报的实例及指标，格式: host=IP,port=PORT,metric1=value1,metric2=value2... 可重复传入")
	flag.Parse()

	if namespace == "" || service == "" {
		flag.Usage()
		log.Fatalf("namespace 和 service 参数必填")
	}
	if len(instances) == 0 {
		flag.Usage()
		log.Fatalf("至少需要传入一个 -instance")
	}
	if metricTimeMs == 0 {
		metricTimeMs = time.Now().UnixMilli()
	}

	// 解析实例参数
	nodes, err := parseInstances(instances, uint32(defaultPort), metricTimeMs)
	if err != nil {
		log.Fatalf("解析 -instance 参数失败: %v", err)
	}

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

	// 创建客户端
	client := schedule.NewMetricsServiceClient(conn)

	// 构造请求
	req := &schedule.MetricsReportRequest{
		Namespace: namespace,
		Service:   service,
		Nodes:     nodes,
	}

	// 发起调用
	callCtx, callCancel := context.WithTimeout(context.Background(), callTimeout)
	defer callCancel()

	resp, err := client.Report(callCtx, req)
	if err != nil {
		log.Fatalf("report 调用失败: %v", err)
	}

	// 打印响应
	fmt.Println("========== report 请求 ==========")
	fmt.Printf("Server      : %s\n", serverAddr)
	fmt.Printf("Namespace   : %s\n", namespace)
	fmt.Printf("Service     : %s\n", service)
	fmt.Printf("MetricTimeMs: %d\n", metricTimeMs)
	fmt.Printf("Nodes (%d):\n", len(nodes))
	for i, n := range nodes {
		fmt.Printf("  [%d] %s:%d metrics=%v\n", i, n.GetHost(), n.GetPort(), n.GetMetrics())
	}
	fmt.Println("========== report 响应 ==========")
	fmt.Printf("Code: %d\n", resp.GetCode())
	fmt.Printf("Info: %s\n", resp.GetInfo())
	fmt.Println("==================================")
}

// parseInstances 解析所有 -instance 字符串为 NodeMetrics 列表
func parseInstances(items []string, defaultPort uint32, metricTimeMs int64) ([]*schedule.NodeMetrics, error) {
	nodes := make([]*schedule.NodeMetrics, 0, len(items))
	for idx, raw := range items {
		node, err := parseOneInstance(raw, defaultPort, metricTimeMs)
		if err != nil {
			return nil, fmt.Errorf("第 %d 个 instance 解析失败: %w", idx+1, err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// parseOneInstance 解析单个实例字符串：host=IP,port=PORT,metric1=v1,metric2=v2
func parseOneInstance(raw string, defaultPort uint32, metricTimeMs int64) (*schedule.NodeMetrics, error) {
	host := ""
	port := defaultPort
	metrics := make(map[string]string)

	for _, kv := range strings.Split(raw, ",") {
		kv = strings.TrimSpace(kv)
		if kv == "" {
			continue
		}
		eq := strings.Index(kv, "=")
		if eq <= 0 {
			return nil, fmt.Errorf("非法 key=value 片段: %q", kv)
		}
		key := strings.TrimSpace(kv[:eq])
		value := strings.TrimSpace(kv[eq+1:])
		switch key {
		case "host":
			host = value
		case "port":
			p, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("port 解析失败: %v", err)
			}
			port = uint32(p)
		default:
			metrics[key] = value
		}
	}
	if host == "" {
		return nil, fmt.Errorf("缺少 host 字段, raw=%q", raw)
	}
	return &schedule.NodeMetrics{
		Host:         host,
		Port:         port,
		Metrics:      metrics,
		MetricTimeMs: metricTimeMs,
	}, nil
}
