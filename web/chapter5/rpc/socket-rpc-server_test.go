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
	"net/http"
	"net/rpc"
	"testing"
)

type Algorithm int

type Args struct {
	X, Y int
}

func (a *Algorithm) Sum(args *Args, reply *int) error {
	*reply = args.X + args.Y
	fmt.Println("sum is ", *reply)
	return nil
}

func TestSocketRPCServer(t *testing.T) {
	alg := new(Algorithm)
	//注册服务
	err := rpc.Register(alg)
	if err != nil {
		fmt.Println(err)
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println(err)
	}
}
