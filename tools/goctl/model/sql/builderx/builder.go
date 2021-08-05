package builderx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-xorm/builder"
)

const dbTag = "db"

// NewEq wraps builder.Eq
func NewEq(in interface{}) builder.Eq {
	return builder.Eq(ToMap(in))
}

// NewGt wraps builder.Gt
func NewGt(in interface{}) builder.Gt {
	return builder.Gt(ToMap(in))
}

// ToMap converts interface into map
func ToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
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
			// set key of map to value in struct field
			val := v.Field(i)
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
			out[tagv] = current
		}
	}

	return out
}

// FieldNames returns field names from given in.
// deprecated: use RawFieldNames instead automatically while model generating after goctl version v1.1.0
func FieldNames(in interface{}) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
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
			out = append(out, tagv)
		} else {
			out = append(out, fi.Name)
		}
	}

	return out
}

// RawFieldNames converts golang struct field into slice string
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
				out = append(out, fmt.Sprintf("%s", tagv))
			} else {
				out = append(out, fmt.Sprintf("`%s`", tagv))
			}
		} else {
			if pg {
				out = append(out, fmt.Sprintf("%s", fi.Name))
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
			}
		}
	}

	return out
}

func PostgreSqlJoin(elems []string) string {
	var b = new(strings.Builder)
	for index, e := range elems {
		b.WriteString(fmt.Sprintf("%s = $%d, ", e, index+1))
	}

	if b.Len() == 0 {
		return b.String()
	}

	return b.String()[0 : b.Len()-2]
}
