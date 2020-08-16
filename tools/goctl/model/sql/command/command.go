package command

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/urfave/cli"
)

func Mysql(ctx *cli.Context) error {
	src := ctx.String("src")
	dir := ctx.String("dir")
	cache := ctx.Bool("cache")
	idea := ctx.Bool("idea")
	var log console.Console
	if idea {
		log = console.NewIdeaConsole()
	} else {
		log = console.NewColorConsole()
	}
	generator := gen.NewDefaultGenerator(src, dir, gen.WithConsoleOption(log))
	err := generator.Start(cache)
	if err != nil {
		log.Error("%v", err)
	}
	return nil
}
