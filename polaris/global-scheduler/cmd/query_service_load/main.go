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
	"time"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		serverAddr  string
		namespace   string
		service     string
		revision    string
		dialTimeout time.Duration
		callTimeout time.Duration
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8085", "global-scheduler gRPC 服务地址")
	//flag.StringVar(&serverAddr, "addr", "9.141.112.151:8082", "global-scheduler gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "", "北极星命名空间（必填）")
	flag.StringVar(&service, "service", "", "北极星服务名（必填）")
	flag.StringVar(&revision, "revision", "", "客户端缓存的服务 revision（可选）")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second, "gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 3*time.Second, "QueryServiceLoad 调用超时")
	flag.Parse()

	if namespace == "" || service == "" {
		flag.Usage()
		log.Fatalf("namespace 和 service 参数必填")
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
	client := schedule.NewScheduleServiceClient(conn)

	// 构造请求
	req := &schedule.ServiceLoadRequest{
		Namespace: namespace,
		Service:   service,
		Revision:  revision,
	}

	// 发起调用
	callCtx, callCancel := context.WithTimeout(context.Background(), callTimeout)
	defer callCancel()

	resp, err := client.QueryServiceLoad(callCtx, req)
	if err != nil {
		log.Fatalf("QueryServiceLoad 调用失败: %v", err)
	}

	// 打印响应
	fmt.Println("========== QueryServiceLoad 请求 ==========")
	fmt.Printf("Server   : %s\n", serverAddr)
	fmt.Printf("Namespace: %s\n", namespace)
	fmt.Printf("Service  : %s\n", service)
	fmt.Printf("Revision : %s\n", revision)
	fmt.Println("========== QueryServiceLoad 响应 ==========")
	fmt.Printf("Code    : %d\n", resp.GetCode())
	fmt.Printf("Info    : %s\n", resp.GetInfo())
	fmt.Printf("Revision: %s\n", resp.GetRevision())

	if instances := resp.GetInstances(); len(instances) > 0 {
		fmt.Printf("Instances (%d):\n", len(instances))
		for i, ins := range instances {
			fmt.Printf("  [%d] Host=%s, Port=%d, Load=%d, IsOverload=%v\n",
				i, ins.GetHost(), ins.GetPort(), ins.GetLoad(), ins.GetIsOverload())
		}
	} else {
		fmt.Println("Instances: (空)")
	}
	fmt.Println("============================================")
}
