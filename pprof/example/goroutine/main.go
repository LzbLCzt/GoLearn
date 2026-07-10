// 演示：goroutine 泄漏场景，用于 /debug/pprof/goroutine 分析
//
// 设计三种典型的泄漏模式，每种都会持续产生"永远回不来"的 goroutine：
//
//  1. leakOnUnbufferedSend：往无缓冲 channel send 但没人 recv，
//     goroutine 永久阻塞在 `ch <- v`（channel send）。
//
//  2. leakOnRecvNoClose：从 channel recv 但生产者忘了 close 也不再发送，
//     goroutine 永久阻塞在 `<-ch`（channel recv）。
//
//  3. leakOnForgottenTimer：起了 goroutine 等 `time.After` 但外部逻辑先返回，
//     等到 timer 到期前 goroutine 一直挂着（模拟"上下文没传，取消不了"）。
//
// 每隔 500ms 触发一批新的泄漏 goroutine，方便对比不同时刻的 goroutine profile。
// 同时提供 /debug/pprof/goroutine 用来抓栈定位泄漏点。
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

// leakOnUnbufferedSend 每次调用产生 1 个泄漏 goroutine，卡在 send
func leakOnUnbufferedSend() {
	ch := make(chan int) // 无缓冲，没人 recv
	go func() {
		// 永远发不出去
		ch <- 42
	}()
}

// leakOnRecvNoClose 每次调用产生 1 个泄漏 goroutine，卡在 recv
func leakOnRecvNoClose() {
	ch := make(chan string)
	go func() {
		// 生产者从未 send 也未 close，这里永远阻塞
		v := <-ch
		_ = v
	}()
}

// leakOnForgottenTimer 每次调用产生 1 个泄漏 goroutine，卡在 time.After
// 真实业务里对应的常见问题：起了 goroutine 做超时等待，但没接 context，
// 外层已经返回也无法取消这个 goroutine。
func leakOnForgottenTimer() {
	go func() {
		// 用一个远大于程序生命周期的超时，模拟"实际上等不到"的场景
		<-time.After(1 * time.Hour)
	}()
}

func spawnLeakers() {
	// 每种泄漏各产生一批，制造出可观的样本便于 pprof 定位
	for i := 0; i < 5; i++ {
		leakOnUnbufferedSend()
	}
	for i := 0; i < 3; i++ {
		leakOnRecvNoClose()
	}
	for i := 0; i < 2; i++ {
		leakOnForgottenTimer()
	}
}

func main() {
	// pprof HTTP server
	go func() {
		log.Println("pprof listening on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("leaker started, press Ctrl+C to stop")
	round := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		round++
		spawnLeakers()
		// 每轮打印一次当前 goroutine 数，方便肉眼观察增长趋势
		fmt.Printf("round %d done, runtime.NumGoroutine=%d\n", round, runtime.NumGoroutine())
	}
}
