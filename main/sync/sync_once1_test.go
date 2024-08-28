package sync

import (
	"fmt"
	"sync"
	"testing"
)

func Test_1(t *testing.T) {
	var once sync.Once
	onceBody := func() {
		fmt.Println("test only once, 只打印一次")
	}
	c := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			once.Do(onceBody) //确保只执行一次
			c <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-c
	}
}
