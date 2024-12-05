package chapter3

import (
	"fmt"
	"net/http"
	"testing"
)

// todo 自定义多路复用器
func TestHandler2(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/cn", WelcomeHandler{Language: "你好，世界"})
	mux.Handle("/en", WelcomeHandler{Language: "Hello, World"})

	server := http.Server{Addr: "127.0.0.1:8088", Handler: mux}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type WelcomeHandler struct {
	Language string
}

// todo 处理器处理请求的方法
func (h WelcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", h.Language)
}
