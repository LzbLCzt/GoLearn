package practice

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func get(url string) (res string, err error) {
	resp, err1 := http.Get(url)
	if err != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4*1024)
	for true {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件读取完毕")
				break
			} else {
				fmt.Println("resp.Body.Read err = ", err)
			}
		}
		res += string(buf[:n])
	}
	return
}

func Test_web(t *testing.T) {
	res, err := get("https://juejin.cn/post/7170962566842155022")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
