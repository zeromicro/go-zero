package golang

import (
	"os"
	"os/exec"
)

func Install(git string) error {
	cmd := exec.Command("go", "install", git)
	env := append([]string{
		"GO111MODULE", "on",
		"GOPROXY", "https://goproxy.cn,direct",
	}, os.Environ()...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
