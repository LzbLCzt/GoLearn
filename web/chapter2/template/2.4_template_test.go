package template

import (
	"fmt"
	"html/template"
	"net/http"
	"testing"
)

func TestTemplate(t *testing.T) {
	http.HandleFunc("/", helloHandleFunc)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template_example.tmpl")
	if err != nil {
		fmt.Println(err)
	}
	name := "我爱go语言"
	t.Execute(w, name)
}
