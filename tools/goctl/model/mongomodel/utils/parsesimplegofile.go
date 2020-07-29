package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"zero/tools/goctl/api/spec"
)

const (
	StructArr = "struct"
	ImportArr = "import"
	Unknown   = "unknown"
)

type Struct struct {
	Name   string
	Fields []spec.Member
}

func readFile(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ParseNetworkGoFile(io string) ([]Struct, []string, error) {
	fset := token.NewFileSet() // 位置是相对于节点

	f, err := parser.ParseFile(fset, "", io, 0)
	if err != nil {
		return nil, nil, err
	}
	return parse(f)
}

func ParseGoFile(pathOrStr string) ([]Struct, []string, error) {
	var goFileStr string
	var err error
	goFileStr, err = readFile(pathOrStr)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet() // 位置是相对于节点

	f, err := parser.ParseFile(fset, "", goFileStr, 0)
	if err != nil {
		return nil, nil, err
	}
	return parse(f)
}

func ParseGoFileByNetwork(io string) ([]Struct, []string, error) {
	fset := token.NewFileSet() // 位置是相对于节点

	f, err := parser.ParseFile(fset, "", io, 0)
	if err != nil {
		return nil, nil, err
	}
	return parse(f)
}

//使用ast包解析golang文件
func parse(f *ast.File) ([]Struct, []string, error) {
	if len(f.Decls) == 0 {
		return nil, nil, fmt.Errorf("you should provide as least 1 struct")
	}
	var structList []Struct
	var importList []string
	for _, decl := range f.Decls {
		structs, imports, err := getStructAndImportInfo(decl)
		if err != nil {
			return nil, nil, err
		}
		structList = append(structList, structs...)
		importList = append(importList, imports...)
	}
	return structList, importList, nil
}

func getStructAndImportInfo(decl ast.Decl) (structs []Struct, imports []string, err error) {
	var structList []Struct
	var importList []string
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil, nil, fmt.Errorf("please input right file")
	}
	for _, tpyObj := range genDecl.Specs {
		switch tpyObj.(type) {
		case *ast.ImportSpec: // import
			importSpec := tpyObj.(*ast.ImportSpec)
			s := importSpec.Path.Value
			importList = append(importList, s)
		case *ast.TypeSpec: //type
			typeSpec := tpyObj.(*ast.TypeSpec)
			switch typeSpec.Type.(type) {
			case *ast.StructType: // struct
				struct1, err := parseStruct(typeSpec)
				if err != nil {
					return nil, nil, err
				}
				structList = append(structList, *struct1)
			}
		default:

		}
	}
	return structList, importList, nil
}

func parseStruct(typeSpec *ast.TypeSpec) (*Struct, error) {
	var result Struct
	structName := typeSpec.Name.Name
	result.Name = structName
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("not struct")
	}

	for _, item := range structType.Fields.List {
		var member spec.Member
		var err error
		member.Name = parseFiledName(item.Names)
		member.Type, err = parseFiledType(item.Type)
		if err != nil {
			return nil, err
		}
		if item.Tag != nil {
			member.Tag = item.Tag.Value
		}
		result.Fields = append(result.Fields, member)
	}
	return &result, nil
}

func parseFiledType(expr ast.Expr) (string, error) {
	switch expr.(type) {
	case *ast.Ident:
		return expr.(*ast.Ident).Name, nil
	case *ast.SelectorExpr:
		selectorExpr := expr.(*ast.SelectorExpr)
		return selectorExpr.X.(*ast.Ident).Name + "." + selectorExpr.Sel.Name, nil
	default:
		return "", fmt.Errorf("can't parse type")
	}
}

func parseFiledName(idents []*ast.Ident) string {
	for _, name := range idents {
		return name.Name
	}
	return ""
}

func UpperCamelToLower(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToLower(name[:1]) + name[1:]
}
