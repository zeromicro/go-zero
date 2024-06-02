package ast

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zeromicro/antlr"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

type (
	// TokenStream defines a token
	TokenStream interface {
		GetStart() antlr.Token
		GetStop() antlr.Token
		GetParser() antlr.Parser
	}

	// ApiVisitor wraps api.BaseApiParserVisitor to call methods which has prefix Visit to
	// visit node from the api syntax
	ApiVisitor struct {
		*api.BaseApiParserVisitor
		debug    bool
		log      console.Console
		prefix   string
		infoFlag bool
	}

	// VisitorOption defines a function with argument ApiVisitor
	VisitorOption func(v *ApiVisitor)

	// Spec describes api spec
	Spec interface {
		Doc() []Expr
		Comment() Expr
		Format() error
		Equal(v any) bool
	}

	// Expr describes ast expression
	Expr interface {
		Prefix() string
		Line() int
		Column() int
		Text() string
		SetText(text string)
		Start() int
		Stop() int
		Equal(expr Expr) bool
		IsNotNil() bool
	}
)

// NewApiVisitor creates an instance for ApiVisitor
func NewApiVisitor(options ...VisitorOption) *ApiVisitor {
	v := &ApiVisitor{
		log: console.NewColorConsole(),
	}
	for _, opt := range options {
		opt(v)
	}
	return v
}

func (v *ApiVisitor) panic(expr Expr, msg string) {
	errString := fmt.Sprintf("%s line %d:%d  %s", v.prefix, expr.Line(), expr.Column(), msg)
	if v.debug {
		fmt.Println(errString)
	}

	panic(errString)
}

// WithVisitorPrefix returns a VisitorOption wrap with specified prefix
func WithVisitorPrefix(prefix string) VisitorOption {
	return func(v *ApiVisitor) {
		v.prefix = prefix
	}
}

// WithVisitorDebug returns a debug VisitorOption
func WithVisitorDebug() VisitorOption {
	return func(v *ApiVisitor) {
		v.debug = true
	}
}

type defaultExpr struct {
	prefix, v    string
	line, column int
	start, stop  int
}

// NewTextExpr creates a default instance for Expr
func NewTextExpr(v string) *defaultExpr {
	return &defaultExpr{
		v: v,
	}
}

func (v *ApiVisitor) newExprWithTerminalNode(node antlr.TerminalNode) *defaultExpr {
	if node == nil {
		return nil
	}

	token := node.GetSymbol()
	return v.newExprWithToken(token)
}

func (v *ApiVisitor) newExprWithToken(token antlr.Token) *defaultExpr {
	if token == nil {
		return nil
	}

	instance := &defaultExpr{}
	instance.prefix = v.prefix
	instance.v = token.GetText()
	instance.line = token.GetLine()
	instance.column = token.GetColumn()
	instance.start = token.GetStart()
	instance.stop = token.GetStop()

	return instance
}

func (v *ApiVisitor) newExprWithText(text string, line, column, start, stop int) *defaultExpr {
	instance := &defaultExpr{}
	instance.prefix = v.prefix
	instance.v = text
	instance.line = line
	instance.column = column
	instance.start = start
	instance.stop = stop

	return instance
}

func (e *defaultExpr) Prefix() string {
	if e == nil {
		return ""
	}

	return e.prefix
}

func (e *defaultExpr) Line() int {
	if e == nil {
		return 0
	}

	return e.line
}

func (e *defaultExpr) Column() int {
	if e == nil {
		return 0
	}

	return e.column
}

func (e *defaultExpr) Text() string {
	if e == nil {
		return ""
	}

	return e.v
}

func (e *defaultExpr) SetText(text string) {
	if e == nil {
		return
	}

	e.v = text
}

func (e *defaultExpr) Start() int {
	if e == nil {
		return 0
	}

	return e.start
}

func (e *defaultExpr) Stop() int {
	if e == nil {
		return 0
	}

	return e.stop
}

func (e *defaultExpr) Equal(expr Expr) bool {
	if e == nil {
		return expr == nil
	}

	if expr == nil {
		return false
	}

	return e.v == expr.Text()
}

func (e *defaultExpr) IsNotNil() bool {
	return e != nil
}

// EqualDoc compares whether the element literals in two Spec are equal
func EqualDoc(spec1, spec2 Spec) bool {
	if spec1 == nil {
		return spec2 == nil
	}

	if spec2 == nil {
		return false
	}

	var expectDoc, actualDoc []Expr
	expectDoc = append(expectDoc, spec2.Doc()...)
	actualDoc = append(actualDoc, spec1.Doc()...)
	sort.Slice(expectDoc, func(i, j int) bool {
		return expectDoc[i].Line() < expectDoc[j].Line()
	})

	for index, each := range actualDoc {
		if !each.Equal(actualDoc[index]) {
			return false
		}
	}

	if spec1.Comment() != nil {
		if spec2.Comment() == nil {
			return false
		}

		if !spec1.Comment().Equal(spec2.Comment()) {
			return false
		}
	} else {
		if spec2.Comment() != nil {
			return false
		}
	}

	return true
}

func (v *ApiVisitor) getDoc(t TokenStream) []Expr {
	return v.getHiddenTokensToLeft(t, api.COMMENTS, false)
}

func (v *ApiVisitor) getComment(t TokenStream) Expr {
	list := v.getHiddenTokensToRight(t, api.COMMENTS)
	if len(list) == 0 {
		return nil
	}

	commentExpr := list[0]
	stop := t.GetStop()
	text := stop.GetText()
	nlCount := strings.Count(text, "\n")

	if commentExpr.Line() != stop.GetLine()+nlCount {
		return nil
	}

	return commentExpr
}

func (v *ApiVisitor) getHiddenTokensToLeft(t TokenStream, channel int, containsCommentOfDefaultChannel bool) []Expr {
	ct := t.GetParser().GetTokenStream().(*antlr.CommonTokenStream)
	tokens := ct.GetHiddenTokensToLeft(t.GetStart().GetTokenIndex(), channel)

	var list []Expr
	for _, each := range tokens {
		if !containsCommentOfDefaultChannel {
			index := each.GetTokenIndex() - 1

			if index > 0 {
				allTokens := ct.GetAllTokens()
				flag := false
				for i := index; i >= 0; i-- {
					tk := allTokens[i]
					if tk.GetChannel() == antlr.LexerDefaultTokenChannel {
						if tk.GetLine() == each.GetLine() {
							flag = true
							break
						}
					}
				}

				if flag {
					continue
				}
			}
		}

		list = append(list, v.newExprWithToken(each))
	}

	return list
}

func (v *ApiVisitor) getHiddenTokensToRight(t TokenStream, channel int) []Expr {
	ct := t.GetParser().GetTokenStream().(*antlr.CommonTokenStream)
	tokens := ct.GetHiddenTokensToRight(t.GetStop().GetTokenIndex(), channel)
	var list []Expr
	for _, each := range tokens {
		list = append(list, v.newExprWithToken(each))
	}

	return list
}
