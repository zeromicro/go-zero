package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type InfoStmt struct {
	Info   token.Token
	LParen token.Token
	Values []*KVExpr
	RParen token.Token
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Position
}

func (i *InfoStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, i.Info, i.LParen)
	if len(i.Values) > 0 {
		w.NewLine()
	}
	tw := w.UseTabWriter()
	for _, v := range i.Values {
		tw.WriteWithInfixIndentln(prefix+indent, v.Key.Text+v.Colon.Text, v.Value.Text)
	}
	tw.Flush()
	if len(i.Values) > 0 {
		w.Write(prefix, i.RParen)
	} else {
		w.Write(i.RParen)
	}
	return w.String()
}

func (i *InfoStmt) stmtNode() {}
