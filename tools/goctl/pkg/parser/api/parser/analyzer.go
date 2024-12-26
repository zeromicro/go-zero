package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/importstack"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

// Analyzer analyzes the ast and converts it to spec.
type Analyzer struct {
	api  *API
	spec *spec.ApiSpec
}

func (a *Analyzer) astTypeToSpec(in ast.DataType) (spec.Type, error) {
	isLiteralType := func(dt ast.DataType) bool {
		if _, ok := dt.(*ast.BaseDataType); ok {
			return true
		}

		_, ok := dt.(*ast.AnyDataType)
		return ok
	}

	switch v := (in).(type) {
	case *ast.BaseDataType:
		raw := v.RawText()
		if IsBaseType(raw) {
			return spec.PrimitiveType{
				RawName: raw,
			}, nil
		}

		return spec.DefineStruct{RawName: raw}, nil
	case *ast.AnyDataType:
		return nil, ast.SyntaxError(v.Pos(), "unsupported any type")
	case *ast.StructDataType:
		var members []spec.Member
		for _, item := range v.Elements {
			m, err := a.fieldToMember(item)
			if err != nil {
				return nil, err
			}
			members = append(members, m)
		}
		if v.RawText() == "{}" {
			return nil, ast.SyntaxError(v.Pos(), "unsupported empty struct")
		}

		return spec.NestedStruct{
			RawName: v.RawText(),
			Members: members,
		}, nil
	case *ast.InterfaceDataType:
		return spec.InterfaceType{RawName: v.RawText()}, nil
	case *ast.MapDataType:
		if !isLiteralType(v.Key) {
			return nil, ast.SyntaxError(v.Pos(), "expected literal type, got <%T>", v.Key)
		}
		if !v.Key.CanEqual() {
			return nil, ast.SyntaxError(v.Pos(), "map key <%T> must be equal data type", v.Key)
		}
		if v.Value.ContainsStruct() {
			return nil, ast.SyntaxError(v.Pos(), "map value unsupported nested struct")
		}

		value, err := a.astTypeToSpec(v.Value)
		if err != nil {
			return nil, err
		}

		return spec.MapType{
			RawName: v.RawText(),
			Key:     v.Key.RawText(),
			Value:   value,
		}, nil
	case *ast.PointerDataType:
		raw := v.DataType.RawText()
		if IsBaseType(raw) {
			return spec.PointerType{RawName: v.RawText(), Type: spec.PrimitiveType{RawName: raw}}, nil
		}

		value, err := a.astTypeToSpec(v.DataType)
		if err != nil {
			return nil, err
		}

		return spec.PointerType{
			RawName: v.RawText(),
			Type:    value,
		}, nil
	case *ast.ArrayDataType:
		if v.Length.Token.Type == token.ELLIPSIS {
			return nil, ast.SyntaxError(v.Pos(), "array length unsupported dynamic length")
		}
		if v.ContainsStruct() {
			return nil, ast.SyntaxError(v.Pos(), "array elements unsupported nested struct")
		}
		value, err := a.astTypeToSpec(v.DataType)
		if err != nil {
			return nil, err
		}

		return spec.ArrayType{
			RawName: v.RawText(),
			Value:   value,
		}, nil
	case *ast.SliceDataType:
		if v.ContainsStruct() {
			return nil, ast.SyntaxError(v.Pos(), "slice elements unsupported nested struct")
		}

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
	a.fillInfo()

	if err := a.fillTypes(); err != nil {
		return err
	}

	if err := a.fillService(); err != nil {
		return err
	}

	sort.SliceStable(a.spec.Types, func(i, j int) bool {
		return a.spec.Types[i].Name() < a.spec.Types[j].Name()
	})

	groups := make([]spec.Group, 0, len(a.spec.Service.Groups))
	for _, v := range a.spec.Service.Groups {
		sort.SliceStable(v.Routes, func(i, j int) bool {
			return v.Routes[i].Path < v.Routes[j].Path
		})
		groups = append(groups, v)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Annotation.Properties[groupKeyText] < groups[j].Annotation.Properties[groupKeyText]
	})
	a.spec.Service.Groups = groups

	return nil
}

func (a *Analyzer) convertAtDoc(atDoc ast.AtDocStmt) spec.AtDoc {
	var ret spec.AtDoc
	switch val := atDoc.(type) {
	case *ast.AtDocLiteralStmt:
		ret.Text = val.Value.Token.Text
	case *ast.AtDocGroupStmt:
		ret.Properties = a.convertKV(val.Values)
	}
	return ret
}

func (a *Analyzer) convertKV(kv []*ast.KVExpr) map[string]string {
	var ret = map[string]string{}
	for _, v := range kv {
		key := strings.TrimSuffix(v.Key.Token.Text, ":")
		ret[key] = v.Value.RawText()
	}

	return ret
}

func (a *Analyzer) fieldToMember(field *ast.ElemExpr) (spec.Member, error) {
	var name []string
	for _, v := range field.Name {
		name = append(name, v.Token.Text)
	}

	tp, err := a.astTypeToSpec(field.DataType)
	if err != nil {
		return spec.Member{}, err
	}

	head, leading := field.CommentGroup()
	m := spec.Member{
		Name:     strings.Join(name, ", "),
		Type:     tp,
		Docs:     head.List(),
		Comment:  leading.String(),
		IsInline: field.IsAnonymous(),
	}
	if field.Tag != nil {
		m.Tag = field.Tag.Token.Text
	}

	return m, nil
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
			head, leading := astRoute.CommentGroup()
			route := spec.Route{
				Method:  astRoute.Route.Method.Token.Text,
				Path:    astRoute.Route.Path.Format(""),
				Doc:     head.List(),
				Comment: leading.List(),
			}
			if astRoute.AtDoc != nil {
				route.AtDoc = a.convertAtDoc(astRoute.AtDoc)
			}
			if astRoute.AtHandler != nil {
				route.AtDoc = a.convertAtDoc(astRoute.AtDoc)
				route.Handler = astRoute.AtHandler.Name.Token.Text
				head, leading := astRoute.AtHandler.CommentGroup()
				route.HandlerDoc = head.List()
				route.HandlerComment = leading.List()
			}

			if astRoute.Route.Request != nil && astRoute.Route.Request.Body != nil {
				requestType, err := a.getType(astRoute.Route.Request, true)
				if err != nil {
					return err
				}
				route.RequestType = requestType
			}
			if astRoute.Route.Response != nil && astRoute.Route.Response.Body != nil {
				responseType, err := a.getType(astRoute.Route.Response, false)
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

func (a *Analyzer) fillInfo() {
	properties := make(map[string]string)
	if a.api.info != nil {
		for _, kv := range a.api.info.Values {
			key := kv.Key.Token.Text
			properties[strings.TrimSuffix(key, ":")] = kv.Value.RawText()
		}
	}
	a.spec.Info.Properties = properties
	infoKeyValue := make(map[string]string)
	for key, value := range properties {
		titleKey := strings.Title(strings.TrimSuffix(key, ":"))
		infoKeyValue[titleKey] = value
	}
	a.spec.Info.Title = infoKeyValue[infoTitleKey]
	a.spec.Info.Desc = infoKeyValue[infoDescKey]
	a.spec.Info.Version = infoKeyValue[infoVersionKey]
	a.spec.Info.Author = infoKeyValue[infoAuthorKey]
	a.spec.Info.Email = infoKeyValue[infoEmailKey]
}

func (a *Analyzer) fillTypes() error {
	for _, item := range a.api.TypeStmt {
		switch v := (item).(type) {
		case *ast.TypeLiteralStmt:
			if err := a.fillTypeExpr(v.Expr); err != nil {
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
	head, _ := expr.CommentGroup()
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
			RawName: expr.Name.Token.Text,
			Members: members,
			Docs:    head.List(),
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

func (a *Analyzer) getType(expr *ast.BodyStmt, req bool) (spec.Type, error) {
	body := expr.Body
	if req && body.IsArrayType() {
		return nil, ast.SyntaxError(body.Pos(), "request body must be struct")
	}

	var tp spec.Type
	var err error
	var rawText = body.Format("")
	if IsBaseType(body.Value.Token.Text) {
		tp = spec.PrimitiveType{RawName: body.Value.Token.Text}
	} else {
		tp, err = a.findDefinedType(body.Value.Token.Text)
		if err != nil {
			return nil, err
		}
	}
	if body.LBrack != nil {
		if body.Star != nil {
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
	if body.Star != nil {
		return spec.PointerType{
			RawName: rawText,
			Type:    tp,
		}, nil
	}
	return tp, nil
}

// Parse parses the given file and returns the parsed spec.
func Parse(filename string, src interface{}) (*spec.ApiSpec, error) {
	p := New(filename, src)
	ast := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return nil, err
	}

	is := importstack.New()
	err := is.Push(ast.Filename)
	if err != nil {
		return nil, err
	}

	importSet := map[string]lang.PlaceholderType{}
	api, err := convert2API(ast, importSet, is)
	if err != nil {
		return nil, err
	}
	if err := api.SelfCheck(); err != nil {
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
	"any":        placeholder.PlaceHolder,
}

// IsBaseType returns true if the given type is a base type.
func IsBaseType(text string) bool {
	_, ok := kind[text]
	return ok
}
