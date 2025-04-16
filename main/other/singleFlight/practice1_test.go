package sync

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ----------------------------------------------------
// todo 基础用法 - 合并重复请求
func TestPractice1(t *testing.T) {
	var sg singleflight.Group
	var wg sync.WaitGroup
	var executed int32 // 原子计数器标记实际执行次数

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ret, err, shared := sg.Do("key", func() (interface{}, error) {
				atomic.AddInt32(&executed, 1) // 记录实际执行次数
				res := "value"                //mock一个数据库的结果
				fmt.Printf("请求 %d 实际执行了操作\n", id)
				time.Sleep(time.Second)
				return res, nil
			})

			time.Sleep(3 * time.Second)
			fmt.Printf("请求 %d 结果: %v, 错误: %v, 是否是共享过来的结果: %v\n",
				id, ret, err, shared)
		}(i)
	}
	wg.Wait()
	fmt.Println("实际执行次数:", executed) // 预期为 1
}

//----------------------------------------------------
//----------------------------------------------------
