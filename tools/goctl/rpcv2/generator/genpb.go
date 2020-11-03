package generator

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
)

func (g *defaultGenerator) GenPb(ctx DirContext, IPATH []string, dir Dir, proto parser.Proto) error {
	cw := new(bytes.Buffer)
	base := filepath.Dir(proto.Src)
	cw.WriteString("protoc ")
	for _, ip := range IPATH {
		cw.WriteString(" -I=" + ip)
	}
	cw.WriteString(" -I=" + base)
	cw.WriteString(" " + proto.Name)
	if strings.Contains(proto.GoPackage, string(filepath.Separator)) {
		cw.WriteString(" --go_out=plugins=grpc:" + ctx.GetInternal().Filename)
	} else {
		cw.WriteString(" --go_out=plugins=grpc:" + dir.Filename)
	}
	command := cw.String()
	g.log.Debug(command)
	_, err := execx.Run(command, "")
	return err
}
