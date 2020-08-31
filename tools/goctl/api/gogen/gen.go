package gogen

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/logx"
	apiformat "github.com/tal-tech/go-zero/tools/goctl/api/format"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const tmpFile = "%s-%d"

var tmpDir = path.Join(os.TempDir(), "goctl")

func GoCommand(c *cli.Context) error {
	apiFile := c.String("api")
	dir := c.String("dir")
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	p, err := parser.NewParser(apiFile)
	if err != nil {
		return err
	}
	api, err := p.Parse()
	if err != nil {
		return err
	}

	logx.Must(util.MkdirIfNotExist(dir))
	logx.Must(genEtc(dir, api))
	logx.Must(genConfig(dir))
	logx.Must(genMain(dir, api))
	logx.Must(genServiceContext(dir, api))
	logx.Must(genTypes(dir, api))
	logx.Must(genHandlers(dir, api))
	logx.Must(genRoutes(dir, api))
	logx.Must(genLogic(dir, api))
	// it does not work
	format(dir)
	createGoModFileIfNeed(dir)

	if err := backupAndSweep(apiFile); err != nil {
		return err
	}

	if err = apiformat.ApiFormat(apiFile, false); err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}

func backupAndSweep(apiFile string) error {
	var err error
	var wg sync.WaitGroup

	wg.Add(2)
	_ = os.MkdirAll(tmpDir, os.ModePerm)

	go func() {
		_, fileName := filepath.Split(apiFile)
		_, e := apiutil.Copy(apiFile, fmt.Sprintf(path.Join(tmpDir, tmpFile), fileName, time.Now().Unix()))
		if e != nil {
			err = e
		}
		wg.Done()
	}()
	go func() {
		if e := sweep(); e != nil {
			err = e
		}
		wg.Done()
	}()
	wg.Wait()

	return err
}

func format(dir string) {
	cmd := exec.Command("go", "fmt", "./"+dir+"...")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func sweep() error {
	keepTime := time.Now().AddDate(0, 0, -7)
	return filepath.Walk(tmpDir, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		pos := strings.LastIndexByte(info.Name(), '-')
		if pos > 0 {
			timestamp := info.Name()[pos+1:]
			seconds, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				// print error and ignore
				fmt.Println(aurora.Red(fmt.Sprintf("sweep ignored file: %s", fpath)))
				return nil
			}

			tm := time.Unix(seconds, 0)
			if tm.Before(keepTime) {
				if err := os.Remove(fpath); err != nil {
					fmt.Println(aurora.Red(fmt.Sprintf("failed to remove file: %s", fpath)))
					return err
				}
			}
		}

		return nil
	})
}

func createGoModFileIfNeed(dir string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	_, hasGoMod := util.FindGoModPath(dir)
	if hasGoMod {
		return
	}

	gopath := os.Getenv("GOPATH")
	parent := path.Join(gopath, "src")
	pos := strings.Index(absDir, parent)
	if pos >= 0 {
		return
	}

	moduleName := absDir[len(filepath.Dir(absDir))+1:]
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		fmt.Println(err.Error())
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf(outStr + "\n" + errStr)
}
