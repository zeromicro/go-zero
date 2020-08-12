package template

var Types = `
type (
	{{.upperStartCamelObject}}Model struct {
		sqlc.CachedConn
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
