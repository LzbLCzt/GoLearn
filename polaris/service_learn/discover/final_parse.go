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
	fmt.Println("=== 开始解析 main.data1 文件 ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	
	// 分析文件内容
	content := string(data)
	
	fmt.Println("\n=== 文件内容分析 ===")
	
	// 检查文件是否包含protobuf文本格式
	if strings.Contains(content, "mmexptcoretrafficproxy") {
		fmt.Println("✓ 包含服务名称: mmexptcoretrafficproxy")
	}
	
	if strings.Contains(content, "Production") {
		fmt.Println("✓ 包含环境信息: Production")
	}
	
	if strings.Contains(content, "internal-enable-nearby") {
		fmt.Println("✓ 包含配置项: internal-enable-nearby")
	}
	
	// 尝试解析为二进制protobuf
	fmt.Println("\n=== 尝试二进制protobuf解析 ===")
	tryBinaryParse(data)
	
	// 分析文件结构
	fmt.Println("\n=== 文件结构分析 ===")
	lines := strings.Split(content, "\n")
	fmt.Printf("文件总行数: %d\n", len(lines))
	
	// 显示前几行内容
	fmt.Println("\n前5行内容:")
	for i := 0; i < min(5, len(lines)); i++ {
		fmt.Printf("%d: %s\n", i+1, lines[i])
	}
	
	// 检查文件是否包含二进制protobuf数据
	fmt.Println("\n=== 二进制数据检查 ===")
	if len(data) > 10 {
		fmt.Printf("文件前10字节(十六进制): %x\n", data[:10])
		fmt.Printf("文件前10字节(ASCII): %q\n", string(data[:10]))
	}
	
	// 检查文件是否包含protobuf消息标识
	if len(data) > 0 {
		firstByte := data[0]
		fmt.Printf("文件第一个字节: 0x%x\n", firstByte)
		if firstByte == 0x0A {
			fmt.Println("✓ 可能是protobuf消息（以0x0A开头）")
		} else {
			fmt.Println("✗ 不是标准protobuf消息格式")
		}
	}
	
	fmt.Println("\n=== 解析完成 ===")
}

func tryBinaryParse(data []byte) {
	// 尝试v2.DiscoverResponse
	rsp := &v2.DiscoverResponse{}
	if err := proto.Unmarshal(data, rsp); err != nil {
		fmt.Printf("✗ v2.DiscoverResponse 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ v2.DiscoverResponse 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", rsp)
	}

	// 尝试apiV1Model.CircuitBreakerV2
	cb := &apiV1Model.CircuitBreakerV2{}
	if err := proto.Unmarshal(data, cb); err != nil {
		fmt.Printf("✗ apiV1Model.CircuitBreakerV2 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ apiV1Model.CircuitBreakerV2 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", cb)
	}

	// 尝试apiV1Model.DiscoverResponse
	dr := &apiV1Model.DiscoverResponse{}
	if err := proto.Unmarshal(data, dr); err != nil {
		fmt.Printf("✗ apiV1Model.DiscoverResponse 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ apiV1Model.DiscoverResponse 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", dr)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}