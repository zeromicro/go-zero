package vben

import (
	"reflect"
	"testing"
)

func TestConvertTagToRules(t *testing.T) {
	type args struct {
		tagString string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "testString",
			args:    args{"omitempty,min=1,max=50"},
			want:    []string{"min: 1", "max: 50"},
			wantErr: false,
		},
		{
			name:    "testLen",
			args:    args{"omitempty,len=50"},
			want:    []string{"len: 50"},
			wantErr: false,
		},
		{
			name:    "testNum",
			args:    args{"omitempty,gte=1,lte=50"},
			want:    []string{"min: 1", "max: 50"},
			wantErr: false,
		},
		{
			name:    "testFloat",
			args:    args{"omitempty,gte=1.1,lte=50.0"},
			want:    []string{"min: 1.1", "max: 50.0"},
			wantErr: false,
		},
		{
			name:    "testFloat2",
			args:    args{"omitempty,gt=1.11,lt=50.01"},
			want:    []string{"min: 1.12", "max: 50.00"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTagToRules(tt.args.tagString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTagToRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTagToRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}
