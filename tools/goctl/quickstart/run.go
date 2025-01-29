package quickstart

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func goStart(dir string) {
	execCommand(dir, "go run .")
}

func goModTidy(dir string) int {
	log.Debug(">> go mod tidy")
	return execCommand(dir, "go mod tidy")
}

func execCommand(dir, arg string, envArgs ...string) int {
	cmd := exec.Command("sh", "-c", arg)
	if runtime.GOOS == vars.OsWindows {
		cmd = exec.Command("cmd.exe", "/c", arg)
	}
	env := append([]string(nil), os.Environ()...)
	env = append(env, envArgs...)
	cmd.Env = env
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return cmd.ProcessState.ExitCode()
}
