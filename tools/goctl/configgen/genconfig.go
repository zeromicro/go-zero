package configgen

import (
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
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

// GenConfigCommand provides the entry of goctl config
func GenConfigCommand(c *cli.Context) error {
	path, err := filepath.Abs(c.String("path"))
	if err != nil {
		errorx.Must(errors.New("abs failed: " + c.String("path")))
	}

	goModPath, found := util.FindGoModPath(path)
	if !found {
		errorx.Must(errors.New("go mod not initial"))
	}

	path = strings.TrimSuffix(path, "/config.go")
	location := filepath.Join(path, "tmp")
	errorx.Must(os.MkdirAll(location, os.ModePerm))

	goPath := filepath.Join(location, "config.go")
	fp, err := os.Create(goPath)
	errorx.Must(err)
	defer fp.Close()
	defer os.RemoveAll(location)

	t := template.Must(template.New("template").Parse(configTemplate))
	errorx.Must(t.Execute(fp, map[string]string{
		"import": filepath.Dir(goModPath),
	}))

	gen := exec.Command("go", "run", "config.go")
	gen.Dir = filepath.Dir(goPath)
	gen.Stderr = os.Stderr
	gen.Stdout = os.Stdout
	errorx.Must(gen.Run())

	path, err = os.Getwd()
	errorx.Must(err)

	errorx.Must(os.Rename(filepath.Dir(goPath)+"/config.yaml", path+"/config.yaml"))
	fmt.Println(aurora.Green("Done."))
	return nil
}
