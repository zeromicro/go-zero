package template

import "fmt"

// Vars defines a template for var block in model
var Vars = fmt.Sprintf(`
var (
	{{.lowerStartCamelObject}}FieldNames          = builderx.RawFieldNames(&{{.upperStartCamelObject}}{})
	{{.lowerStartCamelObject}}Rows                = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
	{{.lowerStartCamelObject}}RowsExpectAutoSet   = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "%screate_time%s", "%supdate_time%s"), ",")
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = {{if .postgreSql}}builderx.PostgreSqlJoin(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s")){{else}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s"), "=?,") + "=?"{{end}}

	{{if .withCache}}{{.cacheKeys}}{{end}}
)
`, "`", "`", "`", "`", "`", "`", "`", "`")
