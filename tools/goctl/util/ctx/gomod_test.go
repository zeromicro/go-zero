package ctx

import (
	"bytes"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func TestProjectFromGoMod(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := pathx.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	_, err = execx.Run("go mod init "+projectName, dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	ctx, err := projectFromGoMod(dir)
	assert.Nil(t, err)
	assert.Equal(t, projectName, ctx.Path)
	assert.Equal(t, dir, ctx.Dir)
}

func Test_getRealModule(t *testing.T) {
	type args struct {
		workDir string
		execRun execx.RunFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *Module
		wantErr bool
	}{
		{
			name: "single module",
			args: args{
				workDir: "/home/foo",
				execRun: func(arg, dir string, in ...*bytes.Buffer) (string, error) {
					return `{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.16"}`, nil
				},
			},
			want: &Module{
				Path:      "foo",
				Dir:       "/home/foo",
				GoMod:     "/home/foo/go.mod",
				GoVersion: "go1.16",
			},
		},
		{
			name: "go work multiple modules",
			args: args{
				workDir: "/home/bar",
				execRun: func(arg, dir string, in ...*bytes.Buffer) (string, error) {
					return `
					{
						"Path":"foo",
						"Dir":"/home/foo",
						"GoMod":"/home/foo/go.mod",
						"GoVersion":"go1.18"
					}
					{
						"Path":"bar",
						"Dir":"/home/bar",
						"GoMod":"/home/bar/go.mod",
						"GoVersion":"go1.18"
					}`, nil
				},
			},
			want: &Module{
				Path:      "bar",
				Dir:       "/home/bar",
				GoMod:     "/home/bar/go.mod",
				GoVersion: "go1.18",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRealModule(tt.args.workDir, tt.args.execRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRealModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRealModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
