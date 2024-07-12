package server

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/config"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/handler"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/util/open"
	"net"
	"time"
)

func Run(port int) error {
	if port <= 0 {
		port = 8080
	}
	ctx := svc.NewServiceContext(&config.Config{
		Port: port,
	})

	server := rest.MustNewServer(rest.RestConf{
		Port: port,
	}, rest.WithFileServer("/", "build"))
	handler.RegisterHandlers(server, ctx)
	defer server.Stop()

	fmt.Printf("serve on http://127.0.0.1:%d, goctl will automatically open the default browser and access it. \nIf it is not opened, please manually click  http://127.0.0.1:%d to access it.\n", port, port)
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
	server.Start()
	return nil
}
