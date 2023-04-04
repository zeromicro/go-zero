package retry

import (
	"fmt"
	"testing"
)

func TestMergeRetryConfig(t *testing.T) {
	type args struct {
		config string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default retry",
			args: args{config: "[]"},
			want: fmt.Sprintf("[%s]", defaultRetry),
		},
		{
			name: "default retry 2",
			args: args{config: ""},
			want: fmt.Sprintf("[%s]", defaultRetry),
		},
		{
			name: "add one config",
			args: args{config: `{"name":[{"service":""}]}`},
			want: fmt.Sprintf(`[{"name":[{"service":""}]},%s]`, defaultRetry),
		},
		{
			name: "add multiple config",
			args: args{config: `[{"name":[{"service":""}]}]`},
			want: fmt.Sprintf(`[{"name":[{"service":""}]},%s]`, defaultRetry),
		},
		{
			name: "add multiple config 2",
			args: args{config: `[{"name":[{"service":""}]},{"name":[{"service":""}]}]`},
			want: fmt.Sprintf(`[{"name":[{"service":""}]},{"name":[{"service":""}]},%s]`, defaultRetry),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeRetryConfig(tt.args.config); got != tt.want {
				t.Errorf("MergeRetryConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
