package template

var Error = `package {{.pkg}}

import "github.com/3Rivers/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
