package bug

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
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
	goVersion    = "GO_VERSION"
)

var openCmd = map[string]string{
	windows: windowsOpen,
	darwin:  darwinOpen,
}

func runE(_ *cobra.Command, _ []string) error {
	env := getEnv()
	content := fmt.Sprintf(issueTemplate, version.BuildVersion, env.string())
	content = url.QueryEscape(content)
	url := fmt.Sprintf("https://github.com/zeromicro/go-zero/issues/new?body=%s", content)

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
