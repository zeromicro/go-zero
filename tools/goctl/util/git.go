package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func CloneIntoGitHome(url string) (dir string, err error) {
	gitHome, err := pathx.GetGitHome()
	if err != nil {
		return "", err
	}
	os.RemoveAll(gitHome)
	ext := filepath.Ext(url)
	repo := strings.TrimSuffix(filepath.Base(url), ext)
	dir = filepath.Join(gitHome, repo)
	path, err := env.LookPath("git")
	if err != nil {
		return "", err
	}
	if !env.CanExec() {
		return "", fmt.Errorf("os %q can not call 'exec' command", runtime.GOOS)
	}
	cmd := exec.Command(path, "clone", url, dir)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
