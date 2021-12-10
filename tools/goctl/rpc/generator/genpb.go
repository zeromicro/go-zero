package generator

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
)

const googleProtocGenGoErr = `--go_out: protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC`

// GenPb generates the pb.go file, which is a layer of packaging for protoc to generate gprc,
// but the commands and flags in protoc are not completely joined in goctl. At present, proto_path(-I) is introduced
func (g *DefaultGenerator) GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, _ *conf.Config, goOptions ...string) error {
	dir := ctx.GetPb()
	cw := new(bytes.Buffer)
	directory, base := filepath.Split(proto.Src)
	directory = filepath.Clean(directory)
	cw.WriteString("protoc ")
	protoImportPathSet := collection.NewSet()
	isSamePackage := true
	for _, ip := range protoImportPath {
		pip := " --proto_path=" + ip
		if protoImportPathSet.Contains(pip) {
			continue
		}

		abs, err := filepath.Abs(ip)
		if err != nil {
			return err
		}

		if abs == directory {
			isSamePackage = true
		} else {
			isSamePackage = false
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

	// Compatible with version 1.4.0，github.com/golang/protobuf/protoc-gen-go@v1.4.0
	// --go_opt usage please see https://developers.google.com/protocol-buffers/docs/reference/go-generated#package
	optSet := collection.NewSet()
	for _, op := range goOptions {
		opt := " --go_opt=" + op
		if optSet.Contains(opt) {
			continue
		}

		optSet.AddStr(op)
		cw.WriteString(" --go_opt=" + op)
	}

	var currentFileOpt string
	if !isSamePackage || (len(proto.GoPackage) > 0 && proto.GoPackage != proto.Package.Name) {
		if filepath.IsAbs(proto.GoPackage) {
			currentFileOpt = " --go_opt=M" + base + "=" + proto.GoPackage
		} else if strings.Contains(proto.GoPackage, string(filepath.Separator)) {
			currentFileOpt = " --go_opt=M" + base + "=./" + proto.GoPackage
		} else {
			currentFileOpt = " --go_opt=M" + base + "=../" + proto.GoPackage
		}
	} else {
		currentFileOpt = " --go_opt=M" + base + "=."
	}

	if !optSet.Contains(currentFileOpt) {
		cw.WriteString(currentFileOpt)
	}

	command := cw.String()
	g.log.Debug(command)
	_, err := execx.Run(command, "")
	if err != nil {
		if strings.Contains(err.Error(), googleProtocGenGoErr) {
			return errors.New(`unsupported plugin protoc-gen-go which installed from the following source:
google.golang.org/protobuf/cmd/protoc-gen-go, 
github.com/protocolbuffers/protobuf-go/cmd/protoc-gen-go;

Please replace it by the following command, we recommend to use version before v1.3.5:
go get -u github.com/golang/protobuf/protoc-gen-go`)
		}

		return err
	}
	return nil
}
