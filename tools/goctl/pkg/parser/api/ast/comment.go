package ast

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

// CommentGroup represents a list of comments.
type CommentGroup []*CommentStmt

// List returns the list of comments.
func (cg CommentGroup) List() []string {
	var list = make([]string, 0, len(cg))
	for _, v := range cg {
		comment := v.Comment.Text
		if util.IsEmptyStringOrWhiteSpace(comment) {
			continue
		}
		list = append(list, comment)
	}
	return list
}

// String joins and returns the comment text.
func (cg CommentGroup) String() string {
	return cg.Join(" ")
}

// Join joins the comments with the given separator.
func (cg CommentGroup) Join(sep string) string {
	if !cg.Valid() {
		return ""
	}
	list := cg.List()
	return strings.Join(list, sep)
}

// Valid returns true if the comment is valid.
func (cg CommentGroup) Valid() bool {
	return len(cg) > 0
}

// CommentStmt represents a comment statement.
type CommentStmt struct {
	// Comment is the comment token.
	Comment token.Token
}

func (c *CommentStmt) HasHeadCommentGroup() bool {
	return false
}

func (c *CommentStmt) HasLeadingCommentGroup() bool {
	return false
}

func (c *CommentStmt) CommentGroup() (head, leading CommentGroup) {
	return
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
