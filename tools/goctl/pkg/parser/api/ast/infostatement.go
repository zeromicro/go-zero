package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type InfoStmt struct {
	Info   *TokenNode
	LParen *TokenNode
	Values []*KVExpr
	RParen *TokenNode

	fw *Writer
}

func (i *InfoStmt) Format(prefix ...string) (result string) {
	if i.fw == nil {
		return
	}

	p := peekOne(prefix)
	i.fw.Skip(i)
	if len(i.Values) == 0 {
		i.fw.Skip(i.Info, i.LParen, i.RParen)
		return
	}

	i.fw.WriteInOneLine(p, i.Info, i.LParen)
	var line = i.fw.lastWriteNode.End().Line
	for _, kv := range i.Values {
		i.fw.Skip(kv)
		if kv.Pos().Line == line {
			i.fw.NewLine()
		}
		i.fw.WriteBetween(Indent, kv.Key, kv.Value)
		line = kv.Value.Pos().Line
	}

	if i.RParen.Pos().Line == line {
		i.fw.NewLine()
	}
	i.fw.Write(p, i.RParen)
	i.fw.NewLine()

	return
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.End()
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Pos()
}

func (i *InfoStmt) stmtNode() {}
