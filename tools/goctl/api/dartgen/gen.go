package dartgen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

// DartCommand create dart network request code
func DartCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	isLegacy := c.Bool("legacy")
	hostname := c.String("hostname")
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}
	if len(hostname) == 0 {
		fmt.Println("you could use '-hostname' flag to specify your server hostname")
		hostname = "go-zero.dev"
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	api.Service = api.Service.JoinPrefix()
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	api.Info.Title = strings.Replace(apiFile, ".api", "", -1)
	logx.Must(genData(dir+"data/", api, isLegacy))
	logx.Must(genApi(dir+"api/", api, isLegacy))
	logx.Must(genVars(dir+"vars/", isLegacy, hostname))
	if err := formatDir(dir); err != nil {
		logx.Errorf("failed to format, %v", err)
	}
	return nil
}
