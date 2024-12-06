package chapter3

import (
	"GoLearn/web/chapter3/controller"
	"net/http"
	"testing"
)

//todo 注册GetUser服务，GetUser完成了数据库查询、模版渲染、响应客户端的工作

func TestServer4(t *testing.T) {
	http.HandleFunc("/getUser", controller.UserController{}.GetUser)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		t.Fatal(err)
	}
}

//client request: http://127.0.0.1:8088/getUser?uid=1
