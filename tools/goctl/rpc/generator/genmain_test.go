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

func TestGenerateMain(t *testing.T) {
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

	projectCtx, err := ctx.Prepare(dir)
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

	err = g.GenMain(dirCtx, proto)
	assert.Nil(t, err)
}
