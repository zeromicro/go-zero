package ast

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/g4gen"
)

type (
	Parser struct {
		content string
		lexer   *parser.ApiLexer
		*parser.ApiParser
	}

	option func(p *Parser)
)

func NewParser(content string, options ...option) *Parser {
	inputStream := antlr.NewInputStream(content)
	lexer := parser.NewApiLexer(inputStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	p := parser.NewApiParser(tokens)
	instance := &Parser{
		content:   content,
		lexer:     lexer,
		ApiParser: p,
	}
	instance.AddErrorCallback(nil)
	for _, opt := range options {
		opt(instance)
	}
	return instance
}

func (p *Parser) AddErrorCallback(callback ErrCallback) {
	p.RemoveErrorListeners()
	errListener := NewErrorListener(callback)
	p.AddErrorListener(errListener)
}

func WithErrorCallback(callback ErrCallback) option {
	return func(p *Parser) {
		p.AddErrorCallback(callback)
	}
}
