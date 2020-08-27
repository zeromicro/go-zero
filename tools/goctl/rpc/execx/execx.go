package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Run(arg string) (string, error) {
	goos := runtime.GOOS
	var cmd *exec.Cmd
	switch goos {
	case "darwin":
		cmd = exec.Command("sh", "-c", arg)
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", arg)
	default:
		return "", fmt.Errorf("unexpected os: %v", goos)
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
	stderr := new(bytes.Buffer)
	cmd.Stdout = os.Stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return err
	}
	return nil
}
