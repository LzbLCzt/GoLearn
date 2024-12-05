package httpRouter

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestHttpRouter1(t *testing.T) {
	router := httprouter.New()
	router.GET("/default", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("default get"))
	})
	router.POST("/default", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("default post"))
	})
	//精确匹配
	router.GET("/user/name", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("user name:" + p.ByName("name")))
	})
	//匹配所有
	router.GET("/user/profile/*name", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("user name:" + p.ByName("name")))
	})
	http.ListenAndServe(":8083", router)
}
