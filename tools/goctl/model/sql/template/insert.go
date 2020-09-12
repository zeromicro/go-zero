package template

var Insert = `
func (m *{{.upperStartCamelObject}}Model) Insert(data {{.upperStartCamelObject}}) (sql.Result, error) {
	query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerStartCamelObject}}RowsExpectAutoSet` + " + `) values ({{.expression}})` " + `
	return m.{{if .withCache}}ExecNoCache{{else}}conn.Exec{{end}}(query, {{.expressionValues}})
}
`
