package cookie

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCookie1(t *testing.T) {
	http.HandleFunc("/", testCookie)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		t.Error(err)
	}
}

func testCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("test_cookie")
	fmt.Printf("cookie:%#v, err:%v\n", c, err)

	cookie := &http.Cookie{
		Name:   "test_cookie",
		Value:  "krrsklHhefUUUFSSKLAkaLlJGGQEXZLJP",
		MaxAge: 3600,
		Domain: "localhost",
		Path:   "/",
	}

	http.SetCookie(w, cookie)

	w.Write([]byte("hello"))
}
