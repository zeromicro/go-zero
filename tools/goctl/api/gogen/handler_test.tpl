package {{.PkgName}}

import (
	"bytes"
    {{if .HasRequest}}"encoding/json"{{end}}
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/stretchr/testify/assert"
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func Test{{.HandlerName}}(t *testing.T)  {
	// 创建一个ServiceContext实例
    	c := config.Config{}
    	svcCtx := svc.NewServiceContext(c)

	// 创建一个HTTP请求
	reqBody := []byte{}
	{{if .HasRequest}}
	reqObj:= types.{{.RequestType}}{
	    //TODO: add fields here
	}
	reqBody, _ = json.Marshal(reqObj)
	{{end}}
	req, err := http.NewRequest("POST", "/unittest",  bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 创建一个HTTP响应记录器
	rr := httptest.NewRecorder()

	// 创建一个LoginHandler实例
	handler := {{.HandlerName}}(svcCtx)

	// 调用LoginHandler
	handler.ServeHTTP(rr, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, rr.Code)

	// 检查响应体
	t.Log(rr.Body.String())
}
