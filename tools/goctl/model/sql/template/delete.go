package sqltemplate

var Delete = `
func (m *{{.upperObject}}Model) Delete({{.lowerPrimaryKey}} {{.dataType}}) error {
	{{if .containsCache}}{{if .isNotPrimaryKey}}data,err:=m.FindOne({{.lowerPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}
	{{.keys}}
    _, err {{if .isNotPrimaryKey}}={{else}}:={{end}} m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.snakePrimaryKey}} = ?` + "`" + `
		return conn.Exec(query, {{.lowerPrimaryKey}})
	}, {{.keyValues}}){{else}}query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.snakePrimaryKey}} = ?` + "`" + `
		_,err:=m.ExecNoCache(query, {{.lowerPrimaryKey}}){{end}}
	return err
}
`
