package restful

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)
import "github.com/emicklei/go-restful"

var defaultLisConf net.Listener

func init() {
	var err error
	defaultLisConf, err = net.Listen("tcp", ":8082")
	if err != nil {
		panic(fmt.Sprintf("Failed to create listener: %v", err))
	}
}

func TestExample(t *testing.T) {
	h := &Handler{}
	ws := newWebService("v1/migration", restful.MIME_JSON, restful.MIME_JSON)
	ws.Route(ws.POST("/service/getCluster").To(h.GetCluster))

	httpHandler := restful.NewContainer()
	httpHandler.Add(ws)

	server := &http.Server{Addr: ":8081", Handler: httpHandler}

	fmt.Println("start server")
	if err := server.Serve(defaultLisConf); err != nil {
		fmt.Printf("http server start error: %v\n", err)
	}

}

func newWebService(prefix, accepts, contentTypes string) *restful.WebService {
	ws := new(restful.WebService)
	// Consumes和Produces作用：限制当前service可以接收的请求类型和返回的类型
	ws.Path(prefix).
		Consumes(accepts).
		Produces(contentTypes)
	return ws
}

type Handler struct {
}

func (h *Handler) GetCluster(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	fmt.Printf("header: %v\n", header)

	if err := resp.WriteHeaderAndJson(http.StatusOK, &Response{Code: 200000, Message: "success"}, restful.MIME_JSON); err != nil {
		fmt.Printf("write header and json error: %v\n", err)
	}
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
