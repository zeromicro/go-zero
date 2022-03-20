package env

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
)

func Action(c *cli.Context) error {
	write := c.StringSlice("write")
	if len(write) > 0 {
		return env.WriteEnv(write)
	}
	args := c.Args()
	var values []string
	for _, arg := range args {
		if !env.Exists(arg) {
			continue
		}
		values = append(values, env.GetOr(arg, ""))
	}
	if len(values) > 0 {
		fmt.Println(strings.Join(values, "\n"))
		return nil
	}
	fmt.Println(env.Print())
	return nil
}
