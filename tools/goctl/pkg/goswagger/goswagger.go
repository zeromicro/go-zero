package goswagger

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/goctl"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
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
	_, err := env.LookUpSwagger()
	return err == nil
}
