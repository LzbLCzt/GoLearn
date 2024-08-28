package practice

import (
	"fmt"
	"testing"
)

func Test_channel_quit(t *testing.T) {
	ch, quit := make(chan int), make(chan int)
	go func(i int) {
		ch <- i
		quit <- i //发送完成信号
	}(10)

	for isFinished := false; !isFinished; {
		select {
		case v := <-ch:
			fmt.Println(v)
		case <-quit:
			isFinished = true
		}
	}
}
