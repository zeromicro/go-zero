package feature

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
)

var feature = `
1、增加goctl model支持
`

func Feature(_ *cli.Context) error {
	fmt.Println(aurora.Blue("\nFEATURE:"))
	fmt.Println(aurora.Blue(feature))
	return nil
}
