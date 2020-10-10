package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
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
	flagDot                 = "."
	suffixServer            = "Server"
	referenceContext        = "context"
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

	anyTypeTemplate = "Any struct {\n\tTypeUrl string `json:\"typeUrl\"`\n\tValue   []byte `json:\"value\"`\n}"

	objectM = make(map[string]*Struct)
)

type (
	astParser struct {
		filterStruct map[string]lang.PlaceholderType
		filterEnum   map[string]*Enum
		console.Console
		fileSet *token.FileSet
		proto   *Proto
	}
	Field struct {
		Name     stringx.String
		Type     Type
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
	ConstLit struct {
		Name     stringx.String
		Document []string
		Comment  []string
		Lit      []*Lit
	}
	Lit struct {
		Key   string
		Value int
	}
	Type struct {
		// eg:context.Context
		Expression string
		// eg: *context.Context
		StarExpression string
		// Invoke Type Expression
		InvokeTypeExpression string
		// eg:context
		Package string
		// eg:Context
		Name string
	}
	Func struct {
		Name         stringx.String
		ParameterIn  Type
		ParameterOut Type
		Document     []string
	}
	RpcService struct {
		Name  stringx.String
		Funcs []*Func
	}
	// parsing for rpc
	PbAst struct {
		ContainsAny bool
		Imports     map[string]string
		Structure   map[string]*Struct
		Service     []*RpcService
		*Proto
	}
)

func MustNewAstParser(proto *Proto, log console.Console) *astParser {
	return &astParser{
		filterStruct: proto.Message,
		filterEnum:   proto.Enum,
		Console:      log,
		fileSet:      token.NewFileSet(),
		proto:        proto,
	}
}
func (a *astParser) Parse() (*PbAst, error) {
	var pbAst PbAst
	pbAst.ContainsAny = a.proto.ContainsAny
	pbAst.Proto = a.proto
	pbAst.Structure = make(map[string]*Struct)
	pbAst.Imports = make(map[string]string)
	structure, imports, services, err := a.parse(a.proto.PbSrc)
	if err != nil {
		return nil, err
	}
	dependencyStructure, err := a.parseExternalDependency()
	if err != nil {
		return nil, err
	}
	for k, v := range structure {
		pbAst.Structure[k] = v
	}
	for k, v := range dependencyStructure {
		pbAst.Structure[k] = v
	}
	for key, path := range imports {
		pbAst.Imports[key] = path
	}
	pbAst.Service = append(pbAst.Service, services...)
	return &pbAst, nil
}

func (a *astParser) parse(pbSrc string) (structure map[string]*Struct, imports map[string]string, services []*RpcService, retErr error) {
	structure = make(map[string]*Struct)
	imports = make(map[string]string)
	data, err := ioutil.ReadFile(pbSrc)
	if err != nil {
		retErr = err
		return
	}
	fSet := a.fileSet
	f, err := parser.ParseFile(fSet, "", data, parser.ParseComments)
	if err != nil {
		retErr = err
		return
	}
	commentMap := ast.NewCommentMap(fSet, f, f.Comments)
	f.Comments = commentMap.Filter(f).Comments()
	strucs, function := a.mustScope(f.Scope, a.mustGetIndentName(f.Name))
	for k, v := range strucs {
		if v == nil {
			continue
		}
		structure[k] = v
	}
	importList := f.Imports
	for _, item := range importList {
		name := a.mustGetIndentName(item.Name)
		if item.Path != nil {
			imports[name] = item.Path.Value
		}
	}
	services = append(services, function...)
	return
}
func (a *astParser) parseExternalDependency() (map[string]*Struct, error) {
	m := make(map[string]*Struct)
	for _, impo := range a.proto.Import {
		ret, _, _, err := a.parse(impo.OriginalPbPath)
		if err != nil {
			return nil, err
		}
		for k, v := range ret {
			m[k] = v
		}
	}
	return m, nil
}

func (a *astParser) mustScope(scope *ast.Scope, sourcePackage string) (map[string]*Struct, []*RpcService) {
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
			st, err := a.parseObject(name, v, sourcePackage)
			a.Must(err)
			structs[st.Name.Lower()] = st

		case *ast.InterfaceType:
			if !strings.HasSuffix(name, suffixServer) {
				continue
			}
			list := a.mustServerFunctions(v, sourcePackage)
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

func (a *astParser) mustServerFunctions(v *ast.InterfaceType, sourcePackage string) []*Func {
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
			inList, err := a.parseFields(params.List, true, sourcePackage)
			a.Must(err)

			for _, data := range inList {
				if data.Type.Package == referenceContext {
					continue
				}
				item.ParameterIn = data.Type
				break
			}
		}
		results := v.Results
		if results != nil {
			outList, err := a.parseFields(results.List, true, sourcePackage)
			a.Must(err)

			for _, data := range outList {
				if data.Type.Package == referenceContext {
					continue
				}
				item.ParameterOut = data.Type
				break
			}
		}
		funcs = append(funcs, &item)
	}
	return funcs
}

func (a *astParser) getFieldType(v string, sourcePackage string) Type {
	var pkg, name, expression, starExpression, invokeTypeExpression string

	if strings.Contains(v, ".") {
		starExpression = v
		if strings.Contains(v, "*") {
			leftIndex := strings.Index(v, "*")
			rightIndex := strings.Index(v, ".")
			if leftIndex >= 0 {
				invokeTypeExpression = v[0:leftIndex+1] + v[rightIndex+1:]
			} else {
				invokeTypeExpression = v[rightIndex+1:]
			}
		} else {
			if strings.HasPrefix(v, "map[") || strings.HasPrefix(v, "[]") {
				leftIndex := strings.Index(v, "]")
				rightIndex := strings.Index(v, ".")
				invokeTypeExpression = v[0:leftIndex+1] + v[rightIndex+1:]
			} else {
				rightIndex := strings.Index(v, ".")
				invokeTypeExpression = v[rightIndex+1:]
			}
		}
	} else {
		expression = strings.TrimPrefix(v, flagStar)
		switch v {
		case "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64",
			"bool", "string", "bytes":
			invokeTypeExpression = v
			break
		default:
			name = expression
			invokeTypeExpression = v
			if strings.HasPrefix(v, "map[") || strings.HasPrefix(v, "[]") {
				starExpression = strings.ReplaceAll(v, flagStar, flagStar+sourcePackage+".")
			} else {
				starExpression = fmt.Sprintf("*%v.%v", sourcePackage, name)
				invokeTypeExpression = v
			}

		}
	}
	expression = strings.TrimPrefix(starExpression, flagStar)
	index := strings.LastIndex(expression, flagDot)
	if index > 0 {
		pkg = expression[0:index]
		name = expression[index+1:]
	} else {
		pkg = sourcePackage
	}

	return Type{
		Expression:           expression,
		StarExpression:       starExpression,
		InvokeTypeExpression: invokeTypeExpression,
		Package:              pkg,
		Name:                 name,
	}
}

func (a *astParser) parseObject(structName string, tp *ast.StructType, sourcePackage string) (*Struct, error) {
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
	members, err := a.parseFields(fieldList, false, sourcePackage)
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		var field Field
		field.Name = m.Name
		field.Type = m.Type
		field.JsonTag = m.JsonTag
		field.Document = m.Document
		field.Comment = m.Comment
		st.Field = append(st.Field, &field)
	}
	objectM[structName] = &st
	return &st, nil
}

func (a *astParser) parseFields(fields []*ast.Field, onlyType bool, sourcePackage string) ([]*Field, error) {
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

		item.Type = a.getFieldType(typeName, sourcePackage)
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

func (f *Func) GetDoc() string {
	return strings.Join(f.Document, util.NL)
}

func (f *Func) HaveDoc() bool {
	return len(f.Document) > 0
}

func (a *PbAst) GenEnumCode() (string, error) {
	var element []string
	for _, item := range a.Enum {
		code, err := item.GenEnumCode()
		if err != nil {
			return "", err
		}
		element = append(element, code)
	}
	return strings.Join(element, util.NL), nil
}

func (a *PbAst) GenTypesCode() (string, error) {
	types := make([]string, 0)
	sts := make([]*Struct, 0)
	for _, item := range a.Structure {
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
	types = append(types, a.genAnyCode())
	for _, item := range a.Enum {
		typeCode, err := item.GenEnumTypeCode()
		if err != nil {
			return "", err
		}
		types = append(types, typeCode)
	}

	buffer, err := util.With("type").Parse(typeTemplate).Execute(map[string]interface{}{
		"types": strings.Join(types, util.NL+util.NL),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (a *PbAst) genAnyCode() string {
	if !a.ContainsAny {
		return ""
	}
	return anyTypeTemplate
}

func (s *Struct) genCode(containsTypeStatement bool) (string, error) {
	fields := make([]string, 0)
	for _, f := range s.Field {
		var comment, doc string
		if len(f.Comment) > 0 {
			comment = f.Comment[0]
		}
		doc = strings.Join(f.Document, util.NL)
		buffer, err := util.With(sx.Rand()).Parse(fieldTemplate).Execute(map[string]interface{}{
			"name":       f.Name.Title(),
			"type":       f.Type.InvokeTypeExpression,
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
		"fields": strings.Join(fields, util.NL),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
