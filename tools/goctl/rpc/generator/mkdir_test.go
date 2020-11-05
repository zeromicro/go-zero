package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

func TestMkDirInGoPath(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := util.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)
	internal := filepath.Join(dir, "internal")
	assert.True(t, true, func() bool {
		return filepath.Join(dir, strings.ToLower(projectName)) == dirCtx.GetCall().Filename && projectName == dirCtx.GetCall().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(dir, "etc") == dirCtx.GetEtc().Filename && filepath.Join(projectName, "etc") == dirCtx.GetEtc().Package
	}())
	assert.True(t, true, func() bool {
		return internal == dirCtx.GetInternal().Filename && filepath.Join(projectName, "internal") == dirCtx.GetInternal().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "config") == dirCtx.GetConfig().Filename && filepath.Join(projectName, "internal", "config") == dirCtx.GetConfig().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "logic") == dirCtx.GetLogic().Filename && filepath.Join(projectName, "internal", "logic") == dirCtx.GetLogic().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "server") == dirCtx.GetServer().Filename && filepath.Join(projectName, "internal", "server") == dirCtx.GetServer().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "svc") == dirCtx.GetSvc().Filename && filepath.Join(projectName, "internal", "svc") == dirCtx.GetSvc().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, strings.ToLower(proto.Service.Name)) == dirCtx.GetPb().Filename && filepath.Join(projectName, "internal", strings.ToLower(proto.Service.Name)) == dirCtx.GetPb().Package
	}())
	assert.True(t, true, func() bool {
		return dir == dirCtx.GetMain().Filename && projectName == dirCtx.GetMain().Package
	}())
}

func TestMkDirInGoMod(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := util.MkdirIfNotExist(dir)
	if err != nil {
		return
	}

	_, err = execx.Run("go mod init "+projectName, dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)
	internal := filepath.Join(dir, "internal")
	assert.True(t, true, func() bool {
		return filepath.Join(dir, strings.ToLower(projectName)) == dirCtx.GetCall().Filename && projectName == dirCtx.GetCall().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(dir, "etc") == dirCtx.GetEtc().Filename && filepath.Join(projectName, "etc") == dirCtx.GetEtc().Package
	}())
	assert.True(t, true, func() bool {
		return internal == dirCtx.GetInternal().Filename && filepath.Join(projectName, "internal") == dirCtx.GetInternal().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "config") == dirCtx.GetConfig().Filename && filepath.Join(projectName, "internal", "config") == dirCtx.GetConfig().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "logic") == dirCtx.GetLogic().Filename && filepath.Join(projectName, "internal", "logic") == dirCtx.GetLogic().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "server") == dirCtx.GetServer().Filename && filepath.Join(projectName, "internal", "server") == dirCtx.GetServer().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, "svc") == dirCtx.GetSvc().Filename && filepath.Join(projectName, "internal", "svc") == dirCtx.GetSvc().Package
	}())
	assert.True(t, true, func() bool {
		return filepath.Join(internal, strings.ToLower(proto.Service.Name)) == dirCtx.GetPb().Filename && filepath.Join(projectName, "internal", strings.ToLower(proto.Service.Name)) == dirCtx.GetPb().Package
	}())
	assert.True(t, true, func() bool {
		return dir == dirCtx.GetMain().Filename && projectName == dirCtx.GetMain().Package
	}())
}
