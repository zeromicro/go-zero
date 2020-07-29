package tsgen

import (
	"errors"
	"fmt"

	"zero/core/lang"
	"zero/tools/goctl/api/parser"
	"zero/tools/goctl/util"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
)

func TsCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	webApi := c.String("webapi")
	caller := c.String("caller")
	unwrapApi := c.Bool("unwrap")
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	p, err := parser.NewParser(apiFile)
	if err != nil {
		return err
	}
	api, err := p.Parse()
	if err != nil {
		fmt.Println(aurora.Red("Failed"))
		return err
	}

	lang.Must(util.MkdirIfNotExist(dir))
	lang.Must(genHandler(dir, webApi, caller, api, unwrapApi))
	lang.Must(genComponents(dir, api))

	fmt.Println(aurora.Green("Done."))
	return nil
}
