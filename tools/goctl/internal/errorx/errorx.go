package errorx

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"os"
	"runtime"
	"strings"
)

func Must(err error, msg ...string) {
	if err == nil {
		return
	}

	buildVersion := os.Getenv("GOCTL_VERSION")
	version := fmt.Sprintf("goctl version: %s %s/%s", buildVersion, runtime.GOOS, runtime.GOARCH)
	colorConsole := console.NewColorConsole()
	output := fmt.Sprintf(`goctl: generation error: %s`, err.Error())
	colorConsole.Error(output)
	colorConsole.Debug(version)
	if len(msg) > 0 {
		colorConsole.Debug(`[message]: %s`, strings.Join(msg, "\n"))
	}

	os.Exit(1)
}
