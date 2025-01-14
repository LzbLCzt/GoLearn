package __1

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

func Test5_1_7(t *testing.T) {
	done := make(chan struct{})
	ch := GenerateIntA(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	close(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch) //todo 对于close的channel，尝试再次读值，会一直返回默认值
	//fmt.Println(<-ch)
	//fmt.Println(<-ch)

	println("NumGoroutine=", runtime.NumGoroutine())
}

func GenerateIntA(done chan struct{}) chan int {
	ch := make(chan int)
	go func() {
	Lable:
		for {
			select {
			case ch <- rand.Int():
			case <-done:
				break Lable
			}
		}
		close(ch)
	}()
	return ch
}
