package mongo

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/mongo/generate"
	file "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// Action provides the entry for goctl mongo code generation.
func Action(ctx *cli.Context) error {
	tp := ctx.StringSlice("type")
	c := ctx.Bool("cache")
	o := strings.TrimSpace(ctx.String("dir"))
	s := ctx.String("style")
	home := ctx.String("home")
	if len(home) > 0 {
		file.RegisterGoctlHome(home)
	}

	if len(tp) == 0 {
		errorx.Must(errors.New("missing type"))
	}

	cfg, err := config.NewConfig(s)
	errorx.Must(err)
	a, err := filepath.Abs(o)
	errorx.Must(err)
	errorx.Must(generate.Do(&generate.Context{
		Types:  tp,
		Cache:  c,
		Output: a,
		Cfg:    cfg,
	}))
	return nil
}
