package template

// Error defines an error template
var Error = `package {{.pkg}}

import "github.com/l306287405/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
