// Copyright (c) 2026 Tencent. All rights reserved.
// Schedule 双 server 压测工具：同时对两台 global-scheduler 发起请求，
// 按秒统计每台 server 的有实例返回 QPS 与限流 QPS，最后输出 20 秒汇总。

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// rateLimitCode 北极星全局限流响应码。
const rateLimitCode = 403006

// serverStats 单台 server 的累计统计（按秒清零）。
type serverStats struct {
	addr           string
	client         schedule.ScheduleServiceClient
	instanceCount  int64 // code != 403006 的响应数（有实例返回）
	rateLimitCount int64 // code == 403006 的响应数（被限流）
	errCount       int64 // 网络/超时等错误
}

func main() {
	var (
		addr1       string
		addr2       string
		namespace   string
		service     string
		qps         int
		duration    time.Duration
		dialTimeout time.Duration
		callTimeout time.Duration
	)
	//flag.StringVar(&addr1, "addr1", "9.141.112.151:8082", "server1 gRPC 地址")
	flag.StringVar(&addr1, "addr1", "9.134.117.127:8082", "server1 gRPC 地址")
	flag.StringVar(&addr2, "addr2", "9.134.117.127:8082", "server2 gRPC 地址")
	flag.StringVar(&namespace, "namespace", "",
		"北极星命名空间（必填，例如 Test）")
	flag.StringVar(&service, "service", "",
		"北极星服务名（必填，例如 lzb_test）")
	flag.IntVar(&qps, "qps", 1000, "单台 server 的目标 QPS")
	flag.DurationVar(&duration, "duration", 20*time.Second,
		"压测持续时间")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second,
		"gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 3*time.Second,
		"单次 Schedule 调用超时")
	flag.Parse()

	if namespace == "" || service == "" {
		flag.Usage()
		log.Fatalf("namespace 和 service 参数必填")
	}

	// 并发拨号两台 server
	stats1 := dialServer(addr1, dialTimeout)
	stats2 := dialServer(addr2, dialTimeout)
	defer closeClient(stats1)
	defer closeClient(stats2)

	fmt.Printf("已连接 %s 和 %s, 目标=%s/%s, "+
		"单台QPS=%d, 持续=%v\n",
		addr1, addr2, namespace, service, qps, duration)
	fmt.Println(strings.Repeat("=", 80))

	stopCh := make(chan struct{})

	// 定时结束
	go func() {
		time.Sleep(duration)
		close(stopCh)
	}()

	// 每秒汇总打印（同时记录每秒快照用于最终汇总）
	type snapshot struct {
		sec        int
		ins1, rl1  int64
		ins2, rl2  int64
		err1, err2 int64
	}
	snapshots := make([]snapshot, 0, int(duration/time.Second)+2)
	var snapMu sync.Mutex

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		startTime := time.Now()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				ins1 := atomic.SwapInt64(&stats1.instanceCount, 0)
				rl1 := atomic.SwapInt64(&stats1.rateLimitCount, 0)
				err1 := atomic.SwapInt64(&stats1.errCount, 0)
				ins2 := atomic.SwapInt64(&stats2.instanceCount, 0)
				rl2 := atomic.SwapInt64(&stats2.rateLimitCount, 0)
				err2 := atomic.SwapInt64(&stats2.errCount, 0)
				sec := int(time.Since(startTime).Seconds())
				fmt.Printf("[%02ds] %s: ins=%d rl=%d err=%d "+
					"| %s: ins=%d rl=%d err=%d "+
					"| TOTAL: ins=%d rl=%d\n",
					sec, stats1.addr, ins1, rl1, err1,
					stats2.addr, ins2, rl2, err2,
					ins1+ins2, rl1+rl2)
				snapMu.Lock()
				snapshots = append(snapshots, snapshot{
					sec: sec, ins1: ins1, rl1: rl1, err1: err1,
					ins2: ins2, rl2: rl2, err2: err2,
				})
				snapMu.Unlock()
			}
		}
	}()

	// 启动两个 sendLoop
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		sendLoop(stopCh, stats1, namespace, service, qps, callTimeout)
	}()
	go func() {
		defer wg.Done()
		sendLoop(stopCh, stats2, namespace, service, qps, callTimeout)
	}()

	wg.Wait()
	// 给最后一个 ticker 一点时间冲刷剩余统计
	time.Sleep(50 * time.Millisecond)

	// 输出 20 秒汇总
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("压测结束，按秒汇总：")
	fmt.Printf("%-6s %-12s %-10s %-12s %-10s %-12s %-10s\n",
		"sec",
		stats1.addr+"-ins", stats1.addr+"-rl",
		stats2.addr+"-ins", stats2.addr+"-rl",
		"total-ins", "total-rl")
	var sumIns1, sumRL1, sumIns2, sumRL2 int64
	snapMu.Lock()
	for _, s := range snapshots {
		fmt.Printf("%-6d %-12d %-10d %-12d %-10d %-12d %-10d\n",
			s.sec, s.ins1, s.rl1, s.ins2, s.rl2,
			s.ins1+s.ins2, s.rl1+s.rl2)
		sumIns1 += s.ins1
		sumRL1 += s.rl1
		sumIns2 += s.ins2
		sumRL2 += s.rl2
	}
	cnt := int64(len(snapshots))
	snapMu.Unlock()
	if cnt == 0 {
		cnt = 1
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("TOTAL  %-12d %-10d %-12d %-10d %-12d %-10d\n",
		sumIns1, sumRL1, sumIns2, sumRL2,
		sumIns1+sumIns2, sumRL1+sumRL2)
	fmt.Printf("AVG/s  %-12d %-10d %-12d %-10d %-12d %-10d\n",
		sumIns1/cnt, sumRL1/cnt, sumIns2/cnt, sumRL2/cnt,
		(sumIns1+sumIns2)/cnt, (sumRL1+sumRL2)/cnt)
	fmt.Println(strings.Repeat("=", 80))
}

// dialServer 与单台 server 建立 gRPC 连接并返回统计上下文。
func dialServer(addr string, dialTimeout time.Duration) *serverStats {
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("连接服务器 %s 失败: %v", addr, err)
	}
	return &serverStats{
		addr:   addr,
		client: schedule.NewScheduleServiceClient(conn),
	}
}

// closeClient 关闭 gRPC 连接（通过类型断言获取底层 conn 不便，
// 这里仅作占位以满足 defer 语义；进程退出时连接会被回收）。
func closeClient(_ *serverStats) {
	// no-op：进程结束时由 runtime 释放
}

// sendLoop 以固定 QPS 持续向某台 server 发送 Schedule 请求，
// 并按响应类型累加 stats。不读取响应中的 FanoutCount。
func sendLoop(stopCh <-chan struct{}, stats *serverStats,
	namespace, service string, qps int, callTimeout time.Duration) {
	if qps <= 0 {
		return
	}
	ticker := time.NewTicker(time.Second / time.Duration(qps))
	defer ticker.Stop()
	req := &schedule.ScheduleRequest{
		Namespace:   namespace,
		Service:     service,
		Loadbalance: schedule.LoadBalanceType_GLOBAL_WRR,
	}
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			go doCall(stats, req, callTimeout)
		}
	}
}

// doCall 发起一次 Schedule 调用并按响应分类累加计数。
func doCall(stats *serverStats, req *schedule.ScheduleRequest,
	callTimeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()
	resp, err := stats.client.Schedule(ctx, req)
	if err != nil {
		atomic.AddInt64(&stats.errCount, 1)
		return
	}
	if resp.GetCode() == rateLimitCode {
		atomic.AddInt64(&stats.rateLimitCount, 1)
		return
	}
	// 按需求：只要 code != 403006 即视为"有实例返回"
	atomic.AddInt64(&stats.instanceCount, 1)
}
