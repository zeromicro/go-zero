package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/tal-tech/go-zero/core/lang"
	sx "github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	flagStar                = "*"
	suffixServer            = "Server"
	referenceContext        = "context."
	unknownPrefix           = "XXX_"
	ignoreJsonTagExpression = `json:"-"`
)

var (
	errorParseError = errors.New("pb parse error")
	typeTemplate    = `type (
	{{.types}}
)`
	structTemplate = `{{if .type}}type {{end}}{{.name}} struct {
	{{.fields}}
}`
	fieldTemplate = `{{if .hasDoc}}{{.doc}}
{{end}}{{.name}} {{.type}} {{.tag}}{{if .hasComment}}{{.comment}}{{end}}`
	objectM = make(map[string]*Struct)
)

type (
	astParser struct {
		golang       []byte
		filterStruct map[string]lang.PlaceholderType
		console.Console
		fileSet *token.FileSet
	}
	Field struct {
		Name     stringx.String
		TypeName string
		JsonTag  string
		Document []string
		Comment  []string
	}
	Struct struct {
		Name     stringx.String
		Document []string
		Comment  []string
		Field    []*Field
	}
	Func struct {
		Name        stringx.String
		InType      string
		InTypeName  string // remove *Context,such as LoginRequest、UserRequest
		OutTypeName string // remove *Context
		OutType     string
		Document    []string
	}
	RpcService struct {
		Name  stringx.String
		Funcs []*Func
	}
	// parsing for rpc
	PbAst struct {
		Package string
		// external reference
		Imports map[string]string
		Strcuts map[string]*Struct
		// rpc server's functions,not all functions
		Service []*RpcService
	}
)

func NewAstParser(golang []byte, filterStruct map[string]lang.PlaceholderType, log console.Console) *astParser {
	return &astParser{
		golang:       golang,
		filterStruct: filterStruct,
		Console:      log,
		fileSet:      token.NewFileSet(),
	}
}
func (a *astParser) Parse() (*PbAst, error) {
	fSet := a.fileSet
	f, err := parser.ParseFile(fSet, "", a.golang, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	commentMap := ast.NewCommentMap(fSet, f, f.Comments)
	f.Comments = commentMap.Filter(f).Comments()
	var pbAst PbAst
	pbAst.Package = a.mustGetIndentName(f.Name)
	imports := make(map[string]string)
	for _, item := range f.Imports {
		if item == nil {
			continue
		}
		if item.Path == nil {
			continue
		}
		key := a.mustGetIndentName(item.Name)
		value := item.Path.Value
		imports[key] = value
	}
	structs, funcs := a.mustScope(f.Scope)
	pbAst.Imports = imports
	pbAst.Strcuts = structs
	pbAst.Service = funcs
	return &pbAst, nil
}

func (a *astParser) mustScope(scope *ast.Scope) (map[string]*Struct, []*RpcService) {
	if scope == nil {
		return nil, nil
	}

	objects := scope.Objects
	structs := make(map[string]*Struct)
	serviceList := make([]*RpcService, 0)
	for name, obj := range objects {
		decl := obj.Decl
		if decl == nil {
			continue
		}
		typeSpec, ok := decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		tp := typeSpec.Type

		switch v := tp.(type) {

		case *ast.StructType:
			st, err := a.parseObject(name, v)
			a.Must(err)
			structs[st.Name.Lower()] = st

		case *ast.InterfaceType:
			if !strings.HasSuffix(name, suffixServer) {
				continue
			}
			list := a.mustServerFunctions(v)
			serviceList = append(serviceList, &RpcService{
				Name:  stringx.From(strings.TrimSuffix(name, suffixServer)),
				Funcs: list,
			})
		}
	}
	targetStruct := make(map[string]*Struct)
	for st := range a.filterStruct {
		lower := strings.ToLower(st)
		targetStruct[lower] = structs[lower]
	}
	return targetStruct, serviceList
}

func (a *astParser) mustServerFunctions(v *ast.InterfaceType) []*Func {
	funcs := make([]*Func, 0)
	methodObject := v.Methods
	if methodObject == nil {
		return nil
	}

	for _, method := range methodObject.List {
		var item Func
		name := a.mustGetIndentName(method.Names[0])
		doc := a.parseCommentOrDoc(method.Doc)
		item.Name = stringx.From(name)
		item.Document = doc
		types := method.Type
		if types == nil {
			funcs = append(funcs, &item)
			continue
		}
		v, ok := types.(*ast.FuncType)
		if !ok {
			continue
		}
		params := v.Params
		if params != nil {
			inList, err := a.parseFields(params.List, true)
			a.Must(err)

			for _, data := range inList {
				if strings.HasPrefix(data.TypeName, referenceContext) {
					continue
				}
				// currently,does not support external references
				item.InTypeName = data.TypeName
				item.InType = strings.TrimPrefix(data.TypeName, flagStar)
				break
			}
		}
		results := v.Results
		if results != nil {
			outList, err := a.parseFields(results.List, true)
			a.Must(err)

			for _, data := range outList {
				if strings.HasPrefix(data.TypeName, referenceContext) {
					continue
				}
				// currently,does not support external references
				item.OutTypeName = data.TypeName
				item.OutType = strings.TrimPrefix(data.TypeName, flagStar)
				break
			}
		}
		funcs = append(funcs, &item)
	}
	return funcs
}

func (a *astParser) parseObject(structName string, tp *ast.StructType) (*Struct, error) {
	if data, ok := objectM[structName]; ok {
		return data, nil
	}
	var st Struct
	st.Name = stringx.From(structName)
	if tp == nil {
		return &st, nil
	}

	fields := tp.Fields
	if fields == nil {
		objectM[structName] = &st
		return &st, nil
	}

	fieldList := fields.List
	members, err := a.parseFields(fieldList, false)
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		var field Field
		field.Name = m.Name
		field.TypeName = m.TypeName
		field.JsonTag = m.JsonTag
		field.Document = m.Document
		field.Comment = m.Comment
		st.Field = append(st.Field, &field)
	}
	objectM[structName] = &st
	return &st, nil
}

func (a *astParser) parseFields(fields []*ast.Field, onlyType bool) ([]*Field, error) {
	ret := make([]*Field, 0)
	for _, field := range fields {
		var item Field
		tag := a.parseTag(field.Tag)
		if tag == "" && !onlyType {
			continue
		}
		if tag == ignoreJsonTagExpression {
			continue
		}

		item.JsonTag = tag
		name := a.parseName(field.Names)
		if strings.HasPrefix(name, unknownPrefix) {
			continue
		}
		item.Name = stringx.From(name)
		typeName, err := a.parseType(field.Type)
		if err != nil {
			return nil, err
		}

		item.TypeName = typeName
		if onlyType {
			ret = append(ret, &item)
			continue
		}
		docs := a.parseCommentOrDoc(field.Doc)
		comments := a.parseCommentOrDoc(field.Comment)

		item.Document = docs
		item.Comment = comments

		isInline := name == ""
		if isInline {
			return nil, a.wrapError(field.Pos(), "unexpected inline type:%s", name)
		}

		ret = append(ret, &item)

	}
	return ret, nil
}

func (a *astParser) parseTag(basicLit *ast.BasicLit) string {
	if basicLit == nil {
		return ""
	}
	value := basicLit.Value
	splits := strings.Split(value, " ")
	if len(splits) == 1 {
		return fmt.Sprintf("`%s`", strings.ReplaceAll(splits[0], "`", ""))
	} else {
		return fmt.Sprintf("`%s`", strings.ReplaceAll(splits[1], "`", ""))
	}
}

// returns
// resp1:type's string expression,like int、string、[]int64、map[string]User、*User
// resp2:error
func (a *astParser) parseType(expr ast.Expr) (string, error) {
	if expr == nil {
		return "", errorParseError
	}

	switch v := expr.(type) {
	case *ast.StarExpr:
		stringExpr, err := a.parseType(v.X)
		if err != nil {
			return "", err
		}

		e := fmt.Sprintf("*%s", stringExpr)
		return e, nil

	case *ast.Ident:
		return a.mustGetIndentName(v), nil
	case *ast.MapType:
		keyStringExpr, err := a.parseType(v.Key)
		if err != nil {
			return "", err
		}

		valueStringExpr, err := a.parseType(v.Value)
		if err != nil {
			return "", err
		}

		e := fmt.Sprintf("map[%s]%s", keyStringExpr, valueStringExpr)
		return e, nil
	case *ast.ArrayType:
		stringExpr, err := a.parseType(v.Elt)
		if err != nil {
			return "", err
		}

		e := fmt.Sprintf("[]%s", stringExpr)
		return e, nil
	case *ast.InterfaceType:
		return "interface{}", nil
	case *ast.SelectorExpr:
		join := make([]string, 0)
		xIdent, ok := v.X.(*ast.Ident)
		xIndentName := a.mustGetIndentName(xIdent)
		if ok {
			join = append(join, xIndentName)
		}
		sel := v.Sel
		join = append(join, a.mustGetIndentName(sel))
		return strings.Join(join, "."), nil
	case *ast.ChanType:
		return "", a.wrapError(v.Pos(), "unexpected type 'chan'")
	case *ast.FuncType:
		return "", a.wrapError(v.Pos(), "unexpected type 'func'")
	case *ast.StructType:
		return "", a.wrapError(v.Pos(), "unexpected inline struct type")
	default:
		return "", a.wrapError(v.Pos(), "unexpected type '%v'", v)
	}
}
func (a *astParser) parseName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	name := names[0]
	return a.mustGetIndentName(name)
}

func (a *astParser) parseCommentOrDoc(cg *ast.CommentGroup) []string {
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

func (a *astParser) mustGetIndentName(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func (a *astParser) wrapError(pos token.Pos, format string, arg ...interface{}) error {
	file := a.fileSet.Position(pos)
	return fmt.Errorf("line %v: %s", file.Line, fmt.Sprintf(format, arg...))
}

func (a *PbAst) GenTypesCode() (string, error) {
	types := make([]string, 0)
	sts := make([]*Struct, 0)
	for _, item := range a.Strcuts {
		sts = append(sts, item)
	}
	sort.Slice(sts, func(i, j int) bool {
		return sts[i].Name.Source() < sts[j].Name.Source()
	})
	for _, s := range sts {
		structCode, err := s.genCode(false)
		if err != nil {
			return "", err
		}

		if structCode == "" {
			continue
		}
		types = append(types, structCode)
	}
	buffer, err := util.With("type").Parse(typeTemplate).Execute(map[string]interface{}{
		"types": strings.Join(types, "\n\n"),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (s *Struct) genCode(containsTypeStatement bool) (string, error) {
	if len(s.Field) == 0 {
		return "", nil
	}
	fields := make([]string, 0)
	for _, f := range s.Field {
		var comment, doc string
		if len(f.Comment) > 0 {
			comment = f.Comment[0]
		}
		doc = strings.Join(f.Document, "\n")
		buffer, err := util.With(sx.Rand()).Parse(fieldTemplate).Execute(map[string]interface{}{
			"name":       f.Name.Title(),
			"type":       f.TypeName,
			"tag":        f.JsonTag,
			"hasDoc":     len(f.Document) > 0,
			"doc":        doc,
			"hasComment": len(f.Comment) > 0,
			"comment":    comment,
		})
		if err != nil {
			return "", err
		}

		fields = append(fields, buffer.String())
	}
	buffer, err := util.With("struct").Parse(structTemplate).Execute(map[string]interface{}{
		"type":   containsTypeStatement,
		"name":   s.Name.Title(),
		"fields": strings.Join(fields, "\n"),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
