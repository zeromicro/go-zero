package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
)

func TestRpcGenerator_Generate_CaseNilImport(t *testing.T) {
	dispatcher := NewDefaultGenerator()
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

func TestRpcGenerator_Generate_Case_Option(t *testing.T) {
	dispatcher := NewDefaultGenerator()
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

func TestRpcGenerator_Generate_Case_Word_Option(t *testing.T) {
	dispatcher := NewDefaultGenerator()
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

// test keyword go
func TestRpcGenerator_Generate_Case_Go_Option(t *testing.T) {
	dispatcher := NewDefaultGenerator()
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

func TestRpcGenerator_Generate_CaseImport(t *testing.T) {
	dispatcher := NewDefaultGenerator()
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
