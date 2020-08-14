package modelgen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/urfave/cli"
)

func FileModelCommand(c *cli.Context) error {
	configFile := c.String("config")
	if len(configFile) == 0 {
		return errors.New("missing config value")
	}
	logx.Must(genModelWithConfigFile(configFile))
	return nil
}

func CmdModelCommand(c *cli.Context) error {
	address := c.String("address")
	force := c.Bool("force")
	schema := c.String("schema")
	redis := c.Bool("redis")
	if len(address) == 0 {
		return errors.New("missing [-address|-a]")
	}
	if len(schema) == 0 {
		return errors.New("missing [--schema|-s]")
	}
	addressArr := strings.Split(address, "@")
	if len(addressArr) < 2 {
		return errors.New("the address format is incorrect")
	}
	user := addressArr[0]
	host := addressArr[1]
	address = fmt.Sprintf("%v@tcp(%v)/information_schema", user, host)
	logx.Must(genModelWithDataSource(address, schema, force, redis, nil))
	return nil
}
