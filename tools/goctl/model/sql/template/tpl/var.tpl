var (
{{.lowerStartCamelObject}}FieldNames = builder.RawFieldNames(&{{.upperStartCamelObject}}{}{{if .postgreSql}}, true{{end}})
{{.lowerStartCamelObject}}Rows = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
{{.lowerStartCamelObject}}RowsExpectAutoSet = {{if .postgreSql}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}", {{end}} {{.ignoreColumns}}), ","){{else}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}", {{end}} {{.ignoreColumns}}), ","){{end}}
{{.lowerStartCamelObject}}RowsWithPlaceHolder = {{if .postgreSql}}builder.PostgreSqlJoin(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", {{.ignoreColumns}})){{else}}strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", {{.ignoreColumns}}), "=?,") + "=?"{{end}}

{{if .withCache}}{{.cacheKeys}}{{end}}
)
