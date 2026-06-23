package generator

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// GenPb generates the pb.go file, which is a layer of packaging for protoc to generate gprc,
// but the commands and flags in protoc are not completely joined in goctl. At present, proto_path(-I) is introduced
func (g *Generator) GenPb(ctx DirContext, c *ZRpcContext) error {
	return g.genPbDirect(ctx, c)
}

func (g *Generator) genPbDirect(ctx DirContext, c *ZRpcContext) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	protocCmd, err := g.buildProtocCmd(c, pwd)
	if err != nil {
		return err
	}

	g.log.Debug("[command]: %s", protocCmd)
	_, err = execx.Run(protocCmd, pwd)
	if err != nil {
		return err
	}
	return g.setPbDir(ctx, c)
}

// buildProtocCmd resolves all transitively imported proto files and appends
// them to the protoc command so that their pb.go files are also generated.
func (g *Generator) buildProtocCmd(c *ZRpcContext, pwd string) (string, error) {
	// Build the full list of proto search paths (absolute).
	protoPaths := make([]string, 0, len(c.ProtoPaths)+1)

	// Always include the directory of the source proto so that imports
	// relative to the source file can be found.
	srcDir := filepath.Dir(c.Src)
	if !filepath.IsAbs(srcDir) {
		srcDir = filepath.Join(pwd, srcDir)
	}
	protoPaths = append(protoPaths, srcDir)

	for _, p := range c.ProtoPaths {
		if !filepath.IsAbs(p) {
			p = filepath.Join(pwd, p)
		}
		protoPaths = append(protoPaths, p)
	}

	importedFiles, err := parser.ResolveImports(c.Src, protoPaths)
	if err != nil {
		return "", err
	}
	if len(importedFiles) == 0 {
		return c.ProtocCmd, nil
	}

	cmd := c.ProtocCmd
	for _, f := range importedFiles {
		// Use the path relative to the best-matching --proto_path entry so that
		// protoc's source_relative output lands in the correct directory.
		// e.g. if --proto_path=./ext and the file is ext/common/types.proto,
		// we pass "common/types.proto" rather than "ext/common/types.proto".
		rel := relativeToProtoPath(f, protoPaths, pwd)
		cmd += " " + rel
	}
	return cmd, nil
}

// relativeToProtoPath returns the path of f relative to the most specific
// (longest) proto_path entry that is a parent of f. Falls back to relative
// to pwd when no proto_path matches.
func relativeToProtoPath(f string, protoPaths []string, pwd string) string {
	bestRel := ""
	bestLen := 0
	for _, pp := range protoPaths {
		prefix := pp + string(filepath.Separator)
		if strings.HasPrefix(f, prefix) && len(pp) > bestLen {
			if rel, err := filepath.Rel(pp, f); err == nil {
				bestRel = rel
				bestLen = len(pp)
			}
		}
	}
	if bestRel != "" {
		return bestRel
	}
	if rel, err := filepath.Rel(pwd, f); err == nil {
		return rel
	}
	return f
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
