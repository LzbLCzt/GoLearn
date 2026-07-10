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
		dialTimeout time.Duration
		callTimeout time.Duration
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8085", "global-scheduler gRPC 服务地址")
	//flag.StringVar(&serverAddr, "addr", "9.141.112.151:8082", "global-scheduler gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "", "北极星命名空间（必填）")
	flag.StringVar(&service, "service", "", "北极星服务名（必填）")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second, "gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 3*time.Second, "Schedule 调用超时")
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
	req := &schedule.ScheduleRequest{
		Namespace:   namespace,
		Service:     service,
		Loadbalance: schedule.LoadBalanceType_GLOBAL_P2C,
		Criteria: &schedule.LoadBalanceCriteria{
			ReplicaCount: 1,
		},
	}

	// 发起调用
	callCtx, callCancel := context.WithTimeout(context.Background(), callTimeout)
	defer callCancel()

	resp, err := client.Schedule(callCtx, req)
	if err != nil {
		log.Fatalf("Schedule 调用失败: %v", err)
	}

	// 打印响应
	fmt.Println("========== Schedule 请求 ==========")
	fmt.Printf("Server     : %s\n", serverAddr)
	fmt.Printf("Namespace  : %s\n", namespace)
	fmt.Printf("Service    : %s\n", service)
	fmt.Printf("LoadBalance: %s\n", schedule.LoadBalanceType_GLOBAL_P2C.String())
	fmt.Println("========== Schedule 响应 ==========")
	fmt.Printf("Code: %d\n", resp.GetCode())
	fmt.Printf("Info: %s\n", resp.GetInfo())

	if nodes := resp.GetNodes(); len(nodes) > 0 {
		fmt.Printf("Nodes (%d):\n", len(nodes))
		for i, node := range nodes {
			fmt.Printf("  [%d] Host=%s, Port=%d, Reused=%v\n",
				i, node.GetHost(), node.GetPort(), node.GetReused())
		}
	} else {
		fmt.Println("Nodes: (空)")
	}

	if bp := resp.GetBackPressure(); bp != nil {
		fmt.Println("BackPressure:")
		fmt.Printf("  Enabled    : %v\n", bp.GetEnabled())
		fmt.Printf("  FanoutCount: %d\n", bp.GetFanoutCount())
		fmt.Printf("  Reason     : %s\n", bp.GetReason())
	} else {
		fmt.Println("BackPressure: (无)")
	}
	fmt.Println("====================================")
}
