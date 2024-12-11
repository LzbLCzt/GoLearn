package example2

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

// server_test.go
func TestServer(t *testing.T) {
	_ = rpc.RegisterName("HelloService", new(HelloService))
	listener, err := net.Listen("tcp", ":1004")
	if err != nil {
		t.Fatal("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		//go rpc.ServeConn(conn)
	}
}

type serverResp struct {
	Result any    `json:"result"` // 服务端返回的结果
	Error  any    `json:"error"`  // 服务端返回的错误
	Id     uint64 `json:"id"`     // 响应的 ID，与请求的 ID 匹配
}

// client_test.go
func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":1005")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	// 构造请求参数
	cr := clientReq{Method: "HelloService.Hello", Params: [1]any{"req content: hello"}, Id: 0}
	reqData, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}
	// 发送请求到服务端
	_, err = conn.Write(reqData)
	if err != nil {
		t.Fatal(err)
	}

	// 接收服务端响应
	respData := make([]byte, 4096)
	n, err := conn.Read(respData)
	if err != nil {
		t.Fatal(err)
	}

	// 解码服务端返回的 JSON 响应
	var resp serverResp
	err = json.Unmarshal(respData[:n], &resp)
	if err != nil {
		t.Fatal(err)
	}

	// 打印服务端返回的结果
	if resp.Error != nil {
		log.Fatalf("Server returned an error: %v", resp.Error)
	} else {
		log.Printf("Server response: %v", resp.Result)
	}
}

type clientReq struct {
	Method string `json:"method"`
	Params [1]any `json:"params"`
	Id     uint64 `json:"id"`
}
