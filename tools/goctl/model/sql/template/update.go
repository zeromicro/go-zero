package sqltemplate

var Update = `
func (m *{{.upperObject}}Model) Update(data {{.upperObject}}) error {
	{{if .containsCache}}{{.primaryCacheKey}}
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerObject}}RowsWithPlaceHolder` + " + `" + ` where {{.primarySnakeCase}} = ?` + "`" + `
		return conn.Exec(query, {{.expressionValues}})
	}, {{.primaryKeyVariable}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerObject}}RowsWithPlaceHolder` + " + `" + ` where {{.primarySnakeCase}} = ?` + "`" + `
    _,err:=m.ExecNoCache(query, {{.expressionValues}}){{end}}
	return err
}
`
