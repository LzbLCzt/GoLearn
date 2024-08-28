package practice

import (
	"fmt"
	"testing"
)

func randGenerator() chan int {
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()
	return ch
}

func Test_rand_generator(t *testing.T) {
	ch := randGenerator()
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
}
