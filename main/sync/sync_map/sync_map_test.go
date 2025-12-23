package sync_map

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
	sync.Map 只能通过Load/Store 来保证并发安全，如果sync.Map 中存储的元素是指针，那么指针指向的内容是不安全的，
	需要额外的锁来保证安全
*/

type Value struct {
	Name string
	Age  int
}

func TestSyncMapNotSafe(t *testing.T) {
	m := sync.Map{}

	m.Store("a", &Value{
		Name: "Alice",
		Age:  18,
	})

	if v, ok := m.Load("a"); ok {
		val := v.(*Value)
		val.Age = 20
		fmt.Printf("v1: %v\n", val)
	}

	if v, ok := m.Load("a"); ok {
		val := v.(*Value)
		fmt.Printf("v2: %v\n", val)
	}
}

type Value2 struct {
	Name string
	Age  int
	lock sync.RWMutex
}

func TestSyncMapSafe(t *testing.T) {
	m := sync.Map{}

	m.Store("a", &Value2{
		Name: "Alice",
		Age:  18,
	})

	if v, ok := m.Load("a"); ok {
		val := v.(*Value2)
		val.lock.Lock()
		val.Age = 20
		val.lock.Unlock()
		fmt.Printf("v1: %v\n", val)
	}

	if v, ok := m.Load("a"); ok {
		val := v.(*Value2)
		fmt.Printf("v2: %v\n", val)
	}
}

var addresses = []string{"a", "b", "c"}

func Test2(t *testing.T) {
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		ticker := time.NewTicker(time.Millisecond * 10)
		defer ticker.Stop()
		for range ticker.C {
			tmp := []string{"a", "b", "c"}
			addresses = tmp
		}
	}()

	i := 0
	for {
		i++
		fmt.Printf("leng of addresses: %d\n", len(addresses))
		if i == 10000 {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	wg.Wait()
}
