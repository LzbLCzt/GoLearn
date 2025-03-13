package unsafe

import (
	"fmt"
	"testing"
	"unsafe"
)
/*
todo unsafe.go主要函数
 */
// ----------------------------------------------
//todo unsafe.Sizeof 是 Go 语言 unsafe 包中的一个函数，它用于 计算变量或类型的大小（以字节为单位）
// func Sizeof(x ArbitraryType) uintptr
// x 是一个 变量或类型，可以是基础类型、结构体、数组等。
// 返回值是 uintptr 类型，表示这个值或类型在内存中占用的字节数。
func Test2Practice1(t *testing.T) {
	var a int
	var b int8
	var c int32
	var d float64
	var e bool

	fmt.Println("Size of int:", unsafe.Sizeof(a))    // 8 (在 64 位系统上)
	fmt.Println("Size of int8:", unsafe.Sizeof(b))   // 1
	fmt.Println("Size of int32:", unsafe.Sizeof(c))  // 4
	fmt.Println("Size of float64:", unsafe.Sizeof(d)) // 8
	fmt.Println("Size of bool:", unsafe.Sizeof(e))   // 1
}

type Person2 struct {
	Age   int32
	Score float64
}
func Test2Practice2(t *testing.T) {
	var p Person2
	fmt.Println("Size of struct Person:", unsafe.Sizeof(p)) // 输出 16: 为什么 unsafe.Sizeof(p) 返回 16 而不是 12？因为 Go 结构体会进行内存对齐，int32 后面会填充 4 字节，以保证 float64 从 8 字节对齐的地址开始

}
func Test2Practice3(t *testing.T) {
	arr := [5]int32{1,2,3,4,5}
	fmt.Println("Size of arr:", unsafe.Sizeof(arr)) // 20, 数组没有额外的填充，所以 unsafe.Sizeof(arr) 返回 20。
}

func Test2Practice4(t *testing.T) {
	var p1 *int
	var p2 *Person
	fmt.Println("Size of pointer:", unsafe.Sizeof(p1)) // 8 (在 64 位系统上), 指针的大小 跟指向的类型无关，在 64 位系统上，指针占 8 字节，在 32 位系统上占 4 字节
	fmt.Println("Size of pointer:", unsafe.Sizeof(p2)) // 8 (在 64 位系统上)
}
// 切片的 unsafe.Sizeof
func Test2Practice5(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("Size of slice:", unsafe.Sizeof(slice)) // 24 (在 64 位系统上)
	/* 切片（slice）本质上是一个结构体, 包含:
	指向底层数组的指针（8 字节）
	长度（len）（8 字节）
	容量（cap）（8 字节）
	type slice struct {
	    array unsafe.Pointer // 底层数组的指针
	    len   int            // 当前切片的长度
	    cap   int            // 当前切片的容量
	}
	但 unsafe.Sizeof(slice) 只计算 slice 结构体的大小，不包括底层数组
	*/
}

func Test2Practice5_1(t *testing.T) {
	s := []int{1,2,3,4,5}
	ptr := unsafe.Pointer(&s)
	sliceHeader := (*struct {
		Data uintptr
		Len int
		Cap int
	})(ptr)
	/*
	Data: 底层数组的起始地址（十六进制表示）。
	Len: 切片的长度，即 len(s) == 5。
	Cap: 切片的容量，即 cap(s) == 5。
	 */
	fmt.Printf("Data: %x, Len: %d, Cap: %d\n", sliceHeader.Data, sliceHeader.Len, sliceHeader.Cap)
}

//todo 通过array指针创建slice
func Test2Practice5_2(t *testing.T) {
	arr := [5]int{1,2,3,4,5}
	ptr := unsafe.Pointer(&arr[0])
	slice := unsafe.Slice((*int)(ptr), 3)
	fmt.Println(slice)
}

func Test2Practice6(t *testing.T) {
	str := "hello"
	fmt.Println("Size of string:", unsafe.Sizeof(str)) // 16 (在 64 位系统上)
	/*
	Go 的 string 本质上是一个结构体：
		type string struct {
			data *byte  // 指向字符串数据的指针 (8 字节)
			len  int    // 字符串长度 (8 字节)
		}
	所以，在 64 位系统 上，string 结构体占 8 + 8 = 16 字节
	unsafe.Sizeof(str) 只计算 string 结构体本身的大小，不包括字符串内容的大小
	 */
}

func Test2Practice6_1(t *testing.T) {
	str := "Hello, Go"
	ptr := unsafe.Pointer(&str)
	strHeader := (*struct {
		Data uintptr
		Len int
	})(ptr)
	/*
	Data：底层字符串的起始内存地址（十六进制表示）。
	Len：字符串的长度，即 len(str) == 10。
	 */
	fmt.Printf("Data: %x, Len: %d\n", strHeader.Data, strHeader.Len)
}
// -----------------------------------------------
