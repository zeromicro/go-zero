package builderx

import (
	"github.com/tal-tech/go-zero/core/stores/builderx"
)

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.RawFieldNames instead.
func FieldNames(in interface{}) []string {
	return builderx.RawFieldNames(in)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.RawFieldNames instead.
func RawFieldNames(in interface{}, postgresSql ...bool) []string {
	return builderx.RawFieldNames(in, postgresSql...)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.PostgreSqlJoin instead.
func PostgreSqlJoin(elems []string) string {
	return builderx.PostgreSqlJoin(elems)
}
