package main

import (
	"bytes"
	"fmt"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"runtime"
	"strconv"
)

type Person interface {
}

type Men struct {
	name string
	age  int
}

func main() {
	var p Person
	m := Men{name: "a", age: 1}
	p = m
	fmt.Println(p)
	if e, ok := p.(*errs.Error); ok {
		fmt.Println(e)
	}
	str := fmt.Sprintf("aa%saa", "z")
	fmt.Println(str)
}

func GetGoroutineID() int {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.Atoi(string(b))
	return n
}

func fibonacci(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

var whitelistApi = map[string]struct{}{
	"/oss/sdk/version":            {},
	"/oss/auth/GetWebAccessByRTX": {},
	"/oss/common/getttl":          {},
	"/oss/common/updatattl":       {},
	"/oss/common/list":            {},
	"/oss/open/queryalarm":        {},
	"/oss/auth/ListCRole":         {},
}

var whitelist = []string{
	"a",
	"b",
	"c",
}
