package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const pluginArg = "_plugin"

type Plugin struct {
	Api         *spec.ApiSpec
	ApiFilePath string
	Style       string
	Dir         string
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

	transferData, err := prepareArgs(c)
	if err != nil {
		return err
	}

	bin, args := getPluginAndArgs(plugin)

	bin, download, err := getCommand(bin)
	if err != nil {
		return err
	}

	if download {
		defer func() {
			_ = os.Remove(bin)
		}()
	}

	content, err := execx.Run(bin+" "+args, filepath.Dir(ex), bytes.NewBuffer(transferData))
	if err != nil {
		return err
	}

	fmt.Println(content)
	return nil
}

func prepareArgs(c *cli.Context) ([]byte, error) {
	apiPath := c.String("api")

	var transferData Plugin
	if len(apiPath) > 0 && util.FileExists(apiPath) {
		api, err := parser.Parse(apiPath)
		if err != nil {
			return nil, err
		}

		transferData.Api = api
	}

	absApiFilePath, err := filepath.Abs(apiPath)
	if err != nil {
		return nil, err
	}

	transferData.ApiFilePath = absApiFilePath
	dirAbs, err := filepath.Abs(c.String("dir"))
	if err != nil {
		return nil, err
	}

	transferData.Dir = dirAbs
	transferData.Style = c.String("style")
	data, err := json.Marshal(transferData)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getCommand(arg string) (string, bool, error) {
	p, err := exec.LookPath(arg)
	if err == nil {
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", false, err
		}
		return abs, false, nil
	}

	var defaultErr = errors.New("invalid plugin value " + arg)
	if strings.HasPrefix(arg, "http") {
		items := strings.Split(arg, "/")
		if len(items) == 0 {
			return "", false, defaultErr
		}

		filename, err := filepath.Abs(pluginArg + items[len(items)-1])
		if err != nil {
			return "", false, err
		}

		err = downloadFile(filename, arg)
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
	defer func() {
		_ = resp.Body.Close()
	}()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)
	return err
}

func NewPlugin() (*Plugin, error) {
	var plugin Plugin
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	var info struct {
		ApiFilePath string
		Style       string
		Dir         string
	}
	err = json.Unmarshal(content, &info)
	if err != nil {
		return nil, err
	}

	plugin.ApiFilePath = info.ApiFilePath
	plugin.Style = info.Style
	plugin.Dir = info.Dir
	api, err := parser.Parse(info.ApiFilePath)
	if err != nil {
		return nil, err
	}

	plugin.Api = api
	return &plugin, nil
}

func getPluginAndArgs(arg string) (string, string) {
	i := strings.Index(arg, "=")
	if i <= 0 {
		return arg, ""
	}

	return trimQuote(arg[:i]), trimQuote(arg[i+1:])
}

func trimQuote(in string) string {
	in = strings.Trim(in, `"`)
	in = strings.Trim(in, `'`)
	in = strings.Trim(in, "`")
	return in
}
