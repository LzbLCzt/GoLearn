package basic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// todo get请求
func TestGet(t *testing.T) {
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	closer := resp.Body
	bytes, err := ioutil.ReadAll(closer)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}
