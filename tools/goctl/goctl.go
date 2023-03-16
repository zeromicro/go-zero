package main

import (
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/tools/goctl/cmd"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

func main() {
	flags.Init()
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
