package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const configTemplate = `package config

import "github.com/tal-tech/go-zero/rpcx"

type Config struct {
	rpcx.RpcServerConf
}
`

func (g *defaultRpcGenerator) genConfig() error {
	configPath := g.dirM[dirConfig]
	fileName := filepath.Join(configPath, fileConfig)
	if util.FileExists(fileName) {
		return nil
	}
	return ioutil.WriteFile(fileName, []byte(configTemplate), os.ModePerm)
}
