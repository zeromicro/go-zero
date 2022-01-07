package ast

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/antlr"
)

type (
	// Parser provides api parsing capabilities
	Parser struct {
		linePrefix string
		debug      bool
		log        console.Console
		antlr.DefaultErrorListener
		src string
	}

	// ParserOption defines an function with argument Parser
	ParserOption func(p *Parser)
)

// NewParser creates an instance for Parser
func NewParser(options ...ParserOption) *Parser {
	p := &Parser{
		log: console.NewColorConsole(),
	}
	for _, opt := range options {
		opt(p)
	}

	return p
}

// Accept can parse any terminalNode of api tree by fn.
// -- for debug
func (p *Parser) Accept(fn func(p *api.ApiParserParser, visitor *ApiVisitor) interface{}, content string) (v interface{}, err error) {
	defer func() {
		p := recover()
		if p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%+v", p)
			}
		}
	}()

	inputStream := antlr.NewInputStream(content)
	lexer := api.NewApiParserLexer(inputStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	apiParser := api.NewApiParserParser(tokens)
	apiParser.RemoveErrorListeners()
	apiParser.AddErrorListener(p)
	var visitorOptions []VisitorOption
	visitorOptions = append(visitorOptions, WithVisitorPrefix(p.linePrefix))
	if p.debug {
		visitorOptions = append(visitorOptions, WithVisitorDebug())
	}
	visitor := NewApiVisitor(visitorOptions...)
	v = fn(apiParser, visitor)
	return
}

// Parse is used to parse the api from the specified file name
func (p *Parser) Parse(filename string) (*Api, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	p.src = abs
	data, err := p.readContent(filename)
	if err != nil {
		return nil, err
	}

	return p.parse(filename, data)
}

// ParseContent is used to parse the api from the specified content
func (p *Parser) ParseContent(content string, filename ...string) (*Api, error) {
	var f string
	if len(filename) > 0 {
		f = filename[0]
		abs, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}

		p.src = abs
	}

	return p.parse(f, content)
}

// parse is used to parse api from the content
// filename is only used to mark the file where the error is located
func (p *Parser) parse(filename, content string) (*Api, error) {
	root, err := p.invoke(filename, content)
	if err != nil {
		return nil, err
	}

	var apiAstList []*Api
	apiAstList = append(apiAstList, root)
	for _, imp := range root.Import {
		dir := filepath.Dir(p.src)
		imp := filepath.Join(dir, imp.Value.Text())
		data, err := p.readContent(imp)
		if err != nil {
			return nil, err
		}

		nestedApi, err := p.invoke(imp, data)
		if err != nil {
			return nil, err
		}

		err = p.valid(root, nestedApi)
		if err != nil {
			return nil, err
		}

		apiAstList = append(apiAstList, nestedApi)
	}

	err = p.checkTypeDeclaration(apiAstList)
	if err != nil {
		return nil, err
	}

	allApi := p.memberFill(apiAstList)
	return allApi, nil
}

func (p *Parser) invoke(linePrefix, content string) (v *Api, err error) {
	defer func() {
		p := recover()
		if p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%+v", p)
			}
		}
	}()

	if linePrefix != "" {
		p.linePrefix = linePrefix
	}

	inputStream := antlr.NewInputStream(content)
	lexer := api.NewApiParserLexer(inputStream)
	lexer.RemoveErrorListeners()
	tokens := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	apiParser := api.NewApiParserParser(tokens)
	apiParser.RemoveErrorListeners()
	apiParser.AddErrorListener(p)
	var visitorOptions []VisitorOption
	visitorOptions = append(visitorOptions, WithVisitorPrefix(p.linePrefix))
	if p.debug {
		visitorOptions = append(visitorOptions, WithVisitorDebug())
	}

	visitor := NewApiVisitor(visitorOptions...)
	v = apiParser.Api().Accept(visitor).(*Api)
	v.LinePrefix = p.linePrefix
	return
}

func (p *Parser) valid(mainApi, nestedApi *Api) error {
	err := p.nestedApiCheck(mainApi, nestedApi)
	if err != nil {
		return err
	}

	mainHandlerMap := make(map[string]PlaceHolder)
	mainRouteMap := make(map[string]PlaceHolder)
	mainTypeMap := make(map[string]PlaceHolder)

	routeMap := func(list []*ServiceRoute) (map[string]PlaceHolder, map[string]PlaceHolder) {
		handlerMap := make(map[string]PlaceHolder)
		routeMap := make(map[string]PlaceHolder)

		for _, g := range list {
			handler := g.GetHandler()
			if handler.IsNotNil() {
				handlerName := handler.Text()
				handlerMap[handlerName] = Holder
				route := fmt.Sprintf("%s://%s", g.Route.Method.Text(), g.Route.Path.Text())
				routeMap[route] = Holder
			}
		}

		return handlerMap, routeMap
	}

	for _, each := range mainApi.Service {
		h, r := routeMap(each.ServiceApi.ServiceRoute)

		for k, v := range h {
			mainHandlerMap[k] = v
		}

		for k, v := range r {
			mainRouteMap[k] = v
		}
	}

	for _, each := range mainApi.Type {
		mainTypeMap[each.NameExpr().Text()] = Holder
	}

	// duplicate route check
	err = p.duplicateRouteCheck(nestedApi, mainHandlerMap, mainRouteMap)
	if err != nil {
		return err
	}

	// duplicate type check
	for _, each := range nestedApi.Type {
		if _, ok := mainTypeMap[each.NameExpr().Text()]; ok {
			return fmt.Errorf("%s line %d:%d duplicate type declaration '%s'",
				nestedApi.LinePrefix, each.NameExpr().Line(), each.NameExpr().Column(), each.NameExpr().Text())
		}
	}

	return nil
}

func (p *Parser) duplicateRouteCheck(nestedApi *Api, mainHandlerMap, mainRouteMap map[string]PlaceHolder) error {
	for _, each := range nestedApi.Service {
		var prefix string
		if each.AtServer != nil {
			p := each.AtServer.Kv.Get(prefixKey)
			if p != nil {
				prefix = p.Text()
			}
		}
		for _, r := range each.ServiceApi.ServiceRoute {
			handler := r.GetHandler()
			if !handler.IsNotNil() {
				return fmt.Errorf("%s handler not exist near line %d", nestedApi.LinePrefix, r.Route.Method.Line())
			}

			if _, ok := mainHandlerMap[handler.Text()]; ok {
				return fmt.Errorf("%s line %d:%d duplicate handler '%s'",
					nestedApi.LinePrefix, handler.Line(), handler.Column(), handler.Text())
			}

			p := fmt.Sprintf("%s://%s", r.Route.Method.Text(), path.Join(prefix, r.Route.Path.Text()))
			if _, ok := mainRouteMap[p]; ok {
				return fmt.Errorf("%s line %d:%d duplicate route '%s'",
					nestedApi.LinePrefix, r.Route.Method.Line(), r.Route.Method.Column(), r.Route.Method.Text()+" "+r.Route.Path.Text())
			}
		}
	}
	return nil
}

func (p *Parser) nestedApiCheck(mainApi, nestedApi *Api) error {
	if len(nestedApi.Import) > 0 {
		importToken := nestedApi.Import[0].Import
		return fmt.Errorf("%s line %d:%d the nested api does not support import",
			nestedApi.LinePrefix, importToken.Line(), importToken.Column())
	}

	if mainApi.Syntax != nil && nestedApi.Syntax != nil {
		if mainApi.Syntax.Version.Text() != nestedApi.Syntax.Version.Text() {
			syntaxToken := nestedApi.Syntax.Syntax
			return fmt.Errorf("%s line %d:%d multiple syntax declaration, expecting syntax '%s', but found '%s'",
				nestedApi.LinePrefix, syntaxToken.Line(), syntaxToken.Column(), mainApi.Syntax.Version.Text(), nestedApi.Syntax.Version.Text())
		}
	}

	if len(mainApi.Service) > 0 {
		mainService := mainApi.Service[0]
		for _, service := range nestedApi.Service {
			if mainService.ServiceApi.Name.Text() != service.ServiceApi.Name.Text() {
				return fmt.Errorf("%s multiple service name declaration, expecting service name '%s', but found '%s'",
					nestedApi.LinePrefix, mainService.ServiceApi.Name.Text(), service.ServiceApi.Name.Text())
			}
		}
	}
	return nil
}

func (p *Parser) memberFill(apiList []*Api) *Api {
	var root Api
	for index, each := range apiList {
		if index == 0 {
			root.Syntax = each.Syntax
			root.Info = each.Info
			root.Import = each.Import
		}

		root.Type = append(root.Type, each.Type...)
		root.Service = append(root.Service, each.Service...)
	}

	return &root
}

// checkTypeDeclaration checks whether a struct type has been declared in context
func (p *Parser) checkTypeDeclaration(apiList []*Api) error {
	types := make(map[string]TypeExpr)

	for _, root := range apiList {
		for _, each := range root.Type {
			types[each.NameExpr().Text()] = each
		}
	}

	for _, apiItem := range apiList {
		linePrefix := apiItem.LinePrefix
		err := p.checkTypes(apiItem, linePrefix, types)
		if err != nil {
			return err
		}

		err = p.checkServices(apiItem, types, linePrefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) checkServices(apiItem *Api, types map[string]TypeExpr, linePrefix string) error {
	for _, service := range apiItem.Service {
		for _, each := range service.ServiceApi.ServiceRoute {
			route := each.Route
			err := p.checkRequestBody(route, types, linePrefix)
			if err != nil {
				return err
			}

			if route.Reply != nil && route.Reply.Name.IsNotNil() && route.Reply.Name.Expr().IsNotNil() {
				reply := route.Reply.Name
				var structName string
				switch tp := reply.(type) {
				case *Literal:
					structName = tp.Literal.Text()
				case *Array:
					switch innerTp := tp.Literal.(type) {
					case *Literal:
						structName = innerTp.Literal.Text()
					case *Pointer:
						structName = innerTp.Name.Text()
					}
				}

				if api.IsBasicType(structName) {
					continue
				}

				_, ok := types[structName]
				if !ok {
					return fmt.Errorf("%s line %d:%d can not found declaration '%s' in context",
						linePrefix, route.Reply.Name.Expr().Line(), route.Reply.Name.Expr().Column(), structName)
				}
			}
		}
	}
	return nil
}

func (p *Parser) checkRequestBody(route *Route, types map[string]TypeExpr, linePrefix string) error {
	if route.Req != nil && route.Req.Name.IsNotNil() && route.Req.Name.Expr().IsNotNil() {
		_, ok := types[route.Req.Name.Expr().Text()]
		if !ok {
			return fmt.Errorf("%s line %d:%d can not found declaration '%s' in context",
				linePrefix, route.Req.Name.Expr().Line(), route.Req.Name.Expr().Column(), route.Req.Name.Expr().Text())
		}
	}
	return nil
}

func (p *Parser) checkTypes(apiItem *Api, linePrefix string, types map[string]TypeExpr) error {
	for _, each := range apiItem.Type {
		tp, ok := each.(*TypeStruct)
		if !ok {
			continue
		}

		for _, member := range tp.Fields {
			err := p.checkType(linePrefix, types, member.DataType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parser) checkType(linePrefix string, types map[string]TypeExpr, expr DataType) error {
	if expr == nil {
		return nil
	}

	switch v := expr.(type) {
	case *Literal:
		name := v.Literal.Text()
		if api.IsBasicType(name) {
			return nil
		}
		_, ok := types[name]
		if !ok {
			return fmt.Errorf("%s line %d:%d can not found declaration '%s' in context",
				linePrefix, v.Literal.Line(), v.Literal.Column(), name)
		}

	case *Pointer:
		name := v.Name.Text()
		if api.IsBasicType(name) {
			return nil
		}
		_, ok := types[name]
		if !ok {
			return fmt.Errorf("%s line %d:%d can not found declaration '%s' in context",
				linePrefix, v.Name.Line(), v.Name.Column(), name)
		}
	case *Map:
		return p.checkType(linePrefix, types, v.Value)
	case *Array:
		return p.checkType(linePrefix, types, v.Literal)
	default:
		return nil
	}
	return nil
}

func (p *Parser) readContent(filename string) (string, error) {
	filename = strings.ReplaceAll(filename, `"`, "")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SyntaxError accepts errors and panic it
func (p *Parser) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	str := fmt.Sprintf(`%s line %d:%d  %s`, p.linePrefix, line, column, msg)
	if p.debug {
		p.log.Error(str)
	}
	panic(str)
}

// WithParserDebug returns a debug ParserOption
func WithParserDebug() ParserOption {
	return func(p *Parser) {
		p.debug = true
	}
}

// WithParserPrefix returns a prefix ParserOption
func WithParserPrefix(prefix string) ParserOption {
	return func(p *Parser) {
		p.linePrefix = prefix
	}
}
