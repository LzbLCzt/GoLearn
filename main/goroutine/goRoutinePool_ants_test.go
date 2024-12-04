package goroutine

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"testing"
	"time"
)

/*
todo 常见的协程池扩展库: ants
*/

func TestAnts(t *testing.T) {
	pool, err := ants.NewPool(10)
	if err != nil {
		panic(err)
	}

	//等待所有任务完成
	defer pool.Release()

	for i := 0; i < 20; i++ {
		i := i
		err = pool.Submit(func() {
			fmt.Printf("task:%d is running...\n", i)
			time.Sleep(1 * time.Second)
		})
		if err != nil {
			fmt.Printf("submit task: %d failed: %v\n", i, err)
		}
	}
}

/*
todo 创建一个带有特定任务函数的 Goroutine 池: ants.NewPoolWithFunc
*/
func myTaskFunc(payload interface{}) {
	num := payload.(int) // 类型断言，将 interface{} 转换为实际的类型
	fmt.Println("Processing number:", num)
}

func TestAnts2(t *testing.T) {
	poolSize := 10
	p, err := ants.NewPoolWithFunc(poolSize, myTaskFunc)
	if err != nil {
		fmt.Println("Failed to create goroutine pool:", err)
		return
	}
	defer p.Release()

	// 提交任务到池中
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		num := i
		_ = p.Invoke(num) // 提交任务
		wg.Done()
	}

	wg.Wait()
	fmt.Println("All tasks are processed.")
}
