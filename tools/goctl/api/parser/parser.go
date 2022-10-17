package parser

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type parserInterface interface {
	parse() (*spec.ApiSpec, error)
}

type parser struct {
	ast  *ast.Api
	spec *spec.ApiSpec
}

// NeedReformatSchemaContent schema 中定义的内容是否需要reformat
func NeedReformatSchemaContent(spec *spec.ApiSpec) bool {
	if strings.HasPrefix(spec.Syntax.Version, openApi) {
		return false
	}
	return true
}

// Parse parses the api file
func Parse(filename string) (*spec.ApiSpec, error) {
	parser, err := newParse(context.Background(), filename)
	if err != nil {
		return nil, err
	}
	return parser.parse()
}

func parseContent(content string, skipCheckTypeDeclaration bool, filename ...string) (*spec.ApiSpec, error) {
	var astParser *ast.Parser
	if skipCheckTypeDeclaration {
		astParser = ast.NewParser(ast.WithParserSkipCheckTypeDeclaration())
	} else {
		astParser = ast.NewParser()
	}
	ast, err := astParser.ParseContent(content, filename...)
	if err != nil {
		return nil, err
	}

	spec := new(spec.ApiSpec)
	p := parser{ast: ast, spec: spec}
	err = p.convert2Spec()
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// ParseContent parses the api content
func ParseContent(content string, filename ...string) (*spec.ApiSpec, error) {
	return parseContent(content, false, filename...)
}

// ParseContentWithParserSkipCheckTypeDeclaration parses the api content with skip check type declaration
func ParseContentWithParserSkipCheckTypeDeclaration(content string, filename ...string) (*spec.ApiSpec, error) {
	return parseContent(content, true, filename...)
}

func newParse(ctx context.Context, filename string) (parserInterface, error) {
	if strings.HasSuffix(filename, ".api") {
		astParser := ast.NewParser(ast.WithParserPrefix(filepath.Base(filename)), ast.WithParserDebug())
		ast, err := astParser.Parse(filename)
		if err != nil {
			return nil, err
		}
		return parser{ast: ast, spec: &spec.ApiSpec{}}, nil
	} else if strings.HasSuffix(filename, ".yaml") {
		parser, err := newOpenApi3Parser(ctx, filename)
		if err == nil {
			return parser, nil
		}
		return newOpenApi2Parser(filename)
	}
	return nil, fmt.Errorf("un expect file type: %s", filename)
}

func (p parser) parse() (*spec.ApiSpec, error) {
	if err := p.convert2Spec(); err != nil {
		return nil, err
	}

	return p.spec, nil
}

func (p parser) convert2Spec() error {
	p.fillInfo()
	p.fillSyntax()
	p.fillImport()
	err := p.fillTypes()
	if err != nil {
		return err
	}

	return p.fillService()
}

func (p parser) fillInfo() {
	properties := make(map[string]string)
	if p.ast.Info != nil {
		for _, kv := range p.ast.Info.Kvs {
			properties[kv.Key.Text()] = kv.Value.Text()
		}
	}
	p.spec.Info.Properties = properties
}

func (p parser) fillSyntax() {
	if p.ast.Syntax != nil {
		p.spec.Syntax = spec.ApiSyntax{
			Version: p.ast.Syntax.Version.Text(),
			Doc:     p.stringExprs(p.ast.Syntax.DocExpr),
			Comment: p.stringExprs([]ast.Expr{p.ast.Syntax.CommentExpr}),
		}
	}
}

func (p parser) fillImport() {
	if len(p.ast.Import) > 0 {
		for _, item := range p.ast.Import {
			p.spec.Imports = append(p.spec.Imports, spec.Import{
				Value:   item.Value.Text(),
				Doc:     p.stringExprs(item.DocExpr),
				Comment: p.stringExprs([]ast.Expr{item.CommentExpr}),
			})
		}
	}
}

func (p parser) fillTypes() error {
	for _, item := range p.ast.Type {
		switch v := (item).(type) {
		case *ast.TypeStruct:
			var members []spec.Member
			for _, item := range v.Fields {
				members = append(members, p.fieldToMember(item))
			}
			p.spec.Types = append(p.spec.Types, spec.DefineStruct{
				RawName: v.Name.Text(),
				Members: members,
				Docs:    p.stringExprs(v.Doc()),
			})
		default:
			return fmt.Errorf("unknown type %+v", v)
		}
	}

	var types []spec.Type
	for _, item := range p.spec.Types {
		switch v := (item).(type) {
		case spec.DefineStruct:
			var members []spec.Member
			for _, member := range v.Members {
				switch v := member.Type.(type) {
				case spec.DefineStruct:
					tp, err := p.findDefinedType(v.RawName)
					if err != nil {
						return err
					}

					member.Type = *tp
				}
				members = append(members, member)
			}
			v.Members = members
			types = append(types, v)
		default:
			return fmt.Errorf("unknown type %+v", v)
		}
	}
	p.spec.Types = types

	return nil
}

func (p parser) findDefinedType(name string) (*spec.Type, error) {
	for _, item := range p.spec.Types {
		if _, ok := item.(spec.DefineStruct); ok {
			if item.Name() == name {
				return &item, nil
			}
		}
	}

	return nil, fmt.Errorf("type %s not defined", name)
}

func (p parser) fieldToMember(field *ast.TypeField) spec.Member {
	name := ""
	tag := ""
	if !field.IsAnonymous {
		name = field.Name.Text()
		if field.Tag == nil {
			panic(fmt.Sprintf("error: line %d:%d field %s has no tag", field.Name.Line(), field.Name.Column(),
				field.Name.Text()))
		}

		tag = field.Tag.Text()
	}
	return spec.Member{
		Name:     name,
		Type:     p.astTypeToSpec(field.DataType),
		Tag:      tag,
		Comment:  p.commentExprs(field.Comment()),
		Docs:     p.stringExprs(field.Doc()),
		IsInline: field.IsAnonymous,
	}
}

func (p parser) astTypeToSpec(in ast.DataType) spec.Type {
	switch v := (in).(type) {
	case *ast.Literal:
		raw := v.Literal.Text()
		if api.IsBasicType(raw) {
			return spec.PrimitiveType{
				RawName: raw,
			}
		}

		return spec.DefineStruct{
			RawName: raw,
		}
	case *ast.Interface:
		return spec.InterfaceType{RawName: v.Literal.Text()}
	case *ast.Map:
		return spec.MapType{RawName: v.MapExpr.Text(), Key: v.Key.Text(), Value: p.astTypeToSpec(v.Value)}
	case *ast.Array:
		return spec.ArrayType{RawName: v.ArrayExpr.Text(), Value: p.astTypeToSpec(v.Literal)}
	case *ast.Pointer:
		raw := v.Name.Text()
		if api.IsBasicType(raw) {
			return spec.PointerType{RawName: v.PointerExpr.Text(), Type: spec.PrimitiveType{RawName: raw}}
		}

		return spec.PointerType{RawName: v.PointerExpr.Text(), Type: spec.DefineStruct{RawName: raw}}
	}

	panic(fmt.Sprintf("unspported type %+v", in))
}

func (p parser) stringExprs(docs []ast.Expr) []string {
	var result []string
	for _, item := range docs {
		if item == nil {
			continue
		}
		result = append(result, item.Text())
	}
	return result
}

func (p parser) commentExprs(comment ast.Expr) string {
	if comment == nil {
		return ""
	}

	return comment.Text()
}

func (p parser) fillService() error {
	var groups []spec.Group
	for _, item := range p.ast.Service {
		var group spec.Group
		p.fillAtServer(item, &group)

		for _, astRoute := range item.ServiceApi.ServiceRoute {
			route := spec.Route{
				AtServerAnnotation: spec.Annotation{},
				Method:             astRoute.Route.Method.Text(),
				Path:               astRoute.Route.Path.Text(),
				Doc:                p.stringExprs(astRoute.Route.DocExpr),
				Comment:            p.stringExprs([]ast.Expr{astRoute.Route.CommentExpr}),
			}
			if astRoute.AtHandler != nil {
				route.Handler = astRoute.AtHandler.Name.Text()
				route.HandlerDoc = append(route.HandlerDoc, p.stringExprs(astRoute.AtHandler.DocExpr)...)
				route.HandlerComment = append(route.HandlerComment, p.stringExprs([]ast.Expr{astRoute.AtHandler.CommentExpr})...)
			}

			err := p.fillRouteAtServer(astRoute, &route)
			if err != nil {
				return err
			}

			if astRoute.Route.Req != nil {
				route.RequestType = p.astTypeToSpec(astRoute.Route.Req.Name)
			}
			if astRoute.Route.Reply != nil {
				route.ResponseType = p.astTypeToSpec(astRoute.Route.Reply.Name)
			}
			if astRoute.AtDoc != nil {
				properties := make(map[string]string)
				for _, kv := range astRoute.AtDoc.Kv {
					properties[kv.Key.Text()] = kv.Value.Text()
				}
				route.AtDoc.Properties = properties
				if astRoute.AtDoc.LineDoc != nil {
					route.AtDoc.Text = astRoute.AtDoc.LineDoc.Text()
				}
			}

			err = p.fillRouteType(&route)
			if err != nil {
				return err
			}

			group.Routes = append(group.Routes, route)

			name := item.ServiceApi.Name.Text()
			if len(p.spec.Service.Name) > 0 && p.spec.Service.Name != name {
				return fmt.Errorf("multiple service names defined %s and %s", name, p.spec.Service.Name)
			}
			p.spec.Service.Name = name
		}
		groups = append(groups, group)
	}
	p.spec.Service.Groups = groups

	return nil
}

func (p parser) fillRouteAtServer(astRoute *ast.ServiceRoute, route *spec.Route) error {
	if astRoute.AtServer != nil {
		properties := make(map[string]string)
		for _, kv := range astRoute.AtServer.Kv {
			properties[kv.Key.Text()] = kv.Value.Text()
		}
		route.AtServerAnnotation.Properties = properties
		if len(route.Handler) == 0 {
			route.Handler = properties["handler"]
		}
		if len(route.Handler) == 0 {
			return fmt.Errorf("missing handler annotation for %q", route.Path)
		}

		for _, char := range route.Handler {
			if !unicode.IsDigit(char) && !unicode.IsLetter(char) {
				return fmt.Errorf("route [%s] handler [%s] invalid, handler name should only contains letter or digit",
					route.Path, route.Handler)
			}
		}
	}
	return nil
}

func (p parser) fillAtServer(item *ast.Service, group *spec.Group) {
	if item.AtServer != nil {
		properties := make(map[string]string)
		for _, kv := range item.AtServer.Kv {
			properties[kv.Key.Text()] = kv.Value.Text()
		}
		group.Annotation.Properties = properties
	}
}

func (p parser) fillRouteType(route *spec.Route) error {
	if route.RequestType != nil {
		switch route.RequestType.(type) {
		case spec.DefineStruct:
			tp, err := p.findDefinedType(route.RequestType.Name())
			if err != nil {
				return err
			}

			route.RequestType = *tp
		}
	}

	if route.ResponseType != nil {
		switch route.ResponseType.(type) {
		case spec.DefineStruct:
			tp, err := p.findDefinedType(route.ResponseType.Name())
			if err != nil {
				return err
			}

			route.ResponseType = *tp
		}
	}

	return nil
}
