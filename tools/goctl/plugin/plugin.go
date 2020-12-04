package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	pluginArg   = "-plugin"
	specJson    = "spec.json"
	contextJson = "context.json"
	specTag     = "-spec"
	contextTag  = "-context"
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

	transferData, err := prepareArgs(c)
	if err != nil {
		return err
	}

	bin, download, err := getCommand(plugin)
	if err != nil {
		return err
	}
	if download {
		defer func() {
			_ = os.Remove(bin)
		}()
	}

	content, err := execx.Run(bin, filepath.Dir(ex), bytes.NewBuffer(transferData))
	if err != nil {
		return err
	}

	fmt.Println(content)
	return nil
}

func prepareArgs(c *cli.Context) ([]byte, error) {
	apiPath := c.String("api")
	dir := c.String("dir")

	var transferData Plugin
	if len(apiPath) > 0 && util.FileExists(apiPath) {
		p, err := parser.NewParser(apiPath)
		if err != nil {
			return nil, err
		}

		api, err := p.Parse()
		if err != nil {
			return nil, err
		}
		transferData.Api = api
	}

	if len(dir) > 0 {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}

		projectCtx, err := ctx.Prepare(abs)
		if err != nil {
			return nil, err
		}
		transferData.Context = projectCtx
	}

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

		var filename = pluginArg + items[len(items)-1]
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
	err = json.Unmarshal(content, &plugin)
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}
