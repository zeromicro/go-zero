package template

var Update = `
func (m *{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.primaryCacheKey}}
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerStartCamelObject}}RowsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return conn.Exec(query, {{.expressionValues}})
	}, {{.primaryKeyVariable}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerStartCamelObject}}RowsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
    ret,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return ret,err
}
`
