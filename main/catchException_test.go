package main

import (
	"fmt"
	"testing"
)

func TestRecoverFunction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil { //todo recover()函数捕获panic
			fmt.Println("Recovered from panic:", r)
		}
	}()

	mayPanic()

	fmt.Println("After mayPanic()")
}

func mayPanic() {
	panic("exception occurred")
}
