package sync

import (
	"fmt"
	"sync"
	"testing"
)

// ----------------------------------------------------
type Task struct {
	ID int
}

var taskPool = sync.Pool{
	New: func() any {
		fmt.Println("Creating new task...")
		return &Task{}
	},
}

func Test_1(t *testing.T) {
	// 第一次从池中拿对象，池是空的，会调用 New()
	task1 := taskPool.Get().(*Task)
	task1.ID = 42
	fmt.Println("Task1:", task1)

	// 放回池中
	taskPool.Put(task1)

	task2 := taskPool.Get().(*Task)
	// 第二次从池中拿对象，复用了 task1
	fmt.Println("Task2:", task2)                                // 会看到 ID = 42（所以你要记得 Reset）
	fmt.Printf("task1 is equal to task2: %t\n", task1 == task2) // true

	task2.ID = 0
	fmt.Printf("task1: %v, task2: %v\n", task1, task2)

	task3 := taskPool.Get().(*Task)
	fmt.Println("Task3:", task3) // 此时池中空了，再次get会返回一个新的对象

	task4 := taskPool.Get().(*Task)
	fmt.Println("Task4:", task4)

	fmt.Printf("task3 is equal to task4: %t\n", task3 == task4) // false
}

//----------------------------------------------------
//----------------------------------------------------
