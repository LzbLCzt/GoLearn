package __1

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test5_2_1(t *testing.T) {
	done := make(chan struct{})
	ch := GenerateInt(done)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	fmt.Println("关闭done")
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func GenerateInt(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
	Label:
		for {
			select {
			case ch <- rand.Int():
			case <-done: // todo 监听done信号，一旦done管道关闭，跳出循环
				break Label
			}
		}
		close(ch)
		fmt.Println("跳出循环")
	}()
	return ch
}

// --------------------------------------------------------
func Test5_2_2(t *testing.T) {
	done := make(chan struct{})
	ch := Generate(done)

	for i := 0; i < 100; i++ {
		fmt.Println(<-ch)
	}
}

func GenerateA(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
	Label:
		for {
			select {
			case ch <- rand.Int():
			case <-done:
				break Label
			}
		}
		close(ch)
	}()
	return ch
}

func GenerateB(done chan struct{}) chan int {
	ch := make(chan int, 10)
	go func() {
	Label:
		for {
			select {
			case ch <- rand.Int():
			case <-done:
				break Label
			}
		}
		close(ch)
	}()
	return ch
}

func Generate(done chan struct{}) chan int {
	ch := make(chan int, 20)
	go func() {
	Label:
		for {
			select {
			case ch <- <-GenerateA(done):
			case ch <- <-GenerateB(done):
			case <-done:
				break Label
			}
		}
		close(ch)
	}()
	return ch
}
