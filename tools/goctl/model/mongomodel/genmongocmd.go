package mongomodel

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/model/mongomodel/gen"
	"github.com/urfave/cli"
)

func ModelCommond(c *cli.Context) error {
	src := c.String("src")
	cache := c.String("cache")

	if len(src) == 0 {
		return errors.New("missing -src")
	}
	var needCache bool
	if cache == "yes" {
		needCache = true
	}

	lang.Must(gen.GenMongoModel(src, needCache))

	fmt.Println(aurora.Green("Done."))
	return nil
}
