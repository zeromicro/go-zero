package dartgen

import (
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func Test_getPropertyFromMember(t *testing.T) {
	tests := []struct {
		name   string
		member spec.Member
		want   string
	}{
		{
			name: "json tag should be ok",
			member: spec.Member{
				Tag:  "`json:\"foo\"`",
				Name: "Foo",
			},
			want: "foo",
		},
		{
			name: "form tag should be ok",
			member: spec.Member{
				Tag:  "`form:\"bar\"`",
				Name: "Bar",
			},
			want: "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPropertyFromMember(tt.member); got != tt.want {
				t.Errorf("getPropertyFromMember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specTypeToDart(t *testing.T) {
	tests := []struct {
		name     string
		specType spec.Type
		want     string
		wantErr  bool
	}{
		{
			name:     "[]string should return List<String>",
			specType: spec.ArrayType{RawName: "[]string", Value: spec.PrimitiveType{RawName: "string"}},
			want:     "List<String>",
		},
		{
			name:     "[]Foo should return List<Foo>",
			specType: spec.ArrayType{RawName: "[]Foo", Value: spec.DefineStruct{RawName: "Foo"}},
			want:     "List<Foo>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := specTypeToDart(tt.specType)
			if (err != nil) != tt.wantErr {
				t.Errorf("specTypeToDart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("specTypeToDart() = %v, want %v", got, tt.want)
			}
		})
	}
}
