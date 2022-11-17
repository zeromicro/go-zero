package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type SyntaxStmt struct {
	Syntax *TokenNode
	Assign *TokenNode
	Value  *TokenNode
}

func (s *SyntaxStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(WithNode(s.Syntax, s.Assign, s.Value), WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
	return w.String()
}

func (s *SyntaxStmt) End() token.Position {
	return s.Value.End()
}

func (s *SyntaxStmt) Pos() token.Position {
	return s.Syntax.Pos()
}

func (s *SyntaxStmt) stmtNode() {}
