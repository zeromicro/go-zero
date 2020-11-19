package generator

import (
	"bytes"
	"path/filepath"
	"strings"

	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
)

func (g *defaultGenerator) GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, _ *conf.Config) error {
	dir := ctx.GetPb()
	cw := new(bytes.Buffer)
	base := filepath.Dir(proto.Src)
	cw.WriteString("protoc ")
	for _, ip := range protoImportPath {
		cw.WriteString(" -I=" + ip)
	}
	cw.WriteString(" -I=" + base)
	cw.WriteString(" " + proto.Name)
	if strings.Contains(proto.GoPackage, "/") {
		cw.WriteString(" --go_out=plugins=grpc:" + ctx.GetMain().Filename)
	} else {
		cw.WriteString(" --go_out=plugins=grpc:" + dir.Filename)
	}
	command := cw.String()
	g.log.Debug(command)
	_, err := execx.Run(command, "")
	return err
}
