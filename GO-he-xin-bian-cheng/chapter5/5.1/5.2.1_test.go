package __1

import (
	"fmt"
	"math/rand"
	"testing"
)

//todo 并发范式-生成器

func Test5_2_1_A(t *testing.T) {
	ch := GenerateInt()
	for i := 0; i < 100; i++ {
		println(<-ch)
	}
}

func GenerateInt() chan int {
	ch := make(chan int, 20)
	go func() {
		for {
			select {
			case ch <- <-GenerateIntAA():
			case ch <- <-GenerateIntBB():
			}
		}
	}()
	return ch
}

func GenerateIntAA() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <- rand.Int()
		}
	}()
	return ch
}

func GenerateIntBB() chan int {
	ch := make(chan int, 10)
	go func() {
		for {
			ch <- rand.Int()
		}
	}()
	return ch
}

// --------------------------------------------------------------------
// todo 退出通知机制
func Test5_2_1_B(t *testing.T) {
	done := make(chan struct{})
	ch := GenerateInt2(done)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	close(done)

	for x := range ch {
		fmt.Println(x)
	}
}

func GenerateInt2(done chan struct{}) chan int {
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

// --------------------------------------------------------------------
func Test5_2_1_C(t *testing.T) {
	done := make(chan struct{})
	ch := GenerateInt3(done)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
	done <- struct{}{} //todo 通知生成器退出
	fmt.Println("stop generate")
	for x := range ch {
		fmt.Printf("after stop generate: %d\n", x)
	}
}

func GenerateInt3(done chan struct{}) chan int {
	ch := make(chan int)
	send := make(chan struct{})
	go func() {
	Label:
		for {
			select {
			case ch <- <-GenerateIntAA3(send):
			case ch <- <-GenerateIntBB3(send):
			case <-done:
				send <- struct{}{}
				send <- struct{}{}
				break Label
			}
		}
		close(ch)
	}()
	return ch
}

func GenerateIntAA3(done chan struct{}) chan int {
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

func GenerateIntBB3(done chan struct{}) chan int {
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
