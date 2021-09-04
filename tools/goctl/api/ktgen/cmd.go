package ktgen

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

// KtCommand the generate kotlin code command entrance
func KtCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	if len(apiFile) == 0 {
		errorx.Must(errors.New("missing -api"))
	}

	if len(dir) == 0 {
		errorx.Must(errors.New("missing -dir"))
	}

	pkg := c.String("pkg")
	if pkg == "" {
		errorx.Must(errors.New("missing -pkg"))
	}

	api, e := parser.Parse(apiFile)
	errorx.Must(e)
	errorx.Must(genBase(dir, pkg, api))
	errorx.Must(genApi(dir, pkg, api))

	return nil
}
