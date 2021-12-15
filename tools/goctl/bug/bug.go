package bug

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

const (
	windows = "windows"
	darwin  = "darwin"

	windowsOpen = "start"
	darwinOpen  = "open"
	linuxOpen   = "xdg-open"

	os           = "OS"
	arch         = "ARCH"
	goctlVersion = "GOCTL_VERSION"
)

var openCmd = map[string]string{
	windows: windowsOpen,
	darwin:  darwinOpen,
}

func Action(_ *cli.Context) error {
	env := getEnv()
	content := fmt.Sprintf(issueTemplate, "<pre>\n"+env.string()+"</pre>")
	content = url.QueryEscape(content)
	url := fmt.Sprintf("https://github.com/zeromicro/go-zero/issues/new?title=TODO&body=%s", content)

	goos := runtime.GOOS
	var cmd string
	var args []string
	cmd, ok := openCmd[goos]
	if !ok {
		cmd = linuxOpen
	}
	if goos == windows {
		args = []string{"/c", "start"}
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
