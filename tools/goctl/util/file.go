package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
)

// NL defines a new line
const (
	NL       = "\n"
	goctlDir = ".goctl"
)

var goctlHome string

func RegisterGoctlHome(home string) {
	goctlHome = home
}

// CreateIfNotExist creates a file if it is not exists
func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}

// RemoveIfExist deletes the specified file if it is exists
func RemoveIfExist(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	return os.Remove(filename)
}

// RemoveOrQuit deletes the specified file if read a permit command from stdin
func RemoveOrQuit(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	fmt.Printf("%s exists, overwrite it?\nEnter to overwrite or Ctrl-C to cancel...",
		aurora.BgRed(aurora.Bold(filename)))
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return os.Remove(filename)
}

// FileExists returns true if the specified file is exists
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// FileNameWithoutExt returns a file name without suffix
func FileNameWithoutExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}

// GetGoctlHome returns the path value of the goctl home where Join $HOME with .goctl
func GetGoctlHome() (string, error) {
	if len(goctlHome) != 0 {
		return goctlHome, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, goctlDir), nil
}

// GetTemplateDir returns the category path value in GoctlHome where could get it by GetGoctlHome
func GetTemplateDir(category string) (string, error) {
	goctlHome, err := GetGoctlHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(goctlHome, category), nil
}

// InitTemplates creates template files GoctlHome where could get it by GetGoctlHome
func InitTemplates(category string, templates map[string]string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}

	if err := MkdirIfNotExist(dir); err != nil {
		return err
	}

	for k, v := range templates {
		if err := createTemplate(filepath.Join(dir, k), v, false); err != nil {
			return err
		}
	}

	return nil
}

// CreateTemplate writes template into file even it is exists
func CreateTemplate(category, name, content string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}
	return createTemplate(filepath.Join(dir, name), content, true)
}

// Clean deletes all templates and removes the parent directory
func Clean(category string) error {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

// LoadTemplate gets template content by the specified file
func LoadTemplate(category, file, builtin string) (string, error) {
	dir, err := GetTemplateDir(category)
	if err != nil {
		return "", err
	}

	file = filepath.Join(dir, file)
	if !FileExists(file) {
		return builtin, nil
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func createTemplate(file, content string, force bool) error {
	if FileExists(file) && !force {
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
