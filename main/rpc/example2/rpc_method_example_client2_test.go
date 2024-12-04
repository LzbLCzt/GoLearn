package example2

import (
	"fmt"
	"net"
	"net/rpc"
	"testing"
)

// todo RPC客户端
// ----------------------------------------------------------------
type HelloServiceClient struct {
	*rpc.Client
}

func (c HelloServiceClient) Say(request string, reply *string) error {
	return c.Client.Call(ServerName+".Say", request, reply)
}

func (c HelloServiceClient) Hello(request string, reply *string) error {
	return c.Client.Call(ServerName+".Hello", request, reply)
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil) // 这行代码的主要目的是在编译时检查 HelloServiceClient 类型是否实现了 HelloServiceInterface 接口。如果 HelloServiceClient 没有实现 HelloServiceInterface 中的所有方法，编译器将报错

func TestHello2(t *testing.T) {
	client, err := DialHelloService("tcp", "9.134.117.127:1001")
	if err != nil {
		t.Fatal(err)
	}
	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("reply: %v \n", reply)

	err = client.Say("say", &reply)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("reply: %v \n", reply)
}

func DialHelloService(network string, address string) (*HelloServiceClient, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	client := rpc.NewClient(conn)
	return &HelloServiceClient{client}, nil
}

// ----------------------------------------------------------------
