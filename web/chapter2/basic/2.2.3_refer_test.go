package basic

import (
	"net/http"
	"testing"
)

func TestRefer(t *testing.T) {
	referer := Refer{handler: http.HandlerFunc(MyHandler), refer: "www.shirdon.com"}
	http.HandleFunc("/hello", hello)
	err := http.ListenAndServe(":8083", referer)
	if err != nil {
		t.Error(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

type Refer struct {
	handler http.Handler
	refer   string
}

func (r Refer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Referer() == r.refer {
		r.handler.ServeHTTP(writer, request)
	} else {
		writer.WriteHeader(403)
	}
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler"))
}
