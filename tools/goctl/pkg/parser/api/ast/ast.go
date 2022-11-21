package ast

import (
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

type Node interface {
	Pos() token.Position
	End() token.Position
	Format(...string) string
	HasHeadCommentGroup() bool
	HasLeadingCommentGroup() bool
	CommentGroup() (head, leading CommentGroup)
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type AST struct {
	Filename     string
	Stmts        []Stmt
	readPosition int
}

type TokenNode struct {
	// HeadCommentGroup are the comments in prev lines.
	HeadCommentGroup CommentGroup
	Token            token.Token
	// LeadingCommentGroup are the tail comments in same line.
	LeadingCommentGroup CommentGroup
}

func (t *TokenNode) CommentGroup() (head, leading CommentGroup) {
	return t.HeadCommentGroup, t.LeadingCommentGroup
}

func NewTokenNode(tok token.Token) *TokenNode {
	return &TokenNode{Token: tok}
}

func (t *TokenNode) IsEmptyString() bool {
	return t.Equal("")
}

func (t *TokenNode) IsZeroString() bool {
	return t.Equal(`""`) || t.Equal("``")
}

func (t *TokenNode) Equal(s string) bool {
	return t.Token.Text == s
}

func (t *TokenNode) SetLeadingCommentGroup(cg CommentGroup) {
	t.LeadingCommentGroup = cg
}

func (t *TokenNode) HasLeadingCommentGroup() bool {
	return t.LeadingCommentGroup.Valid()
}

func (t *TokenNode) HasHeadCommentGroup() bool {
	return t.HeadCommentGroup.Valid()
}

func (t *TokenNode) PeekFirstLeadingComment() *CommentStmt {
	if t.HasLeadingCommentGroup() {
		return t.LeadingCommentGroup[0]
	}
	return nil
}

func (t *TokenNode) PeekFirstHeadComment() *CommentStmt {
	if t.HasHeadCommentGroup() {
		return t.HeadCommentGroup[0]
	}
	return nil
}

func (t *TokenNode) Format(prefix ...string) string {
	p := peekOne(prefix)
	var textList []string
	for _, v := range t.HeadCommentGroup {
		textList = append(textList, v.Format(p))
	}

	var tokenText = p + t.Token.Text
	var validLeadingCommentGroup CommentGroup
	for _, e := range t.LeadingCommentGroup {
		if util.IsEmptyStringOrWhiteSpace(e.Comment.Text) {
			continue
		}
		validLeadingCommentGroup = append(validLeadingCommentGroup, e)
	}

	if len(validLeadingCommentGroup) > 0 {
		tokenText = tokenText + WhiteSpace + t.LeadingCommentGroup.Join(WhiteSpace)
	}

	textList = append(textList, tokenText)
	return strings.Join(textList, NewLine)
}

func (t *TokenNode) Pos() token.Position {
	if t.HasHeadCommentGroup() {
		return t.PeekFirstHeadComment().Pos()
	}
	return t.Token.Position
}

func (t *TokenNode) End() token.Position {
	if t.HasLeadingCommentGroup() {
		return t.LeadingCommentGroup[len(t.LeadingCommentGroup)-1].End()
	}
	return t.Token.Position
}

func (a *AST) Format(w io.Writer) {
	fw := NewWriter(w)
	defer fw.Flush()
	for _, e := range a.Stmts {
		if e.Format() == NilIndent {
			continue
		}

		fw.Write(withNode(e))
		//switch stmt := e.(type) {
		//case *SyntaxStmt:
		//case *ImportGroupStmt:
		//	fw.Write(withNode(stmt))
		//case *ImportLiteralStmt:
		//	fw.Write(withNode(stmt))
		//case *InfoStmt:
		//case *ServiceStmt:
		//case *TypeGroupStmt:
		//case *TypeLiteralStmt:
		//case *RouteStmt:
		//case *CommentStmt:
		//}
	}
}

func (a *AST) Print() {
	_ = Print(a)
}

func SyntaxError(pos token.Position, format string, v ...interface{}) error {
	return fmt.Errorf("syntax error: %s %s", pos.String(), fmt.Sprintf(format, v...))
}

func DuplicateStmtError(pos token.Position, msg string) error {
	return fmt.Errorf("duplicate declaration: %s %s", pos.String(), msg)
}

func peekOne(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[0]
}
