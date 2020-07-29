package sqltemplate

var Imports = `
import (
	{{if .containsCache}}"database/sql"
	"fmt"{{end}}
	"strings"
	"time"

    "zero/core/stores/cache"
	"zero/core/stores/sqlc"
	"zero/core/stores/sqlx"
	"zero/core/stringx"
)
`
