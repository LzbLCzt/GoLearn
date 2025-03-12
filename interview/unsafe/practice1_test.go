package unsafe

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"unsafe"
)
/*
unsafe.Pointer的使用场景
 */
// unsafe.Pointer与任意指针类型(*int、*string)可相互转换
func TestPractice1(t *testing.T) {
	var x int = 42
	ptr := unsafe.Pointer(&x)	//todo 将 *int 转换为 unsafe.Pointer
	p := (*int32)(ptr)	//todo 将 unsafe.Pointer 转换为 *int32
	assert.Equal(t, *p, int32(42))

	*p = 100	//todo 直接修改内存（危险）
	assert.Equal(t, *p, int32(100))
}
// ----------------------------------------------------
//2. 通过unsafe.Pointer + uintptr 计算内存偏移量
type Person struct {
	age int
	name string
}
func TestPractice2(t *testing.T) {
	p := Person{age: 10, name: "mike"}
	ptr := unsafe.Pointer(&p)
	start := uintptr(ptr)	//todo 获取结构体p的起始内存地址
	offset := unsafe.Offsetof(p.age)	//todo 计算 age 字段相对于 p 结构体起始地址的偏移量
	agePtr := (*int)(unsafe.Pointer(start + offset))	//todo 计算age字段的内存地址
	*agePtr = 1000	//todo 直接修改内存
	assert.Equal(t, p.age, 1000)
}

func TestPractice3(t *testing.T) {
	p := Person{age: 1, name: "lzb"}
	ptr := unsafe.Pointer(&p)
	namePtr := (*string)(unsafe.Pointer(uintptr(ptr) + unsafe.Offsetof(p.name)))
	*namePtr = "czt"
	assert.Equal(t, p.name, "czt")
}
// ---------------------------------------------------------------------