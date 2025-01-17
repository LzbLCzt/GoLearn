package sample

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSample(b *testing.T) {
	a := "1"
	t := reflect.TypeOf(a)
	fmt.Println(t.Name())
	fmt.Println(t.Kind())
	fmt.Println(t.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()))
	fmt.Println(t.ConvertibleTo(reflect.TypeOf("1")))
	fmt.Println(t.NumMethod())
	//fmt.Println(t.Method(0))
	fmt.Println(t.PkgPath())
	fmt.Println(t.Size())
}
