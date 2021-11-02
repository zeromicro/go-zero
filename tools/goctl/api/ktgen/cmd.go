package ktgen

import (
	"errors"
	"path"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// KtCommand the generate kotlin code command entrance
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

	api, e := parser.Parse(apiFile)
	if e != nil {
		return e
	}

	service := api.Service
	for _, g := range service.Groups {
		prefix := util.TrimSpace(g.GetAnnotation(spec.RoutePrefixKey))
		for _, r := range g.Routes {
			r.Path = path.Join(prefix, r.Path)
		}
	}
	api.Service = service
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
