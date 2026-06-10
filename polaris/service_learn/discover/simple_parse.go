package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println("=== 简单测试程序开始 ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("✓ 文件读取成功\n")
	fmt.Printf("✓ 文件大小: %d 字节\n", len(data))
	
	// 显示文件开头
	if len(data) > 100 {
		fmt.Printf("✓ 文件前100字节: %s\n", string(data[:100]))
	}
	
	fmt.Println("=== 测试程序结束 ===")
}