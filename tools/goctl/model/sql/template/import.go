package sqltemplate

var Imports = `
import (
	{{if .containsCache}}"database/sql"
	"fmt"{{end}}
	"strings"
	"time"

    "github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)
`
