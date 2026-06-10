package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"google.golang.org/protobuf/encoding/prototext"
)

func main() {
	fmt.Println("=== 开始解析 main.data1 文本格式protobuf文件 ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	
	// 分析文件内容
	content := string(data)
	
	fmt.Println("\n=== 文件内容概览 ===")
	
	// 提取关键信息
	lines := strings.Split(content, "\n")
	fmt.Printf("文件总行数: %d\n", len(lines))
	
	// 显示前20行内容
	fmt.Println("\n前20行内容:")
	for i := 0; i < min(20, len(lines)); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			fmt.Printf("%3d: %s\n", i+1, line)
		}
	}
	
	// 尝试解析为文本格式protobuf
	fmt.Println("\n=== 尝试文本格式protobuf解析 ===")
	
	// 由于文件可能包含多个消息，尝试解析整个文件
	tryTextProtoParse(content, "完整文件")
	
	// 尝试解析文件的一部分（前1000行）
	partialContent := strings.Join(lines[:min(1000, len(lines))], "\n")
	tryTextProtoParse(partialContent, "前1000行")
	
	// 尝试解析单个消息（查找可能的消息边界）
	fmt.Println("\n=== 查找单个protobuf消息 ===")
	findAndParseSingleMessages(content)
	
	fmt.Println("\n=== 解析完成 ===")
}

func tryTextProtoParse(content, desc string) {
	fmt.Printf("\n尝试解析 %s (%d 字节)...\n", desc, len(content))
	
	// 尝试v2.DiscoverResponse
	rsp := &v2.DiscoverResponse{}
	if err := prototext.Unmarshal([]byte(content), rsp); err != nil {
		fmt.Printf("✗ %s - v2.DiscoverResponse 解析失败: %v\n", desc, err)
	} else {
		fmt.Printf("✓ %s - v2.DiscoverResponse 解析成功！\n", desc)
		fmt.Printf("  消息类型: %T\n", rsp)
	}

	// 尝试apiV1Model.CircuitBreakerV2
	cb := &apiV1Model.CircuitBreakerV2{}
	if err := prototext.Unmarshal([]byte(content), cb); err != nil {
		fmt.Printf("✗ %s - apiV1Model.CircuitBreakerV2 解析失败: %v\n", desc, err)
	} else {
		fmt.Printf("✓ %s - apiV1Model.CircuitBreakerV2 解析成功！\n", desc)
		fmt.Printf("  消息类型: %T\n", cb)
	}
}

func findAndParseSingleMessages(content string) {
	lines := strings.Split(content, "\n")
	
	// 查找可能的protobuf消息开始
	var messageStarts []int
	for i, line := range lines {
		if strings.Contains(line, "serviceName") || strings.Contains(line, "workloadName") {
			messageStarts = append(messageStarts, i)
		}
	}
	
	fmt.Printf("发现 %d 个可能的protobuf消息开始位置\n", len(messageStarts))
	
	if len(messageStarts) > 0 {
		// 尝试解析第一个消息
		start := messageStarts[0]
		end := min(start+50, len(lines))
		messageContent := strings.Join(lines[start:end], "\n")
		
		fmt.Printf("\n第一个消息内容（行 %d-%d）:\n", start+1, end)
		fmt.Println(messageContent)
		
		tryTextProtoParse(messageContent, "第一个消息")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}