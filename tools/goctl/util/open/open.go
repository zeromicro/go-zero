package open

import (
	"os/exec"
	"runtime"
)

const (
	windows = "windows"
	darwin  = "darwin"

	windowsOpen = "start"
	darwinOpen  = "open"
	linuxOpen   = "xdg-open"
)

var openCmd = map[string]string{
	windows: windowsOpen,
	darwin:  darwinOpen,
}

func Open(url string) error {
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
