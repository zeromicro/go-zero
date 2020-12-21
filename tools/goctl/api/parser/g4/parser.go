package g4

import (
	"os"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
)

func ParseApi(src string) (bool, error) {
	input, _ := antlr.NewFileStream(os.Args[1])
	lexer := parser.NewApiLexer(input)
	stream := antlr.NewTokenStream(lexer, 0)
	p := parser.NewApiParser(stream)
}
