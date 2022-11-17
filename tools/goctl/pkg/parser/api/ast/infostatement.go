package ast

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type InfoStmt struct {
	Info   *TokenNode
	LParen *TokenNode
	Values []*KVExpr
	RParen *TokenNode
}

func (i *InfoStmt) Format(...string) string {
	if len(i.Values) == 0 {
		return ""
	}
	var textList []string
	for _, v := range i.Values {
		if v.Value.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format(Indent))
	}
	if len(textList) == 0 {
		return ""
	}

	w := NewBufferWriter()
	w.Write(WithNode(i.Info, i.LParen))
	w.NewLine()
	w.WriteText(strings.Join(textList,NewLine))
	w.NewLine()
	w.WriteText(i.RParen.Format())
	return w.String()
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.End()
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Pos()
}

func (i *InfoStmt) stmtNode() {}
