package visitor

import "github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/parser"

type ApiVisitor struct {
	parser.BaseApiParserVisitor
}

func NewApiVisitor() *ApiVisitor {
	return &ApiVisitor{}
}
