package template

import "fmt"

// GenUpdate defines a template for generating update codes
func GenUpdate(dialect SqlDialect) string {
	p1 := dialect.PositionalParameter(0)
	return fmt.Sprintf(`
func (m *default{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
    _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %%s set %%s where {{.originalPrimaryKey}} = %s", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %%s set %%s where {{.originalPrimaryKey}} = %s", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.Exec(query, {{.expressionValues}}){{end}}
	return err
}
`, p1, p1)
}

// UpdateMethod defines an interface method template for generating update codes
var UpdateMethod = `Update(data {{.upperStartCamelObject}}) error`
