package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type SyntaxStmt struct {
	Syntax *TokenNode
	Assign *TokenNode
	Value  *TokenNode
}

func (s *SyntaxStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	syntaxNode := transferTokenNode(s.Syntax,
		withTokenNodePrefix(prefix...), ignoreLeadingComment())
	assignNode := transferTokenNode(s.Assign, ignoreLeadingComment())
	w.Write(withNode(syntaxNode, assignNode, s.Value), withPrefix(prefix...), expectSameLine())
	return w.String()
}

func (s *SyntaxStmt) End() token.Position {
	return s.Value.End()
}

func (s *SyntaxStmt) Pos() token.Position {
	return s.Syntax.Pos()
}

func (s *SyntaxStmt) stmtNode() {}
