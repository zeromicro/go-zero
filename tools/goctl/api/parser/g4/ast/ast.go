package ast

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

type (
	TokenStream interface {
		GetStart() antlr.Token
		GetStop() antlr.Token
		GetParser() antlr.Parser
	}
	ApiVisitor struct {
		api.BaseApiParserVisitor
		debug    bool
		log      console.Console
		prefix   string
		infoFlag bool
	}

	VisitorOption func(v *ApiVisitor)

	Spec interface {
		Doc() []Expr
		Comment() Expr
		Format() error
		Equal(v interface{}) bool
	}

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

func WithVisitorPrefix(prefix string) VisitorOption {
	return func(v *ApiVisitor) {
		v.prefix = prefix
	}
}

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
		if expr != nil {
			return false
		}

		return true
	}

	if expr == nil {
		return false
	}

	return e.v == expr.Text()
}

func (e *defaultExpr) IsNotNil() bool {
	return e != nil
}

func EqualDoc(spec1, spec2 Spec) bool {
	if spec1 == nil {
		if spec2 != nil {
			return false
		}
		return true
	} else {
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
	}
	return true
}

func (v *ApiVisitor) getDoc(t TokenStream) []Expr {
	list := v.getHiddenTokensToLeft(t, api.COMEMNTS, false)
	return list
}

func (v *ApiVisitor) getComment(t TokenStream) Expr {
	list := v.getHiddenTokensToRight(t, api.COMEMNTS)
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
	var tmp []antlr.Token
	for _, each := range tokens {
		tmp = append(tmp, each)
	}

	var list []Expr
	for _, each := range tmp {
		if !containsCommentOfDefaultChannel {
			index := each.GetTokenIndex() - 1

			if index > 0 {
				allTokens := ct.GetAllTokens()
				var flag = false
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

func (v *ApiVisitor) exportCheck(expr Expr) {
	if expr == nil || !expr.IsNotNil() {
		return
	}
	if api.IsBasicType(expr.Text()) {
		return
	}

	if util.UnExport(expr.Text()) {
		v.log.Warning("%s line %d:%d unexported declaration '%s', use %s instead", expr.Prefix(), expr.Line(),
			expr.Column(), expr.Text(), strings.Title(expr.Text()))
	}
}
