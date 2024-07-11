package server

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/config"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/handler"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/util/open"
	"net"
	"net/http"
	"time"
)

func Run(port int) error {
	if port <= 0 {
		port = 8080
	}
	c := &config.Config{Port: port}
	ctx := svc.NewServiceContext(c)
	handler.RegisterCustomHandlers(ctx)

	fmt.Printf("serve on http://127.0.0.1:%d, goctl will automatically open the default browser and access it. \nIf it is not opened, please manually click  http://127.0.0.1:%d to access it.\n", c.Port, c.Port)

	threading.GoSafe(func() {
		for {
			time.Sleep(2 * time.Second)
			conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", c.Port))
			if err != nil {
				fmt.Printf("listening port failed, error: %+v\n", err)
				break
			} else if conn != nil {
				_ = conn.Close()
				_ = open.Open(fmt.Sprintf("http://127.0.0.1:%d", c.Port))
				break
			}
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
}
