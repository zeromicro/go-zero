package mongo

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/mongo/generate"
	"github.com/urfave/cli"
)

// Action provides the entry for goctl mongo code generation.
func Action(ctx *cli.Context) error {
	tp := ctx.StringSlice("type")
	c := ctx.Bool("cache")
	o := strings.TrimSpace(ctx.String("dir"))
	s := ctx.String("style")
	if len(tp) == 0 {
		return errors.New("missing type")
	}

	cfg, err := config.NewConfig(s)
	if err != nil {
		return err
	}

	a, err := filepath.Abs(o)
	if err != nil {
		return err
	}

	return generate.Do(&generate.Context{
		Types:  tp,
		Cache:  c,
		Output: a,
		Cfg:    cfg,
	})
}
