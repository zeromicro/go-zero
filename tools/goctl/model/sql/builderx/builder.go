package builderx

import (
	"github.com/tal-tech/go-zero/core/stores/builderx"
)

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.NewEq instead.
func NewEq(in interface{}) builderx.Eq {
	return builderx.NewEq(in)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.NewGt instead.
func NewGt(in interface{}) builderx.Gt {
	return builderx.NewGt(ToMap(in))
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.ToMap instead.
func ToMap(in interface{}) map[string]interface{} {
	return builderx.ToMap(in)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.FieldNames instead.
func FieldNames(in interface{}) []string {
	return builderx.FieldNames(in)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.RawFieldNames instead.
func RawFieldNames(in interface{}, postgresSql ...bool) []string {
	return builderx.RawFieldNames(in, postgresSql...)
}

// Deprecated: Use github.com/tal-tech/go-zero/core/stores/builderx.PostgreSqlJoin instead.
func PostgreSqlJoin(elems []string) string {
	return builderx.PostgreSqlJoin(elems)
}
