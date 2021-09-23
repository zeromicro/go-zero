package generator

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func TestProtoTmpl(t *testing.T) {
	_ = Clean()
	// exists dir
	err := ProtoTmpl(util.MustTempDir())
	assert.Nil(t, err)

	// not exist dir
	dir := filepath.Join(util.MustTempDir(), "test")
	err = ProtoTmpl(dir)
	assert.Nil(t, err)
}
