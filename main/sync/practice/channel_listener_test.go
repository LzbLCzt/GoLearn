package practice

import (
	"fmt"
	"testing"
	"time"
)

func foo(i int) <-chan int {
	ch := make(chan int)
	go func() { ch <- i }()
	return ch
}
func Test_channel_listener(t *testing.T) {
	ch1, ch2, ch3 := foo(3), foo(6), foo(9)
	ch := make(chan int)
	timeout := time.After(5 * time.Second) //注册定时器，保证for循环5s后break
	go func() {
		for isTimeout := false; !isTimeout; {
			select {
			case v1 := <-ch1:
				fmt.Println("ch1: ", v1)
				ch <- v1
			case v2 := <-ch2:
				fmt.Println("ch2: ", v2)
				ch <- v2
			case v3 := <-ch3:
				fmt.Println("ch3: ", v3)
				ch <- v3
			case <-timeout:
				isTimeout = true
			}
		}
	}()
	for i := 0; i < 3; i++ {
		fmt.Println(<-ch)
	}
}

func Test_2(t *testing.T) {
	ch := make(chan int) //无缓冲channel，发生deadlock，因为没有其他goroutine接收
	ch <- 1
	fmt.Println(111)
}

func Test_3(t *testing.T) {
	ch := make(chan int, 1) //有缓冲channel，程序执行成功
	ch <- 1
	fmt.Println(111)
}
