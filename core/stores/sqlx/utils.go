package sqlx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
)

var errUnbalancedEscape = errors.New("no char after escape char")

func desensitize(datasource string) string {
	// remove account
	pos := strings.LastIndex(datasource, "@")
	if 0 <= pos && pos+1 < len(datasource) {
		datasource = datasource[pos+1:]
	}

	return datasource
}

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
	var argIndex int
	bytes := len(query)

	for i := 0; i < bytes; i++ {
		ch := query[i]
		switch ch {
		case '?':
			if argIndex >= numArgs {
				return "", fmt.Errorf("%d ? in sql, but less arguments provided", argIndex)
			}

			writeValue(&b, args[argIndex])
			argIndex++
		case ':', '$':
			var j int
			for j = i + 1; j < bytes; j++ {
				char := query[j]
				if char < '0' || '9' < char {
					break
				}
			}

			if j > i+1 {
				index, err := strconv.Atoi(query[i+1 : j])
				if err != nil {
					return "", err
				}

				// index starts from 1 for pg or oracle
				if index > argIndex {
					argIndex = index
				}

				index--
				if index < 0 || numArgs <= index {
					return "", fmt.Errorf("wrong index %d in sql", index)
				}

				writeValue(&b, args[index])
				i = j - 1
			}
		case '\'', '"', '`':
			b.WriteByte(ch)

			for j := i + 1; j < bytes; j++ {
				cur := query[j]
				b.WriteByte(cur)

				if cur == '\\' {
					j++
					if j >= bytes {
						return "", errUnbalancedEscape
					}

					b.WriteByte(query[j])
				} else if cur == ch {
					i = j
					break
				}
			}
		default:
			b.WriteByte(ch)
		}
	}

	if argIndex < numArgs {
		return "", fmt.Errorf("%d arguments provided, not matching sql", argIndex)
	}

	return b.String(), nil
}

func logInstanceError(ctx context.Context, datasource string, err error) {
	datasource = desensitize(datasource)
	logx.WithContext(ctx).Errorf("Error on getting sql instance of %s: %v", datasource, err)
}

func logSqlError(ctx context.Context, stmt string, err error) {
	if err != nil && err != ErrNotFound {
		logx.WithContext(ctx).Errorf("stmt: %s, error: %s", stmt, err.Error())
	}
}

func writeValue(buf *strings.Builder, arg any) {
	switch v := arg.(type) {
	case bool:
		if v {
			buf.WriteByte('1')
		} else {
			buf.WriteByte('0')
		}
	case string:
		buf.WriteByte('\'')
		buf.WriteString(escape(v))
		buf.WriteByte('\'')
	case time.Time:
		buf.WriteByte('\'')
		buf.WriteString(v.String())
		buf.WriteByte('\'')
	case *time.Time:
		buf.WriteByte('\'')
		buf.WriteString(v.String())
		buf.WriteByte('\'')
	default:
		buf.WriteString(mapping.Repr(v))
	}
}
