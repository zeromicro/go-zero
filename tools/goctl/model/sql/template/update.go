package template

// Update defines a template for generating update codes
var Update = `
func (m *default{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.primaryCacheKey}}
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ?", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.Exec(query, {{.expressionValues}})
	}, {{.primaryKeyVariable}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = ?", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return err
}
`

// UpdateMethod defines an interface method template for generating update codes
var UpdateMethod = `Update(data {{.upperStartCamelObject}}) error`
