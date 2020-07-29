package sqltemplate

var Types = `
type (
	{{.upperObject}}Model struct {
		sqlc.CachedConn
		table string
	}

	{{.upperObject}} struct {
		{{.fields}}
	}
)
`
