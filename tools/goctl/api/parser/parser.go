package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"unicode"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type parser struct {
	ast  *ast.Api
	spec *spec.ApiSpec
}

func Parse(filename string) (*spec.ApiSpec, error) {
	astParser := ast.NewParser(ast.WithParserPrefix(filepath.Base(filename)))
	ast, err := astParser.Parse(filename)
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

func ParseContent(content string) (*spec.ApiSpec, error) {
	astParser := ast.NewParser()
	ast, err := astParser.ParseContent(content)
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
	var properties = make(map[string]string, 0)
	if p.ast.Info != nil {
		p.spec.Info = spec.Info{}
		for _, kv := range p.ast.Info.Kvs {
			properties[kv.Key.Text()] = kv.Value.Text()
		}
	}
	p.spec.Info.Properties = properties
}

func (p parser) fillSyntax() {
	if p.ast.Syntax != nil {
		p.spec.Syntax = spec.ApiSyntax{Version: p.ast.Syntax.Version.Text()}
	}
}

func (p parser) fillImport() {
	if len(p.ast.Import) > 0 {
		for _, item := range p.ast.Import {
			p.spec.Imports = append(p.spec.Imports, spec.Import{Value: item.Value.Text()})
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
			return errors.New(fmt.Sprintf("unknown type %+v", v))
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
					} else {
						member.Type = *tp
					}
				}
				members = append(members, member)
			}
			v.Members = members
			types = append(types, v)
		default:
			return errors.New(fmt.Sprintf("unknown type %+v", v))
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
	return nil, errors.New(fmt.Sprintf("type %s not defined", name))
}

func (p parser) fieldToMember(field *ast.TypeField) spec.Member {
	var name = ""
	var tag = ""
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
			return spec.PrimitiveType{RawName: raw}
		} else {
			return spec.DefineStruct{RawName: raw}
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
		} else {
			return spec.PointerType{RawName: v.PointerExpr.Text(), Type: spec.DefineStruct{RawName: raw}}
		}
	}

	panic(fmt.Sprintf("unspported type %+v", in))
}

func (p parser) stringExprs(docs []ast.Expr) []string {
	var result []string
	for _, item := range docs {
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
		if item.AtServer != nil {
			var properties = make(map[string]string, 0)
			for _, kv := range item.AtServer.Kv {
				properties[kv.Key.Text()] = kv.Value.Text()
			}
			group.Annotation.Properties = properties
		}

		for _, astRoute := range item.ServiceApi.ServiceRoute {
			route := spec.Route{
				Annotation: spec.Annotation{},
				Method:     astRoute.Route.Method.Text(),
				Path:       astRoute.Route.Path.Text(),
			}
			if astRoute.AtHandler != nil {
				route.Handler = astRoute.AtHandler.Name.Text()
			}

			if astRoute.AtServer != nil {
				var properties = make(map[string]string, 0)
				for _, kv := range astRoute.AtServer.Kv {
					properties[kv.Key.Text()] = kv.Value.Text()
				}
				route.Annotation.Properties = properties
				if len(route.Handler) == 0 {
					route.Handler = properties["handler"]
				}
				if len(route.Handler) == 0 {
					return fmt.Errorf("missing handler annotation for %q", route.Path)
				}

				for _, char := range route.Handler {
					if !unicode.IsDigit(char) && !unicode.IsLetter(char) {
						return errors.New(fmt.Sprintf("route [%s] handler [%s] invalid, handler name should only contains letter or digit",
							route.Path, route.Handler))
					}
				}
			}

			if astRoute.Route.Req != nil {
				route.RequestType = p.astTypeToSpec(astRoute.Route.Req.Name)
			}
			if astRoute.Route.Reply != nil {
				route.ResponseType = p.astTypeToSpec(astRoute.Route.Reply.Name)
			}
			if astRoute.AtDoc != nil {
				var properties = make(map[string]string, 0)
				for _, kv := range astRoute.AtDoc.Kv {
					properties[kv.Key.Text()] = kv.Value.Text()
				}
				route.AtDoc.Properties = properties
				if astRoute.AtDoc.LineDoc != nil {
					route.AtDoc.Text = astRoute.AtDoc.LineDoc.Text()
				}
			}

			err := p.fillRouteType(&route)
			if err != nil {
				return err
			}

			group.Routes = append(group.Routes, route)

			name := item.ServiceApi.Name.Text()
			if len(p.spec.Service.Name) > 0 && p.spec.Service.Name != name {
				return errors.New(fmt.Sprintf("mulit service name defined %s and %s", name, p.spec.Service.Name))
			}
			p.spec.Service.Name = name
		}
		groups = append(groups, group)
	}
	p.spec.Service.Groups = groups

	return nil
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
