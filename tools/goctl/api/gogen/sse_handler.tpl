// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

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

		{{end}}// Buffer size of 16 is chosen as a reasonable default to balance throughput and memory usage.
		// You can change this based on your application's needs.
		// if your go-zero version less than 1.8.1, you need to add 3 lines below.
        // w.Header().Set("Content-Type", "text/event-stream")
        // w.Header().Set("Cache-Control", "no-cache")
        // w.Header().Set("Connection", "keep-alive")
		client := make(chan {{.ResponseType}}, 16)

        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        threading.GoSafeCtx(r.Context(), func() {
            defer close(client)
            err := l.{{.Call}}({{if .HasRequest}}&req, {{end}}client)
            if err != nil {
                logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                return
            }
        })

        for {
            select {
            case data, ok := <-client:
                if !ok {
                    return
                }
                output, err := json.Marshal(data)
                if err != nil {
                    logc.Errorw(r.Context(), "{{.HandlerName}}", logc.Field("error", err))
                    continue
                }

                if _, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err != nil {
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
