package format

import (
	"io"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func Format(source []byte, w io.Writer) error {
	nodeSet := ast.NewNodeSet()
	p := parser.New(nodeSet, "", source)
	result := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return err
	}

	fw := ast.NewWriter(w, nodeSet)
	result.Format(fw)
	return nil
}
