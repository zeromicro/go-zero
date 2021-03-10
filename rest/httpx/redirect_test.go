package httpx

import (
	"net/http"
	"testing"
)

func TestRedirect(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com/", nil)

	w := tracedResponseWriter{headers: make(map[string][]string)}

	Redirect(http.StatusMovedPermanently, "https://www.baidu.com/", r, &w)

}
