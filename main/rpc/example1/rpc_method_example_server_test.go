package example1

import (
	"log"
	"net"
	"net/rpc"
	"testing"
)

// todo RPC服务端
// ----------------------------------------------------------------
type HelloService struct{}

// 写一个rpc方法
// 接口方法必须是func (t *T) MethodName(argType T1, replyType *T2) error
func (p *HelloService) Hello(request string, reply *string) error {
	log.Println("HelloService Hello")
	*reply = "Hello：" + request
	return nil
}

// 注册rpc方法
func TestRegisterRPCMethod(t *testing.T) {
	_ = rpc.RegisterName("HelloService", new(HelloService))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(conn)
	}
}

// ----------------------------------------------------------------
