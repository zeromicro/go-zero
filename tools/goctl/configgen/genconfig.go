package configgen

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const configTemplate = `package main

import (
	"io/ioutil"
	"os"
	"{{.import}}"

	"github.com/ghodss/yaml"
)

func main() {
	var c config.Config
	template, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("config.yaml", template, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
`

func GenConfigCommand(c *cli.Context) error {
	path, err := filepath.Abs(c.String("path"))
	if err != nil {
		return errors.New("abs failed: " + c.String("path"))
	}
	goModPath, hasFound := util.FindGoModPath(path)
	if !hasFound {
		return errors.New("go mod not initial")
	}
	path = strings.TrimSuffix(path, "/config.go")
	location := path + "/tmp"
	err = os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return err
	}

	goPath := filepath.Join(location, "config.go")
	fp, err := os.Create(goPath)
	if err != nil {
		return err
	}
	defer fp.Close()
	defer os.RemoveAll(location)

	t := template.Must(template.New("template").Parse(configTemplate))
	if err := t.Execute(fp, map[string]string{
		"import": filepath.Dir(goModPath),
	}); err != nil {
		return err
	}

	gen := exec.Command("go", "run", "config.go")
	gen.Dir = filepath.Dir(goPath)
	gen.Stderr = os.Stderr
	gen.Stdout = os.Stdout
	err = gen.Run()
	if err != nil {
		panic(err)
	}
	path, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Rename(filepath.Dir(goPath)+"/config.yaml", path+"/config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}
