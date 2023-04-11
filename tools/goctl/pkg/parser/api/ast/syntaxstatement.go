package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

// SyntaxStmt represents a syntax statement.
type SyntaxStmt struct {
	// Syntax is the syntax token.
	Syntax *TokenNode
	// Assign is the assign token.
	Assign *TokenNode
	// Value is the syntax value.
	Value *TokenNode
}

func (s *SyntaxStmt) HasHeadCommentGroup() bool {
	return s.Syntax.HasHeadCommentGroup()
}

func (s *SyntaxStmt) HasLeadingCommentGroup() bool {
	return s.Value.HasLeadingCommentGroup()
}

func (s *SyntaxStmt) CommentGroup() (head, leading CommentGroup) {
	return s.Syntax.HeadCommentGroup, s.Syntax.LeadingCommentGroup
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
