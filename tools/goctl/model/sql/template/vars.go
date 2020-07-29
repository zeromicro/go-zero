package sqltemplate

var Vars = `
var (
	{{.lowerObject}}FieldNames          = builderx.FieldNames(&{{.upperObject}}{})
	{{.lowerObject}}Rows                = strings.Join({{.lowerObject}}FieldNames, ",")
	{{.lowerObject}}RowsExpectAutoSet   = strings.Join(stringx.Remove({{.lowerObject}}FieldNames, "{{.snakePrimaryKey}}", "create_time", "update_time"), ",")
	{{.lowerObject}}RowsWithPlaceHolder = strings.Join(stringx.Remove({{.lowerObject}}FieldNames, "{{.snakePrimaryKey}}", "create_time", "update_time"), "=?,") + "=?"

	{{.keysDefine}}

	{{if .createNotFound}}ErrNotFound               = sqlx.ErrNotFound{{end}}
)
`
