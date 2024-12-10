//++++++++++++++++++++++++++++++++++++++++
// 《Go Web编程实战派从入门到精通》源码
//++++++++++++++++++++++++++++++++++++++++
// Author:廖显东（ShirDon）
// Blog:https://www.shirdon.com/
// 仓库地址：https://goWebCodeFromBook
// 仓库地址：https://github.com/shirdonl/goWebActualCombat
//++++++++++++++++++++++++++++++++++++++++

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() { //todo 输入来自命令行终端，通过go test没法输入
	conn, err := net.Dial("tcp", ":8088")
	if err != nil {
		fmt.Println("连接服务端失败", err)
	}
	reader := bufio.NewReader(os.Stdin) //os.Stdin 代表标准输入[终端]
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("没有更多的输入或用户已断开连接")
				break // 退出循环
			}
			fmt.Println("读取用户输入失败", err)
			continue
		}
		line = strings.Trim(line, "\r\n")

		if line == "exit" {
			fmt.Println("用户推出客户端")
			break
		}

		n, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			panic(err)
		}
		fmt.Printf("客户端发送了 %d 字节的数据到服务端\n", n)
	}
}
