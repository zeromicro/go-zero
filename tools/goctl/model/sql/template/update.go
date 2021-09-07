package template

// Update defines a template for generating update codes
var Update = `
func (m *default{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return err
}
`

// UpdateMethod defines an interface method template for generating update codes
var UpdateMethod = `Update(data {{.upperStartCamelObject}}) error`
