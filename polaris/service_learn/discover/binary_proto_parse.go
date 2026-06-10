package main

import (
	"fmt"
	"io/ioutil"
	"log"

	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("=== 开始解析 main.data1 二进制protobuf文件 ===")
	
	// 读取文件
	data, err := ioutil.ReadFile("main.data1")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", len(data))
	
	// 分析文件头
	fmt.Println("\n=== 文件头分析 ===")
	if len(data) >= 10 {
		fmt.Printf("前10字节(十六进制): %x\n", data[:10])
		fmt.Printf("前10字节(ASCII): %q\n", string(data[:10]))
	}
	
	// 移除UTF-8 BOM标记（如果存在）
	var cleanData []byte
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbf && data[2] == 0xbd {
		cleanData = data[3:]
		fmt.Println("✓ 已移除UTF-8 BOM标记")
	} else {
		cleanData = data
		fmt.Println("✓ 无BOM标记")
	}
	
	// 查找真正的protobuf消息开始位置
	fmt.Println("\n=== 查找protobuf消息 ===")
	
	// protobuf消息通常以0x0A开头
	var messageStart int = -1
	for i := 0; i < len(cleanData)-1; i++ {
		if cleanData[i] == 0x0A {
			messageStart = i
			break
		}
	}
	
	if messageStart == -1 {
		fmt.Println("✗ 未找到protobuf消息开始标记")
		return
	}
	
	fmt.Printf("找到protobuf消息开始位置: %d\n", messageStart)
	
	// 提取protobuf消息数据
	messageData := cleanData[messageStart:]
	fmt.Printf("protobuf消息大小: %d 字节\n", len(messageData))
	
	// 显示消息前几个字节
	if len(messageData) > 50 {
		fmt.Printf("消息前50字节(十六进制): %x\n", messageData[:50])
	}
	
	// 尝试解析protobuf消息
	fmt.Println("\n=== 尝试解析protobuf消息 ===")
	
	// 尝试v2.DiscoverResponse
	rsp := &v2.DiscoverResponse{}
	if err := proto.Unmarshal(messageData, rsp); err != nil {
		fmt.Printf("✗ v2.DiscoverResponse 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ v2.DiscoverResponse 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", rsp)
		// 使用String()方法显示内容摘要
		fmt.Printf("  内容摘要: %s\n", rsp.String())
	}

	// 尝试apiV1Model.CircuitBreakerV2
	cb := &apiV1Model.CircuitBreakerV2{}
	if err := proto.Unmarshal(messageData, cb); err != nil {
		fmt.Printf("✗ apiV1Model.CircuitBreakerV2 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ apiV1Model.CircuitBreakerV2 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", cb)
		fmt.Printf("  内容摘要: %s\n", cb.String())
	}

	// 尝试apiV1Model.DiscoverResponse
	dr := &apiV1Model.DiscoverResponse{}
	if err := proto.Unmarshal(messageData, dr); err != nil {
		fmt.Printf("✗ apiV1Model.DiscoverResponse 解析失败: %v\n", err)
	} else {
		fmt.Printf("✓ apiV1Model.DiscoverResponse 解析成功！\n")
		fmt.Printf("  消息类型: %T\n", dr)
		fmt.Printf("  内容摘要: %s\n", dr.String())
	}
	
	// 检查文件是否包含多个protobuf消息
	fmt.Println("\n=== 检查多个消息 ===")
	checkMultipleMessages(cleanData)
	
	fmt.Println("\n=== 解析完成 ===")
}

func checkMultipleMessages(data []byte) {
	// 统计protobuf消息数量
	var messageCount int
	var lastPos int
	
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x0A && i > lastPos+100 { // 至少间隔100字节
			messageCount++
			lastPos = i
			
			if messageCount >= 5 { // 只检查前5个
				break
			}
		}
	}
	
	fmt.Printf("发现 %d 个可能的protobuf消息\n", messageCount)
	
	if messageCount > 1 {
		fmt.Println("文件可能包含多个protobuf消息")
	}
}