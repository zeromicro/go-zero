package completion

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func Completion(c *cli.Context) error {
	goos := runtime.GOOS
	if goos == vars.OsWindows {
		return fmt.Errorf("%q: only support unix-like OS", goos)
	}

	name := c.String("name")
	if len(name) == 0 {
		name = defaultCompletionFilename
	}
	if filepath.IsAbs(name) {
		return fmt.Errorf("unsupport absolute path: %q", name)
	}

	home, err := pathx.GetAutoCompleteHome()
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	zshF := filepath.Join(home, "zsh", defaultCompletionFilename)
	err = pathx.MkdirIfNotExist(filepath.Dir(zshF))
	if err != nil {
		return err
	}

	bashF := filepath.Join(home, "bash", defaultCompletionFilename)
	err = pathx.MkdirIfNotExist(filepath.Dir(bashF))
	if err != nil {
		return err
	}

	flag := magic
	err = ioutil.WriteFile(zshF, zsh, os.ModePerm)
	if err == nil {
		flag |= flagZsh
	}

	err = ioutil.WriteFile(bashF, bash, os.ModePerm)
	if err == nil {
		flag |= flagBash
	}

	buffer.WriteString(aurora.BrightGreen("generation auto completion success!\n").String())
	buffer.WriteString(aurora.BrightGreen("executes the following script to setting shell:\n").String())
	switch flag {
	case magic | flagZsh:
		buffer.WriteString(aurora.BrightCyan(fmt.Sprintf("echo PROG=goctl source %s >> ~/.zshrc && source ~/.zshrc", zshF)).String())
	case magic | flagBash:
		buffer.WriteString(aurora.BrightCyan(fmt.Sprintf("echo PROG=goctl source %s >> ~/.bashrc && source ~/.bashrc", bashF)).String())
	case magic | flagZsh | flagBash:
		buffer.WriteString(aurora.BrightCyan(fmt.Sprintf(`echo PROG=goctl source %s >> ~/.zshrc && source ~/.zshrc
or
echo PROG=goctl source %s >> ~/.bashrc && source ~/.bashrc`, zshF, bashF)).String())
	default:
		return nil
	}

	fmt.Println(buffer.String())
	return nil
}
