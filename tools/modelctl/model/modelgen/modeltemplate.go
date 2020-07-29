package modelgen

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"zero/tools/modelctl/model"
	"zero/tools/modelctl/util"
)

type (
	Query struct {
		Rows   string
		Args   string
		Values string
	}
	Update struct {
		Rows   string
		Values string
	}
	Template struct {
		Package                  string
		PrimaryKeyField          string
		PrimaryKeyFieldCamel     string
		PrimaryKeyType           string
		ModelCamelWithLowerStart string
		ModelLowerSplitByPound   string
		ModelCamelWithUpperStart string
		Fields                   string
		Abbr                     string
		WithCache                bool
		Insert                   Query
		// 包含where语句
		Update Update
	}

	StructField struct {
		// 字段名称,下划线
		NameWithUnderline string
		// 字段名称,驼峰式,大写开头
		NameCamelWithUpperStart string
		// 字段数据类型
		DataType string
		// 字段注释
		Comment string
		// 是否为主键
		PrimaryKey bool
	}
)

const (
	fieldTemplateText = "{{.NameCamelWithUpperStart}} {{.DataType}} `db:\"{{.NameWithUnderline}}\" json:\"{{.NameWithUnderline}},omitempty\"` {{.Comment}}"
	modelTemplateText = `package {{.Package}}

import (
    "strings"
    "time"

    "zero/core/stores/redis"
    "zero/core/stores/sqlx"
    "zero/service/shared/builderx"
    "zero/service/shared/cache"

)
var (
    {{.ModelCamelWithLowerStart}}QueryRows = strings.Join(builderx.FieldNames(&{{.ModelCamelWithUpperStart}}{}), ",")
{{if .WithCache}}{{.ModelCamelWithLowerStart}}CachePrefix = "xjy#{{.ModelLowerSplitByPound}}#"
    {{.ModelCamelWithLowerStart}}Expire = 7 * 24 * 3600{{end}}
)

type (
    {{.ModelCamelWithUpperStart}}Model struct {
       {{if .WithCache}} *CachedModel{{else}}
		table string
    	conn  sqlx.SqlConn{{end}}
    }
    {{.ModelCamelWithUpperStart}} struct {
        {{.Fields}}
    }
)

func New{{.ModelCamelWithUpperStart}}Model(table string, conn sqlx.SqlConn{{if .WithCache}}, rds *redis.Redis{{end}}) *{{.ModelCamelWithUpperStart}}Model {
	{{if .WithCache}} return &{{.ModelCamelWithUpperStart}}Model{NewCachedModel(conn, table, rds)}
	 {{else}}return &{{.ModelCamelWithUpperStart}}Model{table:table,conn:conn}{{end}}
}

func ({{.Abbr}} *{{.ModelCamelWithUpperStart}}Model) Insert(data *{{.ModelCamelWithUpperStart}}) error {
    querySql:="insert into "+{{.Abbr}}.table+" ({{.Insert.Rows}}) value ({{.Insert.Args}})"
    _, err := {{.Abbr}}.conn.Exec(querySql, {{.Insert.Values}})
    return err
}

func ({{.Abbr}} *{{.ModelCamelWithUpperStart}}Model) Update(data *{{.ModelCamelWithUpperStart}}) error {
{{if .WithCache}}err := {{.Abbr}}.cleanCache(data.{{.PrimaryKeyField}})
    if err != nil {
        return err
    }
    querySql := "update " + {{.Abbr}}.table + " set {{.Update.Rows}}"
    _, err = {{.Abbr}}.conn.Exec(querySql,{{.Update.Values}} )
    return err
{{else}}querySql := "update " + {{.Abbr}}.table + " set {{.Update.Rows}}"
    _, err := {{.Abbr}}.conn.Exec(querySql,{{.Update.Values}} )
    return err{{end}}
}

func ({{.Abbr}} *{{.ModelCamelWithUpperStart}}Model) FindOne({{.PrimaryKeyFieldCamel}} {{.PrimaryKeyType}})(*{{.ModelCamelWithUpperStart}},error){
    querySql:="select "+{{.ModelCamelWithLowerStart}}QueryRows+" from "+{{.Abbr}}.table+" where {{.PrimaryKeyFieldCamel}} = ? limit 1"
    var resp {{.ModelCamelWithUpperStart}}
{{if .WithCache}}key := cache.FormatKey({{.ModelCamelWithLowerStart}}CachePrefix,{{.PrimaryKeyFieldCamel}})
    err := {{.Abbr}}.QueryRow(&resp, key, {{.ModelCamelWithLowerStart}}Expire, func(conn sqlx.Session, v interface{}) error {
        return conn.QueryRow(v, querySql, {{.PrimaryKeyFieldCamel}})
    })
    if err != nil {
        if err == sqlx.ErrNotFound {
            return nil, ErrNotFound
        }
        return nil, err
    }
{{else}}err := {{.Abbr}}.conn.QueryRow(&resp, querySql, {{.PrimaryKeyFieldCamel}})
    if err != nil {
        if err == sqlx.ErrNotFound {
            return nil, ErrNotFound
        }
        return nil, err
    }{{end}}
    return &resp, nil
}

func ({{.Abbr}} *{{.ModelCamelWithUpperStart}}Model) DeleteById({{.PrimaryKeyFieldCamel}} {{.PrimaryKeyType}}) error {
{{if .WithCache}}err := {{.Abbr}}.cleanCache({{.PrimaryKeyFieldCamel}})
    if err != nil {
        return err
    }
    querySql := "delete from " + {{.Abbr}}.table + " where {{.PrimaryKeyFieldCamel}} = ? "
    _, err = {{.Abbr}}.conn.Exec(querySql, {{.PrimaryKeyFieldCamel}})
    return err
{{else}}querySql := "delete from " + {{.Abbr}}.table + " where {{.PrimaryKeyFieldCamel}} = ? "
    _, err := {{.Abbr}}.conn.Exec(querySql, {{.PrimaryKeyFieldCamel}})
    return err{{end}}
}

{{if .WithCache}}
func ({{.Abbr}} *{{.ModelCamelWithUpperStart}}Model) cleanCache({{.PrimaryKeyFieldCamel}} {{.PrimaryKeyType}}) error {
    key := cache.FormatKey({{.ModelCamelWithLowerStart}}CachePrefix,{{.PrimaryKeyFieldCamel}})
    _, err := {{.Abbr}}.rds.Del(key)
    return err
}
{{end}}

`
)

func generateModelTemplate(packageName, table string, fileds []*Field) (*Template, error) {
	list := make([]*StructField, 0)
	var containsPrimaryKey bool
	for _, item := range fileds {
		goType, ok := model.CommonMysqlDataTypeMap[item.Type]
		if !ok {
			return nil, errors.New(fmt.Sprintf("table:%v,the data type %v of mysql does not match", table, item.Type))
		}
		if !containsPrimaryKey {
			containsPrimaryKey = item.Primary == "PRI"
		}
		list = append(list, &StructField{
			NameWithUnderline:       item.Name,
			NameCamelWithUpperStart: util.FmtUnderLine2Camel(item.Name, true),
			DataType:                goType,
			Comment:                 item.Comment,
			PrimaryKey:              item.Primary == "PRI",
		})
	}
	if !containsPrimaryKey {
		return nil, errors.New(fmt.Sprintf("table:%v,primary key does not exist", table))
	}
	var structBuffer, insertRowsBuffer, insertArgBuffer, insertValuesBuffer, updateRowsBuffer, updateValuesBuffer bytes.Buffer
	var primaryField *StructField
	for index, item := range list {
		out, err := convertField(item)
		if err != nil {
			return nil, err
		}
		structBuffer.WriteString(out + "\n")
		if !item.PrimaryKey {
			insertRowsBuffer.WriteString(item.NameWithUnderline)
			insertArgBuffer.WriteString("?")
			insertValuesBuffer.WriteString("data." + item.NameCamelWithUpperStart)

			updateRowsBuffer.WriteString(item.NameWithUnderline + "=?")
			updateValuesBuffer.WriteString("data." + item.NameCamelWithUpperStart)

			if index < len(list)-1 {
				insertRowsBuffer.WriteString(",")
				insertArgBuffer.WriteString(",")
				insertValuesBuffer.WriteString(",")

				updateRowsBuffer.WriteString(",")
			}
			updateValuesBuffer.WriteString(",")
		} else {
			primaryField = item
		}
	}

	updateRowsBuffer.WriteString(" where " + primaryField.NameWithUnderline + "=?")
	updateValuesBuffer.WriteString(" data." + primaryField.NameCamelWithUpperStart)
	modelSplitByPoundArr := strings.Replace(table, "_", "#", -1)
	templateStruct := Template{
		Package:                  packageName,
		PrimaryKeyField:          primaryField.NameCamelWithUpperStart,
		PrimaryKeyFieldCamel:     primaryField.NameWithUnderline,
		PrimaryKeyType:           primaryField.DataType,
		ModelCamelWithLowerStart: util.FmtUnderLine2Camel(table, false),
		ModelLowerSplitByPound:   modelSplitByPoundArr,
		ModelCamelWithUpperStart: util.FmtUnderLine2Camel(table, true),
		Fields:                   structBuffer.String(),
		Abbr:                     util.Abbr(table) + "m",
		Insert: Query{
			Rows:   insertRowsBuffer.String(),
			Args:   insertArgBuffer.String(),
			Values: insertValuesBuffer.String(),
		},
		Update: Update{
			Rows:   updateRowsBuffer.String(),
			Values: updateValuesBuffer.String(),
		},
	}
	return &templateStruct, nil
}

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
