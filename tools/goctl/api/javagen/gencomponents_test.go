package javagen

import (
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestCheckStructUnsupportedType(t *testing.T) {
	var c componentsContext
	_, done, err := c.checkStruct(spec.PrimitiveType{RawName: "int64"})
	if !done {
		t.Fatal("checkStruct should stop for a non-struct type")
	}
	if err == nil || err.Error() != "unsupported type int64" {
		t.Fatalf("checkStruct error = %v, want %q", err, "unsupported type int64")
	}
}
