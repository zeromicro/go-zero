package golang

import (
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func GetParentPackage(dir string) (string, string, error) {
	return GetParentPackageWithModule(dir, "")
}

func GetParentPackageWithModule(dir, moduleName string) (string, string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}

	var projectCtx *ctx.ProjectContext
	if len(moduleName) > 0 {
		projectCtx, err = ctx.PrepareWithModule(abs, moduleName)
	} else {
		projectCtx, err = ctx.Prepare(abs)
	}
	if err != nil {
		return "", "", err
	}

	return buildParentPackage(projectCtx)
}

// buildParentPackage extracts the common logic for building parent package paths
func buildParentPackage(projectCtx *ctx.ProjectContext) (string, string, error) {
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
