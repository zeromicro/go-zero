package generator

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// GenPb generates the pb.go file, which is a layer of packaging for protoc to generate gprc,
// but the commands and flags in protoc are not completely joined in goctl. At present, proto_path(-I) is introduced
func (g *Generator) GenPb(ctx DirContext, c *ZRpcContext) error {
	return g.genPbDirect(ctx, c)
}

func (g *Generator) genPbDirect(ctx DirContext, c *ZRpcContext) error {
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

func (g *Generator) setPbDir(ctx DirContext, c *ZRpcContext) error {
	pbDir, err := findPbFile(c.GoOutput, c.Src, false)
	if err != nil {
		return err
	}
	if len(pbDir) == 0 {
		return fmt.Errorf("pg.go is not found under %q", c.GoOutput)
	}
	grpcDir, err := findPbFile(c.GrpcOutput, c.Src, true)
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

func findPbFile(current string, src string, grpc bool) (string, error) {
	protoName := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
	pbFile := protoName + "." + pbSuffix
	grpcFile := protoName + grpcSuffix

	fileSystem := os.DirFS(current)
	var ret string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, pbSuffix) {
			if grpc {
				if strings.HasSuffix(path, grpcFile) {
					ret = path
					return os.ErrExist
				}
			} else if strings.HasSuffix(path, pbFile) {
				ret = path
				return os.ErrExist
			}
		}
		return nil
	})
	if errors.Is(err, os.ErrExist) {
		return pathx.ReadLink(filepath.Dir(filepath.Join(current, ret)))
	}
	return "", err
}
