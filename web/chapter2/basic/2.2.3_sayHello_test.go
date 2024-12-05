package basic

import (
	"net/http"
	"testing"
)

// todo 创建和解析HTTP服务器端

func TestSayHello(t *testing.T) {
	http.HandleFunc("/helloo", SayHello)
	http.ListenAndServe(":8081", nil)
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
