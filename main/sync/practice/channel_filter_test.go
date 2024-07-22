package practice

import (
	"fmt"
	"testing"
)

const max = 100

func integerGenerator() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

func filter(ch chan int, num int) chan int {
	out_ch := make(chan int)
	go func() {
		for {
			ele := <-ch
			if ele%num != 0 {
				out_ch <- ele
			}
		}
	}()
	return out_ch
}

func Test_channel_filter(t *testing.T) {
	ch := integerGenerator()
	num := <-ch
	for num <= max {
		fmt.Println(num)
		next_ch := filter(ch, num)
		num = <-next_ch
	}
}
