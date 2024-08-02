package ast

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeromicro/antlr"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

type (
	// Parser provides api parsing capabilities
	Parser struct {
		antlr.DefaultErrorListener
		linePrefix               string
		debug                    bool
		log                      console.Console
		skipCheckTypeDeclaration bool
		handlerMap               map[string]PlaceHolder
		routeMap                 map[string]PlaceHolder
		typeMap                  map[string]PlaceHolder
		fileMap                  map[string]PlaceHolder
		importStatck             importStack
		syntax                   *SyntaxExpr
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
	p.handlerMap = make(map[string]PlaceHolder)
	p.routeMap = make(map[string]PlaceHolder)
	p.typeMap = make(map[string]PlaceHolder)
	p.fileMap = make(map[string]PlaceHolder)

	return p
}

// Accept can parse any terminalNode of api tree by fn.
// -- for debug
func (p *Parser) Accept(fn func(p *api.ApiParserParser, visitor *ApiVisitor) any, content string) (v any, err error) {
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

	data, err := p.readContent(filename)
	if err != nil {
		return nil, err
	}

	p.importStatck.push(abs)
	return p.parse(filename, data)
}

// ParseContent is used to parse the api from the specified content
func (p *Parser) ParseContent(content string, filename ...string) (*Api, error) {
	var f, abs string
	if len(filename) > 0 {
		f = filename[0]
		a, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}

		abs = a
	}

	p.importStatck.push(abs)
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
	p.storeVerificationInfo(root)
	p.syntax = root.Syntax
	impApiAstList, err := p.invokeImportedApi(filename, root.Import)
	if err != nil {
		return nil, err
	}
	apiAstList = append(apiAstList, impApiAstList...)

	if !p.skipCheckTypeDeclaration {
		err = p.checkTypeDeclaration(apiAstList)
		if err != nil {
			return nil, err
		}
	}

	allApi := p.memberFill(apiAstList)
	return allApi, nil
}

func (p *Parser) invokeImportedApi(filename string, imports []*ImportExpr) ([]*Api, error) {
	var apiAstList []*Api
	for _, imp := range imports {
		dir := filepath.Dir(filename)
		impPath := strings.ReplaceAll(imp.Value.Text(), "\"", "")
		if !filepath.IsAbs(impPath) {
			impPath = filepath.Join(dir, impPath)
		}
		// import cycle check
		if err := p.importStatck.push(impPath); err != nil {
			return nil, err
		}
		// ignore already imported file
		if p.alreadyImported(impPath) {
			p.importStatck.pop()
			continue
		}
		p.fileMap[impPath] = PlaceHolder{}

		data, err := p.readContent(impPath)
		if err != nil {
			return nil, err
		}

		nestedApi, err := p.invoke(impPath, data)
		if err != nil {
			return nil, err
		}

		err = p.valid(nestedApi)
		if err != nil {
			return nil, err
		}
		p.storeVerificationInfo(nestedApi)
		apiAstList = append(apiAstList, nestedApi)
		list, err := p.invokeImportedApi(impPath, nestedApi.Import)
		p.importStatck.pop()
		apiAstList = append(apiAstList, list...)

		if err != nil {
			return nil, err
		}
	}
	return apiAstList, nil
}

func (p *Parser) alreadyImported(filename string) bool {
	_, ok := p.fileMap[filename]
	return ok
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

// storeVerificationInfo stores information for verification
func (p *Parser) storeVerificationInfo(api *Api) {
	routeMap := func(list []*ServiceRoute, prefix string) {
		for _, g := range list {
			handler := g.GetHandler()
			if handler.IsNotNil() {
				handlerName := handler.Text()
				p.handlerMap[handlerName] = Holder
				route := fmt.Sprintf("%s://%s", g.Route.Method.Text(), path.Join(prefix, g.Route.Path.Text()))
				p.routeMap[route] = Holder
			}
		}
	}

	for _, each := range api.Service {
		var prefix string
		if each.AtServer != nil {
			pExp := each.AtServer.Kv.Get(prefixKey)
			if pExp != nil {
				prefix = pExp.Text()
			}
		}
		routeMap(each.ServiceApi.ServiceRoute, prefix)
	}

	for _, each := range api.Type {
		p.typeMap[each.NameExpr().Text()] = Holder
	}
}

func (p *Parser) valid(nestedApi *Api) error {
	if p.syntax != nil && nestedApi.Syntax != nil {
		if p.syntax.Version.Text() != nestedApi.Syntax.Version.Text() {
			syntaxToken := nestedApi.Syntax.Syntax
			return fmt.Errorf("%s line %d:%d multiple syntax declaration, expecting syntax '%s', but found '%s'",
				nestedApi.LinePrefix, syntaxToken.Line(), syntaxToken.Column(), p.syntax.Version.Text(), nestedApi.Syntax.Version.Text())
		}
	}

	// duplicate route check
	err := p.duplicateRouteCheck(nestedApi)
	if err != nil {
		return err
	}

	// duplicate type check
	for _, each := range nestedApi.Type {
		if _, ok := p.typeMap[each.NameExpr().Text()]; ok {
			return fmt.Errorf("%s line %d:%d duplicate type declaration '%s'",
				nestedApi.LinePrefix, each.NameExpr().Line(), each.NameExpr().Column(), each.NameExpr().Text())
		}
	}

	return nil
}

func (p *Parser) duplicateRouteCheck(nestedApi *Api) error {
	for _, each := range nestedApi.Service {
		var prefix, group string
		if each.AtServer != nil {
			p := each.AtServer.Kv.Get(prefixKey)
			if p != nil {
				prefix = p.Text()
			}
			g := each.AtServer.Kv.Get(groupKey)
			if g != nil {
				group = g.Text()
			}
		}
		for _, r := range each.ServiceApi.ServiceRoute {
			handler := r.GetHandler()
			if !handler.IsNotNil() {
				return fmt.Errorf("%s handler not exist near line %d", nestedApi.LinePrefix, r.Route.Method.Line())
			}

			handlerKey := handler.Text()
			if len(group) > 0 {
				handlerKey = fmt.Sprintf("%s/%s", group, handler.Text())
			}
			if _, ok := p.handlerMap[handlerKey]; ok {
				return fmt.Errorf("%s line %d:%d duplicate handler '%s'",
					nestedApi.LinePrefix, handler.Line(), handler.Column(), handlerKey)
			}

			routeKey := fmt.Sprintf("%s://%s", r.Route.Method.Text(), path.Join(prefix, r.Route.Path.Text()))
			if _, ok := p.routeMap[routeKey]; ok {
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
					return fmt.Errorf("%s line %d:%d can not find declaration '%s' in context",
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
			return fmt.Errorf("%s line %d:%d can not find declaration '%s' in context",
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
			return fmt.Errorf("%s line %d:%d can not find declaration '%s' in context",
				linePrefix, v.Literal.Line(), v.Literal.Column(), name)
		}

	case *Pointer:
		name := v.Name.Text()
		if api.IsBasicType(name) {
			return nil
		}
		_, ok := types[name]
		if !ok {
			return fmt.Errorf("%s line %d:%d can not find declaration '%s' in context",
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
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SyntaxError accepts errors and panic it
func (p *Parser) SyntaxError(_ antlr.Recognizer, _ any, line, column int, msg string, _ antlr.RecognitionException) {
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

func WithParserSkipCheckTypeDeclaration() ParserOption {
	return func(p *Parser) {
		p.skipCheckTypeDeclaration = true
	}
}
