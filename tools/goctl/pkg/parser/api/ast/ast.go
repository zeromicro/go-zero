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
			w.Skip(stmt)
			w.WriteInOneLineBetween(NilIndent, stmt.Syntax, stmt.Value)
			w.NewLine()
		case *ImportGroupStmt:
			w.Skip(stmt)
			var values []string
			for _, v := range stmt.Values {
				if v.Token.IsEmptyString() {
					continue
				}
				values = append(values, v.Token.Text)
			}
			if len(values) == 0 {
				w.Skip(stmt.Import, stmt.LParen, stmt.RParen)
				continue
			}
			w.WriteInOneLineBetween(NilIndent, stmt.Import, stmt.LParen)
			var line = w.lastWriteNode.End().Line
			for _, v := range stmt.Values {
				if v.Pos().Line == line {
					w.NewLine()
				}
				w.Write(Indent, v)
				line = v.End().Line
			}
			if stmt.RParen.Pos().Line == line {
				w.NewLine()
			}
			w.Write(NilIndent, stmt.RParen)
			w.NewLine()
		case *ImportLiteralStmt:
			w.Skip(stmt)
			if stmt.Value.Token.IsEmptyString() {
				w.Skip(stmt.Import, stmt.Value)
				return
			}

			w.WriteInOneLineBetween(NilIndent, stmt.Import, stmt.Value)
		case *InfoStmt:
			w.Skip(stmt)
			if len(stmt.Values) == 0 {
				w.Skip(stmt.Info, stmt.LParen, stmt.RParen)
				return
			}

			w.WriteInOneLine(NilIndent, stmt.Info, stmt.LParen)
			var line = w.lastWriteNode.End().Line
			for _, kv := range stmt.Values {
				w.Skip(kv)
				if kv.Pos().Line == line {
					w.NewLine()
				}
				w.WriteBetween(Indent, kv.Key, kv.Value)
				line = kv.Value.End().Line
			}

			if stmt.RParen.Pos().Line == line {
				w.NewLine()
			}
			w.Write(NilIndent, stmt.RParen)
			w.NewLine()
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
