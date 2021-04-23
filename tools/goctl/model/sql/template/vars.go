package template

import "fmt"

// GenVars produces a template for var block in model
func GenVars(dialect SqlDialect) string {
	quote := dialect.IdentifierQuote()
	escQuote := escapeQuote(quote)
	p1, p2 := dialect.PositionalParameter(0), dialect.PositionalParameter(1)
	return fmt.Sprintf(`
var (
	{{.lowerStartCamelObject}}FieldNames          = rawFieldNames(&{{.upperStartCamelObject}}{})
	{{.lowerStartCamelObject}}Rows                = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
	{{.lowerStartCamelObject}}RowsExpectAutoSet   = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}",{{end}} "%screate_time%s", "%supdate_time%s"), ",")
	{{.lowerStartCamelObject}}RowsWithPlaceHolder = strings.Join(stringx.Remove({{.lowerStartCamelObject}}FieldNames, "{{.originalPrimaryKey}}", "%screate_time%s", "%supdate_time%s"), "=%s,") + "=%s"

	{{if .withCache}}{{.cacheKeys}}{{end}}
)

const dbTag = "db"
func rawFieldNames(in interface{}) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %%T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(dbTag); tagv != "" {
			out = append(out, fmt.Sprintf("%s%%s%s", tagv))
		} else {
			out = append(out, fmt.Sprintf("\"%%s\"", fi.Name))
		}
	}

	return out
}
`, escQuote, escQuote, escQuote, escQuote, escQuote, escQuote, escQuote, escQuote, p1, p2, escQuote, escQuote)
}

func escapeQuote(q string) string {
	if q == `"` {
		return `\"`
	}
	return q
}

type SqlDialect interface {
	IdentifierQuote() string
	PositionalParameter(idx int) string
}
