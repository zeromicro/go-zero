package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"zero/core/conf"
	"zero/core/hash"
	"zero/core/lang"
	"zero/core/logx"
	"zero/core/mr"
	"zero/core/stringx"
	"zero/tools/goctl/api/apigen"
	"zero/tools/goctl/api/dartgen"
	"zero/tools/goctl/api/docgen"
	"zero/tools/goctl/api/format"
	"zero/tools/goctl/api/gogen"
	"zero/tools/goctl/api/javagen"
	"zero/tools/goctl/api/tsgen"
	"zero/tools/goctl/api/validate"
	"zero/tools/goctl/configgen"
	"zero/tools/goctl/docker"
	"zero/tools/goctl/feature"
	"zero/tools/goctl/model/mongomodel"
	"zero/tools/goctl/util"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
)

const (
	autoUpdate     = "GOCTL_AUTO_UPDATE"
	configFile     = ".goctl"
	configTemplate = `url = http://47.97.184.41:7777/`
	toolName       = "goctl"
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
			Name:  "model",
			Usage: "generate sql model",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "the file that contains main function",
				},
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "the target dir",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "mongo",
					Usage: "generate mongoModel files for provided somemongo.go in go file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "src, s",
							Usage: "the src file",
						},
						cli.StringFlag{
							Name:  "cache",
							Usage: "need cache code([yes/no] default value is no)",
						},
					},
					Action: mongomodel.ModelCommond,
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

func genConfigFile(file string) error {
	return ioutil.WriteFile(file, []byte(configTemplate), 0600)
}

func getAbsFile() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(filepath.Dir(exe))
	if err != nil {
		return "", err
	}

	return path.Join(dir, filepath.Base(os.Args[0])), nil
}

func getFilePerm(file string) (os.FileMode, error) {
	info, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}

func update() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	absConfigFile := path.Join(usr.HomeDir, configFile)
	if !util.FileExists(absConfigFile) {
		if err := genConfigFile(absConfigFile); err != nil {
			fmt.Println(err)
			return
		}
	}

	props, err := conf.LoadProperties(absConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	u, err := url.Parse(props.GetString("url"))
	if err != nil {
		fmt.Println(err)
		return
	}

	u.Path = path.Join(u.Path, toolName)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := getAbsFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Md5", hash.Md5Hex(content))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	mode, err := getFilePerm(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if err := ioutil.WriteFile(file, content, mode); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	logx.Disable()

	done := make(chan lang.PlaceholderType)
	mr.FinishVoid(func() {
		if os.Getenv(autoUpdate) != "off" && !stringx.Contains(os.Args, "-iu") {
			update()
		}
		close(done)
	}, func() {
		app := cli.NewApp()
		app.Usage = "a cli tool to generate code"
		app.Version = BuildTime
		app.Commands = commands
		// cli already print error messages
		if err := app.Run(os.Args); err != nil {
			fmt.Println("error:", err)
		}
	}, func() {
		select {
		case <-done:
		case <-time.After(time.Second):
			fmt.Println(aurora.Yellow("Updating goctl, please wait..."))
		}
	})
}
