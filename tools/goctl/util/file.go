package util

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/logx"
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

func CreateTemplateAndExecute(filename, text string, arg map[string]interface{}, forceUpdate bool, disableFormatCodeArgs ...bool) error {
	if FileExists(filename) && !forceUpdate {
		return nil
	}
	var buffer = new(bytes.Buffer)
	templateName := fmt.Sprintf("%d", time.Now().UnixNano())
	t, err := template.New(templateName).Parse(text)
	if err != nil {
		return err
	}
	err = t.Execute(buffer, arg)
	if err != nil {
		return err
	}
	var disableFormatCode bool
	for _, f := range disableFormatCodeArgs {
		disableFormatCode = f
	}
	var bts = buffer.Bytes()
	s := buffer.String()
	logx.Info(s)
	if !disableFormatCode {
		bts, err = format.Source(buffer.Bytes())
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, bts, os.ModePerm)
}

func FormatCodeAndWrite(filename string, code []byte) error {
	if FileExists(filename) {
		return nil
	}
	bts, err := format.Source(code)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bts, os.ModePerm)
}
