package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/handler"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/util/open"
)

const (
	DefaultHost  = "127.0.0.1"
	DefaultPort  = 2000
	networkTCP   = "tcp"
	outputFormat = `serve on %s, goctl will automatically open the default browser and access it. 
If it is not opened, please manually click  %s to access it.
`
)

func Run(port int) error {
	if port <= 0 {
		port = DefaultPort
	}

	url := fmt.Sprintf("%s:%d", DefaultHost, port)
	address := fmt.Sprintf("http://%s", url)
	ctx := svc.NewServiceContext()
	if err := handler.RegisterCustomHandlers(ctx); err != nil {
		return err
	}

	fmt.Printf(outputFormat, address, address)
	// listening port and open browser.
	threading.GoSafe(func() {
		for {
			time.Sleep(500 * time.Millisecond)
			ok, err := ping(url)
			if err != nil {
				fmt.Printf("listening port failed, error: %+v\n", err)
				break
			} else if ok {
				_ = open.Open(address)
				break
			}
		}
	})

	return http.ListenAndServe(fmt.Sprintf("%s:%d", DefaultHost, port), nil)
}

func ping(address string) (bool, error) {
	conn, err := net.Dial(networkTCP, address)
	if err != nil {
		return false, err
	} else if conn != nil {
		_ = conn.Close()
		_ = open.Open(address)
		return true, nil
	}
	return false, nil
}
