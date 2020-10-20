package template

var Insert = `
func (m *{{.upperStartCamelObject}}Model) Insert(data {{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{if .containsIndexCache}}{{.keys}}
    ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerStartCamelObject}}RowsExpectAutoSet` + " + `) values ({{.expression}})` " + `
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerStartCamelObject}}RowsExpectAutoSet` + " + `) values ({{.expression}})` " + `
    ret,err:=m.ExecNoCache(query, {{.expressionValues}})
	{{end}}{{else}}query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerStartCamelObject}}RowsExpectAutoSet` + " + `) values ({{.expression}})` " + `
    ret,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return ret,err
}
`
