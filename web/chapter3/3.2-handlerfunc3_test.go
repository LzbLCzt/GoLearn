// ++++++++++++++++++++++++++++++++++++++++
// 《Go Web编程实战派从入门到精通》源码
// ++++++++++++++++++++++++++++++++++++++++
// Author:廖显东（ShirDon）
// Blog:https://www.shirdon.com/
// 仓库地址：https://gitee.com/shirdonl/goWebActualCombat
// 仓库地址：https://github.com/shirdonl/goWebActualCombat
// ++++++++++++++++++++++++++++++++++++++++
package chapter3

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "欢迎来到Go Web首页！处理器为：indexHandler！")
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "欢迎来到Go Web欢迎页！处理器为：hiHandler！")
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "欢迎来到Go Web欢迎页！处理器为：webHandler！")
}

func TestHandlerfunc3(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/hi", hiHandler)
	mux.HandleFunc("/hi/web", webHandler)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	tmp := TMP{Person: Person{Name: "shirdon", Age: 18}}
	tmp.eat()
	tmp2 := TMP2{Person: Person{Name: "shirdon", Age: 18}}
	tmp2.eat()
}

type TMP2 TMP

type TMP struct {
	Person
	m map[string]string
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) eat() {
	fmt.Println("吃饭")
}
