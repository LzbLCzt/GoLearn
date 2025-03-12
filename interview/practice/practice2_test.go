package practice

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPractice6(t *testing.T) {
	x := 10 // x 在栈上，每个 Goroutine 独立
	go func() {
		x += 1
		fmt.Printf("goroutine中的变量值: %v\n", x)
	}()
	time.Sleep(time.Second * 2)
	fmt.Printf("main线程中的变量值: %v\n", x) //todo x被修改了
}

func TestPractice7(t *testing.T) {
	x := 10 // x 在栈上，每个 Goroutine 独立
	go func(y int) {
		y += 1
		fmt.Printf("goroutine中的变量值: %v\n", y)
	}(x)
	time.Sleep(time.Second * 2)
	fmt.Printf("main线程中的变量值: %v\n", x) //todo x没有被修改，因为值传递
}

// ----------------------------------------------------------------
// todo 竞态问题举例子
var counter int

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000000; i++ {
		counter++ // 竞态条件
	}
}
func TestPractice8(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go increment(&wg)
	go increment(&wg)
	wg.Wait()
	fmt.Println("Counter:", counter) // todo 结果可能不稳定
}

// todo 解决方案1
var counter2 atomic.Int32

func increment2(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000000; i++ {
		counter2.Add(1)
	}
}
func TestPractice9(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go increment2(&wg)
	go increment2(&wg)
	wg.Wait()
	fmt.Printf("Counter2: %v\n", counter2)
}

// todo 解决方案2
var counter3 int32
var mu3 sync.Mutex

func increment3(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000000; i++ {
		mu3.Lock()
		counter3++
		mu3.Unlock()
	}
}
func TestPractice10(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go increment3(&wg)
	go increment3(&wg)
	wg.Wait()
	fmt.Println("Counter3:", counter3)
}
