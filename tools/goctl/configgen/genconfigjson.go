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
	"github.com/urfave/cli"
	"zero/tools/goctl/vars"
)

const configTemplate = `package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"{{.import}}"
)

func main() {
	var c config.Config
	template, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("config.json", template, os.ModePerm)
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
	xi := strings.Index(path, vars.ProjectName)
	if xi <= 0 {
		return errors.New("path should the absolute path of config go file")
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
		"import": path[xi:],
	}); err != nil {
		return err
	}

	cmd := exec.Command("go", "run", goPath)
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(aurora.Green("Done."))
	return nil
}
