package ctx

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func TestGOPATH(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := util.MkdirIfNotExist(dir)
	if err != nil {
		return
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	ctx, err := GOPATH(dir)
	assert.Nil(t, err)
	assert.Equal(t, dir, ctx.Dir)
	assert.Equal(t, projectName, ctx.Path)
}

func TestGOPATH_NotInGoSrc(t *testing.T) {
	dft := build.Default
	gp := dft.GOPATH
	if len(gp) == 0 {
		return
	}
	projectName := stringx.Rand()
	dir := filepath.Join(gp, "src", projectName)
	err := util.MkdirIfNotExist(dir)
	if err != nil {
		return
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	_, err = GOPATH("testPath")
	assert.NotNil(t, err)
}
