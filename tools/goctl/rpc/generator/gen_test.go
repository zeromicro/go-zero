package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

func TestRpcGenerateCaseNilImport(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_stream.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseWordOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_word_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

// test keyword go
func TestRpcGenerateCaseGoOption(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_go_option.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}

func TestRpcGenerateCaseImport(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_import.proto", abs, []string{"./base"})
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.True(t, func() bool {
			return strings.Contains(err.Error(), "package base is not in GOROOT")
		}())
	}
}

func TestRpcGenerateCaseServiceRpcNamingSnake(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	if err := dispatcher.Prepare(); err == nil {
		g := NewRpcGenerator(dispatcher)
		abs, err := filepath.Abs("./test")
		assert.Nil(t, err)

		err = g.Generate("./test_service_rpc_naming_snake.proto", abs, nil)
		defer func() {
			_ = os.RemoveAll(abs)
		}()
		assert.Nil(t, err)

		_, err = execx.Run("go test "+abs, abs)
		assert.Nil(t, err)
	}
}
