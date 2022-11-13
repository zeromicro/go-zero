package format

import (
	"io"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

type formatter struct {
	ast      *ast.AST
	curStmt  ast.Stmt
	peekStmt ast.Stmt

	fw Formatter
}

func newFormatter(source []byte, w io.Writer) (*formatter, error) {
	p := parser.New("", source, parser.SkipComment)
	ast := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return nil, err
	}

	return &formatter{
		ast: ast,
		fw:  NewWriter(w, p.TokenSet),
	}, nil
}

func (f *formatter) format() {
	defer func() {
		f.fw.Flush()
	}()
	for f.curStmt != nil {
		switch val := f.curStmt.(type) {
		case *ast.AtDocGroupStmt:
			f.fw.WriteSpaceInfix(indent, val.AtDoc, val.LParen)
			for _, v := range val.Values {
				f.fw.WriteBetween(indent, v.Key, v.Value)
			}
			f.fw.Write(indent, val.RParen)
		case *ast.AtDocLiteralStmt:
		case *ast.AtHandlerStmt:
		case *ast.AtServerStmt:
		case *ast.BodyStmt:
		case *ast.CommentStmt:
		case *ast.ImportGroupStmt:
			if len(val.Values) == 0 {
				f.fw.Skip(val.Import, val.RParen)
				break
			}
			f.fw.WriteSpaceInfixBetween(nilIndent, val.Import, val.LParen)
			preTok := val.LParen
			for _, v := range val.Values {
				if v.Line() == preTok.Line() {
					f.fw.NewLine()
				}
				f.fw.WriteBetween(indent, v, v)
				preTok = v
			}
			if val.RParen.Line() == preTok.Line() && len(val.Values) > 0 {
				f.fw.NewLine()
			}
			f.fw.Write(nilIndent, val.RParen)
			f.fw.NewLine()
		case *ast.ImportLiteralStmt:
			if val.Value.IsEmptyString() {
				f.fw.Skip(val.Import, val.Value)
				break
			}
			f.fw.WriteSpaceInfixBetween(nilIndent, val.Import, val.Value)
		case *ast.InfoStmt:
			if len(val.Values) == 0 {
				f.fw.Skip(val.Info, val.RParen)
				break
			}
			f.fw.WriteSpaceInfixBetween(nilIndent, val.Info, val.LParen)
			preTok := val.LParen
			for _, v := range val.Values {
				if v.Key.Line() == preTok.Line() {
					f.fw.NewLine()
				}
				f.fw.WriteBetween(indent, v.Key, v.Value)
				preTok = v.Value
			}
			if val.RParen.Line() == preTok.Line() && len(val.Values) > 0 {
				f.fw.NewLine()
			}
			f.fw.Write(nilIndent, val.RParen)
			f.fw.NewLine()
		case *ast.RouteStmt:
		case *ast.ServiceItemStmt:
		case *ast.ServiceStmt:
		case *ast.SyntaxStmt:
			f.fw.WriteSpaceInfixBetween(nilIndent, val.Syntax,val.Value)
			f.fw.NewLine()
		case *ast.TypeGroupStmt:
		case *ast.TypeLiteralStmt:
			switch dt := val.Expr.DataType.(type) {
			case *ast.AnyDataType:
				f.fw.WriteBetween(nilIndent, val.Type, dt.Any)
			case *ast.ArrayDataType:
				f.fw.Write(nilIndent, val.Type, dt.RBrack)
				if !dt.DataType.ContainsStruct() {

				}
			case *ast.BaseDataType:
				f.fw.WriteBetween(nilIndent, val.Type, dt.Base)
			case *ast.InterfaceDataType:
				f.fw.WriteBetween(nilIndent, val.Type, dt.Interface)
			case *ast.MapDataType:
			case *ast.PointerDataType:
			case *ast.SliceDataType:
			case *ast.StructDataType:
			}
		}
		f.nextStmt()
	}
}

func (f *formatter) formatDataType() {

}

func (f *formatter) nextStmt() {
	f.curStmt = f.peekStmt
	f.peekStmt = f.ast.NextStmt()
}

func Format(source []byte, w io.Writer) error {
	f, err := newFormatter(source, w)
	if err != nil {
		return err
	}
	f.nextStmt()
	f.nextStmt()
	f.format()
	return nil
}
