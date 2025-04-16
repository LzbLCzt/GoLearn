package sync

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ----------------------------------------------------
// todo 防止缓存击穿 - 合并重复数据库查询
func TestPractice2(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup

	// 模拟 10 个并发请求
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val, _ := cache.Get("user_123")
			fmt.Printf("协程 %d 获取结果: %s\n", id, val)
		}(i)
		time.Sleep(2 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("实际执行次数:", cache.executed) // todo 实际执行次数: 239
}

type Cache struct {
	mu       sync.Mutex
	data     map[string]string
	group    singleflight.Group
	executed int32 // 原子计数器标记实际执行次数
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// 获取数据（缓存未命中时合并请求）
func (c *Cache) Get(key string) (string, error) {
	//c.mu.Lock()
	//if val, ok := c.data[key]; ok {
	//	c.mu.Unlock()
	//	return val, nil
	//}
	//c.mu.Unlock()

	// 使用 singleflight 防止缓存击穿
	val, err, _ := c.group.Do(key, func() (interface{}, error) {
		atomic.AddInt32(&c.executed, 1)
		fmt.Printf("缓存未命中，实际查询数据库: %s\n", key)
		time.Sleep(100 * time.Millisecond) // 模拟数据库查询
		result := "数据库结果-" + key

		// 更新缓存
		c.mu.Lock()
		defer c.mu.Unlock()
		c.data[key] = result
		return result, nil
	})

	return val.(string), err
}

//----------------------------------------------------
//----------------------------------------------------
