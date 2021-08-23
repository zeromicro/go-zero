package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

// Run provides the execution of shell scripts in golang,
// which can support macOS, Windows, and Linux operating systems.
// Other operating systems are currently not supported
func Run(arg, dir string, in ...*bytes.Buffer) (string, error) {
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
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if len(in) > 0 {
		cmd.Stdin = in[0]
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSuffix(stderr.String(), util.NL))
		}
		return "", err
	}

	return strings.TrimSuffix(stdout.String(), util.NL), nil
}
