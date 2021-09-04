package javagen

import (
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// JavaCommand the generate java code command entrance
func JavaCommand(c *cli.Context) error {
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

	packetName := api.Service.Name
	if strings.HasSuffix(packetName, "-api") {
		packetName = packetName[:len(packetName)-4]
	}

	errorx.Must(util.MkdirIfNotExist(dir))
	errorx.Must(genPacket(dir, packetName, api))
	errorx.Must(genComponents(dir, packetName, api))

	fmt.Println(aurora.Green("Done."))
	return nil
}
