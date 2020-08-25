package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
)

func CreateIfNotExist(file string) (*os.File, error) {
	_, err := os.Stat(file)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s already exist", file)
	}

	return os.Create(file)
}

func RemoveIfExist(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	return os.Remove(filename)
}

func RemoveOrQuit(filename string) error {
	if !FileExists(filename) {
		return nil
	}

	fmt.Printf("%s exists, overwrite it?\nEnter to overwrite or Ctrl-C to cancel...",
		aurora.BgRed(aurora.Bold(filename)))
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return os.Remove(filename)
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func FileNameWithoutExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}
