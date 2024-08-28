package practice

import (
	"fmt"
	"testing"
	"time"
)

func Timer(duration time.Duration) chan bool {
	ch := make(chan bool)
	go func() {
		time.Sleep(duration)
		ch <- true
	}()
	return ch
}

func Test_timer(t *testing.T) {
	timer := Timer(3 * time.Second)
	select {
	case <-timer:
		fmt.Println("already 3s")
		return
	}
}
