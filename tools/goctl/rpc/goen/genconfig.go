package gogen

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var configTemplate = `package config

import "github.com/tal-tech/go-zero/rpcx"

type (
	Config struct {
		rpcx.RpcServerConf
	}
)
`

func (g *defaultRpcGenerator) genConfig() error {
	configPath := g.dirM[dirConfig]
	fileName := filepath.Join(configPath, fileConfig)
	return ioutil.WriteFile(fileName, []byte(configTemplate), os.ModePerm)
}
