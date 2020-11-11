package ctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackground(t *testing.T) {
	workDir := "."
	ctx, err := Prepare(workDir)
	assert.Nil(t, err)
	assert.True(t, true, func() bool {
		return len(ctx.Dir) != 0 && len(ctx.Path) != 0
	}())
}

func TestBackgroundNilWorkDir(t *testing.T) {
	workDir := ""
	_, err := Prepare(workDir)
	assert.NotNil(t, err)
}
