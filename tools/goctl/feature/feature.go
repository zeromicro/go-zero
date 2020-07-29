package feature

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
)

var feature = `
1、新增对rpc错误转换处理
  1.1、目前暂时仅处理not found 和 unknown错误
2、增加feature命令支持，详细使用请通过命令[goctl -feature]查看
`

func Feature(c *cli.Context) error {
	fmt.Println(aurora.Blue("\nFEATURE:"))
	fmt.Println(aurora.Blue(feature))
	return nil
}
