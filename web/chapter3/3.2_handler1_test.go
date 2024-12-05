package chapter3

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandler1(t *testing.T) {
	handle1 := handle1{}
	handle2 := handle2{}
	server := http.Server{Addr: "127.0.0.1:8085", Handler: nil}
	http.Handle("/handle1", &handle1)
	http.Handle("/handle2", &handle2)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

type handle1 struct{}

func (h1 *handle1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi, handle1")
}

type handle2 struct{}

func (h2 *handle2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi, handle2")
}
