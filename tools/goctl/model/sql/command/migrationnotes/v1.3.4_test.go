package migrationnotes

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_needShow1_3_4(t *testing.T) {
	root, err := os.MkdirTemp("", "goctl-model")
	require.NoError(t, err)
	defer os.RemoveAll(root)

	dir1 := path.Join(root, "dir1")
	require.NoError(t, os.Mkdir(dir1, fs.ModePerm))
	os.WriteFile(filepath.Join(dir1, "foo_gen.go"), nil, fs.ModePerm)

	dir2 := path.Join(root, "dir2")
	require.NoError(t, os.Mkdir(dir2, fs.ModePerm))
	os.WriteFile(filepath.Join(dir2, "foomodel.go"), nil, fs.ModePerm)

	dir3 := path.Join(root, "dir3")
	require.NoError(t, os.Mkdir(dir3, fs.ModePerm))
	os.WriteFile(filepath.Join(dir3, "irrelevant.go"), nil, fs.ModePerm)

	type args struct {
		dir   string
		style string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "dir that contains *_gen.go should return false",
			args: args{
				dir: dir1,
			},
			want: false,
		},
		{
			name: "dir that contains *model.go without *_gen.go should return true",
			args: args{
				dir: dir2,
			},
			want: true,
		},
		{
			name: "dir that only contains irrelevant files should return false",
			args: args{
				dir: dir3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := needShow1_3_4(tt.args.dir, tt.args.style)
			if (err != nil) != tt.wantErr {
				t.Errorf("needShow1_3_4() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("needShow1_3_4() = %v, want %v", got, tt.want)
			}
		})
	}
}
