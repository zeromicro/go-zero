package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/apigen"
	"github.com/zeromicro/go-zero/tools/goctl/api/dartgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/docgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/javagen"
	"github.com/zeromicro/go-zero/tools/goctl/api/ktgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/new"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/validate"
	"github.com/zeromicro/go-zero/tools/goctl/bug"
	"github.com/zeromicro/go-zero/tools/goctl/completion"
	"github.com/zeromicro/go-zero/tools/goctl/docker"
	"github.com/zeromicro/go-zero/tools/goctl/env"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/kube"
	"github.com/zeromicro/go-zero/tools/goctl/migrate"
	"github.com/zeromicro/go-zero/tools/goctl/model/mongo"
	model "github.com/zeromicro/go-zero/tools/goctl/model/sql/command"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	rpc "github.com/zeromicro/go-zero/tools/goctl/rpc/cli"
	"github.com/zeromicro/go-zero/tools/goctl/tpl"
	"github.com/zeromicro/go-zero/tools/goctl/upgrade"
)

const codeFailure = 1

var commands = []cli.Command{
	{
		Name:   "bug",
		Usage:  "report a bug",
		Action: bug.Action,
	},
	{
		Name:   "upgrade",
		Usage:  "upgrade goctl to latest version",
		Action: upgrade.Upgrade,
	},
	{
		Name:  "env",
		Usage: "check or edit goctl environment",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "write, w",
				Usage: "edit goctl environment",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:   "install",
				Usage:  "goctl env installation",
				Action: env.Install,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "force,f",
						Usage: "silent installation of non-existent dependencies",
					},
					cli.BoolFlag{
						Name:  "verbose, v",
						Usage: "enable log output",
					},
				},
			},
			{
				Name:  "check",
				Usage: "detect goctl env and dependency tools",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "install, i",
						Usage: "install dependencies if not found",
					},
					cli.BoolFlag{
						Name:  "force, f",
						Usage: "silent installation of non-existent dependencies",
					},
					cli.BoolFlag{
						Name:  "verbose, v",
						Usage: "enable log output",
					},
				},
				Action: env.Check,
			},
		},
		Action: env.Action,
	},
	{
		Name:        "migrate",
		Usage:       "migrate from tal-tech to zeromicro",
		Description: "migrate is a transition command to help users migrate their projects from tal-tech to zeromicro version",
		Action:      migrate.Migrate,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "verbose enables extra logging",
			},
			cli.StringFlag{
				Name:  "version",
				Usage: "the target release version of github.com/zeromicro/go-zero to migrate",
			},
		},
	},
	{
		Name:  "api",
		Usage: "generate api related files",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Usage: "output a sample api file",
			},
			cli.StringFlag{
				Name: "home",
				Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority",
			},
			cli.StringFlag{
				Name: "remote",
				Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
					"https://github.com/zeromicro/go-zero-template directory structure",
			},
			cli.StringFlag{
				Name:  "branch",
				Usage: "the branch of the remote repo, it does work with --remote",
			},
		},
		Action: apigen.ApiCommand,
		Subcommands: []cli.Command{
			{
				Name:      "new",
				Usage:     "fast create api service",
				UsageText: "example: goctl api new [options] service-name",
				Action:    new.CreateServiceCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
					cli.StringFlag{
						Name:  "style",
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]",
					},
				},
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
						Name:  "iu",
						Usage: "ignore update",
					},
					cli.BoolFlag{
						Name:  "stdin",
						Usage: "use stdin to input api doc content, press \"ctrl + d\" to send EOF",
					},
					cli.BoolFlag{
						Name:  "declare",
						Usage: "use to skip check api types already declare",
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
					cli.StringFlag{
						Name:     "o",
						Required: false,
						Usage:    "the output markdown directory",
					},
				},
				Action: docgen.DocCommand,
			},
			{
				Name:  "go",
				Usage: "generate go files for provided api in api file",
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
						Name:  "style",
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
					},
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
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
						Name:  "webapi",
						Usage: "the web api file path",
					},
					cli.StringFlag{
						Name:  "caller",
						Usage: "the web api caller",
					},
					cli.BoolFlag{
						Name:  "unwrap",
						Usage: "unwrap the webapi caller for import",
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
					cli.BoolFlag{
						Name:  "legacy",
						Usage: "legacy generator for flutter v1",
					},
					cli.StringFlag{
						Name:  "hostname",
						Usage: "hostname of the server",
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
				Name:  "plugin",
				Usage: "custom file generator",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "plugin, p",
						Usage: "the plugin file",
					},
					cli.StringFlag{
						Name:  "dir",
						Usage: "the target directory",
					},
					cli.StringFlag{
						Name:  "api",
						Usage: "the api file",
					},
					cli.StringFlag{
						Name:  "style",
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
					},
				},
				Action: plugin.PluginCommand,
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
			cli.StringFlag{
				Name:  "base",
				Usage: "the base image to build the docker image, default scratch",
				Value: "scratch",
			},
			cli.IntFlag{
				Name:  "port",
				Usage: "the port to expose, default none",
				Value: 0,
			},
			cli.StringFlag{
				Name: "home",
				Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority",
			},
			cli.StringFlag{
				Name: "remote",
				Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
					"https://github.com/zeromicro/go-zero-template directory structure",
			},
			cli.StringFlag{
				Name:  "branch",
				Usage: "the branch of the remote repo, it does work with --remote",
			},
			cli.StringFlag{
				Name:  "version",
				Usage: "the goctl builder golang image version",
			},
			cli.StringFlag{
				Name:  "tz",
				Usage: "the timezone of the container",
				Value: "Asia/Shanghai",
			},
		},
		Action: docker.DockerCommand,
	},
	{
		Name:  "kube",
		Usage: "generate kubernetes files",
		Subcommands: []cli.Command{
			{
				Name:  "deploy",
				Usage: "generate deployment yaml file",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "name",
						Usage:    "the name of deployment",
						Required: true,
					},
					cli.StringFlag{
						Name:     "namespace",
						Usage:    "the namespace of deployment",
						Required: true,
					},
					cli.StringFlag{
						Name:     "image",
						Usage:    "the docker image of deployment",
						Required: true,
					},
					cli.StringFlag{
						Name:  "secret",
						Usage: "the secret to image pull from registry",
					},
					cli.IntFlag{
						Name:  "requestCpu",
						Usage: "the request cpu to deploy",
						Value: 500,
					},
					cli.IntFlag{
						Name:  "requestMem",
						Usage: "the request memory to deploy",
						Value: 512,
					},
					cli.IntFlag{
						Name:  "limitCpu",
						Usage: "the limit cpu to deploy",
						Value: 1000,
					},
					cli.IntFlag{
						Name:  "limitMem",
						Usage: "the limit memory to deploy",
						Value: 1024,
					},
					cli.StringFlag{
						Name:     "o",
						Usage:    "the output yaml file",
						Required: true,
					},
					cli.IntFlag{
						Name:  "replicas",
						Usage: "the number of replicas to deploy",
						Value: 3,
					},
					cli.IntFlag{
						Name:  "revisions",
						Usage: "the number of revision history to limit",
						Value: 5,
					},
					cli.IntFlag{
						Name:     "port",
						Usage:    "the port of the deployment to listen on pod",
						Required: true,
					},
					cli.IntFlag{
						Name:  "nodePort",
						Usage: "the nodePort of the deployment to expose",
						Value: 0,
					},
					cli.IntFlag{
						Name:  "minReplicas",
						Usage: "the min replicas to deploy",
						Value: 3,
					},
					cli.IntFlag{
						Name:  "maxReplicas",
						Usage: "the max replicas of deploy",
						Value: 10,
					},
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
					cli.StringFlag{
						Name:  "serviceAccount",
						Usage: "the ServiceAccount for the deployment",
					},
				},
				Action: kube.DeploymentCommand,
			},
		},
	},
	{
		Name:  "rpc",
		Usage: "generate rpc code",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Usage: "output a sample proto file",
			},
			cli.StringFlag{
				Name: "home",
				Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority",
			},
			cli.StringFlag{
				Name: "remote",
				Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
					"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
					"https://github.com/zeromicro/go-zero-template directory structure",
			},
			cli.StringFlag{
				Name:  "branch",
				Usage: "the branch of the remote repo, it does work with --remote",
			},
		},
		Action: rpc.RPCTemplate,
		Subcommands: []cli.Command{
			{
				Name:      "new",
				Usage:     `generate rpc demo service`,
				UsageText: "example: goctl rpc new [options] service-name",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:   "go_opt",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "go-grpc_opt",
						Hidden: true,
					},
					cli.StringFlag{
						Name:  "style",
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
					},
					cli.BoolFlag{
						Name:  "idea",
						Usage: "whether the command execution environment is from idea plugin. [optional]",
					},
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
					cli.BoolFlag{
						Name:  "verbose, v",
						Usage: "enable log output",
					},
				},
				Action: rpc.RPCNew,
			},
			{
				Name:  "template",
				Usage: `generate proto template`,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "out, o",
						Usage: "the target path of proto (deprecated)",
					},
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time," +
							" if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
				},
				Action: rpc.RPCTemplate,
			},
			{
				Name:        "protoc",
				Usage:       "generate grpc code",
				UsageText:   `example: goctl rpc protoc xx.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.`,
				Description: "for details, see https://go-zero.dev/cn/goctl-rpc.html",
				Action:      rpc.ZRPC,
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:   "go_out",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "go-grpc_out",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "go_opt",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "go-grpc_opt",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "plugin",
						Hidden: true,
					},
					cli.StringSliceFlag{
						Name:   "proto_path,I",
						Hidden: true,
					},
					cli.StringFlag{
						Name:  "zrpc_out",
						Usage: "the zrpc output directory",
					},
					cli.StringFlag{
						Name:  "style",
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
					},
					cli.StringFlag{
						Name:  "home",
						Usage: "the goctl home path of the template",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
					cli.BoolFlag{
						Name:  "verbose, v",
						Usage: "enable log output",
					},
				},
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
								Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
							},
							cli.BoolFlag{
								Name:  "cache, c",
								Usage: "generate code with cache [optional]",
							},
							cli.BoolFlag{
								Name:  "idea",
								Usage: "for idea plugin [optional]",
							},
							cli.StringFlag{
								Name:  "database, db",
								Usage: "the name of database [optional]",
							},
							cli.StringFlag{
								Name: "home",
								Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority",
							},
							cli.StringFlag{
								Name: "remote",
								Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
									"https://github.com/zeromicro/go-zero-template directory structure",
							},
							cli.StringFlag{
								Name:  "branch",
								Usage: "the branch of the remote repo, it does work with --remote",
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
								Usage: `the data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`,
							},
							cli.StringSliceFlag{
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
								Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
							},
							cli.BoolFlag{
								Name:  "idea",
								Usage: "for idea plugin [optional]",
							},
							cli.StringFlag{
								Name: "home",
								Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority",
							},
							cli.StringFlag{
								Name: "remote",
								Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
									"https://github.com/zeromicro/go-zero-template directory structure",
							},
							cli.StringFlag{
								Name:  "branch",
								Usage: "the branch of the remote repo, it does work with --remote",
							},
						},
						Action: model.MySqlDataSource,
					},
				},
			},
			{
				Name:  "pg",
				Usage: `generate postgresql model`,
				Subcommands: []cli.Command{
					{
						Name:  "datasource",
						Usage: `generate model from datasource`,
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "url",
								Usage: `the data source of database,like "postgres://root:password@127.0.0.1:5432/database?sslmode=disable"`,
							},
							cli.StringFlag{
								Name:  "table, t",
								Usage: `the table or table globbing patterns in the database`,
							},
							cli.StringFlag{
								Name:  "schema, s",
								Usage: `the table schema, default is [public]`,
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
								Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
							},
							cli.BoolFlag{
								Name:  "idea",
								Usage: "for idea plugin [optional]",
							},
							cli.StringFlag{
								Name: "home",
								Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority",
							},
							cli.StringFlag{
								Name: "remote",
								Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
									"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
									"https://github.com/zeromicro/go-zero-template directory structure",
							},
							cli.StringFlag{
								Name:  "branch",
								Usage: "the branch of the remote repo, it does work with --remote",
							},
						},
						Action: model.PostgreSqlDataSource,
					},
				},
			},
			{
				Name:  "mongo",
				Usage: `generate mongo model`,
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "type, t",
						Usage: "specified model type name",
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
						Usage: "the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]",
					},
					cli.StringFlag{
						Name: "home",
						Usage: "the goctl home path of the template, --home and --remote cannot be set at the same time," +
							" if they are, --remote has higher priority",
					},
					cli.StringFlag{
						Name: "remote",
						Usage: "the remote git repo of the template, --home and --remote cannot be set at the same time, " +
							"if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the " +
							"https://github.com/zeromicro/go-zero-template directory structure",
					},
					cli.StringFlag{
						Name:  "branch",
						Usage: "the branch of the remote repo, it does work with --remote",
					},
				},
				Action: mongo.Action,
			},
		},
	},
	{
		Name:  "template",
		Usage: "template operation",
		Subcommands: []cli.Command{
			{
				Name:  "init",
				Usage: "initialize the all templates(force update)",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "home",
						Usage: "the goctl home path of the template",
					},
				},
				Action: tpl.GenTemplates,
			},
			{
				Name:  "clean",
				Usage: "clean the all cache templates",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "home",
						Usage: "the goctl home path of the template",
					},
				},
				Action: tpl.CleanTemplates,
			},
			{
				Name:  "update",
				Usage: "update template of the target category to the latest",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "category,c",
						Usage: "the category of template, enum [api,rpc,model,docker,kube]",
					},
					cli.StringFlag{
						Name:  "home",
						Usage: "the goctl home path of the template",
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
						Usage: "the category of template, enum [api,rpc,model,docker,kube]",
					},
					cli.StringFlag{
						Name:  "name,n",
						Usage: "the target file name of template",
					},
					cli.StringFlag{
						Name:  "home",
						Usage: "the goctl home path of the template",
					},
				},
				Action: tpl.RevertTemplates,
			},
		},
	},
	{
		Name:   "completion",
		Usage:  "generation completion script, it only works for unix-like OS",
		Action: completion.Completion,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "the filename of auto complete script, default is [goctl_autocomplete]",
			},
		},
	},
}

func main() {
	logx.Disable()
	load.Disable()

	cli.BashCompletionFlag = cli.BoolFlag{
		Name:   completion.BashCompletionFlag,
		Hidden: true,
	}
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands

	// cli already print error messages.
	if err := app.Run(os.Args); err != nil {
		fmt.Println(aurora.Red(err.Error()))
		os.Exit(codeFailure)
	}
}
