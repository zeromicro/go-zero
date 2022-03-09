package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

const googleProtocGenGoErr = `--go_out: protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC`

// GenPb generates the pb.go file, which is a layer of packaging for protoc to generate gprc,
// but the commands and flags in protoc are not completely joined in goctl. At present, proto_path(-I) is introduced
func (g *DefaultGenerator) GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, _ *conf.Config, c *ZRpcContext, goOptions ...string) error {
	if c != nil {
		return g.genPbDirect(ctx, c)
	}

	// deprecated: use genPbDirect instead.
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

func (g *DefaultGenerator) genPbDirect(ctx DirContext, c *ZRpcContext) error {
	g.log.Debug("[command]: %s", c.ProtocCmd)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = execx.Run(c.ProtocCmd, pwd)
	if err != nil {
		return err
	}
	return g.setPbDir(ctx, c)
}

func (g *DefaultGenerator) setPbDir(ctx DirContext, c *ZRpcContext) error {
	pbDir, err := findPbFile(c.GoOutput, false)
	if err != nil {
		return err
	}
	if len(pbDir) == 0 {
		return fmt.Errorf("pg.go is not found under %q", c.GoOutput)
	}
	grpcDir, err := findPbFile(c.GrpcOutput, true)
	if err != nil {
		return err
	}
	if len(grpcDir) == 0 {
		return fmt.Errorf("_grpc.pb.go is not found in %q", c.GrpcOutput)
	}
	if pbDir != grpcDir {
		return fmt.Errorf("the pb.go and _grpc.pb.go must under the same dir: "+
			"\n pb.go: %s\n_grpc.pb.go: %s", pbDir, grpcDir)
	}
	if pbDir == c.Output {
		return fmt.Errorf("the output of pb.go and _grpc.pb.go must not be the same "+
			"with --zrpc_out:\npb output: %s\nzrpc out: %s", pbDir, c.Output)
	}
	ctx.SetPbDir(pbDir, grpcDir)
	return nil
}

const (
	pbSuffix   = "pb.go"
	grpcSuffix = "_grpc.pb.go"
)

func findPbFile(current string, grpc bool) (string, error) {
	fileSystem := os.DirFS(current)
	var ret string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, pbSuffix) {
			if grpc {
				if strings.HasSuffix(path, grpcSuffix) {
					ret = path
					return os.ErrExist
				}
			} else if !strings.HasSuffix(path, grpcSuffix) {
				ret = path
				return os.ErrExist
			}
		}
		return nil
	})
	if err == os.ErrExist {
		return filepath.Dir(filepath.Join(current, ret)), nil
	}
	return "", err
}
