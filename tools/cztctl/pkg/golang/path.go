package golang

import (
	"path/filepath"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/util/ctx"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
)

func GetParentPackage(dir string) (string, string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return "", "", err
	}

	// fix https://github.com/zeromicro/go-zero/issues/1058
	wd := projectCtx.WorkDir
	d := projectCtx.Dir
	same, err := pathx.SameFile(wd, d)
	if err != nil {
		return "", "", err
	}

	trim := strings.TrimPrefix(projectCtx.WorkDir, projectCtx.Dir)
	if same {
		trim = strings.TrimPrefix(strings.ToLower(projectCtx.WorkDir), strings.ToLower(projectCtx.Dir))
	}

	return filepath.ToSlash(filepath.Join(projectCtx.Path, trim)), projectCtx.Path, nil
}
