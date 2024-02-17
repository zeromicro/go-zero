package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var testTypes = `
	type User struct{}
	type Class struct{}
`

func TestDo(t *testing.T) {
	t.Run("should generate model", func(t *testing.T) {
		cfg, err := config.NewConfig(config.DefaultFormat)
		assert.Nil(t, err)

		tempDir := pathx.MustTempDir()
		typesfile := filepath.Join(tempDir, "types.go")
		err = os.WriteFile(typesfile, []byte(testTypes), 0o666)
		assert.Nil(t, err)

		err = Do(&Context{
			Types:  []string{"User", "Class"},
			Cache:  false,
			Output: tempDir,
			Cfg:    cfg,
		})

		assert.Nil(t, err)
	})

	t.Run("missing config", func(t *testing.T) {
		tempDir := t.TempDir()
		typesfile := filepath.Join(tempDir, "types.go")
		err := os.WriteFile(typesfile, []byte(testTypes), 0o666)
		assert.Nil(t, err)

		err = Do(&Context{
			Types:  []string{"User", "Class"},
			Cache:  false,
			Output: tempDir,
			Cfg:    nil,
		})
		assert.Error(t, err)
	})

	t.Run("invalid config", func(t *testing.T) {
		cfg := &config.Config{NamingFormat: "foo"}
		tempDir := t.TempDir()
		typesfile := filepath.Join(tempDir, "types.go")
		err := os.WriteFile(typesfile, []byte(testTypes), 0o666)
		assert.Nil(t, err)

		err = Do(&Context{
			Types:  []string{"User", "Class"},
			Cache:  false,
			Output: tempDir,
			Cfg:    cfg,
		})
		assert.Error(t, err)
	})
}
