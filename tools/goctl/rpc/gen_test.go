package base

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

func TestParseImport(t *testing.T) {
	src, _ := filepath.Abs("./test.proto")
	base, _ := filepath.Abs("./base.proto")
	imports, containsAny, err := parser.ParseImport(src)
	assert.Nil(t, err)
	assert.Equal(t, true, containsAny)
	assert.Equal(t, 1, len(imports))
	assert.Equal(t, "github.com/tal-tech/go-zero/tools/goctl/rpc", imports[0].PbImportName)
	assert.Equal(t, base, imports[0].OriginalProtoPath)
}

func TestTransfer(t *testing.T) {
	src, _ := filepath.Abs("./test.proto")
	abs, _ := filepath.Abs("./test")
	imports, _, _ := parser.ParseImport(src)
	proto, err := parser.Transfer(src, abs, imports, console.NewConsole(false))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(proto.Service))
	assert.Equal(t, "Greeter", proto.Service[0].Name.Source())
	assert.Equal(t, 5, len(proto.Structure))
	data, ok := proto.Structure["map"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "M", data.Field[0].Name.Source())
}
