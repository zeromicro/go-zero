package javagen

import (
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestSpecTypeToJavaIntegerWidths(t *testing.T) {
	for _, name := range []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
		if _, err := specTypeToJava(spec.PrimitiveType{RawName: name}); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
