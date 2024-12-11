package rpcWatch

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//todo 在很多系统中都提供了Watch监视功能的接口，当系统满足某种条件时Watch方法返回监控的结果。在这里我们可以尝试通过RPC框架实现一个基本的Watch功能

type KVStoreService struct {
	//m成员是一个map类型，用于存储KV数据。filter成员对应每个Watch调用时定义的过滤器函数列表。而mu成员为互斥锁，用于在多个Goroutine访问或修改时对其它成员提供保护
	m      map[string]string
	filter map[string]func(key string)
	mu     sync.Mutex
}

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

func (p *KVStoreService) Get(key string, value *string) error {
	fmt.Println("method get is called")
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}

	return fmt.Errorf("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	fmt.Println("method set is called")
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]

	if oldValue := p.m[key]; oldValue != value {
		for _, fn := range p.filter {
			fn(key)
		}
	}

	p.m[key] = value
	return nil
}

func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	fmt.Println("method watch is called")
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	ch := make(chan string, 10) // buffered

	p.mu.Lock()
	p.filter[id] = func(key string) { ch <- key }
	p.mu.Unlock()

	select {
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		return fmt.Errorf("timeout")
	case key := <-ch:	//todo 监控KVStoreService.m中的key的value是否发生变化
		*keyChanged = key
		return nil
	}

	return nil
}