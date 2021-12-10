package generator

import (
	"path/filepath"
	"strings"

	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	wd       = "wd"
	etc      = "etc"
	internal = "internal"
	config   = "config"
	logic    = "logic"
	server   = "server"
	svc      = "svc"
	pb       = "pb"
	call     = "call"
)

type (
	// DirContext defines a rpc service directories context
	DirContext interface {
		GetCall() Dir
		GetEtc() Dir
		GetInternal() Dir
		GetConfig() Dir
		GetLogic() Dir
		GetServer() Dir
		GetSvc() Dir
		GetPb() Dir
		GetMain() Dir
		GetServiceName() stringx.String
	}

	// Dir defines a directory
	Dir struct {
		Base     string
		Filename string
		Package  string
	}

	defaultDirContext struct {
		inner       map[string]Dir
		serviceName stringx.String
	}
)

func mkdir(ctx *ctx.ProjectContext, proto parser.Proto, _ *conf.Config) (DirContext, error) {
	inner := make(map[string]Dir)
	etcDir := filepath.Join(ctx.WorkDir, "etc")
	internalDir := filepath.Join(ctx.WorkDir, "internal")
	configDir := filepath.Join(internalDir, "config")
	logicDir := filepath.Join(internalDir, "logic")
	serverDir := filepath.Join(internalDir, "server")
	svcDir := filepath.Join(internalDir, "svc")
	pbDir := filepath.Join(ctx.WorkDir, proto.GoPackage)
	callDir := filepath.Join(ctx.WorkDir, strings.ToLower(stringx.From(proto.Service.Name).ToCamel()))
	if strings.EqualFold(proto.Service.Name, proto.GoPackage) {
		callDir = filepath.Join(ctx.WorkDir, strings.ToLower(stringx.From(proto.Service.Name+"_client").ToCamel()))
	}

	inner[wd] = Dir{
		Filename: ctx.WorkDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(ctx.WorkDir, ctx.Dir))),
		Base:     filepath.Base(ctx.WorkDir),
	}
	inner[etc] = Dir{
		Filename: etcDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(etcDir, ctx.Dir))),
		Base:     filepath.Base(etcDir),
	}
	inner[internal] = Dir{
		Filename: internalDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(internalDir, ctx.Dir))),
		Base:     filepath.Base(internalDir),
	}
	inner[config] = Dir{
		Filename: configDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(configDir, ctx.Dir))),
		Base:     filepath.Base(configDir),
	}
	inner[logic] = Dir{
		Filename: logicDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(logicDir, ctx.Dir))),
		Base:     filepath.Base(logicDir),
	}
	inner[server] = Dir{
		Filename: serverDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(serverDir, ctx.Dir))),
		Base:     filepath.Base(serverDir),
	}
	inner[svc] = Dir{
		Filename: svcDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(svcDir, ctx.Dir))),
		Base:     filepath.Base(svcDir),
	}
	inner[pb] = Dir{
		Filename: pbDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(pbDir, ctx.Dir))),
		Base:     filepath.Base(pbDir),
	}
	inner[call] = Dir{
		Filename: callDir,
		Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(callDir, ctx.Dir))),
		Base:     filepath.Base(callDir),
	}
	for _, v := range inner {
		err := util.MkdirIfNotExist(v.Filename)
		if err != nil {
			return nil, err
		}
	}
	serviceName := strings.TrimSuffix(proto.Name, filepath.Ext(proto.Name))
	return &defaultDirContext{
		inner:       inner,
		serviceName: stringx.From(strings.ReplaceAll(serviceName, "-", "")),
	}, nil
}

func (d *defaultDirContext) GetCall() Dir {
	return d.inner[call]
}

func (d *defaultDirContext) GetEtc() Dir {
	return d.inner[etc]
}

func (d *defaultDirContext) GetInternal() Dir {
	return d.inner[internal]
}

func (d *defaultDirContext) GetConfig() Dir {
	return d.inner[config]
}

func (d *defaultDirContext) GetLogic() Dir {
	return d.inner[logic]
}

func (d *defaultDirContext) GetServer() Dir {
	return d.inner[server]
}

func (d *defaultDirContext) GetSvc() Dir {
	return d.inner[svc]
}

func (d *defaultDirContext) GetPb() Dir {
	return d.inner[pb]
}

func (d *defaultDirContext) GetMain() Dir {
	return d.inner[wd]
}

func (d *defaultDirContext) GetServiceName() stringx.String {
	return d.serviceName
}

// Valid returns true if the directory is valid
func (d *Dir) Valid() bool {
	return len(d.Filename) > 0 && len(d.Package) > 0
}
