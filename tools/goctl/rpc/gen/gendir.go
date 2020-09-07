package gen

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

//  target
//	├── etc
//	├── internal
//	│   ├── config
//	│   ├── handler
//	│   ├── logic
//	│   ├── pb
//	│   └── svc
func (g *defaultRpcGenerator) createDir() error {
	ctx := g.Ctx
	m := make(map[string]string)
	m[dirTarget] = ctx.TargetDir
	m[dirEtc] = filepath.Join(ctx.TargetDir, dirEtc)
	m[dirInternal] = filepath.Join(ctx.TargetDir, dirInternal)
	m[dirConfig] = filepath.Join(ctx.TargetDir, dirInternal, dirConfig)
	m[dirServer] = filepath.Join(ctx.TargetDir, dirInternal, dirServer)
	m[dirLogic] = filepath.Join(ctx.TargetDir, dirInternal, dirLogic)
	m[dirPb] = filepath.Join(ctx.TargetDir, dirPb)
	m[dirSvc] = filepath.Join(ctx.TargetDir, dirInternal, dirSvc)
	for _, d := range m {
		err := util.MkdirIfNotExist(d)
		if err != nil {
			return err
		}
	}
	g.dirM = m
	return nil
}

func (g *defaultRpcGenerator) mustGetPackage(dir string) string {
	target := g.dirM[dir]
	projectPath := g.Ctx.ProjectPath
	relativePath := strings.TrimPrefix(target, projectPath)
	os := runtime.GOOS
	switch os {
	case "windows":
		relativePath = filepath.ToSlash(relativePath)
	case "darwin", "linux":
	default:
		g.Ctx.Fatalln("unexpected os: %s", os)
	}
	return g.Ctx.Module + relativePath
}
