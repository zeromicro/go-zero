package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/apigen"
	"github.com/tal-tech/go-zero/tools/goctl/api/dartgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/docgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/format"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/api/javagen"
	"github.com/tal-tech/go-zero/tools/goctl/api/ktgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/new"
	"github.com/tal-tech/go-zero/tools/goctl/api/tsgen"
	"github.com/tal-tech/go-zero/tools/goctl/api/validate"
	"github.com/tal-tech/go-zero/tools/goctl/configgen"
	"github.com/tal-tech/go-zero/tools/goctl/docker"
	model "github.com/tal-tech/go-zero/tools/goctl/model/sql/command"
	rpc "github.com/tal-tech/go-zero/tools/goctl/rpc/cli"
	"github.com/tal-tech/go-zero/tools/goctl/tpl"
	"github.com/urfave/cli"
)

var (
	BuildVersion = "20201108"
	commands     = []cli.Command{
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
					Name:   "new",
					Usage:  "fast create api service",
					Action: new.NewService,
				},
				{
					Name:  "format",
					Usage: "format api files",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the format target dir",
						},
						cli.BoolFlag{
							Name:     "iu",
							Usage:    "ignore update",
							Required: false,
						},
						cli.BoolFlag{
							Name:     "stdin",
							Usage:    "use stdin to input api doc content, press \"ctrl + d\" to send EOF",
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
			},
		},
		{
			Name:  "docker",
			Usage: "generate Dockerfile",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "go",
					Usage: "the file that contains main function",
				},
			},
			Action: docker.DockerCommand,
		},
		{
			Name:  "rpc",
			Usage: "generate rpc code",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: `generate rpc demo service`,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [optional]",
						},
					},
					Action: rpc.RpcNew,
				},
				{
					Name:  "template",
					Usage: `generate proto template`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "out, o",
							Usage: "the target path of proto",
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
						cli.StringSliceFlag{
							Name:  "proto_path, I",
							Usage: `native command of protoc, specify the directory in which to search for imports. [optional]`,
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: `the target path of the code`,
						},
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [optional]",
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
									Usage: "the path or path globbing patterns of the ddl",
								},
								cli.StringFlag{
									Name:  "dir, d",
									Usage: "the target dir",
								},
								cli.StringFlag{
									Name:  "style",
									Usage: "the file naming style, lower|camel|underline,default is lower",
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
									Usage: `the table or table globbing patterns in the database`,
								},
								cli.BoolFlag{
									Name:  "cache, c",
									Usage: "generate code with cache [optional]",
								},
								cli.StringFlag{
									Name:  "dir, d",
									Usage: "the target dir",
								},
								cli.StringFlag{
									Name:  "style",
									Usage: "the file naming style, lower|camel|snake, default is lower",
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
			Name:  "template",
			Usage: "template operation",
			Subcommands: []cli.Command{
				{
					Name:   "init",
					Usage:  "initialize the all templates(force update)",
					Action: tpl.GenTemplates,
				},
				{
					Name:   "clean",
					Usage:  "clean the all cache templates",
					Action: tpl.CleanTemplates,
				},
				{
					Name:  "update",
					Usage: "update template of the target category to the latest",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "category,c",
							Usage: "the category of template, enum [api,rpc,model]",
						},
					},
					Action: tpl.UpdateTemplates,
				},
				{
					Name:  "revert",
					Usage: "revert the target template to the latest",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "category,c",
							Usage: "the category of template, enum [api,rpc,model]",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "the target file name of template",
						},
					},
					Action: tpl.RevertTemplates,
				},
			},
		},
	}
)

func main() {
	logx.Disable()

	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", BuildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	// cli already print error messages
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
