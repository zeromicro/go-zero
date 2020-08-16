package template

var New = `
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn,{{if .withCache}} c cache.CacheConf,{{end}} table string) *{{.upperStartCamelObject}}Model {
	return &{{.upperStartCamelObject}}Model{
		{{if .withCache}}CachedConn: sqlc.NewConn(conn, c){{else}}conn:conn{{end}},
		table:      table,
	}
}
`
