package ast

import (
	"fmt"
	"reflect"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type Formatter interface {
	// Format formats Node into string, do not end with '\n'.
	Format(prefix string) string
}

type Node interface {
	Formatter
	Pos() token.Position
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
	Filename string
	Stmts    []Stmt
}

func (a *AST) Format() string {
	w := NewWriter()
	for _, s := range a.Stmts {
		w.Writeln(s.Format(noLead))
		switch val := s.(type) {
		case *SyntaxStmt, *TypeGroupStmt, *ServiceStmt, *InfoStmt:
			w.NewLine()
		case *TypeLiteralStmt:
			if val.Expr.isStruct() {
				w.NewLine()
			}
		}
	}
	return w.String()
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

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	vo := reflect.ValueOf(v)
	if vo.Kind() == reflect.Ptr {
		return vo.IsNil()
	}
	return false
}
