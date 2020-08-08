package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genVars(table *InnerTable) (string, error) {
	t, err := template.New("vars").Parse(sqltemplate.Vars)
	if err != nil {
		return "", err
	}
	varBuffer := new(bytes.Buffer)
	m, err := genCacheKeys(table)
	if err != nil {
		return "", err
	}
	keys := make([]string, 0)
	for _, v := range m {
		keys = append(keys, v.Expression)
	}
	err = t.Execute(varBuffer, map[string]interface{}{
		"lowerObject":     table.LowerCamelCase,
		"upperObject":     table.UpperCamelCase,
		"createNotFound":  table.CreateNotFound,
		"keysDefine":      strings.Join(keys, "\r\n"),
		"snakePrimaryKey": table.PrimaryField.SnakeCase,
	})
	if err != nil {
		return "", err
	}
	return varBuffer.String(), nil
}
