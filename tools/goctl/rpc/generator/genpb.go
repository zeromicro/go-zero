package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
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

	goctlHomeBin := filepath.Join(goctlHome, "bin")
	err = util.MkdirIfNotExist(goctlHomeBin)
	if err != nil {
		return err
	}

	protocGenGo := filepath.Join(goctlHome, "bin", "protoc-gen-go")
	g.log.Debug("checking protoc-gen-go state ...")
	goGetCmd := "\ngo install github.com/golang/protobuf/protoc-gen-go@v1.3.2"

	if util.FileExists(protocGenGo) {
		g.log.Success("protoc-gen-go exists ...")
		goGetCmd = ""
	} else {
		g.log.Error("missing protoc-gen-go: downloading ...")
	}

	goos := runtime.GOOS
	switch goos {
	case vars.OsLinux, vars.OsMac:
		cmd = getUnixLikeCmd(goctlHome, goctlHomeBin, goGetCmd, cmd)
		g.log.Debug("%s", cmd)
	case vars.OsWindows:
		cmd = getWindowsCmd(goctlHome, goctlHomeBin, goGetCmd, cmd)
		// Do not support to execute commands in context, the solution is created
		// a batch file to execute it on Windows.
		batFile, err := createBatchFile(goctlHome, cmd)
		if err != nil {
			return err
		}

		g.log.Debug("%s", cmd)
		cmd = batFile
	default:
		return fmt.Errorf("unsupported os: %s", goos)
	}

	_, err = execx.Run(cmd, "")
	return err
}

func getUnixLikeCmd(goctlHome, goctlHomeBin, goGetCmd, cmd string) string {
	return fmt.Sprintf(`export GOPATH=%s 
export GOBIN=%s 
export PATH=$PATH:$GOPATH:$GOBIN
export GO111MODULE=on
export GOPROXY=https://goproxy.cn %s
%s`, goctlHome, goctlHomeBin, goGetCmd, cmd)
}

func getWindowsCmd(goctlHome, goctlHomeBin, goGetCmd, cmd string) string {
	return fmt.Sprintf(`set GOPATH=%s
set GOBIN=%s
set path=%s
set GO111MODULE=on
set GOPROXY=https://goproxy.cn %s
%s`, goctlHome, goctlHomeBin, "%path%;"+goctlHome+";"+goctlHomeBin, goGetCmd, cmd)
}

func createBatchFile(goctlHome, cmd string) (string, error) {
	batFile := filepath.Join(goctlHome, ".generate.bat")
	if !util.FileExists(batFile) {
		err := ioutil.WriteFile(batFile, []byte(cmd), os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return batFile, nil
}
