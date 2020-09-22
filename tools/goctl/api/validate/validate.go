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

	p, err := parser.NewParser(apiFile)
	if err != nil {
		return err
	}
	_, err = p.Parse()
	if err == nil {
		fmt.Println(aurora.Green("api format ok"))
	}
	return err
}
