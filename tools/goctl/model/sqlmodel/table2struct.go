package sqlmodel

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	util2 "github.com/tal-tech/go-zero/tools/modelctl/util"
)

var (
	commonMysqlDataTypeMap = map[string]string{
		"tinyint":    "int64",
		"smallint":   "int64",
		"mediumint":  "int64",
		"int":        "int64",
		"integer":    "int64",
		"bigint":     "int64",
		"float":      "float64",
		"double":     "float64",
		"decimal":    "float64",
		"date":       "time.Time",
		"time":       "string",
		"year":       "int64",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"char":       "string",
		"varchar":    "string",
		"tinyblob":   "string",
		"tinytext":   "string",
		"blob":       "string",
		"text":       "string",
		"mediumblob": "string",
		"mediumtext": "string",
		"longblob":   "string",
		"longtext":   "string",
	}
)

var modelTemplate = `
type ( 

	{{.ModelCamelWithUpperStart}}Model struct {
		table string
		conn  sqlx.SqlConn
	}

	{{.ModelCamelWithUpperStart}} struct {
		{{.Fields}}
	}
)

func New{{.ModelCamelWithUpperStart}}Model(table string, conn sqlx.SqlConn) *{{.ModelCamelWithUpperStart}}Model {
	return &{{.ModelCamelWithUpperStart}}Model{table: table, conn: conn}
}
`

const (
	fieldTemplateText = "{{.NameCamelWithUpperStart}} {{.DataType}} `db:\"{{.NameWithUnderline}}\"` {{.Comment}}"
)

type (
	Template struct {
		ModelCamelWithUpperStart string
		Fields                   string
	}

	StructField struct {
		// 字段名称,下划线
		NameWithUnderline string
		// 字段名称,驼峰式,大写开头
		NameCamelWithUpperStart string
		// 字段名称,驼峰式,小写开头
		NameCamelWithLowerStart string
		// 字段数据类型
		DataType string
		// 字段注释
		Comment string
	}

	GenStruct struct {
		// 表对应的struct
		TableStruct string `json:"tableStruct"`
		// 表对应生成的model信息，参考模板${modelTemplate}
		TableModel string `json:"tableModel"`
		// 主键
		PrimaryKey string `json:"primaryKey"`
	}
)

func generateTypeModel(table string, fields []*Column) (*GenStruct, error) {
	var resp GenStruct
	structString, fieldsString, err := convertStruct(table, fields)
	if err != nil {
		return nil, err
	}
	templateStruct := Template{
		ModelCamelWithUpperStart: util2.FmtUnderLine2Camel(table, true),
		Fields:                   fieldsString,
	}
	tl, err := template.New("").Parse(modelTemplate)
	if err != nil {
		return nil, err
	}
	var resultBuffer = bytes.NewBufferString("")
	err = tl.Execute(resultBuffer, templateStruct)
	if err != nil {
		return nil, err
	}
	resp.TableStruct = structString
	resp.TableModel = resultBuffer.String()
	return &resp, nil
}

// returns struct、fields、error
func convertStruct(table string, columns []*Column) (string, string, error) {
	var structBuffer, fieldsBuffer bytes.Buffer
	structBuffer.WriteString("package model \n\n")
	structBuffer.WriteString("type " + fmtUnderLine2Camel(table, true) + " struct {\n")
	for index, item := range columns {
		goType, ok := commonMysqlDataTypeMap[item.DataType]
		if !ok {
			return "", "", errors.New(fmt.Sprintf("table: %s,the data type %s of %s does not match", table, item.DataType, item.Name))
		}
		out, err := convertField(&StructField{
			NameWithUnderline:       item.Name,
			NameCamelWithUpperStart: fmtUnderLine2Camel(item.Name, true),
			NameCamelWithLowerStart: fmtUnderLine2Camel(item.Name, false),
			DataType:                goType,
			Comment:                 item.Comment,
		})
		if err != nil {
			return "", "", err
		}
		structBuffer.WriteString("\t" + out)
		structBuffer.WriteString("\n")

		if index == 0 {
			fieldsBuffer.WriteString(out)
		} else {
			fieldsBuffer.WriteString("\t\t" + out)
		}
		if index < len(columns)-1 {
			fieldsBuffer.WriteString("\n")
		}

	}
	structBuffer.WriteString("}")
	return structBuffer.String(), fieldsBuffer.String(), nil
}

// column转换成struct field
func convertField(field *StructField) (string, error) {
	if strings.TrimSpace(field.Comment) != "" {
		field.Comment = "// " + field.Comment
	}
	tl, err := template.New("").Parse(fieldTemplateText)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	err = tl.Execute(buf, field)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 简单的下划线转驼峰格式
func fmtUnderLine2Camel(in string, upperStart bool) string {
	if strings.TrimSpace(in) == "" {
		return ""
	}
	var words []string
	if strings.Contains(in, "_") {
		words = strings.Split(in, "_")
		if len(words) == 0 {
			return ""
		}
	}
	if len(words) == 0 {
		if !upperStart {
			bts := []byte(in)
			r := bytes.ToLower([]byte{bts[0]})
			bts[0] = r[0]
			return string(bts)
		} else {
			return strings.Title(in)
		}
	}
	var buffer bytes.Buffer
	for index, word := range words {
		if strings.TrimSpace(word) == "" {
			continue
		}
		bts := []byte(word)
		if index == 0 && !upperStart {
			bts[0] = bytes.ToLower([]byte{bts[0]})[0]
			buffer.Write(bts)
			continue
		}
		bts = bytes.Title(bts)
		buffer.Write(bts)
	}
	return buffer.String()
}
