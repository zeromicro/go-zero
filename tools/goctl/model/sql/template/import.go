package template

var (
	Imports = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/3Rivers/go-zero/core/stores/cache"
	"github.com/3Rivers/go-zero/core/stores/sqlc"
	"github.com/3Rivers/go-zero/core/stores/sqlx"
	"github.com/3Rivers/go-zero/core/stringx"
	"github.com/3Rivers/go-zero/tools/goctl/model/sql/builderx"
)
`
	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/3Rivers/go-zero/core/stores/sqlc"
	"github.com/3Rivers/go-zero/core/stores/sqlx"
	"github.com/3Rivers/go-zero/core/stringx"
	"github.com/3Rivers/go-zero/tools/goctl/model/sql/builderx"
)
`
)
