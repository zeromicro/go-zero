package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const configTemplate = `package config

import "github.com/tal-tech/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
}
`

func (g *defaultGenerator) GenConfig(ctx DirContext, _ parser.Proto) error {
	dir := ctx.GetConfig()
	fileName := filepath.Join(dir.Filename, formatFilename("config")+".go")
	if util.FileExists(fileName) {
		return nil
	}

	text, err := util.LoadTemplate(category, configTemplateFileFile, configTemplate)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(text), os.ModePerm)
}
