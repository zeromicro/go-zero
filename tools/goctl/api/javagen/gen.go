package javagen

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// JavaCommand the generate java code command entrance
func JavaCommand(c *cli.Context) error {
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
	packetName := strings.TrimSuffix(api.Service.Name, "-api")
	logx.Must(util.MkdirIfNotExist(dir))
	logx.Must(genPacket(dir, packetName, api))
	logx.Must(genComponents(dir, packetName, api))

	fmt.Println(aurora.Green("Done."))
	return nil
}
