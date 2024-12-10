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
	"log"
	"net/rpc"
)

// 参数结构体
type ArgsTwo struct {
	X, Y int
}

func main() {
	fmt.Println("客户端启动")
	client, err := rpc.DialHTTP("tcp", ":8088")
	if err != nil {
		log.Fatal("dialing: %v", err)
		return
	}
	//获取第一个输入值
	//i1, _ := strconv.Atoi(os.Args[1])
	////获取第二个输入值
	//i2, _ := strconv.Atoi(os.Args[2])
	//args := ArgsTwo{i1, i2}
	args := ArgsTwo{1, 2}
	var reply int
	//调用命名函数，等待它完成，并返回其错误状态。
	err = client.Call("Algorithm.Sum", args, &reply)
	if err != nil {
		log.Fatal("arith error: %v", err)
		return
	}
	fmt.Printf("Algorithm 和为: %d+%d=%d\n", args.X, args.Y, reply)
}
