package dartgen

import (
	"errors"
	"path"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// DartCommand create dart network request code
func DartCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	service := api.Service
	for _, g := range service.Groups {
		prefix := util.TrimSpace(g.GetAnnotation(spec.RoutePrefixKey))
		for _, r := range g.Routes {
			r.Path = path.Join(prefix, r.Path)
		}
	}
	api.Service = service
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	api.Info.Title = strings.Replace(apiFile, ".api", "", -1)
	logx.Must(genData(dir+"data/", api))
	logx.Must(genApi(dir+"api/", api))
	logx.Must(genVars(dir + "vars/"))
	return nil
}
