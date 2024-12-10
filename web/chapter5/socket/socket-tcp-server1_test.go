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
	"fmt"
	_ "io"
	"log"
	"net"
	"testing"
)

func TestTcpServer1(t *testing.T) {
	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		t.Error(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("访问客户端信息： con=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	for {
		buf := make([]byte, 512)
		n, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Print(string(buf[:n]))
	}
}
