package test

import (
	"fmt"
	"testing"
	"unsafe"
)

type Person struct {
	name string
	age int
}

func TestPrac(t *testing.T) {
	x := Person{name: "aaa", age: 10}
	ptr := unsafe.Pointer(&x)
	off := uintptr(ptr) + unsafe.Offsetof(x.name)
	agePtr := (*string)(unsafe.Pointer(off))
	*agePtr = "ccc"
	fmt.Println(x)
}
