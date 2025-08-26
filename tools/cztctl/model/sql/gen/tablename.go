package gen

import (
	"github.com/lerity-yao/go-zero/tools/cztctl/model/sql/template"
	"github.com/lerity-yao/go-zero/tools/cztctl/util"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
)

func genTableName(table Table) (string, error) {
	text, err := pathx.LoadTemplate(category, tableNameTemplateFile, template.TableName)
	if err != nil {
		return "", err
	}

	output, err := util.With("tableName").
		Parse(text).
		Execute(map[string]any{
			"tableName":             table.Name.Source(),
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
