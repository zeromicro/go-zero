package sqlmodel

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

var (
	commonTemplate = `
var (
	{{.modelWithLowerStart}}FieldNames = builderx.FieldNames(&{{.model}}{})
	{{.modelWithLowerStart}}Rows = strings.Join({{.modelWithLowerStart}}FieldNames, ",")
	{{.modelWithLowerStart}}RowsExpectAutoSet = strings.Join(stringx.Remove({{.modelWithLowerStart}}FieldNames, {{.expected}}), ",")
	{{.modelWithLowerStart}}RowsWithPlaceHolder = strings.Join(stringx.Remove({{.modelWithLowerStart}}FieldNames, {{.expected}}), "=?,") + "=?"
)
`
)

type (
	fieldExp struct {
		name string
		tag  string
		ty   string
	}
	structExp struct {
		genStruct       GenStruct
		name            string
		idAutoIncrement bool
		Fields          []fieldExp
		primaryKey      string
		conditions      []string
		ignoreFields    []string
	}
)

func NewStructExp(s GenStruct) (*structExp, error) {
	src := s.TableStruct
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)

	if err != nil {
		return nil, err
	}
	if len(file.Decls) == 0 {
		return nil, sqlError("no struct inspect")
	}
	typeDecl := file.Decls[0].(*ast.GenDecl)

	if len(typeDecl.Specs) == 0 {
		return nil, sqlError("no type specs")
	}
	typeSpec := typeDecl.Specs[0].(*ast.TypeSpec)
	structDecl := typeSpec.Type.(*ast.StructType)

	var strExp structExp
	strExp.genStruct = s
	strExp.primaryKey = s.PrimaryKey
	strExp.name = typeSpec.Name.Name
	fields := structDecl.Fields.List
	for _, field := range fields {
		typeExpr := field.Type

		start := typeExpr.Pos() - 1
		end := typeExpr.End() - 1

		typeInSource := src[start:end]
		if len(field.Names) == 0 {
			return nil, sqlError("field name empty")
		}
		var name = field.Names[0].Name
		var tag = util.Untitle(name)
		if field.Tag != nil {
			tag = src[field.Tag.Pos():field.Tag.End()]
		}
		dbTag := getDBColumnName(tag, name)
		strExp.Fields = append(strExp.Fields, fieldExp{
			name: name,
			ty:   typeInSource,
			tag:  dbTag,
		})
		if dbTag == strExp.primaryKey && typeInSource == "int64" {
			strExp.ignoreFields = append(strExp.ignoreFields, dbTag)
			strExp.idAutoIncrement = true
		}
		if name == "UpdateTime" || name == "CreateTime" {
			strExp.ignoreFields = append(strExp.ignoreFields, dbTag)
		}
	}
	return &strExp, nil
}

func (s *structExp) genMysqlCRUD() (string, error) {
	commonStr, err := s.genCommon()
	if err != nil {
		return "", err
	}
	insertStr, err := s.genInsert()
	if err != nil {
		return "", err
	}
	updateStr, err := s.genUpdate()
	if err != nil {
		return "", err
	}

	deleteStr, err := s.genDelete()
	if err != nil {
		return "", err
	}

	queryOneStr, err := s.genQueryOne()
	if err != nil {
		return "", err
	}

	queryListStr, err := s.genQueryList()
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"package model \n", commonStr, s.genStruct.TableModel, queryOneStr, queryListStr, deleteStr, insertStr, updateStr}, "\n"), nil
}

func getDBColumnName(tag, name string) string {
	matches := spec.TagRe.FindStringSubmatch(tag)
	for i := range matches {
		name := spec.TagSubNames[i]
		if name == "name" {
			return matches[i]
		}
	}

	return util.Untitle(name)
}

func (s *structExp) genCommon() (string, error) {
	templateName := commonTemplate
	t := template.Must(template.New("commonTemplate").Parse(templateName))
	var tmplBytes bytes.Buffer
	var ignoreFieldsQuota []string
	for _, item := range s.ignoreFields {
		ignoreFieldsQuota = append(ignoreFieldsQuota, fmt.Sprintf("\"%s\"", item))
	}
	err := t.Execute(&tmplBytes, map[string]string{
		"model":               s.name,
		"expected":            strings.Join(ignoreFieldsQuota, ", "),
		"modelWithLowerStart": fmtUnderLine2Camel(s.name, false),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}

func (s *structExp) buildCondition() (string, string) {
	var conditionExp []string
	var valueConditions []string
	for _, field := range s.Fields {
		if stringx.Contains(s.conditions, strings.ToLower(field.tag)) ||
			stringx.Contains(s.conditions, strings.ToLower(field.name)) {
			conditionExp = append(conditionExp, fmt.Sprintf("%s %s", util.Untitle(field.name), field.ty))
			valueConditions = append(valueConditions, fmt.Sprintf("%s = ?", field.tag))
		}
	}
	return strings.Join(conditionExp, ", "), strings.Join(valueConditions, " and ")
}

func (s *structExp) conditionNames() []string {
	var conditionExp []string
	for _, field := range s.Fields {
		if stringx.Contains(s.conditions, strings.ToLower(field.tag)) ||
			stringx.Contains(s.conditions, strings.ToLower(field.name)) {
			conditionExp = append(conditionExp, util.Untitle(field.name))
		}
	}
	return conditionExp
}
