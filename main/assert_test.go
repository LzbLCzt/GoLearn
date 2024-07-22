package main

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	var p Person
	var men = Men{name: "zhengbangli", age: 18}
	p = men
	switch v := p.(type) {
	case Men:
		fmt.Println(v)
		p1 := p.(Men)
		fmt.Println(p1)
		fmt.Println("p1和p是同一个变量，地址值一样: ", p1 == p)
	default:
	}
}
