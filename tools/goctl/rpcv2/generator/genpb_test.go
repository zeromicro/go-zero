package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func TestDefaultGenerator_GenPb_CaseNilImport(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		//_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Background(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_stream.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_stream.pb.go")
	err = g.GenPb(dirCtx, nil, dirCtx.GetPb(), proto)
	assert.Nil(t, err)
	assert.True(t, func() bool {
		return util.FileExists(targetPb)
	}())
}

func TestDefaultGenerator_GenPb_CaseImport(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Background(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_stream.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	err = g.Prepare()
	if err != nil {
		return
	}

	err = g.GenPb(dirCtx, nil, dirCtx.GetPb(), proto)
	assert.Nil(t, err)

	targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_stream.pb.go")
	assert.True(t, func() bool {
		return util.FileExists(targetPb)
	}())
}

func TestDefaultGenerator_GenPb_Case_Path_Option(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Background(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	err = g.Prepare()
	if err != nil {
		return
	}

	err = g.GenPb(dirCtx, nil, dirCtx.GetPb(), proto)
	assert.Nil(t, err)

	targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_option.pb.go")
	assert.True(t, func() bool {
		return util.FileExists(targetPb)
	}())
}

func TestDefaultGenerator_GenPb_Case_Word_Option(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Background(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_word_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	err = g.Prepare()
	if err != nil {
		return
	}

	err = g.GenPb(dirCtx, nil, dirCtx.GetPb(), proto)
	assert.Nil(t, err)

	targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_word_option.pb.go")
	assert.True(t, func() bool {
		return util.FileExists(targetPb)
	}())
}

// test keyword go
func TestDefaultGenerator_GenPb_Case_Go_Option(t *testing.T) {
	project := "stream"
	abs, err := filepath.Abs("./test")
	assert.Nil(t, err)

	dir := filepath.Join(abs, project)
	err = util.MkdirIfNotExist(dir)
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(abs)
	}()

	projectCtx, err := ctx.Background(dir)
	assert.Nil(t, err)

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse("./test_go_option.proto")
	assert.Nil(t, err)

	dirCtx, err := mkdir(projectCtx, proto)
	assert.Nil(t, err)

	g := NewDefaultGenerator()
	err = g.Prepare()
	if err != nil {
		return
	}

	err = g.GenPb(dirCtx, nil, dirCtx.GetPb(), proto)
	assert.Nil(t, err)

	targetPb := filepath.Join(dirCtx.GetPb().Filename, "test_go_option.pb.go")
	assert.True(t, func() bool {
		return util.FileExists(targetPb)
	}())
}
