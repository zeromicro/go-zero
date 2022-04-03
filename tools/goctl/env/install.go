package env

import "github.com/urfave/cli"

func Install(c *cli.Context) error {
	force := c.Bool("force")
	verbose := c.Bool("verbose")
	return Prepare(true, force, verbose)
}
