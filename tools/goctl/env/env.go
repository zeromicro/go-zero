package env

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
)

func Action(c *cli.Context) error {
	write := c.StringSlice("write")
	if len(write) > 0 {
		return env.WriteEnv(write)
	}
	fmt.Println(env.Print())
	return nil
}
