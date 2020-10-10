package base

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/gen"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
)

func TestParseImport(t *testing.T) {
	src, _ := filepath.Abs("./test.proto")
	base, _ := filepath.Abs("./base.proto")
	imports, containsAny, err := parser.ParseImport(src)
	assert.Nil(t, err)
	assert.Equal(t, false, containsAny)
	assert.Equal(t, 1, len(imports))
	assert.Equal(t, "github.com/tal-tech/go-zero/tools/goctl/rpc", imports[0].PbImportName)
	assert.Equal(t, base, imports[0].OriginalProtoPath)
}

func TestParseImport2(t *testing.T) {
	src, _ := filepath.Abs("./test.proto")
	base, _ := filepath.Abs("./base.proto")
	imports, containsAny, err := parser.ParseImport(src)
	assert.Nil(t, err)
	assert.Equal(t, false, containsAny)
	assert.Equal(t, 1, len(imports))
	assert.Equal(t, "github.com/tal-tech/go-zero/tools/goctl/rpc", imports[0].PbImportName)
	assert.Equal(t, base, imports[0].OriginalProtoPath)
}

func TestTransfer(t *testing.T) {
	src, _ := filepath.Abs("./test.proto")
	abs, _ := filepath.Abs("./test")
	rpcCtx := ctx.MustCreateRpcContext(src, abs, "", false)
	g := gen.NewDefaultRpcGenerator(rpcCtx)
	err := g.Generate()
	assert.Nil(t, err)
}
