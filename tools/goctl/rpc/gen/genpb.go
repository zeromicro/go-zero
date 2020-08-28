package gen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dsymonds/gotoc/parser"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	astParser "github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func (g *defaultRpcGenerator) genPb() error {
	importPath, filename := filepath.Split(g.Ctx.ProtoFileSrc)
	tree, err := parser.ParseFiles([]string{filename}, []string{importPath})
	if err != nil {
		return err
	}

	if len(tree.Files) == 0 {
		return errors.New("proto ast parse failed")
	}

	file := tree.Files[0]
	if len(file.Package) == 0 {
		return errors.New("expected package, but nothing found")
	}

	targetStruct := make(map[string]lang.PlaceholderType)
	for _, item := range file.Messages {
		if len(item.Messages) > 0 {
			return fmt.Errorf(`line %v: unexpected inner message near: "%v""`, item.Messages[0].Position.Line, item.Messages[0].Name)
		}

		name := stringx.From(item.Name)
		if _, ok := targetStruct[name.Lower()]; ok {
			return fmt.Errorf("line %v: duplicate %v", item.Position.Line, name)
		}
		targetStruct[name.Lower()] = lang.Placeholder
	}

	pbPath := g.dirM[dirPb]
	protoFileName := filepath.Base(g.Ctx.ProtoFileSrc)
	err = g.protocGenGo(pbPath)
	if err != nil {
		return err
	}

	pbGo := strings.TrimSuffix(protoFileName, ".proto") + ".pb.go"
	pbFile := filepath.Join(pbPath, pbGo)
	bts, err := ioutil.ReadFile(pbFile)
	if err != nil {
		return err
	}

	aspParser := astParser.NewAstParser(bts, targetStruct, g.Ctx.Console)
	ast, err := aspParser.Parse()
	if err != nil {
		return err
	}

	if len(ast.Service) == 0 {
		return fmt.Errorf("service not found")
	}
	g.ast = ast
	return nil
}

func (g *defaultRpcGenerator) protocGenGo(target string) error {
	src := filepath.Dir(g.Ctx.ProtoFileSrc)
	sh := fmt.Sprintf(`protoc -I=%s --go_out=plugins=grpc:%s %s`, src, target, g.Ctx.ProtoFileSrc)
	stdout, err := execx.Run(sh)
	if err != nil {
		return err
	}

	if len(stdout) > 0 {
		g.Ctx.Info(stdout)
	}

	return nil
}
