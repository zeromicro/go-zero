package ast

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type Node interface {
	Pos() token.Position
	End() token.Position
	Format(prefix ...string) string
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
	Token token.Token
}

func NewTokenNode(tok token.Token) *TokenNode {
	return &TokenNode{Token: tok}
}

func (t *TokenNode) Format(...string) string {
	return t.Token.Text
}

func (t *TokenNode) Pos() token.Position {
	return t.Token.Position
}

func (t *TokenNode) End() token.Position {
	return t.Token.Position
}

func (a *AST) Format(w *Writer) {
	defer func() {
		w.WriteTailGaps()
		w.Flush()
	}()
	for _, e := range a.Stmts {
		switch stmt := e.(type) {
		case *SyntaxStmt:
			stmt.fw = w
			stmt.Format(NilIndent)
		case *ImportGroupStmt:
		case *ImportLiteralStmt:
		case *InfoStmt:
		case *ServiceStmt:
		case *TypeGroupStmt:
		case *TypeLiteralStmt:
		case *RouteStmt:
		}
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
