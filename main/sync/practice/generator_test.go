package practice

import (
	"fmt"
	"testing"
)

/*
todo 开发一个自增整数生成器
*/
func addElementToChannel(ch chan<- int) {
	go func() {
		for i := 0; ; i++ {
			ch <- i //当主线程需要数据时，我们加进去chan
		}
	}()
}

func Test_generator1(t *testing.T) {
	var ch = make(chan int)
	addElementToChannel(ch)

	for i := 0; i < 100; i++ {
		fmt.Println(<-ch) //主线程从chan索要数据
	}
}
