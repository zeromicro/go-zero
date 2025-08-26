// copy from core/stores/sqlx/utils.go

package mocksql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
)

// ErrNotFound is the alias of sql.ErrNoRows
var ErrNotFound = sql.ErrNoRows

func escape(input string) string {
	var b strings.Builder

	for _, ch := range input {
		switch ch {
		case '\x00':
			b.WriteString(`\x00`)
		case '\r':
			b.WriteString(`\r`)
		case '\n':
			b.WriteString(`\n`)
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '"':
			b.WriteString(`\"`)
		case '\x1a':
			b.WriteString(`\x1a`)
		default:
			b.WriteRune(ch)
		}
	}

	return b.String()
}

func format(query string, args ...any) (string, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return query, nil
	}

	var b strings.Builder
	argIndex := 0

	for _, ch := range query {
		if ch == '?' {
			if argIndex >= numArgs {
				return "", fmt.Errorf("%d ? in sql, but less arguments provided", argIndex)
			}

			arg := args[argIndex]
			argIndex++

			switch v := arg.(type) {
			case bool:
				if v {
					b.WriteByte('1')
				} else {
					b.WriteByte('0')
				}
			case string:
				b.WriteByte('\'')
				b.WriteString(escape(v))
				b.WriteByte('\'')
			default:
				b.WriteString(mapping.Repr(v))
			}
		} else {
			b.WriteRune(ch)
		}
	}

	if argIndex < numArgs {
		return "", fmt.Errorf("%d ? in sql, but more arguments provided", argIndex)
	}

	return b.String(), nil
}

func logSqlError(stmt string, err error) {
	if err != nil && err != ErrNotFound {
		logx.Errorf("stmt: %s, error: %s", stmt, err.Error())
	}
}
