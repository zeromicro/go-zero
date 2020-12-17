package template

import "fmt"

var Vars = fmt.Sprintf(`
var (
	{{.lowerStartCamelObject}}FieldNames          = builderx.FieldNames(&{{.upperStartCamelObject}}{})
	{{.lowerStartCamelObject}}Rows                = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
	{{.lowerStartCamelObject}}RowsExpectAutoSet   = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "%screate_time%s", "%supdate_time%s"), ",")
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s"), "=?,") + "=?"

	{{if .withCache}}{{.cacheKeys}}{{end}}
)
`, "`", "`", "`", "`", "`", "`", "`", "`")
