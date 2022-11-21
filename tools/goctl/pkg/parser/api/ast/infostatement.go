package ast

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type InfoStmt struct {
	Info   *TokenNode
	LParen *TokenNode
	Values []*KVExpr
	RParen *TokenNode
}

func (i *InfoStmt) Format(prefix ...string) string {
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
	infoNode := transferTokenNode(i.Info, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	w.Write(withNode(infoNode, i.LParen))
	w.NewLine()
	for _, v := range i.Values {
		node := transferTokenNode(v.Key, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment())
		w.Write(withNode(node, v.Value), expectIndentInfix(), expectSameLine())
		w.NewLine()
	}
	w.Write(withNode(transferTokenNode(i.RParen, withTokenNodePrefix(prefix...))))
	return w.String()
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.End()
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Pos()
}

func (i *InfoStmt) stmtNode() {}
