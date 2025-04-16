package time

import (
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestSchedule(t *testing.T) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case h := <-ticker.C:
			fmt.Println("time added: ", h)
			fmt.Println("cur time: ", time.Now())
			fmt.Println()
		}
	}
}

func TestSchedule2(t *testing.T) {
	ch := make(chan int)
	go func(i int) {
		ch <- i
	}(10)
	timeout := time.After(5 * time.Second) //返回一个只发送channel，到了时间就会发信号出来
	for isTimeout := false; !isTimeout; {  //保证for不会持续循环，最多5s后跳出
		select {
		case <-ch:
		case <-timeout:
			isTimeout = true
		}
	}
}

func Test(t *testing.T) {
	pool := sync.Pool{
		New: func() interface{} {
			return &Task{}
		},
	}

	task1 := pool.Get().(*Task)
	task1.id = 1
	fmt.Println(task1)

	pool.Put(task1)

	task2 := pool.Get().(*Task)
	fmt.Println(task2)
	fmt.Println(task2 == task1)

	task3 := pool.Get().(*Task)
	fmt.Println(task3)
	fmt.Println(task3 == task1)
	fmt.Println(task3 == task2)
}

func TestPractice1(t *testing.T) {
	var x int = 42
	ptr := unsafe.Pointer(&x) //todo 将 *int 转换为 unsafe.Pointer
	p := (*int32)(ptr)        //todo 将 unsafe.Pointer 转换为 *int32
	assert.Equal(t, *p, int32(42))

	*p = 100 //todo 直接修改内存（危险）
	assert.Equal(t, *p, int32(100))
}

type Task struct {
	id int
}

type Person struct {
	age  int
	name string
}

func TestPractice2(t *testing.T) {
	p := Person{age: 10, name: "mike"}
	ptr := unsafe.Pointer(&p)
	start := uintptr(ptr)                            //todo 获取结构体p的起始内存地址
	offset := unsafe.Offsetof(p.age)                 //todo 计算 age 字段相对于 p 结构体起始地址的偏移量
	agePtr := (*int)(unsafe.Pointer(start + offset)) //todo 计算age字段的内存地址
	*agePtr = 1000                                   //todo 直接修改内存
	assert.Equal(t, p.age, 1000)
}
