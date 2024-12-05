package httpRouter

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"testing"
)

//todo Handler包可以处理不同的二级域名，它先根据域名获取对应的Handler路由，然后调用处理（分发机制）

func TestHttprouter2(t *testing.T) {
	userRouter := httprouter.New()
	userRouter.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("sub1"))
	})

	dataRouter := httprouter.New()
	dataRouter.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("sub2"))
	})

	//分别用于处理不同的二级域名(也叫做子域名),主域名必须是子域名的后缀
	hs := make(HostMap)
	hs["sub1.localhost:8888"] = userRouter
	hs["sub2.localhost:8888"] = dataRouter

	log.Fatal(http.ListenAndServe(":8888", hs))
}

type HostMap map[string]http.Handler

func (hs HostMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//根据域名获取对应的Handler路由，然后调用处理（分发机制）
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}
