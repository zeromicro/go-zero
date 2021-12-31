package builderx

import (
	"github.com/z-micro/go-zero/core/stores/builder"
)

// Deprecated: Use github.com/z-micro/go-zero/core/stores/builder.RawFieldNames instead.
func FieldNames(in interface{}) []string {
	return builder.RawFieldNames(in)
}

// Deprecated: Use github.com/z-micro/go-zero/core/stores/builder.RawFieldNames instead.
func RawFieldNames(in interface{}, postgresSql ...bool) []string {
	return builder.RawFieldNames(in, postgresSql...)
}

// Deprecated: Use github.com/z-micro/go-zero/core/stores/builderx.PostgreSqlJoin instead.
func PostgreSqlJoin(elems []string) string {
	return builder.PostgreSqlJoin(elems)
}
