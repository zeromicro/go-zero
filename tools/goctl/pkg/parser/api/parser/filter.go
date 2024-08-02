package parser

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type filterBuilder struct {
	filename      string
	m             map[string]token.Position
	checkExprName string
	errorManager  *errorManager
}

func (b *filterBuilder) check(nodes ...*ast.TokenNode) {
	for _, node := range nodes {
		fileNodeText := fmt.Sprintf("%s/%s", b.filename, node.Token.Text)
		if pos, ok := b.m[fileNodeText]; ok && pos != node.Token.Position {
			b.errorManager.add(ast.DuplicateStmtError(node.Pos(), "duplicate "+b.checkExprName))
		} else {
			b.m[fileNodeText] = node.Token.Position
		}
	}
}

func (b *filterBuilder) checkNodeWithPrefix(prefix string, nodes ...*ast.TokenNode) {
	for _, node := range nodes {
		joinText := fmt.Sprintf("%s/%s", prefix, node.Token.Text)
		if pos, ok := b.m[joinText]; ok && pos != node.Token.Position {
			b.errorManager.add(ast.DuplicateStmtError(node.Pos(), "duplicate "+b.checkExprName))
		} else {
			b.m[joinText] = node.Token.Position
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

func (f *filter) addCheckItem(filename, checkExprName string) *filterBuilder {
	b := &filterBuilder{
		filename:      filename,
		m:             make(map[string]token.Position),
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
