package tsgen

import (
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// TsCommand provides the entry to generate typescript codes
func TsCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	webAPI := c.String("webapi")
	caller := c.String("caller")
	unwrapAPI := c.Bool("unwrap")
	if len(apiFile) == 0 {
		errorx.Must(errors.New("missing -api"))
	}

	if len(dir) == 0 {
		errorx.Must(errors.New("missing -dir"))
	}

	api, err := parser.Parse(apiFile)
	errorx.Must(err)
	errorx.Must(util.MkdirIfNotExist(dir))
	errorx.Must(genHandler(dir, webAPI, caller, api, unwrapAPI))
	errorx.Must(genComponents(dir, api))

	fmt.Println(aurora.Green("Done."))
	return nil
}
