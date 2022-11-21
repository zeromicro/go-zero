package ast

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

type CommentGroup []*CommentStmt

func (cg CommentGroup) Join(sep string) string {
	if !cg.Valid() {
		return ""
	}
	var list = make([]string, 0, len(cg))
	for _, v := range cg {
		comment := v.Format(NilIndent)
		if util.IsEmptyStringOrWhiteSpace(comment) {
			continue
		}
		list = append(list, comment)
	}
	return strings.Join(list, sep)
}

func (cg CommentGroup) Valid() bool {
	return len(cg) > 0
}

type CommentStmt struct {
	Comment token.Token
}

func (c *CommentStmt) stmtNode() {}

func (c *CommentStmt) Pos() token.Position {
	return c.Comment.Position
}

func (c *CommentStmt) End() token.Position {
	return c.Comment.Position
}

func (c *CommentStmt) Format(prefix ...string) string {
	return peekOne(prefix) + c.Comment.Text
}
