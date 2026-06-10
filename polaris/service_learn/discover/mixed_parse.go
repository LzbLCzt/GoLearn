package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("=== 开始解析 main.data1 混合格式文件 ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	
	// 分析文件结构
	fmt.Println("\n=== 文件结构分析 ===")
	
	// 查找可能的protobuf消息开始位置
	// protobuf消息通常以0x0A开头
	var protoStartPositions []int
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x0A && i > 0 {
			// 检查前面是否有合理的字段长度
			if i > 1 && isValidProtoField(data[i-1]) {
				protoStartPositions = append(protoStartPositions, i)
			}
		}
		if len(protoStartPositions) >= 10 { // 只找前10个
			break
		}
	}
	
	fmt.Printf("发现 %d 个可能的protobuf消息开始位置\n", len(protoStartPositions))
	
	// 尝试解析每个可能的protobuf消息
	for i, pos := range protoStartPositions {
		if i >= 3 { // 只解析前3个
			break
		}
		
		fmt.Printf("\n=== 尝试解析第 %d 个消息（位置 %d） ===\n", i+1, pos)
		
		// 提取消息数据（假设消息长度不超过100KB）
		end := min(pos+100000, len(data))
		messageData := data[pos:end]
		
		tryParseMessage(messageData, fmt.Sprintf("消息%d", i+1))
	}
	
	// 分析文本内容
	fmt.Println("\n=== 文本内容分析 ===")
	content := string(data)
	
	// 提取关键信息
	lines := strings.Split(content, "\n")
	fmt.Printf("文件总行数: %d\n", len(lines))
	
	// 查找服务相关信息
	for i, line := range lines {
		if i >= 20 { // 只检查前20行
			break
		}
		if strings.Contains(line, "mmexptcoretrafficproxy") {
			fmt.Printf("第%d行: 服务名称 - %s\n", i+1, strings.TrimSpace(line))
		}
		if strings.Contains(line, "Production") {
			fmt.Printf("第%d行: 环境信息 - %s\n", i+1, strings.TrimSpace(line))
		}
		if strings.Contains(line, "internal-") {
			fmt.Printf("第%d行: 配置项 - %s\n", i+1, strings.TrimSpace(line))
		}
	}
	
	fmt.Println("\n=== 解析完成 ===")
}

func tryParseMessage(data []byte, messageName string) {
	fmt.Printf("消息大小: %d 字节\n", len(data))
	
	// 尝试v2.DiscoverResponse
	rsp := &v2.DiscoverResponse{}
	if err := proto.Unmarshal(data, rsp); err != nil {
		fmt.Printf("✗ %s - v2.DiscoverResponse 解析失败: %v\n", messageName, err)
	} else {
		fmt.Printf("✓ %s - v2.DiscoverResponse 解析成功！\n", messageName)
		fmt.Printf("  消息类型: %T\n", rsp)
	}

	// 尝试apiV1Model.CircuitBreakerV2
	cb := &apiV1Model.CircuitBreakerV2{}
	if err := proto.Unmarshal(data, cb); err != nil {
		fmt.Printf("✗ %s - apiV1Model.CircuitBreakerV2 解析失败: %v\n", messageName, err)
	} else {
		fmt.Printf("✓ %s - apiV1Model.CircuitBreakerV2 解析成功！\n", messageName)
		fmt.Printf("  消息类型: %T\n", cb)
	}
}

func isValidProtoField(b byte) bool {
	// protobuf字段标签通常在1-15范围内
	return b >= 1 && b <= 15
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}