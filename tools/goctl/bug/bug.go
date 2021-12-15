package bug

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/tal-tech/go-zero/tools/goctl/internal/version"
	"github.com/urfave/cli"
)

type env map[string]string

func (e env) string() string {
	if e == nil {
		return ""
	}
	w := bytes.NewBuffer(nil)
	for k, v := range e {
		w.WriteString(fmt.Sprintf("%s = %q\n", k, v))
	}
	return w.String()
}

func Action(_ *cli.Context) error {
	env := getEnv()
	content := fmt.Sprintf(issueTemplate, "<pre>\n"+env.string()+"</pre>")
	content= url.QueryEscape(content)
	url := fmt.Sprintf("https://github.com/zeromicro/go-zero/issues/new?title=TODO&body=%s", content)
	os := runtime.GOOS
	var cmd string
	var args []string
	switch os {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func getEnv() env {
	e := make(env)
	e["OS"] = runtime.GOOS
	e["ARCH"] = runtime.GOARCH
	e["GOCTL_VERSION"] = version.BuildVersion
	return e
}
