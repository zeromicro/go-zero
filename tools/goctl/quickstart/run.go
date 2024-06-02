package quickstart

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

const goproxy = "GOPROXY=https://goproxy.cn,direct"

func goStart(dir string) {
	execCommand(dir, "go run .", prepareGoProxyEnv()...)
}

func goModTidy(dir string) int {
	log.Debug(">> go mod tidy")
	return execCommand(dir, "go mod tidy", prepareGoProxyEnv()...)
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

func prepareGoProxyEnv(envArgs ...string) []string {
	if env.InChina() {
		return append(envArgs, goproxy)
	}

	return envArgs
}
