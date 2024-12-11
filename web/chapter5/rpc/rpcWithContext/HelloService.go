package rpcWithContext

import (
	"fmt"
	"log"
	"net"
)

type HelloService struct {
	conn net.Conn
	isLogin bool //todo 基于上下文信息，我们可以方便地为RPC服务增加简单的登陆状态的验证
}

func (p *HelloService) Hello(request string, reply *string) error {
	if !p.isLogin {
		return fmt.Errorf("please login")
	}
	*reply = "hello:" + request + ", from" + p.conn.RemoteAddr().String()
	return nil
}

func (p *HelloService) Login(request string, reply *string) error {
	if request != "user:password" {
		return fmt.Errorf("auth failed")
	}
	log.Println("login ok")
	p.isLogin = true
	return nil
}