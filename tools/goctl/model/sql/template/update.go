package template

var Update = `
func (m *{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.primaryCacheKey}}
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerStartCamelObject}}RowsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return conn.Exec(query, {{.expressionValues}})
	}, {{.primaryKeyVariable}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerStartCamelObject}}RowsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
    _,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return err
}
`
