package template

import (
	"html/template"
	"net/http"
	"testing"
)

type PageData struct {
	Title string
	Name  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "My Page",
		Name:  "Alice",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Test_1(t *testing.T) {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
