package gogen

import (
	"path/filepath"
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
	m[dirHandler] = filepath.Join(ctx.TargetDir, dirInternal, dirHandler)
	m[dirLogic] = filepath.Join(ctx.TargetDir, dirInternal, dirLogic)
	m[dirPb] = filepath.Join(ctx.TargetDir, dirPb)
	m[dirSvc] = filepath.Join(ctx.TargetDir, dirInternal, dirSvc)
	m[dirShared] = g.Ctx.SharedDir
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
	return g.Ctx.Module + relativePath
}
