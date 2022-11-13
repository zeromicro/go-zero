package ast

import (
	"fmt"
	"reflect"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type Node interface {
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
	Filename     string
	Stmts        []Stmt
	readPosition int
}

func (a *AST) NextStmt() Stmt {
	if a.readPosition == len(a.Stmts) {
		return nil
	}
	defer func() {
		a.readPosition += 1
	}()
	return a.Stmts[a.readPosition]
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
