package builder

import (
	"fmt"
	"reflect"
	"strings"
)

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func RawFieldNames(in any, postgresSql ...bool) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var pg bool
	if len(postgresSql) > 0 {
		pg = postgresSql[0]
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			if pg {
				out = append(out, fi.Name)
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
			}
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			if pg {
				out = append(out, tagv)
			} else {
				out = append(out, fmt.Sprintf("`%s`", tagv))
			}
		}
	}

	return out
}

// PostgreSqlJoin concatenates the given elements into a string.
func PostgreSqlJoin(elems []string) string {
	b := new(strings.Builder)
	for index, e := range elems {
		b.WriteString(fmt.Sprintf("%s = $%d, ", e, index+2))
	}

	if b.Len() == 0 {
		return b.String()
	}

	return b.String()[0 : b.Len()-2]
}
