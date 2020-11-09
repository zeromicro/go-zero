package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

func TestGenerateCaseNilImport(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		//_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_stream.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	if err := g.Prepare(); err == nil {
		targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_stream.pb.go")
		err = g.GenPb(dirCtx, nil, proto)
		assert.Nil(t, err)
		assert.True(t, func() bool {
			return util.FileExists(targetPb)
		}())
	}
}

func TestGenerateCaseImport(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_stream.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	if err := g.Prepare(); err == nil {
		err = g.GenPb(dirCtx, nil, proto)
		assert.Nil(t, err)

		targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_stream.pb.go")
		assert.True(t, func() bool {
			return util.FileExists(targetPb)
		}())
	}
}

func TestGenerateCasePathOption(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	if err := g.Prepare(); err == nil {
		err = g.GenPb(dirCtx, nil, proto)
		assert.Nil(t, err)

		targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_option.pb.go")
		assert.True(t, func() bool {
			return util.FileExists(targetPb)
		}())
	}
}

func TestGenerateCaseWordOption(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_word_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	if err := g.Prepare(); err == nil {

		err = g.GenPb(dirCtx, nil, proto)
		assert.Nil(t, err)

		targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_word_option.pb.go")
		assert.True(t, func() bool {
			return util.FileExists(targetPb)
		}())
	}
}

// test keyword go
func TestGenerateCaseGoOption(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Prepare(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_go_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	if err := g.Prepare(); err == nil {

		err = g.GenPb(dirCtx, nil, proto)
		assert.Nil(t, err)

		targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_go_option.pb.go")
		assert.True(t, func() bool {
			return util.FileExists(targetPb)
		}())
	}
}
