package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestExampleClient(t *testing.T) {
	// 生成有效的请求ID（长度≥16）
	requestID := uuid.New().String() // 通常为36个字符

	// 示例1: GET请求 - 获取用户
	getUser(requestID)

	// 示例2: POST请求 - 创建用户
	createUser(requestID)
}

func getUser(requestID string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8080/users/123", nil)
	if err != nil {
		panic(err)
	}

	// 设置有效的请求ID
	req.Header.Set("X-Request-ID", requestID)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("GET响应: %s\n", body)
	fmt.Printf("服务端返回的Request-ID: %s\n", resp.Header.Get("X-Request-ID"))
}

func createUser(requestID string) {
	client := &http.Client{}

	userData := map[string]string{
		"name":  "李四",
		"email": "lisi@example.com",
	}

	jsonData, _ := json.Marshal(userData)

	req, err := http.NewRequest("POST", "http://localhost:8080/users",
		bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID) // 使用有效的请求ID

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("POST响应: %s\n", body)
}
