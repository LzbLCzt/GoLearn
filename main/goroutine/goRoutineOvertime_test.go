package goroutine

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

/*
goroutine设置超时机制
*/
func TestGoRoutineOvertime(t *testing.T) {
	for i := 0; i < 1000; i++ {
		startRoutineWithOvertime(doSomething)
	}
	time.Sleep(2 * time.Second)
	t.Log("当前goroutine数量: ", runtime.NumGoroutine())	//查看运行中的goroutine数量
}

func doSomething(done chan bool) {
	time.Sleep(time.Second)
	done <- true
}

func startRoutineWithOvertime(f func(chan bool)) {
	c := make(chan bool, 1)
	go f(c)
	select {
	case <- c:
		fmt.Println("函数doSomething正常执行完成，无超时发生")
	case <- time.After(time.Millisecond):	//设置1ms的超时时间
		fmt.Println("函数doSomething执行超时")
	}
}