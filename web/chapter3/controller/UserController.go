package controller

import (
	"GoLearn/web/chapter3/model"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type UserController struct{}

func (c UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid, _ := strconv.Atoi(query["uid"][0])

	user, err := model.GetUser(uid)
	fmt.Printf("user: %v\n", user)
	t, err := template.ParseFiles("view/t3.html")
	if err != nil {
		fmt.Printf("parse template failed, err:%v\n", err)
	}
	userInfo := []string{user.Name, user.Phone}
	err = t.Execute(w, userInfo)
	if err != nil {
		fmt.Printf("execute template failed, err:%v\n", err)
	}
}
