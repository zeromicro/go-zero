package generator

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
)

func (g *defaultGenerator) GenPb(ctx DirContext, IPATH string, dir Dir, proto parser.Proto) error {
	cw := new(bytes.Buffer)
	base := filepath.Dir(proto.Src)
	cw.WriteString("protoc ")
	cw.WriteString(" -I=" + base)
	if len(IPATH) > 0 {
		cw.WriteString(" -I=" + IPATH)
	}
	cw.WriteString(" " + proto.Src)
	if strings.Contains(proto.GoPackage, string(filepath.Separator)) {
		cw.WriteString(" --go_out=plugins=grpc:" + ctx.GetInternal().Filename)
	} else {
		cw.WriteString(" --go_out=plugins=grpc:" + dir.Filename)
	}
	command := cw.String()
	g.log.Debug(command)
	_, err := execx.Run(command, ctx.GetWorkDir().Filename)
	return err
}
