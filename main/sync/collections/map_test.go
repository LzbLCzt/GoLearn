package collections

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_Collections1(t *testing.T) {
	m := make(map[int]string)
	m[1] = "a"
	m[2] = "b"
	m[3] = "c"
	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go func() {
			fmt.Println("ele: ", m[i+1]) //todo 这里i没有被拷贝出一个新的变量，被goroutine和主线程共享的，会随着for循环执行发生变化
			ch <- i                      //todo 正确的做法参考Test_Collections2
		}()
	}
	for i := 0; i < 3; i++ { //保证goroutine全部执行完成
		<-ch
	}
}

func Test_Collections2(t *testing.T) {
	m := make(map[int]string)
	m[1] = "a"
	m[2] = "b"
	m[3] = "c"
	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go func(key int) {
			fmt.Println("ele: ", m[key])
			ch <- key
		}(i + 1)
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
}

func Test_Collections3(t *testing.T) {
	//todo 一个普通的map在不同goroutine是线程不安全的
	//todo 通过加锁解决，参考Test_Collections4，或使用sync.Map，参考Test_Collections5
	m := make(map[int]int)
	m[1] = 1
	m[2] = 2
	m[3] = 3
	ch := make(chan int)
	//var mutex sync.Mutex

	for i := 0; i < 3; i++ {
		go func(key int) {
			//mutex.Lock()
			//defer mutex.Unlock()
			time.Sleep(10 * time.Millisecond)
			fmt.Println("ele before modify: ", m[key])
			m[key] = key + 1
			fmt.Println("ele after modify: ", m[key])
			ch <- key
		}(i + 1)
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
}

func Test_Collections4(t *testing.T) {
	m := make(map[int]int)
	m[1] = 1
	m[2] = 2
	m[3] = 3
	ch := make(chan int)
	var mutex sync.Mutex

	for i := 0; i < 3; i++ {
		go func(key int) {
			mutex.Lock()
			defer mutex.Unlock()
			time.Sleep(10 * time.Millisecond)
			fmt.Println("ele before modify: ", m[key])
			m[key] = key + 1
			fmt.Println("ele after modify: ", m[key])
			ch <- key
		}(i + 1)
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
}

func Test_Collections5(t *testing.T) {
	m := sync.Map{}
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	ch := make(chan int)
	var mutex sync.Mutex

	for i := 0; i < 3; i++ {
		go func(key int) {
			mutex.Lock()
			defer mutex.Unlock()
			time.Sleep(10 * time.Millisecond)
			v, _ := m.Load(key)
			fmt.Println("ele before modify: ", v)
			m.Store(key, key+1)
			v2, _ := m.Load(key)
			fmt.Println("ele after modify: ", v2)
			ch <- key
		}(i + 1)
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
}
