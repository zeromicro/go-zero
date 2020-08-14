package ktgen

import (
	"errors"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

func KtCommand(c *cli.Context) error {
	apiFile := c.String("api")
	if apiFile == "" {
		return errors.New("missing -api")
	}
	dir := c.String("dir")
	if dir == "" {
		return errors.New("missing -dir")
	}
	pkg := c.String("pkg")
	if pkg == "" {
		return errors.New("missing -pkg")
	}

	p, e := parser.NewParser(apiFile)
	if e != nil {
		return e
	}
	api, e := p.Parse()
	if e != nil {
		return e
	}

	e = genBase(dir, pkg, api)
	if e != nil {
		return e
	}
	e = genApi(dir, pkg, api)
	if e != nil {
		return e
	}
	return nil
}
