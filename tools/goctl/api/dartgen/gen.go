package dartgen

import (
	"errors"
	"strings"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

func DartCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
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
		return err
	}

	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	api.Info.Title = strings.Replace(apiFile, ".api", "", -1)
	lang.Must(genData(dir+"data/", api))
	lang.Must(genApi(dir+"api/", api))
	lang.Must(genVars(dir + "vars/"))
	return nil
}
