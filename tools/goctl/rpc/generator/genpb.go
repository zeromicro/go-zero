package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

// GenPb generates the pb.go file, which is a layer of packaging for protoc to generate gprc,
// but the commands and flags in protoc are not completely joined in goctl. At present, proto_path(-I) is introduced
func (g *DefaultGenerator) GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, _ *conf.Config, goOptions ...string) error {
	dir := ctx.GetPb()
	cw := new(bytes.Buffer)
	directory, _ := filepath.Split(proto.Src)
	directory = filepath.Clean(directory)
	cw.WriteString("protoc ")
	protoImportPathSet := collection.NewSet()
	for _, ip := range protoImportPath {
		pip := " --proto_path=" + ip
		if protoImportPathSet.Contains(pip) {
			continue
		}

		protoImportPathSet.AddStr(pip)
		cw.WriteString(pip)
	}
	currentPath := " --proto_path=" + directory
	if !protoImportPathSet.Contains(currentPath) {
		cw.WriteString(currentPath)
	}

	cw.WriteString(" " + proto.Name)
	if strings.Contains(proto.GoPackage, "/") {
		cw.WriteString(" --go_out=plugins=grpc:" + ctx.GetMain().Filename)
	} else {
		cw.WriteString(" --go_out=plugins=grpc:" + dir.Filename)
	}

	return g.generatePbWithVersion132(cw.String())
}

// generatePbWithVersion132 generates pb.go by specifying protoc-gen-go@1.3.2 version
func (g *DefaultGenerator) generatePbWithVersion132(cmd string) error {
	goctlHome, err := util.GetGoctlHome()
	if err != nil {
		return err
	}

	err = util.MkdirIfNotExist(goctlHome)
	if err != nil {
		return err
	}

	protocGenGo := filepath.Join(goctlHome, "bin", "protoc-gen-go")
	goGetCmd := "\ngo get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2"
	if util.FileExists(protocGenGo) {
		goGetCmd = ""
	}
	goos := runtime.GOOS
	switch goos {
	case vars.OsLinux, vars.OsMac:
		cmd = fmt.Sprintf(`export GOPATH=%s %s
%s`, goctlHome, goGetCmd, cmd)
	case vars.OsWindows:
		cmd = fmt.Sprintf(`set GOPATH=%s %s
%s`, goctlHome, goGetCmd, cmd)
	default:
		return fmt.Errorf("unsupported os: %s", goos)
	}

	g.log.Debug(cmd)
	_, err = execx.Run(cmd, "")
	return err
}
