package tsgen

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// TsCommand provides the entry to generate typescript codes
func TsCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	webAPI := c.String("webapi")
	caller := c.String("caller")
	unwrapAPI := c.Bool("unwrap")
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}

	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		fmt.Println(aurora.Red("Failed"))
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	api.Service = api.Service.JoinPrefix()
	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(genHandler(dir, webAPI, caller, api, unwrapAPI))
	logx.Must(genComponents(dir, api))

	fmt.Println(aurora.Green("Done."))
	return nil
}
