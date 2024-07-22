package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func get() {
	//创建get请求
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Println("err", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
