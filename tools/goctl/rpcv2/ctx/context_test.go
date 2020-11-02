package ctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackground(t *testing.T) {
	workDir := "."
	ctx, err := Background(workDir)
	assert.Nil(t, err)
	assert.True(t, true, func() bool {
		return len(ctx.Dir) != 0 && len(ctx.Path) != 0
	}())
}

func TestBackground_NilWorkDir(t *testing.T) {
	workDir := ""
	_, err := Background(workDir)
	assert.NotNil(t, err)
}
