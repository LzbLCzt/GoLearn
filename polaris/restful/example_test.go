package restful

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// 请求ID中间件
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取Request-ID
		requestID := r.Header.Get("X-Request-ID")

		// 如果客户端未提供或长度不足16，则生成新的UUID
		if len(requestID) < 16 {
			if requestID == "" {
				requestID = uuid.New().String()
			} else {
				requestID = fmt.Sprintf("%s_auto_%s", requestID, uuid.New().String())
			}
		}

		// 将Request-ID设置到响应头
		w.Header().Set("X-Request-ID", requestID)

		// 将Request-ID存入上下文，供后续处理使用
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 用户资源处理器
type UserHandler struct{}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	vars := mux.Vars(r)
	userID := vars["id"]

	response := map[string]interface{}{
		"user_id":    userID,
		"name":       "张三",
		"request_id": requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("requestID").(string)

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message":    "用户创建成功",
		"user":       user,
		"request_id": requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	userHandler := &UserHandler{}

	// 应用请求ID中间件
	router.Use(RequestIDMiddleware)

	// 定义RESTful路由
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	log.Println("RESTful服务启动在 :8080 端口")
	log.Fatal(http.ListenAndServe(":8080", router))
}
