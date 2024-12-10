//++++++++++++++++++++++++++++++++++++++++
// 《Go Web编程实战派从入门到精通》源码
//++++++++++++++++++++++++++++++++++++++++
// Author:廖显东（ShirDon）
// Blog:https://www.shirdon.com/
// 仓库地址：https://goWebCodeFromBook
// 仓库地址：https://github.com/shirdonl/goWebActualCombat
//++++++++++++++++++++++++++++++++++++++++

package main

import (
	"net/http"
	"testing"
)

func TestRequestHeader(t *testing.T) {
	http.HandleFunc("/redirect", Redirect)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		t.Fatal("ListenAndServe: ", err)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	//todo 通过Header 设置一个 301 重定向
	w.Header().Set("Location", "https://www.shirdon.com")
	w.WriteHeader(301)
}
