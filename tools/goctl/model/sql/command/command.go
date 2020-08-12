package command

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/urfave/cli"
)

func Mysql(ctx *cli.Context) error {
	src := ctx.String("src")
	dir := ctx.String("dir")
	cache := ctx.Bool("cache")
	generator := gen.NewDefaultGenerator(src, dir)
	return generator.Start(cache)
}
