package util

import "testing"

func TestConvertValidateTagToSwagger(t *testing.T) {
	type args struct {
		tagData string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{"json:\"path,optional\" validate:\"omitempty,min=1,max=50\""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertValidateTagToSwagger(tt.args.tagData); got != tt.want {
				t.Errorf("ConvertValidateTagToSwagger() = %v, want %v", got, tt.want)
			}
		})
	}
}
