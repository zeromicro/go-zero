package ast

import (
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

// Node represents a node in the AST.
type Node interface {
	// Pos returns the position of the first character belonging to the node.
	Pos() token.Position
	// End returns the position of the first character immediately after the node.
	End() token.Position
	// Format returns the node's text after format.
	Format(...string) string
	// HasHeadCommentGroup returns true if the node has head comment group.
	HasHeadCommentGroup() bool
	// HasLeadingCommentGroup returns true if the node has leading comment group.
	HasLeadingCommentGroup() bool
	// CommentGroup returns the node's head comment group and leading comment group.
	CommentGroup() (head, leading CommentGroup)
}

// Stmt represents a statement in the AST.
type Stmt interface {
	Node
	stmtNode()
}

// Expr represents an expression in the AST.
type Expr interface {
	Node
	exprNode()
}

// AST represents a parsed API file.
type AST struct {
	Filename     string
	Stmts        []Stmt
	readPosition int
}

// TokenNode represents a token node in the AST.
type TokenNode struct {
	// HeadCommentGroup are the comments in prev lines.
	HeadCommentGroup CommentGroup
	// Token is the token of the node.
	Token token.Token
	// LeadingCommentGroup are the tail comments in same line.
	LeadingCommentGroup CommentGroup

	// headFlag and leadingFlag is a comment flag only used in transfer another Node to TokenNode,
	// headFlag's value is true do not represent HeadCommentGroup is not empty,
	// leadingFlag's values is true do not represent LeadingCommentGroup is not empty.
	headFlag, leadingFlag bool
}

// NewTokenNode creates and returns a new TokenNode.
func NewTokenNode(tok token.Token) *TokenNode {
	return &TokenNode{Token: tok}
}

// IsEmptyString returns true if the node is empty string.
func (t *TokenNode) IsEmptyString() bool {
	return t.Equal("")
}

// IsZeroString returns true if the node is zero string.
func (t *TokenNode) IsZeroString() bool {
	return t.Equal(`""`) || t.Equal("``")
}

// Equal returns true if the node's text is equal to the given text.
func (t *TokenNode) Equal(s string) bool {
	return t.Token.Text == s
}

// SetLeadingCommentGroup sets the node's leading comment group.
func (t *TokenNode) SetLeadingCommentGroup(cg CommentGroup) {
	t.LeadingCommentGroup = cg
}

// RawText returns the node's raw text.
func (t *TokenNode) RawText() string {
	text := t.Token.Text
	if strings.HasPrefix(text, "`") {
		text = strings.TrimPrefix(text, "`")
		text = strings.TrimSuffix(text, "`")
	} else if strings.HasPrefix(text, `"`) {
		text = strings.TrimPrefix(text, `"`)
		text = strings.TrimSuffix(text, `"`)
	}

	return text
}

func (t *TokenNode) HasLeadingCommentGroup() bool {
	return t.LeadingCommentGroup.Valid() || t.leadingFlag
}

func (t *TokenNode) HasHeadCommentGroup() bool {
	return t.HeadCommentGroup.Valid() || t.headFlag
}

func (t *TokenNode) CommentGroup() (head, leading CommentGroup) {
	return t.HeadCommentGroup, t.LeadingCommentGroup
}

// PeekFirstLeadingComment returns the first leading comment of the node.
func (t *TokenNode) PeekFirstLeadingComment() *CommentStmt {
	if len(t.LeadingCommentGroup) > 0 {
		return t.LeadingCommentGroup[0]
	}
	return nil
}

// PeekFirstHeadComment returns the first head comment of the node.
func (t *TokenNode) PeekFirstHeadComment() *CommentStmt {
	if len(t.HeadCommentGroup) > 0 {
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
	if len(t.HeadCommentGroup) > 0 {
		return t.PeekFirstHeadComment().Pos()
	}
	return t.Token.Position
}

func (t *TokenNode) End() token.Position {
	if len(t.LeadingCommentGroup) > 0 {
		return t.LeadingCommentGroup[len(t.LeadingCommentGroup)-1].End()
	}
	return t.Token.Position
}

// Format formats the AST.
func (a *AST) Format(w io.Writer) {
	fw := NewWriter(w)
	defer fw.Flush()
	for idx, e := range a.Stmts {
		if e.Format() == NilIndent {
			continue
		}

		fw.Write(withNode(e))
		fw.NewLine()
		switch e.(type) {
		case *SyntaxStmt:
			fw.NewLine()
		case *ImportGroupStmt:
			fw.NewLine()
		case *ImportLiteralStmt:
			if idx < len(a.Stmts)-1 {
				_, ok := a.Stmts[idx+1].(*ImportLiteralStmt)
				if !ok {
					fw.NewLine()
				}
			}
		case *InfoStmt:
			fw.NewLine()
		case *ServiceStmt:
			fw.NewLine()
		case *TypeGroupStmt:
			fw.NewLine()
		case *TypeLiteralStmt:
			fw.NewLine()
		case *CommentStmt:
		}
	}
}

// FormatForUnitTest formats the AST for unit test.
func (a *AST) FormatForUnitTest(w io.Writer) {
	fw := NewWriter(w)
	defer fw.Flush()
	for _, e := range a.Stmts {
		text := e.Format()
		if text == NilIndent {
			continue
		}

		fw.WriteText(text)
	}
}

// Print prints the AST.
func (a *AST) Print() {
	_ = Print(a)
}

// SyntaxError represents a syntax error.
func SyntaxError(pos token.Position, format string, v ...interface{}) error {
	return fmt.Errorf("syntax error: %s %s", pos.String(), fmt.Sprintf(format, v...))
}

// DuplicateStmtError represents a duplicate statement error.
func DuplicateStmtError(pos token.Position, msg string) error {
	return fmt.Errorf("duplicate declaration: %s %s", pos.String(), msg)
}

func peekOne(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[0]
}
