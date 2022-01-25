package bug

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
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

	return strings.TrimSuffix(w.String(), "\n")
}

func getEnv() env {
	e := make(env)
	e[os] = runtime.GOOS
	e[arch] = runtime.GOARCH
	e[goctlVersion] = version.BuildVersion
	e[goVersion] = runtime.Version()
	return e
}
