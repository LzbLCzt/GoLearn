// 演示：多个 channel 等待场景，用于 /debug/pprof/block 分析
//
// 场景说明：
//  1. slowProducer：生产者往一个无缓冲 channel 发送数据，间隔较长，
//     导致大量 consumer goroutine 在 <-ch 处阻塞。
//  2. fanIn：多个 goroutine 在 select 多路等待，等待时间较久。
//  3. barrier：一批 worker 同时等待同一个 done channel 关闭（广播场景）。
//
// 这些阻塞事件都会被 /debug/pprof/block 记录，通过采样栈可以定位到
// 具体阻塞在哪一行 channel 操作上。
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

// slowProducer 慢生产者：每 50ms 才发一个，导致 N 个 consumer 阻塞在 <-ch
func slowProducer(ch chan<- int, n int) {
	for i := 0; i < n; i++ {
		time.Sleep(50 * time.Millisecond)
		ch <- i
	}
	close(ch)
}

// consumer 从 channel 接收数据，会阻塞在 <-ch 上
func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		_ = v
		_ = id
	}
}

// fanInWorker 在多个 channel 上 select 等待
func fanInWorker(id int, a, b, c <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-a:
	case <-b:
	case <-c:
	}
	_ = id
}

// barrierWorker 等待广播信号（done 关闭）
func barrierWorker(id int, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-done
	_ = id
}

func runWorkload() {
	// ---------- 场景 1：慢生产者 + 多消费者 ----------
	ch := make(chan int) // 无缓冲
	var wg1 sync.WaitGroup
	consumerN := 50
	for i := 0; i < consumerN; i++ {
		wg1.Add(1)
		go consumer(i, ch, &wg1)
	}
	go slowProducer(ch, 20)

	// ---------- 场景 2：fan-in select 多路等待 ----------
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	var wg2 sync.WaitGroup
	fanN := 30
	for i := 0; i < fanN; i++ {
		wg2.Add(1)
		go fanInWorker(i, a, b, c, &wg2)
	}
	// 300ms 后才随便往一个 channel 发信号，让所有 worker 都阻塞一段时间
	go func() {
		time.Sleep(300 * time.Millisecond)
		for i := 0; i < fanN; i++ {
			a <- i
		}
	}()

	// ---------- 场景 3：广播 barrier ----------
	done := make(chan struct{})
	var wg3 sync.WaitGroup
	barrierN := 40
	for i := 0; i < barrierN; i++ {
		wg3.Add(1)
		go barrierWorker(i, done, &wg3)
	}
	// 500ms 后广播释放
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(done)
	}()

	wg1.Wait()
	wg2.Wait()
	wg3.Wait()
}

func main() {
	// ★ 关键：显式开启 block profile 采样，否则 /debug/pprof/block 为空
	// 传 1 表示所有阻塞事件都采样（定位问题最准，开销最大）
	runtime.SetBlockProfileRate(1)

	// pprof HTTP server
	go func() {
		log.Println("pprof listening on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// 持续跑 workload，方便随时抓 profile
	log.Println("workload started, press Ctrl+C to stop")
	round := 0
	for {
		round++
		start := time.Now()
		runWorkload()
		fmt.Printf("round %d done in %v\n", round, time.Since(start))
	}
}
