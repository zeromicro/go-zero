type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
        WithSession(session sqlx.Session) *default{{.upperStartCamelObject}}Model
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}sqlc.CachedConn{{else}}conn sqlx.SqlConn{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
