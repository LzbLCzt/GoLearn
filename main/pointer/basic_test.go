package pointer

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

// go中可以修改类型为基础变量的实参的值，java则不行（因为java中基础类型变量默认是值传递）
func TestBasic(t *testing.T) {
	a := 10
	modify(&a)
	assert.Equal(t, a, 110)
}

func modify(a *int) {
	*a = *a + 100
}
