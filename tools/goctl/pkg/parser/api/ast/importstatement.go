package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type ImportStmt interface {
	Stmt
	importNode()
}

type ImportLiteralStmt struct {
	Import *TokenNode
	Value  *TokenNode

	fw *Writer
}

func (i *ImportLiteralStmt) Format(prefix ...string) (result string) {
	if i.fw == nil {
		return
	}

	i.fw.Skip(i)
	if i.Value.Token.IsEmptyString() {
		i.fw.Skip(i.Import, i.Value)
		return
	}

	i.fw.WriteInOneLine(peekOne(prefix), i.Import, i.Value)
	return
}

func (i *ImportLiteralStmt) End() token.Position {
	return i.Value.End()
}

func (i *ImportLiteralStmt) importNode() {}

func (i *ImportLiteralStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportLiteralStmt) stmtNode() {}

type ImportGroupStmt struct {
	Import *TokenNode
	LParen *TokenNode
	Values []*TokenNode
	RParen *TokenNode

	fw *Writer
}

func (i *ImportGroupStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (i *ImportGroupStmt) End() token.Position {
	return i.RParen.End()
}

func (i *ImportGroupStmt) importNode() {}

func (i *ImportGroupStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportGroupStmt) stmtNode() {}
