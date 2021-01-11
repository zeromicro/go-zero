package validate

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

func GoValidateApi(c *cli.Context) error {
	apiFile := c.String("api")

	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}

	_, err := parser.Parse(apiFile)
	if err == nil {
		fmt.Println(aurora.Green("api format ok"))
	}
	return err
}
