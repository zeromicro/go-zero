package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func TestGenerateConfig(t *testing.T) {
	_ = Clean()
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
	err = g.GenConfig(dirCtx, proto)
	assert.Nil(t, err)

	// test file exists
	err = g.GenConfig(dirCtx, proto)
	assert.Nil(t, err)
}
