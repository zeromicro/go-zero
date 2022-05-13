package quickstart

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

func goStart(dir string) {
	goproxy := "GOPROXY=https://goproxy.cn"
	execCommand(dir, "go run .", goproxy)
}

func goModTidy(dir string) int {
	goproxy := "GOPROXY=https://goproxy.cn"
	log.Debug(">> go mod tidy")
	return execCommand(dir, "go mod tidy", goproxy)
}

func execCommand(dir string, arg string, envArgs ...string) int {
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
