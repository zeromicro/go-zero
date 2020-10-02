package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

func Run(arg string, dir string) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case vars.OsMac, vars.OsLinux:
		cmd = exec.Command("sh", "-c", arg)
	case vars.OsWindows:
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
	}
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	dtsout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = dtsout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(stderr.String())
		}
		return "", err
	}

	return dtsout.String(), nil
}
