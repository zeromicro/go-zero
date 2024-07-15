package server

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/handler"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/util/open"
)

//go:embed static/*
var assets embed.FS

func Run(port int) error {
	if port <= 0 {
		port = 8080
	}

	ctx := svc.NewServiceContext(assets)
	handler.RegisterCustomHandlers(ctx)

	fmt.Printf(`serve on http://127.0.0.1:%d, goctl will automatically open the default browser and access it. 
If it is not opened, please manually click  http://127.0.0.1:%d to access it.
`, port, port)
	threading.GoSafe(func() {
		for {
			time.Sleep(2 * time.Second)
			conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
			if err != nil {
				fmt.Printf("listening port failed, error: %+v\n", err)
				break
			} else if conn != nil {
				_ = conn.Close()
				_ = open.Open(fmt.Sprintf("http://127.0.0.1:%d", port))
				break
			}
		}
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
