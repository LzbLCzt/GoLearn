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
	"encoding/json"
	"net/http"
	"testing"
)

type Greeting struct {
	Message string `json:"message"`
}

func TestWrite3(t *testing.T) {
	http.HandleFunc("/", Hello)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	greeting := Greeting{
		"欢迎一起学习《Go Web编程实战派从入门到精通》",
	}
	bytes, err := json.Marshal(greeting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
