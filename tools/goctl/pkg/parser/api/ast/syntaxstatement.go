package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type SyntaxStmt struct {
	Syntax *TokenNode
	Assign *TokenNode
	Value  *TokenNode

	fw *Writer
}

func (s *SyntaxStmt) Format(prefix ...string) string {
	if s.fw == nil {
		return ""
	}

	s.fw.Skip(s)
	s.fw.WriteSpaceInfixBetween(peekOne(prefix), s.Syntax, s.Value)
	s.fw.NewLine()
	return ""
}

func (s *SyntaxStmt) End() token.Position {
	return s.Value.End()
}

func (s *SyntaxStmt) Pos() token.Position {
	return s.Syntax.Pos()
}

func (s *SyntaxStmt) stmtNode() {}
