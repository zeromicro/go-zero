package sqltemplate

var New = `
func New{{.upperObject}}Model(conn sqlx.SqlConn, c cache.CacheConf, table string) *{{.upperObject}}Model {
	return &{{.upperObject}}Model{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}
`
