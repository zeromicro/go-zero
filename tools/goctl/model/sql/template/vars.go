package template

import "fmt"

// Vars defines a template for var block in model
var Vars = fmt.Sprintf(`
var (
	{{.lowerStartCamelObject}}FieldNames          = builderx.RawFieldNames(&{{.upperStartCamelObject}}{}{{if .postgreSql}},true{{end}})
	{{.lowerStartCamelObject}}Rows                = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
	{{.lowerStartCamelObject}}RowsExpectAutoSet   = {{if .postgreSql}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "%screate_time%s", "%supdate_time%s"), ","){{else}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "%screate_time%s", "%supdate_time%s"), ","){{end}}
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = {{if .postgreSql}}builderx.PostgreSqlJoin(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s")){{else}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s"), "=?,") + "=?"{{end}}

	{{if .withCache}}{{.cacheKeys}}{{end}}
)
`, "", "", "", "", // postgreSql mode
	"`", "`", "`", "`",
	"", "", "", "", // postgreSql mode
	"`", "`", "`", "`")
