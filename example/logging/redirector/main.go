package main

import (
	"fmt"
	"time"

	"zero/core/logx"
)

func main() {
	logx.MustSetup(logx.LogConf{
		Mode: "console",
	})
	logx.CollectSysLog()

	line := "asdkg"
	logx.Info(line)
	fmt.Print(line)
	time.Sleep(time.Second)
}
