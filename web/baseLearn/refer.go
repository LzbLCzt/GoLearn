package main

import "net/http"

func refer() {
	//自定义Handle，实现拦截器的功能：在发送http请求时，只有带上指定的refer参数，该请求才能调用成功，否则返回403
	referer := &Refer{
		handler: http.HandlerFunc(myHandler),
		refer:   "www.shirdon.com",
	}
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", referer)
}

type Refer struct {
	handler http.Handler
	refer   string
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is handler"))
}

func (this *Refer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Referer() == this.refer {
		this.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}
