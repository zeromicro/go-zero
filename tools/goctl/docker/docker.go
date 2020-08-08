package docker

import (
	"errors"

	"github.com/tal-tech/go-zero/tools/goctl/gen"
	"github.com/urfave/cli"
)

func DockerCommand(c *cli.Context) error {
	goFile := c.String("go")
	namespace := c.String("namespace")
	if len(goFile) == 0 || len(namespace) == 0 {
		return errors.New("-go and -namespace can't be empty")
	}

	if err := gen.GenerateDockerfile(goFile, "-f", "etc/config.json"); err != nil {
		return err
	}

	return gen.GenerateMakefile(goFile, namespace)
}
