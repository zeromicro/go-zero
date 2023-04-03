package util

import (
	"reflect"
	"testing"
)

func TestConvertValidateTagToSwagger(t *testing.T) {
	type args struct {
		tagData string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "testString",
			args:    args{"json:\"path,optional\" validate:\"omitempty,min=1,max=50\""},
			want:    []string{"// min length : 1\n", "// max length : 50\n"},
			wantErr: false,
		},
		{
			name:    "testLen",
			args:    args{"json:\"path,optional\" validate:\"omitempty,len=50\""},
			want:    []string{"// max length : 50\n", "// min length : 50\n"},
			wantErr: false,
		},
		{
			name:    "testNum",
			args:    args{"json:\"path,optional\" validate:\"omitempty,gte=1,lte=50\""},
			want:    []string{"// min : 1\n", "// max : 50\n"},
			wantErr: false,
		},
		{
			name:    "testFloat",
			args:    args{"json:\"path,optional\" validate:\"omitempty,gte=1.1,lte=50.0\""},
			want:    []string{"// min : 1.1\n", "// max : 50.0\n"},
			wantErr: false,
		},
		{
			name:    "testFloat2",
			args:    args{"json:\"path,optional\" validate:\"omitempty,gt=1.11,lt=50.01\""},
			want:    []string{"// min : 1.12\n", "// max : 50.00\n"},
			wantErr: false,
		},
		{
			name:    "testRequired",
			args:    args{"json:\"path,optional\" validate:\"required,gte=1.1,lte=50.0\""},
			want:    []string{"// required : true\n", "// min : 1.1\n", "// max : 50.0\n"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertValidateTagToSwagger(tt.args.tagData)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertValidateTagToSwagger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertValidateTagToSwagger() got = %v, want %v", got, tt.want)
			}
		})
	}
}
