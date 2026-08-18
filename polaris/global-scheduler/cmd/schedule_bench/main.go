package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.woa.com/polaris/polaris-server-api/api/schedule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		serverAddr     string
		namespace      string
		service        string
		qps            int
		duration       time.Duration
		dialTimeout    time.Duration
		callTimeout    time.Duration
		simulateFanout bool
	)
	flag.StringVar(&serverAddr, "addr", "9.134.117.127:8085",
		"global-scheduler gRPC 服务地址")
	flag.StringVar(&namespace, "namespace", "", "北极星命名空间（必填）")
	flag.StringVar(&service, "service", "", "北极星服务名（必填）")
	flag.IntVar(&qps, "qps", 300, "目标 QPS（可运行时通过 stdin 动态调整）")
	flag.DurationVar(&duration, "duration", 60*time.Second,
		"压测持续时间（默认60s，0表示无限）")
	flag.DurationVar(&dialTimeout, "dial-timeout", 5*time.Second,
		"gRPC 拨号超时")
	flag.DurationVar(&callTimeout, "call-timeout", 3*time.Second,
		"单次 Schedule 调用超时")
	flag.BoolVar(&simulateFanout, "simulate-fanout", false,
		"模拟SDK fanout分散：收到fanout=N后实际发送QPS=target/N")
	flag.Parse()

	if namespace == "" || service == "" {
		flag.Usage()
		log.Fatalf("namespace 和 service 参数必填")
	}

	// 建立 gRPC 连接
	dialCtx, dialCancel := context.WithTimeout(
		context.Background(), dialTimeout)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("连接服务器 %s 失败: %v", serverAddr, err)
	}
	defer conn.Close()

	client := schedule.NewScheduleServiceClient(conn)
	fmt.Printf("已连接 %s, 目标: %s/%s, QPS=%d, 持续=%v\n",
		serverAddr, namespace, service, qps, duration)
	fmt.Println("提示: 输入数字回车可动态调整 QPS，输入 q 退出")
	fmt.Println(strings.Repeat("=", 70))

	// 原子变量控制目标 QPS
	var targetQPS int64
	atomic.StoreInt64(&targetQPS, int64(qps))

	// 统计计数器
	var (
		sentCount      int64
		errCount       int64
		rateLimitCount int64
		lastFanout     uint32 = 1 // 初始均摊数为1
		fanoutChange   int64
	)

	// 停止信号
	stopCh := make(chan struct{})

	// stdin 读取协程：动态调整 QPS
	go readStdin(&targetQPS, stopCh)

	// 定时结束
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			close(stopCh)
		}()
	}

	// 每秒汇总打印协程
	go printStats(stopCh, &sentCount, &errCount, &rateLimitCount,
		&lastFanout, &fanoutChange, &targetQPS)

	// 发送请求主循环
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sendLoop(stopCh, client, namespace, service, callTimeout,
			&targetQPS, &sentCount, &errCount, &rateLimitCount,
			&lastFanout, &fanoutChange, simulateFanout)
	}()

	wg.Wait()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("压测结束")
}

// readStdin 从标准输入读取 QPS 调整指令
func readStdin(targetQPS *int64, stopCh chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "q" || line == "quit" {
			close(stopCh)
			return
		}
		newQPS, parseErr := strconv.Atoi(line)
		if parseErr != nil || newQPS < 0 {
			fmt.Printf("[!] 无效输入: %q, 请输入正整数或 q\n", line)
			continue
		}
		atomic.StoreInt64(targetQPS, int64(newQPS))
		fmt.Printf("[*] QPS 已调整为 %d\n", newQPS)
	}
}

// printStats 每秒汇总打印统计信息
func printStats(stopCh chan struct{}, sentCount, errCount,
	rateLimitCount *int64, lastFanout *uint32,
	fanoutChange, targetQPS *int64) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	startTime := time.Now()
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			sent := atomic.SwapInt64(sentCount, 0)
			errs := atomic.SwapInt64(errCount, 0)
			rl := atomic.SwapInt64(rateLimitCount, 0)
			fc := atomic.LoadUint32(lastFanout)
			changed := atomic.SwapInt64(fanoutChange, 0)
			elapsed := time.Since(startTime).Truncate(time.Second)
			tgt := atomic.LoadInt64(targetQPS)
			marker := ""
			if changed > 0 {
				marker = " ← CHANGED!"
			}
			fmt.Printf("[%v] target=%d sent=%d err=%d "+
				"rateLimit=%d fanout=%d%s\n",
				elapsed, tgt, sent, errs, rl, fc, marker)
		}
	}
}

// sendLoop 持续发送 Schedule 请求
func sendLoop(stopCh chan struct{},
	client schedule.ScheduleServiceClient,
	namespace, service string, callTimeout time.Duration,
	targetQPS *int64, sentCount, errCount, rateLimitCount *int64,
	lastFanout *uint32, fanoutChange *int64,
	simulateFanout bool) {
	// 使用动态 ticker 控制发送速率
	var ticker *time.Ticker
	var currentEffectiveQPS int64
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		tgt := atomic.LoadInt64(targetQPS)
		if tgt <= 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// 模拟 fanout 分散：实际发送 QPS = target / fanout
		effectiveQPS := tgt
		if simulateFanout {
			fc := int64(atomic.LoadUint32(lastFanout))
			if fc > 1 {
				effectiveQPS = tgt / fc
				if effectiveQPS < 1 {
					effectiveQPS = 1
				}
			}
		}
		// QPS 变化时重建 ticker
		if effectiveQPS != currentEffectiveQPS {
			if ticker != nil {
				ticker.Stop()
			}
			ticker = time.NewTicker(time.Second / time.Duration(effectiveQPS))
			currentEffectiveQPS = effectiveQPS
		}
		select {
		case <-stopCh:
			if ticker != nil {
				ticker.Stop()
			}
			return
		case <-ticker.C:
			// 异步发送，不阻塞 ticker
			go func() {
				callCtx, cancel := context.WithTimeout(
					context.Background(), callTimeout)
				resp, callErr := client.Schedule(callCtx,
					&schedule.ScheduleRequest{
						Namespace:   namespace,
						Service:     service,
						Loadbalance: schedule.LoadBalanceType_GLOBAL_WRR,
					})
				cancel()
				if callErr != nil {
					atomic.AddInt64(errCount, 1)
					return
				}
				// 检测限流响应（code=403006）
				if resp.GetCode() == 403006 {
					atomic.AddInt64(rateLimitCount, 1)
					// 限流响应中也可能携带 BackPressure，需要解析以更新 fanout
					if bp := resp.GetBackPressure(); bp != nil {
						fc := bp.GetFanoutCount()
						old := atomic.SwapUint32(lastFanout, fc)
						if old != fc {
							fmt.Printf("  >> fanout 变更: %d → %d (enabled=%v, reason=%s)\n", old, fc, bp.GetEnabled(), bp.GetReason())
						}
					}
					return
				}
				atomic.AddInt64(sentCount, 1)
				if bp := resp.GetBackPressure(); bp != nil {
					fc := bp.GetFanoutCount()
					old := atomic.SwapUint32(lastFanout, fc)
					if old != fc {
						atomic.AddInt64(fanoutChange, 1)
						fmt.Printf("  >> fanout 变更: %d → %d"+
							" (enabled=%v, reason=%s)\n",
							old, fc, bp.GetEnabled(),
							bp.GetReason())
					}
				}
			}()
		}
	}
}
