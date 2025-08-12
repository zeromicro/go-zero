package {{.PkgName}}

import (
    "encoding/json"
    "fmt"
	"net/http"

    "github.com/zeromicro/go-zero/core/logc"
    "github.com/zeromicro/go-zero/core/threading"
	{{if .HasRequest}}"github.com/zeromicro/go-zero/rest/httpx"{{end}}
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		{{end}}client := make(chan {{.ResponseType}}, 16)
        defer func() {
            close(client)
        }()
        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        threading.GoSafeCtx(r.Context(), func() {
            err := l.{{.Call}}({{if .HasRequest}}&req, {{end}}client)
            if err != nil {
                logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                return
            }
        })

        for {
            select {
            case data := <-client:
                output, err := json.Marshal(data)
                if err != nil {
                    logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                    continue
                }

                if, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err!=nil{
                    logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                    return
                }
               if flusher, ok := w.(http.Flusher); ok {
                   flusher.Flush()
               }
            case <-r.Context().Done():
                return
            }
        }
	}
}
