//go:build linux || darwin || freebsd

package migrate

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lerity-yao/go-zero/tools/cztctl/util/console"
	"github.com/zeromicro/go-zero/core/syncx"
)

func cancelOnSignals() {
	doneChan := syncx.NewDoneChan()
	defer doneChan.Close()

	go func(dc *syncx.DoneChan) {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT)
		select {
		case <-c:
			console.Error(`
migrate failed, reason: "User Canceled"`)
			os.Exit(0)
		case <-dc.Done():
			return
		}
	}(doneChan)
}
