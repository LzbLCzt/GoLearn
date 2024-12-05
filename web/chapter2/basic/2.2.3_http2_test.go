package basic

import (
	"log"
	"net/http"
	"testing"
)

func TestHttp2(t *testing.T) {
	server := http.Server{Addr: ":8085", Handler: http.HandlerFunc(handle)}
	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got connection: %s", r.Proto) //打印请求协议
	w.Write([]byte("Hello this is a HTTP 2 message"))
}
