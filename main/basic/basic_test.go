package basic

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestBasic(t *testing.T) {
	var p Person
	h := Human{Name: "shennong"}
	p = h
	fmt.Printf("type of p: %T \n", p)

	tt := reflect.TypeOf(p)
	fmt.Printf("type of p: %v \n", tt)

}

type Person interface {
	Eat()
}

type Human struct {
	Name string
}

func (h Human) Eat() {
	fmt.Printf("%s eating", h.Name)
}

func Test2(t *testing.T) {
	println("GOMAXPROCS=", runtime.GOMAXPROCS(0))
}
