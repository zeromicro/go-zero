package mongomodel

import (
	"errors"
	"fmt"

	"zero/core/lang"
	"zero/tools/goctl/model/mongomodel/gen"

	"github.com/logrusorgru/aurora"
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
