package template

var New = `
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn, c cache.CacheConf, table string) *{{.upperStartCamelObject}}Model {
	return &{{.upperStartCamelObject}}Model{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}
`
