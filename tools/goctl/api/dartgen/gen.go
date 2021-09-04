package dartgen

import (
	"errors"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

// DartCommand create dart network request code
func DartCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	if len(apiFile) == 0 {
		errorx.Must(errors.New("missing -api"))
	}

	if len(dir) == 0 {
		errorx.Must(errors.New("missing -dir"))
	}

	api, err := parser.Parse(apiFile)
	errorx.Must(err)

	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	api.Info.Title = strings.Replace(apiFile, ".api", "", -1)
	errorx.Must(genData(dir+"data/", api))
	errorx.Must(genApi(dir+"api/", api))
	errorx.Must(genVars(dir + "vars/"))
	return nil
}
