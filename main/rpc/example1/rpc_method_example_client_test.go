package example1

import (
	"fmt"
	"net"
	"net/rpc"
	"testing"
)

// todo RPC客户端
// ----------------------------------------------------------------
func TestHello(t *testing.T) {
	conn, err := net.Dial("tcp", "9.134.117.127:1234")
	if err != nil {
		t.Fatal("dialing:", err)
	}
	client := rpc.NewClient(conn)
	var reply string
	err = client.Call("HelloService.Hello", "--request input--", &reply)
	if err != nil {
		t.Fatal("arith error:", err)
	}
	fmt.Println(reply)
}

// ----------------------------------------------------------------
