package basic

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	h := &Human{Name: "zhang"}
	Call(h)

	d := &Dog{Name: "wang"}
	Call(d)
}

type Animal interface {
	Eat()
}

type Human struct {
	Name string
}

type Dog struct {
	Name string
}

func (h *Human) Eat() {
	fmt.Printf("%s eating", h.Name)
}

func (h *Dog) Eat() {
	fmt.Printf("%s eating", h.Name)
}

func Call(p Animal) {
	p.Eat()
}
