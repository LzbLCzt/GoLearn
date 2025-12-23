package map_struction

import (
	"fmt"
	"testing"
)

// 定义一个包含多个字段的结构体，用于作为map的键
type ComplexKey struct {
	ID      int64
	Name    string
	Data    [100]int // 一个较大的数组字段，用于增加结构体大小
	Counter float64
}

func TestMap(t *testing.T) {
	// 创建一个键为ComplexKey结构体的map
	valueMap := make(map[ComplexKey]string)

	// 创建一个实例
	key := ComplexKey{
		ID:      1001,
		Name:    "test_key",
		Counter: 1.0,
	}

	// 1. 写入操作：会发生值拷贝
	valueMap[key] = "第一次写入的值"
	fmt.Printf("写入操作完成\n")

	// 2. 读取操作：也会发生值拷贝
	retrievedValue := valueMap[key]
	fmt.Printf("读取到的值: %s\n", retrievedValue)

	// 3. 演示修改临时副本不影响原map
	tempKey := key // 这里会发生一次拷贝
	tempKey.Name = "修改后的临时键"
	tempKey.Data[0] = 999 // 修改拷贝副本中的数据

	// 尝试用修改后的键访问map
	_, exists := valueMap[tempKey]
	fmt.Printf("用修改后的键访问map: 存在=%t\n", exists) // 输出: false

	// 原始键仍然可以访问到原始值
	originalValue := valueMap[key]
	fmt.Printf("原始键访问到的值: %s\n", originalValue) // 输出: 第一次写入的值
}
