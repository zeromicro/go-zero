package template

var Insert = `
func (m *{{.upperStartCamelObject}}Model) Insert(data {{.upperStartCamelObject}}) error {
	query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "`(` + " + `{{.lowerStartCamelObject}}RowsExpectAutoSet` + " + `) value ({{.expression}})` " + `
	_, err := m.ExecNoCache(query, {{.expressionValues}})
	return err
}
`
