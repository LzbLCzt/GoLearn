package time

import (
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case h := <-ticker.C:
			fmt.Println("time added: ", h)
			fmt.Println("cur time: ", time.Now())
			fmt.Println()
		}
	}
}

func TestSchedule2(t *testing.T) {
	ch := make(chan int)
	go func(i int) {
		ch <- i
	}(10)
	timeout := time.After(5 * time.Second) //返回一个只发送channel，到了时间就会发信号出来
	for isTimeout := false; !isTimeout; {  //保证for不会持续循环，最多5s后跳出
		select {
		case <-ch:
		case <-timeout:
			isTimeout = true
		}
	}
}
