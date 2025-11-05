// server/main.go
package main

import (
	greeter "GoLearn/polaris/request/trpc/api/Greeter"
	svc "GoLearn/polaris/request/trpc/server/Greeter"
	"git.code.oa.com/trpc-go/trpc-go"
	"log"
)

func main() {
	// 创建服务器实例
	s := trpc.NewServer()

	// 注册服务实现
	greeter.RegisterGreeterService(s, &svc.GreeterServer{})

	// 启动服务
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
