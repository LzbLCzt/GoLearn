package sync

import (
	"fmt"
	"sync"
	"testing"
)
/*
	复用字节切片
 */
// 声明一个全局的 sync.Pool 用于存放字节切片
var bufferPool = sync.Pool{
	New: func() interface{} {
		// 当池中无可用对象时，New函数会被调用来创建新对象
		return make([]byte, 0, 1024) // 初始容量为1024的空切片
	},
}

func getBuffer() []byte {
	return bufferPool.Get().([]byte) // 从池中获取一个字节切片
}

func putBuffer(buf []byte) {
	buf = buf[:0]       // 重置切片长度，清空内容（重要！）
	bufferPool.Put(buf) // 将切片放回池中
}

func TestSyncPool2(t *testing.T) {
	// 模拟处理数据：假设需要频繁使用缓冲区
	data := []byte("Hello, World!")

	// 获取一个缓冲区（可能是新的，也可能是复用的）
	buf := getBuffer()
	defer putBuffer(buf) // 确保使用完毕后归还

	// 使用缓冲区
	buf = append(buf, data...)
	fmt.Println("Processed data:", string(buf))
}
