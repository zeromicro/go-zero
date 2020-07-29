package main

import (
	"fmt"
	"os"

	"zero/core/logx"
	"zero/tools/modelctl/model/configtemplategen"
	"zero/tools/modelctl/model/modelgen"

	"github.com/urfave/cli"
)

var commands = []cli.Command{
	{
		Name:  "model",
		Usage: "generate model files",
		Subcommands: []cli.Command{
			{
				Name:   "template",
				Usage:  "generate the json config template",
				Action: configgen.ConfigCommand,
			},
			{
				Name:   "file",
				Usage:  "generated from a configuration file",
				Action: modelgen.FileModelCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "config,c",
						Usage: "the file path of config",
					},
				},
			},
			{
				Name:   "cmd",
				Usage:  "generated from the command line,it will be generated from ALL TABLES",
				Action: modelgen.CmdModelCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "address,a",
						Usage: "database connection address,format:\"[username]:[password]@[address]\"",
					},
					cli.StringFlag{
						Name:  "schema,s",
						Usage: "the target database name",
					},
					cli.BoolFlag{
						Name:  "force,f",
						Usage: "whether to force the generation, if it is, it may cause the source file to be lost,[default:false]",
					},
					cli.BoolFlag{
						Name:  "redis,r",
						Usage: "whether to generate with redis cache when generating files,[default:false]",
					},
				},
			},
		},
	},
}

func main() {
	logx.Disable()
	app := cli.NewApp()
	app.Usage = "a cli tool to generate model"
	app.Version = "0.0.1"
	app.Commands = commands
	// cli already print error messages
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
