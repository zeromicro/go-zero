package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/importstack"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const (
	atServerGroupKey  = "group"
	atServerPrefixKey = "prefix"
)

// API is the parsed api file.
type API struct {
	Filename      string
	Syntax        *ast.SyntaxStmt
	info          *ast.InfoStmt    // Info block does not participate in code generation.
	importStmt    []ast.ImportStmt // ImportStmt block does not participate in code generation.
	TypeStmt      []ast.TypeStmt
	ServiceStmts  []*ast.ServiceStmt
	importManager *importstack.ImportStack
	importSet     map[string]lang.PlaceholderType
}

func convert2API(a *ast.AST, importSet map[string]lang.PlaceholderType, is *importstack.ImportStack) (*API, error) {
	var api = new(API)
	api.importManager = is
	api.importSet = importSet
	api.Filename = a.Filename
	one := a.Stmts[0]
	syntax, ok := one.(*ast.SyntaxStmt)
	if !ok {
		syntax = &ast.SyntaxStmt{
			Syntax: ast.NewTokenNode(
				token.Token{
					Type: token.IDENT,
					Text: token.Syntax,
				},
			),
			Assign: ast.NewTokenNode(
				token.Token{
					Type: token.ASSIGN,
					Text: "=",
				},
			),
			Value: ast.NewTokenNode(
				token.Token{
					Type: token.STRING,
					Text: `"v1"`,
				},
			),
		}
	}

	api.Syntax = syntax
	var hasSyntax, hasInfo bool
	for i := 0; i < len(a.Stmts); i++ {
		one := a.Stmts[i]
		switch val := one.(type) {
		case *ast.SyntaxStmt:
			if hasSyntax {
				return nil, ast.DuplicateStmtError(val.Pos(), "duplicate syntax statement")
			} else {
				hasSyntax = true
			}
		case *ast.InfoStmt:
			if api.info != nil {
				if hasInfo {
					return nil, ast.DuplicateStmtError(val.Pos(), "duplicate info statement")
				}
			} else {
				hasInfo = true
			}
			api.info = val
		case ast.ImportStmt:
			api.importStmt = append(api.importStmt, val)
		case ast.TypeStmt:
			api.TypeStmt = append(api.TypeStmt, val)
		case *ast.ServiceStmt:
			api.ServiceStmts = append(api.ServiceStmts, val)
		}
	}

	return api, nil
}

func (api *API) checkImportStmt() error {
	f := newFilter()
	b := f.addCheckItem(api.Filename, "import value expression")
	for _, v := range api.importStmt {
		switch val := v.(type) {
		case *ast.ImportLiteralStmt:
			b.check(val.Value)
		case *ast.ImportGroupStmt:
			b.check(val.Values...)
		}
	}
	return f.error()
}

func (api *API) checkInfoStmt() error {
	if api.info == nil {
		return nil
	}
	f := newFilter()
	b := f.addCheckItem(api.Filename, "info key expression")
	for _, v := range api.info.Values {
		b.check(v.Key)
	}
	return f.error()
}

func (api *API) checkServiceStmt() error {
	f := newFilter()
	serviceNameChecker := f.addCheckItem(api.Filename, "service name expression")
	handlerChecker := f.addCheckItem(api.Filename, "handler expression")
	pathChecker := f.addCheckItem(api.Filename, "path expression")
	var serviceName = map[string]string{}
	for _, v := range api.ServiceStmts {
		name := strings.TrimSuffix(v.Name.Format(""), "-api")
		if sn, ok := serviceName[name]; ok {
			if sn != name {
				serviceNameChecker.errorManager.add(ast.SyntaxError(v.Name.Pos(), "multiple service name"))
			}
		} else {
			serviceName[name] = name
		}
		var (
			prefix = api.getAtServerValue(v.AtServerStmt, atServerPrefixKey)
			group  = api.getAtServerValue(v.AtServerStmt, atServerGroupKey)
		)
		for _, item := range v.Routes {
			handlerChecker.checkNodeWithPrefix(group, item.AtHandler.Name)
			path := fmt.Sprintf("[%s]:%s", prefix, item.Route.Format(""))
			pathChecker.check(
				ast.NewTokenNode(
					token.Token{
						Text:     path,
						Position: item.Route.Pos(),
					},
				),
			)
		}
	}
	return f.error()
}

func (api *API) checkTypeStmt() error {
	f := newFilter()
	b := f.addCheckItem(api.Filename, "type expression")
	for _, v := range api.TypeStmt {
		switch val := v.(type) {
		case *ast.TypeLiteralStmt:
			b.check(val.Expr.Name)
		case *ast.TypeGroupStmt:
			for _, expr := range val.ExprList {
				b.check(expr.Name)
			}
		}
	}
	return f.error()
}

func (api *API) checkTypeDeclareContext() error {
	var typeMap = map[string]placeholder.Type{}
	for _, v := range api.TypeStmt {
		switch tp := v.(type) {
		case *ast.TypeLiteralStmt:
			typeMap[tp.Expr.Name.Token.Text] = placeholder.PlaceHolder
		case *ast.TypeGroupStmt:
			for _, v := range tp.ExprList {
				typeMap[v.Name.Token.Text] = placeholder.PlaceHolder
			}
		}
	}

	return api.checkTypeContext(typeMap)
}

func (api *API) checkTypeContext(declareContext map[string]placeholder.Type) error {
	var em = newErrorManager()
	for _, v := range api.TypeStmt {
		switch tp := v.(type) {
		case *ast.TypeLiteralStmt:
			em.add(api.checkTypeExprContext(declareContext, tp.Expr.DataType))
		case *ast.TypeGroupStmt:
			for _, v := range tp.ExprList {
				em.add(api.checkTypeExprContext(declareContext, v.DataType))
			}
		}
	}
	return em.error()
}

func (api *API) checkTypeExprContext(declareContext map[string]placeholder.Type, tp ast.DataType) error {
	switch val := tp.(type) {
	case *ast.ArrayDataType:
		return api.checkTypeExprContext(declareContext, val.DataType)
	case *ast.BaseDataType:
		if IsBaseType(val.Base.Token.Text) {
			return nil
		}
		_, ok := declareContext[val.Base.Token.Text]
		if !ok {
			return ast.SyntaxError(val.Base.Pos(), "unresolved type <%s>", val.Base.Token.Text)
		}
		return nil
	case *ast.MapDataType:
		var manager = newErrorManager()
		manager.add(api.checkTypeExprContext(declareContext, val.Key))
		manager.add(api.checkTypeExprContext(declareContext, val.Value))
		return manager.error()
	case *ast.PointerDataType:
		return api.checkTypeExprContext(declareContext, val.DataType)
	case *ast.SliceDataType:
		return api.checkTypeExprContext(declareContext, val.DataType)
	case *ast.StructDataType:
		var manager = newErrorManager()
		for _, e := range val.Elements {
			manager.add(api.checkTypeExprContext(declareContext, e.DataType))
		}
		return manager.error()
	}
	return nil
}

func (api *API) getAtServerValue(atServer *ast.AtServerStmt, key string) string {
	if atServer == nil {
		return ""
	}

	for _, val := range atServer.Values {
		if val.Key.Token.Text == key {
			return val.Value.Token.Text
		}
	}

	return ""
}

func (api *API) mergeAPI(in *API) error {
	if api.Syntax.Value.Format() != in.Syntax.Value.Format() {
		return ast.SyntaxError(
			in.Syntax.Value.Pos(),
			"multiple syntax value expression, expected <%s>, got <%s>",
			api.Syntax.Value.Format(),
			in.Syntax.Value.Format(),
		)
	}
	api.TypeStmt = append(api.TypeStmt, in.TypeStmt...)
	api.ServiceStmts = append(api.ServiceStmts, in.ServiceStmts...)
	return nil
}

func (api *API) parseImportedAPI(imports []ast.ImportStmt) ([]*API, error) {
	var list []*API
	if len(imports) == 0 {
		return list, nil
	}

	var importValueSet = map[string]token.Token{}
	for _, imp := range imports {
		switch val := imp.(type) {
		case *ast.ImportLiteralStmt:
			importValueSet[strings.ReplaceAll(val.Value.Token.Text, `"`, "")] = val.Value.Token
		case *ast.ImportGroupStmt:
			for _, v := range val.Values {
				importValueSet[strings.ReplaceAll(v.Token.Text, `"`, "")] = v.Token
			}
		}
	}

	dir := filepath.Dir(api.Filename)
	for impPath, tok := range importValueSet {
		if !filepath.IsAbs(impPath) {
			impPath = filepath.Join(dir, impPath)
		}
		// import cycle check
		if err := api.importManager.Push(impPath); err != nil {
			return nil, ast.SyntaxError(tok.Position, err.Error())
		}

		if _, ok := api.importSet[impPath]; ok {
			api.importManager.Pop()
			continue
		}
		api.importSet[impPath] = lang.Placeholder

		p := New(impPath, "")
		ast := p.Parse()
		if err := p.CheckErrors(); err != nil {
			return nil, err
		}

		nestedApi, err := convert2API(ast, api.importSet, api.importManager)
		if err != nil {
			return nil, err
		}

		if err = nestedApi.parseReverse(); err != nil {
			return nil, err
		}

		api.importManager.Pop()
		list = append(list, nestedApi)

		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (api *API) parseReverse() error {
	list, err := api.parseImportedAPI(api.importStmt)
	if err != nil {
		return err
	}
	for _, e := range list {
		if err = api.mergeAPI(e); err != nil {
			return err
		}
	}
	return nil
}

func (api *API) SelfCheck() error {
	if err := api.parseReverse(); err != nil {
		return err
	}
	if err := api.checkImportStmt(); err != nil {
		return err
	}
	if err := api.checkInfoStmt(); err != nil {
		return err
	}
	if err := api.checkTypeStmt(); err != nil {
		return err
	}
	if err := api.checkServiceStmt(); err != nil {
		return err
	}
	return api.checkTypeDeclareContext()
}
