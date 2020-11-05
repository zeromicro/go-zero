package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtoTmpl(t *testing.T) {
	out, err := filepath.Abs("./test/test.proto")
	assert.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(filepath.Dir(out))
	}()
	err = ProtoTmpl(out)
	assert.Nil(t, err)
	_, err = os.Stat(out)
	assert.Nil(t, err)
}
