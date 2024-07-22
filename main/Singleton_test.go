package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
todo 本文测试单例模式，用多个协程初始化Cache，是否是同一个Cache实例
*/
type Cache struct {
	sync.RWMutex
	data map[string]int
}

var (
	instance *Cache
	once     sync.Once
)

// 初始化单例对象
func GetInstance() *Cache {
	once.Do(func() {
		instance = &Cache{
			data: make(map[string]int),
		}
		instance.start()
		go instance.scheduleUpdate()
	})
	return instance
}

func (c *Cache) start() {
	fmt.Println("starting...")
}

// 从数据库更新数据
func (c *Cache) updateData() {
	// 模拟从数据库拉取数据
	newData := map[string]int{
		"service1": 100,
		"service2": 200,
	}
	c.Lock()
	c.data = newData
	c.Unlock()

	fmt.Println("Data updated:", c.data)
}

// 定时更新数据
func (c *Cache) scheduleUpdate() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			c.updateData()
		}
	}
}

// 测试以上代码正常运行
func TestA(t *testing.T) {
	cache := GetInstance()
	time.Sleep(2 * time.Hour) // 模拟程序运行，保持程序不退出
	fmt.Println(cache.data)
}

func TestSingleton(t *testing.T) {
	var wg sync.WaitGroup
	instances := make([]*Cache, 100)

	for i := 0; i < 100; i++ {
		fmt.Println(i)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			instances[i] = GetInstance()
		}(i)
	}
	wg.Wait()

	// 检查所有实例是否相同
	first := instances[0]
	for _, instance := range instances {
		if instance != first {
			t.Errorf("GetInstance returned different instances, test failed")
		}
	}
	fmt.Println("instance is singleton, test succeed")
}
