package restful

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"
)

func TestExample2(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/api/users/{id}", getUser).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8083", r))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("get user id %s", id)
	// 查询数据库，获取用户信息
	user := "xxx"
	json.NewEncoder(w).Encode(user)
}
