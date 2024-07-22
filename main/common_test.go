package main

import (
	"fmt"
	"testing"
	"time"
)

func Echo(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func Test_1(t *testing.T) {
	go Echo("协程执行")
	Echo("主线程执行")
}

func Test_2(t *testing.T) {
	ch := make(chan string)
	go func() {
		a := <-ch
		fmt.Println(a)
	}()
	ch <- "sleep"

}
