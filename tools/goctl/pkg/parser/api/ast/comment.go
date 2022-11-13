package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type CommentStmt struct {
	Comment token.Token
}

func (c *CommentStmt) Pos() token.Position {
	return c.Comment.Position
}

func (c *CommentStmt) stmtNode() {}
