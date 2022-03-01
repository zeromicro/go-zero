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
