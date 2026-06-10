package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("=== 开始解析 main.data1 文件（文本格式protobuf） ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	
	// 移除UTF-8 BOM标记
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbf && data[2] == 0xbd {
		data = data[3:]
		fmt.Println("✓ 已移除UTF-8 BOM标记")
	}
	
	// 显示文件开头内容
	content := string(data)
	firstLines := strings.SplitN(content, "\n", 10)
	fmt.Println("\n=== 文件前10行内容 ===")
	for i, line := range firstLines {
		if i >= 10 {
			break
		}
		fmt.Printf("第%d行: %s\n", i+1, line)
	}
	
	fmt.Println("\n=== 尝试文本格式protobuf解析 ===")
	
	// 尝试解析为文本格式protobuf
	tryTextParse(data, "v2.DiscoverResponse", &v2.DiscoverResponse{})
	tryTextParse(data, "apiV1Model.CircuitBreakerV2", &apiV1Model.CircuitBreakerV2{})
	tryTextParse(data, "apiV1Model.DiscoverResponse", &apiV1Model.DiscoverResponse{})
	
	fmt.Println("\n=== 分析文件结构 ===")
	
	// 文件可能包含多个protobuf消息
	lines := strings.Split(content, "\n")
	fmt.Printf("文件总行数: %d\n", len(lines))
	
	// 查找可能的protobuf消息分隔符
	var messageStarts []int
	for i, line := range lines {
		if strings.Contains(line, "mmexptcoretrafficproxy") {
			messageStarts = append(messageStarts, i)
		}
	}
	
	fmt.Printf("发现 %d 个可能的protobuf消息开始位置\n", len(messageStarts))
	
	if len(messageStarts) > 0 {
		fmt.Println("\n=== 提取第一个消息进行分析 ===")
		
		// 提取第一个消息
		start := messageStarts[0]
		end := len(lines)
		if len(messageStarts) > 1 {
			end = messageStarts[1]
		}
		
		firstMessage := strings.Join(lines[start:end], "\n")
		fmt.Printf("第一个消息大小: %d 字节\n", len(firstMessage))
		fmt.Printf("第一个消息前200字符:\n%s\n", firstMessage[:min(200, len(firstMessage))])
		
		// 尝试解析第一个消息
		tryTextParse([]byte(firstMessage), "第一个消息(v2.DiscoverResponse)", &v2.DiscoverResponse{})
		tryTextParse([]byte(firstMessage), "第一个消息(apiV1Model.CircuitBreakerV2)", &apiV1Model.CircuitBreakerV2{})
	}
	
	fmt.Println("\n=== 解析完成 ===")
}

func tryTextParse(data []byte, typeName string, msg proto.Message) {
	if err := prototext.Unmarshal(data, msg); err != nil {
		fmt.Printf("✗ %s 文本格式解析失败: %v\n", typeName, err)
	} else {
		fmt.Printf("✓ %s 文本格式解析成功！\n", typeName)
		fmt.Printf("  消息类型: %T\n", msg)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}