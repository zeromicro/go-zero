type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
        withSession(session sqlx.Session) {{.lowerStartCamelObject}}Model
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}sqlc.CachedConn{{else}}conn sqlx.SqlConn{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
