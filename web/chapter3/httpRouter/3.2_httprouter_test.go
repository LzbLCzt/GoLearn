package httpRouter

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"testing"
)

func TestHttprouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", Index) //注册get方法
	//router.POST() //注册post方法
	//router.DELETE() //注册delete方法
	log.Fatal(http.ListenAndServe(":8088", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Index"))
}
