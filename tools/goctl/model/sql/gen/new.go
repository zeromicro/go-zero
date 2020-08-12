package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genNew(table Table) (string, error) {
	output, err := templatex.With("new").
		Parse(template.New).
		Execute(map[string]interface{}{
			"upperStartCamelObject": table.Name.Snake2Camel(),
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
