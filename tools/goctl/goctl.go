package main

import (
	"fmt"
	"os"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/apigen"
	"github.com/tal-tech/go-zero/tools/goctl/api/dartgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/docgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/format"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/api/javagen"
	"github.com/tal-tech/go-zero/tools/goctl/api/jsgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/ktgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/tsgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/validate"
	"github.com/tal-tech/go-zero/tools/goctl/configgen"
	"github.com/tal-tech/go-zero/tools/goctl/docker"
	"github.com/tal-tech/go-zero/tools/goctl/feature"
	model "github.com/tal-tech/go-zero/tools/goctl/model/sql/command"
	rpc "github.com/tal-tech/go-zero/tools/goctl/rpc/command"
	"github.com/urfave/cli"
)

var (
	BuildTime = "not set"
	commands  = []cli.Command{
		{
			Name:  "api",
			Usage: "generate api related files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "o",
					Usage: "the output api file",
				},
			},
			Action: apigen.ApiCommand,
			Subcommands: []cli.Command{
				{
					Name:  "format",
					Usage: "format api files",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the format target dir",
						},
						cli.BoolFlag{
							Name:  "p",
							Usage: "print result to console",
						},
						cli.BoolFlag{
							Name:     "iu",
							Usage:    "ignore update",
							Required: false,
						},
					},
					Action: format.GoFormatApi,
				},
				{
					Name:  "validate",
					Usage: "validate api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "api",
							Usage: "validate target api file",
						},
					},
					Action: validate.GoValidateApi,
				},
				{
					Name:  "doc",
					Usage: "generate doc files",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
					},
					Action: docgen.DocCommand,
				},
				{
					Name:  "go",
					Usage: "generate go files for provided api in yaml file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: gogen.GoCommand,
				},
				{
					Name:  "java",
					Usage: "generate java files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: javagen.JavaCommand,
				},
				{
					Name:  "ts",
					Usage: "generate ts files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
						cli.StringFlag{
							Name:     "webapi",
							Usage:    "the web api file path",
							Required: false,
						},
						cli.StringFlag{
							Name:     "caller",
							Usage:    "the web api caller",
							Required: false,
						},
						cli.BoolFlag{
							Name:     "unwrap",
							Usage:    "unwrap the webapi caller for import",
							Required: false,
						},
					},
					Action: tsgen.TsCommand,
				},
				{
					Name:  "dart",
					Usage: "generate dart files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: dartgen.DartCommand,
				},
				{
					Name:  "kt",
					Usage: "generate kotlin code for provided api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target directory",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
						cli.StringFlag{
							Name:  "pkg",
							Usage: "define package name for kotlin file",
						},
					},
					Action: ktgen.KtCommand,
				},
				{
					Name:  "js",
					Usage: "generate javascript code for provided api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target directory",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: jsgen.JsCommand,
				},
			},
		},
		{
			Name:  "docker",
			Usage: "generate Dockerfile and Makefile",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "go",
					Usage: "the file that contains main function",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Usage: "which namespace of kubernetes to deploy the service",
				},
			},
			Action: docker.DockerCommand,
		},
		{
			Name:  "rpc",
			Usage: "generate rpc code",
			Subcommands: []cli.Command{
				{
					Name:  "template",
					Usage: `generate proto template`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "out, o",
							Usage: "the target path of proto",
						},
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [option]",
						},
					},
					Action: rpc.RpcTemplate,
				},
				{
					Name:  "proto",
					Usage: `generate rpc from proto`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "src, s",
							Usage: "the file path of the proto source file",
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: `the target path of the code,default path is "${pwd}". [option]`,
						},
						cli.StringFlag{
							Name:  "service, srv",
							Usage: `the name of rpc service. [option]`,
						},
						cli.StringFlag{
							Name:  "shared",
							Usage: `the dir of the shared file,default path is "${pwd}/shared. [option]`,
						},
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [option]",
						},
					},
					Action: rpc.Rpc,
				},
			},
		},
		{
			Name:  "model",
			Usage: "generate model code",
			Subcommands: []cli.Command{
				{
					Name:  "mysql",
					Usage: `generate mysql model`,
					Subcommands: []cli.Command{
						{
							Name:  "ddl",
							Usage: `generate mysql model from ddl`,
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "src, s",
									Usage: "the file path of the ddl source file",
								},
								cli.StringFlag{
									Name:  "dir, d",
									Usage: "the target dir",
								},
								cli.BoolFlag{
									Name:  "cache, c",
									Usage: "generate code with cache [optional]",
								},
								cli.BoolFlag{
									Name:  "idea",
									Usage: "for idea plugin [optional]",
								},
							},
							Action: model.MysqlDDL,
						},
						{
							Name:  "datasource",
							Usage: `generate model from datasource`,
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "url",
									Usage: `the data source of database,like "root:password@tcp(127.0.0.1:3306)/database`,
								},
								cli.StringFlag{
									Name:  "table, t",
									Usage: `source table,tables separated by commas,like "user,course`,
								},
								cli.BoolFlag{
									Name:  "cache, c",
									Usage: "generate code with cache [optional]",
								},
								cli.StringFlag{
									Name:  "dir, d",
									Usage: "the target dir",
								},
								cli.BoolFlag{
									Name:  "idea",
									Usage: "for idea plugin [optional]",
								},
							},
							Action: model.MyDataSource,
						},
					},
				},
			},
		},
		{
			Name:  "config",
			Usage: "generate config json",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Usage: "the target config go file",
				},
			},
			Action: configgen.GenConfigCommand,
		},
		{
			Name:   "feature",
			Usage:  "the features of the latest version",
			Action: feature.Feature,
		},
	}
)

func main() {
	logx.Disable()

	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = BuildTime
	app.Commands = commands
	// cli already print error messages
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
