package plugin

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

const pluginNamePrefix = "plugin_"

func Do(args []string, osArgs []string) error {
	var plugin = flag.String("plugin", "", "")
	pluginArgs, err := prepareArgs()
	if err != nil {
		return err
	}

	if len(*plugin) == 0 {
		return errors.New("missing plugin")
	}

	bin, download, err := getCommand(args[1])
	if err != nil {
		return err
	}
	if download {
		defer os.Remove(bin)
	}

	var commands []string
	commands = append(commands, bin)
	commands = append(commands, args[2:]...)
	commands = append(commands, osArgs[1:]...)
	commands = append(commands, pluginArgs...)

	_, err = execx.Run(strings.Join(commands, " "), "")
	return err
}

func prepareArgs() ([]string, error) {
	var pluginArgs []string
	var apiPath = flag.String("api", "", "")
	if len(*apiPath) > 0 && util.FileExists(*apiPath) {
		p, err := parser.NewParser(*apiPath)
		if err != nil {
			return nil, err
		}

		api, err := p.Parse()
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(api)
		if err != nil {
			return nil, err
		}
		pluginArgs = append(pluginArgs, "-spec")
		pluginArgs = append(pluginArgs, string(data))
	}

	var dir = flag.String("-dir", "", "")
	if len(*dir) > 0 {
		abs, err := filepath.Abs(*dir)
		if err != nil {
			return nil, err
		}

		projectCtx, err := ctx.Prepare(abs)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(projectCtx)
		if err != nil {
			return nil, err
		}
		pluginArgs = append(pluginArgs, "-context")
		pluginArgs = append(pluginArgs, string(data))
	}

	return pluginArgs, nil
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

		var filename = pluginNamePrefix + items[len(items)-1]
		err := downloadFile(filename, arg)
		if err != nil {
			return "", false, err
		}

		os.Chmod(filename, os.ModePerm)
		return filename, true, nil
	}
	return "", false, defaultErr
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
