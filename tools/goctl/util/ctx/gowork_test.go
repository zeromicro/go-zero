package ctx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func Test_isGoWork(t *testing.T) {
	dir := filepath.Join("/tmp", stringx.Rand())

	err := pathx.MkdirIfNotExist(dir)
	assert.Nil(t, err)

	defer os.RemoveAll(dir)

	gowork, err := isGoWork(dir)
	assert.False(t, gowork)
	assert.Nil(t, err)

	_, err = execx.Run("go work init", dir)
	assert.Nil(t, err)

	gowork, err = isGoWork(dir)
	assert.True(t, gowork)
	assert.Nil(t, err)

	subDir := filepath.Join(dir, stringx.Rand())
	err = pathx.MkdirIfNotExist(subDir)
	assert.Nil(t, err)

	gowork, err = isGoWork(subDir)
	assert.True(t, gowork)
	assert.Nil(t, err)
}
