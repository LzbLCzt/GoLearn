package example2

import (
	"log"
	"net"
	"net/rpc"
	"testing"
)

// todo server
const ServerName = "HelloService"

type HelloServiceInterface interface {
	Hello(request string, reply *string) error
	Say(request string, reply *string) error
}

func RegisterHelloService(srv HelloServiceInterface) error {
	return rpc.RegisterName(ServerName, srv)
}

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	log.Println("HelloService Hello")
	*reply = "Hello：" + request
	return nil
}

func (p *HelloService) Say(request string, reply *string) error {
	log.Println("Say Hello")
	*reply = "Hello：" + request
	return nil
}

func TestRegisterRPCMethod2(t *testing.T) {
	_ = RegisterHelloService(&HelloService{})
	listener, err := net.Listen("tcp", ":1001")
	if err != nil {
		log.Fatal("net.Listen err:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("listener.Accept err:", err)
		}
		go rpc.ServeConn(conn)
	}
}
