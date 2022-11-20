func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where {{.originalPrimaryField}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	return conn.QueryRowCtx(ctx, v, query, primary)
}
