package builder

import (
	"fmt"
	"reflect"
	"strings"
)

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func RawFieldNames(in interface{}, postgresSql ...bool) []string {
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
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			if pg {
				out = append(out, tagv)
			} else {
				out = append(out, fmt.Sprintf("`%s`", tagv))
			}
		} else {
			if pg {
				out = append(out, fi.Name)
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
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
