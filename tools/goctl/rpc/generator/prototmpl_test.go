package generator

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtoTmpl(t *testing.T) {
	_ = Clean()
	// exists dir
	err := ProtoTmpl(t.TempDir())
	assert.Nil(t, err)

	// not exist dir
	dir := filepath.Join(t.TempDir(), "test")
	err = ProtoTmpl(dir)
	assert.Nil(t, err)
}
