package goswagger

import (
	"strings"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/goctl"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
)

const (
	Name = "swagger"
	url  = "github.com/go-swagger/go-swagger/cmd/swagger@latest"
)

func Install(cacheDir string) (string, error) {
	return goctl.Install(cacheDir, Name, func(dest string) (string, error) {
		err := golang.Install(url)
		return dest, err
	})
}

func Exists() bool {
	ver, err := Version()
	if err != nil {
		return false
	}
	return len(ver) > 0
}

// Version is used to get the version of the protoc-gen-go plugin. For older versions, protoc-gen-go does not support
// version fetching, so if protoc-gen-go --version is executed, it will cause the process to block, so it is controlled
// by a timer to prevent the older version process from blocking.
func Version() (string, error) {
	path, err := env.LookUpSwagger()
	if err != nil {
		return "", err
	}
	versionC := make(chan string)
	go func(c chan string) {
		version, _ := execx.Run(path+" --version", "")
		fields := strings.Fields(version)
		if len(fields) > 1 {
			c <- fields[1]
		}
	}(versionC)
	t := time.NewTimer(time.Second)
	select {
	case <-t.C:
		return "", nil
	case version := <-versionC:
		return version, nil
	}
}
