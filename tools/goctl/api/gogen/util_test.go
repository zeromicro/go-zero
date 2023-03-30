package gogen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertRoutePathToSwagger(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "testPath",
			args: args{
				data: "/init/test/:name/:age",
			},
			want: "/init/test/{name}/{age}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ConvertRoutePathToSwagger(tt.args.data), "ConvertRoutePathToSwagger(%v)", tt.args.data)
		})
	}
}
