package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

var (
	ErrStructNotFound      = errors.New("struct not found")
	ErrUnSupportInlineType = errors.New("unsupport inline type")
	interfaceExpr          = `interface{}`
	objectM                = make(map[string]*spec.Type)
)

const (
	golangF = `package ast
	%s
`
	pkgPrefix = "package"
)

type StructParser struct {
	Src string
}

func (sp *StructParser) Parse() ([]spec.Type, error) {
	if !strings.HasPrefix(sp.Src, pkgPrefix) {
		sp.Src = fmt.Sprintf(golangF, sp.Src)
	}

	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", sp.Src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	commentMap := ast.NewCommentMap(fSet, f, f.Comments)
	f.Comments = commentMap.Filter(f).Comments()
	scope := f.Scope
	if scope == nil {
		return nil, ErrStructNotFound
	}
	objects := scope.Objects
	structs := make([]*spec.Type, 0)
	for structName, obj := range objects {
		st, err := sp.parseObject(structName, obj)
		if err != nil {
			return nil, err
		}
		structs = append(structs, st)
	}
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Name < structs[j].Name
	})
	resp := make([]spec.Type, 0)
	for _, item := range structs {
		resp = append(resp, *item)
	}
	return resp, nil
}

func (sp *StructParser) parseObject(structName string, obj *ast.Object) (*spec.Type, error) {
	if data, ok := objectM[structName]; ok {
		return data, nil
	}
	var st spec.Type
	st.Name = structName
	if obj.Decl == nil {
		objectM[structName] = &st
		return &st, nil
	}
	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		objectM[structName] = &st
		return &st, nil
	}
	if decl.Type == nil {
		objectM[structName] = &st
		return &st, nil
	}
	tp, ok := decl.Type.(*ast.StructType)
	if !ok {
		objectM[structName] = &st
		return &st, nil
	}
	fields := tp.Fields
	if fields == nil {
		objectM[structName] = &st
		return &st, nil
	}
	fieldList := fields.List
	members, err := sp.parseFields(fieldList)
	if err != nil {
		return nil, err
	}
	st.Members = members
	objectM[structName] = &st
	return &st, nil
}

func (sp *StructParser) parseFields(fields []*ast.Field) ([]spec.Member, error) {
	members := make([]spec.Member, 0)
	for _, field := range fields {
		docs := parseCommentOrDoc(field.Doc)
		comments := parseCommentOrDoc(field.Comment)
		name := parseName(field.Names)
		tp, stringExpr, err := sp.parseType(field.Type)
		if err != nil {
			return nil, err
		}
		tag := parseTag(field.Tag)
		isInline := name == ""
		if isInline {
			var err error
			name, err = sp.getInlineName(tp)
			if err != nil {
				return nil, err
			}
		}
		members = append(members, spec.Member{
			Name:     name,
			Type:     stringExpr,
			Expr:     tp,
			Tag:      tag,
			Comments: comments,
			Docs:     docs,
			IsInline: isInline,
		})

	}
	return members, nil
}

func (sp *StructParser) getInlineName(tp interface{}) (string, error) {
	switch v := tp.(type) {
	case *spec.Type:
		return v.Name, nil
	case *spec.PointerType:
		return sp.getInlineName(v.Star)
	case *spec.StructType:
		return v.StringExpr, nil
	default:
		return "", ErrUnSupportInlineType
	}
}

func (sp *StructParser) getInlineTypePrefix(tp interface{}) (string, error) {
	if tp == nil {
		return "", nil
	}
	switch tp.(type) {
	case *ast.Ident:
		return "", nil
	case *ast.StarExpr:
		return "*", nil
	case *ast.TypeSpec:
		return "", nil
	default:
		return "", ErrUnSupportInlineType
	}
}

func parseTag(basicLit *ast.BasicLit) string {
	if basicLit == nil {
		return ""
	}
	return basicLit.Value
}

// returns
// resp1: type can convert to *spec.PointerType|*spec.BasicType|*spec.MapType|*spec.ArrayType|*spec.InterfaceType
// resp2: type's string expression,like int、string、[]int64、map[string]User、*User
// resp3: error
func (sp *StructParser) parseType(expr ast.Expr) (interface{}, string, error) {
	if expr == nil {
		return nil, "", errors.New("parse error " + sp.Src)
	}
	exprStr := sp.Src[expr.Pos():expr.End()]
	switch v := expr.(type) {
	case *ast.StarExpr:
		star, stringExpr, err := sp.parseType(v.X)
		if err != nil {
			return nil, "", err
		}
		e := fmt.Sprintf("*%s", stringExpr)
		return &spec.PointerType{Star: star, StringExpr: e}, e, nil
	case *ast.Ident:
		if isBasicType(v.Name) {
			return &spec.BasicType{Name: v.Name, StringExpr: v.Name}, v.Name, nil
		} else if v.Obj != nil {
			obj := v.Obj
			if obj.Name != v.Name { // 防止引用自己而无限递归
				specType, err := sp.parseObject(v.Name, v.Obj)
				if err != nil {
					return nil, "", err
				} else {
					return specType, v.Obj.Name, nil
				}
			} else {
				inlineType, err := sp.getInlineTypePrefix(obj.Decl)
				if err != nil {
					return nil, "", err
				}
				return &spec.StructType{
					StringExpr: fmt.Sprintf("%s%s", inlineType, v.Name),
				}, v.Name, nil
			}
		} else {
			return nil, "", fmt.Errorf(" [%s] - member is not exist, expr is %s", v.Name, exprStr)
		}
	case *ast.MapType:
		key, keyStringExpr, err := sp.parseType(v.Key)
		if err != nil {
			return nil, "", err
		}

		value, valueStringExpr, err := sp.parseType(v.Value)
		if err != nil {
			return nil, "", err
		}

		keyType, ok := key.(*spec.BasicType)
		if !ok {
			return nil, "", fmt.Errorf("[%+v] - unsupported type of map key, expr is  %s", v.Key, exprStr)
		}

		e := fmt.Sprintf("map[%s]%s", keyStringExpr, valueStringExpr)
		return &spec.MapType{
			Key:        keyType.Name,
			Value:      value,
			StringExpr: e,
		}, e, nil
	case *ast.ArrayType:
		arrayType, stringExpr, err := sp.parseType(v.Elt)
		if err != nil {
			return nil, "", err
		}

		e := fmt.Sprintf("[]%s", stringExpr)
		return &spec.ArrayType{ArrayType: arrayType, StringExpr: e}, e, nil
	case *ast.InterfaceType:
		return &spec.InterfaceType{StringExpr: interfaceExpr}, interfaceExpr, nil
	case *ast.ChanType:
		return nil, "", errors.New("[chan] - unsupported type, expr is " + exprStr)
	case *ast.FuncType:
		return nil, "", errors.New("[func] - unsupported type, expr is " + exprStr)
	case *ast.StructType: // todo can optimize
		return nil, "", errors.New("[struct] - unsupported inline struct type, expr is " + exprStr)
	case *ast.SelectorExpr:
		x := v.X
		sel := v.Sel
		xIdent, ok := x.(*ast.Ident)
		if ok {
			name := xIdent.Name
			if name != "time" && sel.Name != "Time" {
				return nil, "", fmt.Errorf("[outter package] - package: %s, unsupport type", exprStr)
			}

			tm := fmt.Sprintf("time.Time")
			return &spec.TimeType{
				StringExpr: tm,
			}, tm, nil
		}
		return nil, "", errors.New("parse error " + exprStr)
	default:
		return nil, "", errors.New("parse error " + exprStr)
	}
}

func isBasicType(tp string) bool {
	switch tp {
	case
		"bool",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"complex64",
		"complex128",
		"string",
		"int",
		"uint",
		"uintptr",
		"byte",
		"rune",
		"Type",
		"Type1",
		"IntegerType",
		"FloatType",
		"ComplexType":
		return true
	default:
		return false
	}
}
func parseName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	name := names[0]
	return parseIdent(name)
}

func parseIdent(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func parseCommentOrDoc(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}
	comments := make([]string, 0)
	for _, comment := range cg.List {
		if comment == nil {
			continue
		}
		text := strings.TrimSpace(comment.Text)
		if text == "" {
			continue
		}
		comments = append(comments, text)
	}
	return comments
}
