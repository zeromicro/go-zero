package plugin

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

const (
	pluginArg   = "-plugin"
	specJson    = "spec.json"
	contextJson = "context.json"
)

func Do(args []string, osArgs []string) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fs := flag.NewFlagSet("goctl", flag.ContinueOnError)
	var plugin = fs.String("plugin", "", "")
	for index, item := range os.Args {
		if strings.HasPrefix(item, pluginArg) {
			fs.Parse(os.Args[index:])
			break
		}
	}

	pluginArgs, tempFile, err := prepareArgs()
	if err != nil {
		return err
	}

	defer func() {
		for _, item := range tempFile {
			os.Remove(item)
		}
	}()

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
	fmt.Printf("%+v", commands)
	_, err = execx.Run(strings.Join(commands, " "), filepath.Dir(ex))
	return err
}

func prepareArgs() ([]string, []string, error) {
	var pluginArgs []string
	var tempFiles []string
	fs := flag.NewFlagSet("goctl", flag.ContinueOnError)
	var apiPath = fs.String("api", "", "")
	var dir = fs.String("dir", "", "")
	var args []string
	for index, item := range os.Args {
		if item == "-api" || item == "-dir" {
			args = append(args, item, os.Args[index+1])
		}
	}
	_ = fs.Parse(args)

	var timestamp = fmt.Sprint(time.Now().Unix())
	if len(*apiPath) > 0 && util.FileExists(*apiPath) {
		p, err := parser.NewParser(*apiPath)
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

		pluginArgs = append(pluginArgs, "-spec")
		absFile, _ := filepath.Abs(filename)
		pluginArgs = append(pluginArgs, absFile)
		tempFiles = append(tempFiles, absFile)
	}

	if len(*dir) > 0 {
		abs, err := filepath.Abs(*dir)
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

		pluginArgs = append(pluginArgs, "-context")
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
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
