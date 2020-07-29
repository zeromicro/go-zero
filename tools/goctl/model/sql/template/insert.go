package sqltemplate

var Insert = `
func (m *{{.upperObject}}Model) Insert(data {{.upperObject}}) error {
	query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "`(` + " + `{{.lowerObject}}RowsExpectAutoSet` + " + `) value ({{.expression}})` " + `
	_, err := m.ExecNoCache(query, {{.expressionValues}})
	return err
}
`
