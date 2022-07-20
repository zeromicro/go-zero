package generate

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func TestTemplate(t *testing.T) {
	tempDir := t.TempDir()
	pathx.RegisterGoctlHome(tempDir)
	t.Cleanup(func() {
		pathx.RegisterGoctlHome("")
	})

	t.Run("Category", func(t *testing.T) {
		assert.Equal(t, category, Category())
	})

	t.Run("Clean", func(t *testing.T) {
		err := Clean()
		assert.NoError(t, err)
	})

	t.Run("Templates", func(t *testing.T) {
		err := Templates()
		assert.NoError(t, err)
		assert.True(t, pathx.FileExists(filepath.Join(tempDir, category, modelTemplateFile)))
	})

	t.Run("RevertTemplate", func(t *testing.T) {
		assert.NoError(t, RevertTemplate(modelTemplateFile))
		assert.Error(t, RevertTemplate("foo"))
	})

	t.Run("Update", func(t *testing.T) {
		assert.NoError(t, Update())
	})
}
