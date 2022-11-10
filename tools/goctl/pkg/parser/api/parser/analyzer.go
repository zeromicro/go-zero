package parser

import (
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
)

type Analyzer struct {
	api  *API
	spec *spec.ApiSpec
}

func (a *Analyzer) astTypeToSpec(in ast.DataType) (spec.Type, error) {
	isLiteralType := func(dt ast.DataType) bool {
		_, ok := dt.(*ast.BaseDataType)
		if ok {
			return true
		}
		_, ok = dt.(*ast.AnyDataType)
		return ok
	}
	switch v := (in).(type) {
	case *ast.BaseDataType:
		return spec.PrimitiveType{RawName: v.RawText()}, nil
	case *ast.AnyDataType:
		return spec.PrimitiveType{RawName: v.RawText()}, nil
	case *ast.StructDataType:
		// TODO(keson) feature: can be extended
	case *ast.InterfaceDataType:
		return spec.InterfaceType{RawName: v.RawText()}, nil
	case *ast.MapDataType:
		if !isLiteralType(v.Key) {
			return nil, ast.SyntaxError(v.Pos(), "expected literal type, got <%T>", v)
		}
		if !v.Key.CanEqual() {
			return nil, ast.SyntaxError(v.Pos(), "map key <%T> must be equal data type", v)
		}
		value, err := a.astTypeToSpec(v.Value)
		if err != nil {
			return nil, err
		}
		return spec.MapType{
			RawName: v.RawText(),
			Key:     v.RawText(),
			Value:   value,
		}, nil
	case *ast.PointerDataType:
		value, err := a.astTypeToSpec(v.DataType)
		if err != nil {
			return nil, err
		}
		return spec.PointerType{
			RawName: v.RawText(),
			Type:    value,
		}, nil
	case *ast.ArrayDataType:
		// TODO(keson) feature: can be extended
	case *ast.SliceDataType:
		value, err := a.astTypeToSpec(v.DataType)
		if err != nil {
			return nil, err
		}
		return spec.ArrayType{
			RawName: v.RawText(),
			Value:   value,
		}, nil
	}

	return nil, ast.SyntaxError(in.Pos(), "unsupported type <%T>", in)
}

func (a *Analyzer) convert2Spec() error {
	if err := a.fillTypes(); err != nil {
		return err
	}

	return a.fillService()
}

func (a *Analyzer) convertAtDoc(atDoc ast.AtDocStmt) spec.AtDoc {
	var ret spec.AtDoc
	switch val := atDoc.(type) {
	case *ast.AtDocLiteralStmt:
		ret.Text = val.Value.Text
	case *ast.AtDocGroupStmt:
		ret.Properties = a.convertKV(val.Values)
	}
	return ret
}

func (a *Analyzer) convertKV(kv []*ast.KVExpr) map[string]string {
	var ret = map[string]string{}
	for _, v := range kv {
		ret[v.Key.Text] = v.Value.Text
	}
	return ret
}

// break changes: not support anonymous structure.
func (a *Analyzer) fieldToMember(field *ast.ElemExpr) (spec.Member, error) {
	var name []string
	for _, v := range field.Name {
		name = append(name, v.Text)
	}
	tp, err := a.astTypeToSpec(field.DataType)
	if err != nil {
		return spec.Member{}, err
	}
	return spec.Member{
		Name: strings.Join(name, ", "),
		Type: tp,
		Tag:  field.Tag.Text,
	}, nil
}

func (a *Analyzer) fillRouteType(route *spec.Route) error {
	if route.RequestType != nil {
		switch route.RequestType.(type) {
		case spec.DefineStruct:
			tp, err := a.findDefinedType(route.RequestType.Name())
			if err != nil {
				return err
			}

			route.RequestType = tp
		}
	}

	if route.ResponseType != nil {
		switch route.ResponseType.(type) {
		case spec.DefineStruct:
			tp, err := a.findDefinedType(route.ResponseType.Name())
			if err != nil {
				return err
			}

			route.ResponseType = tp
		}
	}

	return nil
}

func (a *Analyzer) fillService() error {
	var groups []spec.Group
	for _, item := range a.api.ServiceStmts {
		var group spec.Group
		if item.AtServerStmt != nil {
			group.Annotation.Properties = a.convertKV(item.AtServerStmt.Values)
		}

		for _, astRoute := range item.Routes {
			route := spec.Route{
				Method: astRoute.Route.Method.Text,
				Path:   astRoute.Route.Path.Format(""),
			}
			if astRoute.AtDoc != nil {
				route.AtDoc = a.convertAtDoc(astRoute.AtDoc)
			}
			if astRoute.AtHandler != nil {
				route.AtDoc = a.convertAtDoc(astRoute.AtDoc)
				route.Handler = astRoute.AtHandler.Name.Text
			}

			if astRoute.Route.Request != nil {
				requestType, err := a.getType(astRoute.Route.Request)
				if err != nil {
					return err
				}
				route.RequestType = requestType
			}
			if astRoute.Route.Response != nil {
				responseType, err := a.getType(astRoute.Route.Response)
				if err != nil {
					return err
				}
				route.ResponseType = responseType
			}

			if err := a.fillRouteType(&route); err != nil {
				return err
			}

			group.Routes = append(group.Routes, route)

			name := item.Name.Format("")
			if len(a.spec.Service.Name) > 0 && a.spec.Service.Name != name {
				return ast.SyntaxError(item.Name.Pos(), "multiple service names defined <%s> and <%s>", name, a.spec.Service.Name)
			}
			a.spec.Service.Name = name
		}
		groups = append(groups, group)
	}

	a.spec.Service.Groups = groups
	return nil
}

func (a *Analyzer) fillTypes() error {
	for _, item := range a.api.TypeStmt {
		switch v := (item).(type) {
		case *ast.TypeLiteralStmt:
			err := a.fillTypeExpr(v.Expr)
			if err != nil {
				return err
			}
		case *ast.TypeGroupStmt:
			for _, expr := range v.ExprList {
				err := a.fillTypeExpr(expr)
				if err != nil {
					return err
				}
			}
		}
	}

	var types []spec.Type
	for _, item := range a.spec.Types {
		switch v := (item).(type) {
		case spec.DefineStruct:
			var members []spec.Member
			for _, member := range v.Members {
				switch v := member.Type.(type) {
				case spec.DefineStruct:
					tp, err := a.findDefinedType(v.RawName)
					if err != nil {
						return err
					}

					member.Type = tp
				}
				members = append(members, member)
			}
			v.Members = members
			types = append(types, v)
		default:
			return fmt.Errorf("unknown type %+v", v)
		}
	}
	a.spec.Types = types

	return nil
}

func (a *Analyzer) fillTypeExpr(expr *ast.TypeExpr) error {
	switch val := expr.DataType.(type) {
	case *ast.StructDataType:
		var members []spec.Member
		for _, item := range val.Elements {
			m, err := a.fieldToMember(item)
			if err != nil {
				return err
			}
			members = append(members, m)
		}
		a.spec.Types = append(a.spec.Types, spec.DefineStruct{
			RawName: expr.Name.Text,
			Members: members,
		})
		return nil
	default:
		return ast.SyntaxError(expr.Pos(), "expected <struct> expr, got <%T>", expr.DataType)
	}
}

func (a *Analyzer) findDefinedType(name string) (spec.Type, error) {
	for _, item := range a.spec.Types {
		if _, ok := item.(spec.DefineStruct); ok {
			if item.Name() == name {
				return item, nil
			}
		}
	}

	return nil, fmt.Errorf("type %s not defined", name)
}

func (a *Analyzer) getType(expr *ast.BodyStmt) (spec.Type, error) {
	body := expr.Body
	var tp spec.Type
	var err error
	var rawText = body.Format("")
	if IsBaseType(body.Value.Text) {
		tp = spec.PrimitiveType{RawName: body.Value.Text}
	} else {
		tp, err = a.findDefinedType(body.Value.Text)
		if err != nil {
			return nil, err
		}
	}
	if body.LBrack.Valid() {
		if body.Star.Valid() {
			return spec.PointerType{
				RawName: rawText,
				Type:    tp,
			}, nil
		}
		return spec.ArrayType{
			RawName: rawText,
			Value:   tp,
		}, nil
	}
	if body.Star.Valid() {
		return spec.PointerType{
			RawName: rawText,
			Type:    tp,
		}, nil
	}
	return tp, nil
}

func Parse(filename string, src interface{}) (*spec.ApiSpec, error) {
	p := New(filename, src, SkipComment)
	ast := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return nil, err
	}

	var importManager = make(map[string]placeholder.Type)
	importManager[ast.Filename]=placeholder.PlaceHolder
	api, err := convert2API(ast, importManager)
	if err != nil {
		return nil, err
	}

	var result = new(spec.ApiSpec)
	analyzer := Analyzer{
		api:  api,
		spec: result,
	}

	err = analyzer.convert2Spec()
	if err != nil {
		return nil, err
	}

	return result, nil
}

var kind = map[string]placeholder.Type{
	"bool":       placeholder.PlaceHolder,
	"int":        placeholder.PlaceHolder,
	"int8":       placeholder.PlaceHolder,
	"int16":      placeholder.PlaceHolder,
	"int32":      placeholder.PlaceHolder,
	"int64":      placeholder.PlaceHolder,
	"uint":       placeholder.PlaceHolder,
	"uint8":      placeholder.PlaceHolder,
	"uint16":     placeholder.PlaceHolder,
	"uint32":     placeholder.PlaceHolder,
	"uint64":     placeholder.PlaceHolder,
	"uintptr":    placeholder.PlaceHolder,
	"float32":    placeholder.PlaceHolder,
	"float64":    placeholder.PlaceHolder,
	"complex64":  placeholder.PlaceHolder,
	"complex128": placeholder.PlaceHolder,
	"string":     placeholder.PlaceHolder,
	"byte":       placeholder.PlaceHolder,
	"rune":       placeholder.PlaceHolder,
	"any":       placeholder.PlaceHolder,
}

func IsBaseType(text string) bool {
	_, ok := kind[text]
	return ok
}
