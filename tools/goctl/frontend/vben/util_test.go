package vben

import (
	"fmt"
	"testing"
)

func TestFindBeginEndOfLocaleField(t *testing.T) {
	type args struct {
		data   string
		target string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{
			name: "test",
			args: args{data: "export default {\n  login: '登录',\n  errorLogList: '错误日志列表',\n  test: {\n    hello: 'hello',\n  },\n};\n",
				target: "test"},
			want:  76,
			want1: 108,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindBeginEndOfLocaleField(tt.args.data, tt.args.target)
			fmt.Println(tt.args.data[got:got1])
			if got != tt.want {
				t.Errorf("FindBeginEndOfLocaleField() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindBeginEndOfLocaleField() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
