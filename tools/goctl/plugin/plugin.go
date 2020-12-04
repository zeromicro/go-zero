package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
	"github.com/urfave/cli"
)

const (
	pluginPrefix = "plugin-"
	specJson     = "spec.json"
	contextJson  = "context.json"
	specTag      = "-spec"
	contextTag   = "-context"
)

type Plugin struct {
	Api     *spec.ApiSpec
	Context *ctx.ProjectContext
}

func PluginCommand(c *cli.Context) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	var plugin = c.String("plugin")
	if len(plugin) == 0 {
		return errors.New("missing plugin")
	}

	pluginArgs, tempFile, err := prepareArgs(c)
	if err != nil {
		return err
	}

	defer func() {
		for _, item := range tempFile {
			os.Remove(item)
		}
	}()

	bin, download, err := getCommand(plugin)
	if err != nil {
		return err
	}
	if download {
		defer os.Remove(bin)
	}

	var commands []string
	commands = append(commands, bin)
	commands = append(commands, pluginArgs...)

	content, err := execx.Run(strings.Join(commands, " "), filepath.Dir(ex))
	if err != nil {
		return err
	}

	fmt.Println(content)
	return nil
}

func prepareArgs(c *cli.Context) ([]string, []string, error) {
	var pluginArgs []string
	var tempFiles []string
	apiPath := c.String("api")
	dir := c.String("dir")

	var timestamp = fmt.Sprint(time.Now().Unix())
	if len(apiPath) > 0 && util.FileExists(apiPath) {
		p, err := parser.NewParser(apiPath)
		if err != nil {
			return nil, nil, err
		}

		api, err := p.Parse()
		if err != nil {
			return nil, nil, err
		}

		data, err := json.Marshal(api)
		if err != nil {
			return nil, nil, err
		}

		filename := timestamp + "-" + specJson
		err = ioutil.WriteFile(filename, data, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}

		pluginArgs = append(pluginArgs, specTag)
		absFile, _ := filepath.Abs(filename)
		pluginArgs = append(pluginArgs, absFile)
		tempFiles = append(tempFiles, absFile)
	}

	if len(dir) > 0 {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, nil, err
		}

		projectCtx, err := ctx.Prepare(abs)
		if err != nil {
			return nil, nil, err
		}

		data, err := json.Marshal(projectCtx)
		if err != nil {
			return nil, nil, err
		}

		filename := timestamp + "-" + contextJson
		err = ioutil.WriteFile(filename, data, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}

		pluginArgs = append(pluginArgs, contextTag)
		absFile, _ := filepath.Abs(filename)
		pluginArgs = append(pluginArgs, absFile)
		tempFiles = append(tempFiles, absFile)
	}

	return pluginArgs, tempFiles, nil
}

func getCommand(arg string) (string, bool, error) {
	if util.FileExists(arg) {
		return arg, false, nil
	}

	var defaultErr = errors.New("invalid plugin value " + arg)
	if strings.HasPrefix(arg, "http") {
		items := strings.Split(arg, "/")
		if len(items) == 0 {
			return "", false, defaultErr
		}

		var filename = pluginPrefix + items[len(items)-1]
		err := downloadFile(filename, arg)
		if err != nil {
			return "", false, err
		}

		os.Chmod(filename, os.ModePerm)
		return filename, true, nil
	}
	return arg, false, nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func NewPlugin(args []string) (*Plugin, error) {
	var plugin Plugin
	for index, item := range args {
		if item == specTag {
			specFile := args[index+1]
			var api spec.ApiSpec
			content, err := ioutil.ReadFile(specFile)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(content, &api)
			if err != nil {
				return nil, err
			}

			plugin.Api = &api
		}

		if item == contextTag {
			contextFile := args[index+1]
			var context ctx.ProjectContext
			content, err := ioutil.ReadFile(contextFile)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(content, &context)
			if err != nil {
				return nil, err
			}

			plugin.Context = &context
		}
	}
	return &plugin, nil
}
