package template

import "fmt"

// GenDelete defines a delete template
func GenDelete(dialect SqlDialect) string {
	p1 := dialect.PositionalParameter(0)
	return fmt.Sprintf(`
func (m *default{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %%s where {{.originalPrimaryKey}} = %s", m.table)
		return conn.Exec(query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %%s where {{.originalPrimaryKey}} = %s", m.table)
		_,err:=m.conn.Exec(query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
`, p1, p1)
}

// DeleteMethod defines a delete template for interface method
var DeleteMethod = `Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error`
