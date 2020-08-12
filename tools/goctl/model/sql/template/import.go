package template

var Imports = `
import (
	"database/sql"{{if .withCache}}
	"fmt"
	{{end}}"strings"
	"time"

  "github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
  "github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)
`
