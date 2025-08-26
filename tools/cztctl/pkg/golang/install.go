package golang

import (
	"os"
	"os/exec"
)

func Install(git string) error {
	cmd := exec.Command("go", "install", git)
	env := os.Environ()
	env = append(env, "GO111MODULE=on")
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
