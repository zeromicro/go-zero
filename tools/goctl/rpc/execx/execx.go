package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Run(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if stderr.Len() != 0 {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}

func RunShOrBat(arg string) error {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case "darwin":
		cmd = exec.Command("sh", "-c", arg)
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return fmt.Errorf("unexpected os: %v", goos)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
