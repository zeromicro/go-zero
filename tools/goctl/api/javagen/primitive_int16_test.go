package javagen

import (
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestSpecTypeToJavaIntegerWidths(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
	}{
		{name: "int", want: "int"},
		{name: "int8", want: "int"},
		{name: "int16", want: "int"},
		{name: "int32", want: "int"},
		{name: "int64", want: "long"},
		{name: "uint", want: "int"},
		{name: "uint8", want: "int"},
		{name: "uint16", want: "int"},
		{name: "uint32", want: "int"},
		{name: "uint64", want: "long"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := specTypeToJava(spec.PrimitiveType{RawName: test.name})
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("specTypeToJava(%q) = %q, want %q", test.name, got, test.want)
			}
		})
	}
}
