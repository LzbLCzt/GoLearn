package rpcWatch

import (
	"fmt"
	"net/http"
	"net/rpc"
	"testing"
)

func TestWatchServer(t *testing.T) {
	service := NewKVStoreService()
	err := rpc.Register(service)
	if err != nil {
		fmt.Println("err=====", err.Error())
		return
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("err=====", err.Error())
	}
}
