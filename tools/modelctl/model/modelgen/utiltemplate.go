package modelgen

const (
	utilTemplateText = `package {{.Package}}

import (
    "errors"
	
    {{if .WithCache}}"github.com/tal-tech/go-zero/core/stores/redis"
    "github.com/tal-tech/go-zero/core/stores/sqlc"
    "github.com/tal-tech/go-zero/core/stores/sqlx"{{end}}
)
{{if .WithCache}}
type CachedModel struct {
    table string
    conn  sqlx.SqlConn
    rds   *redis.Redis
    sqlc.CachedConn
}

func NewCachedModel(conn sqlx.SqlConn, table string, rds *redis.Redis) *CachedModel {
    return &CachedModel{
        table:      table,
        conn:       conn,
        rds:        rds,
        CachedConn: sqlc.NewCachedConn(conn, rds),
    }
}
{{end}}
var (
    ErrNotFound = errors.New("not found")
)
`
)
