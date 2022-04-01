package template

const (
	// Update defines a template for generating update codes
	Update = `
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    _,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return err
}
`

	// UpdateMethod defines an interface method template for generating update codes
	UpdateMethod = `Update(ctx context.Context, data *{{.upperStartCamelObject}}) error`
)
