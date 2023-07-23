package parser

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
)

type filterBuilder struct {
	m             map[string]placeholder.Type
	checkExprName string
	errorManager  *errorManager
}

func (b *filterBuilder) check(nodes ...*ast.TokenNode) {
	for _, node := range nodes {
		if _, ok := b.m[node.Token.Text]; ok {
			b.errorManager.add(ast.DuplicateStmtError(node.Pos(), "duplicate "+b.checkExprName))
		} else {
			b.m[node.Token.Text] = placeholder.PlaceHolder
		}
	}
}

func (b *filterBuilder) error() error {
	return b.errorManager.error()
}

type filter struct {
	builders []*filterBuilder
}

func newFilter() *filter {
	return &filter{}
}

func (f *filter) addCheckItem(checkExprName string) *filterBuilder {
	b := &filterBuilder{
		m:             make(map[string]placeholder.Type),
		checkExprName: checkExprName,
		errorManager:  newErrorManager(),
	}
	f.builders = append(f.builders, b)
	return b
}

func (f *filter) error() error {
	if len(f.builders) == 0 {
		return nil
	}
	var errorManager = newErrorManager()
	for _, b := range f.builders {
		errorManager.add(b.error())
	}
	return errorManager.error()
}
