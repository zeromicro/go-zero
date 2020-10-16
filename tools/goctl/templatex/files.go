package templatex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const goctlDir = ".goctl"

func InitTemplates(category string, templates map[string]string) error {
	dir, err := getTemplateDir(category)
	if err != nil {
		return err
	}

	if err := util.MkdirIfNotExist(dir); err != nil {
		return err
	}

	for k, v := range templates {
		if err := createTemplate(filepath.Join(dir, k), v); err != nil {
			return err
		}
	}

	fmt.Printf("Templates are generated in %s, %s\n", aurora.Green(dir),
		aurora.Red("edit on your risk!"))

	return nil
}

func LoadTemplate(category, file, builtin string) (string, error) {
	dir, err := getTemplateDir(category)
	if err != nil {
		return "", err
	}

	file = filepath.Join(dir, file)
	if !util.FileExists(file) {
		return builtin, nil
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func createTemplate(file, content string) error {
	if util.FileExists(file) {
		println(1)
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func getTemplateDir(category string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, goctlDir, category), nil
}
