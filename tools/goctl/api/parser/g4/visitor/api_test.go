package visitor

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/parser"
)

func TestApiParser_Api(t *testing.T) {
	inputStream := antlr.NewInputStream(`syntax = "v2`)
	lexer := parser.NewApiLexer(inputStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	p := parser.NewApiParser(tokens)
	p.RemoveErrorListeners()
	errorListener := NewErrorListener()
	p.AddErrorListener(errorListener)
	v := NewApiVisitor()
	p.Api().Accept(v)
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{}
}

func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	lineHeader := "line " + strconv.Itoa(line) + ":" + strconv.Itoa(column)
	fmt.Println(lineHeader, msg)
}
