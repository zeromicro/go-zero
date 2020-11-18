package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

func TestRpcGenerate(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	err := dispatcher.Prepare()
	if err != nil {
		logx.Error(err)
		return
	}
	projectName := stringx.Rand()
	g := NewRpcGenerator(dispatcher, namingLower)

	// case go path
	src := filepath.Join(build.Default.GOPATH, "src")
	_, err = os.Stat(src)
	if err != nil {
		return
	}

	projectDir := filepath.Join(src, projectName)
	srcDir := projectDir
	defer func() {
		_ = os.RemoveAll(srcDir)
	}()
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// case go mod
	workDir := t.TempDir()
	name := filepath.Base(workDir)
	_, err = execx.Run("go mod init "+name, workDir)
	if err != nil {
		logx.Error(err)
		return
	}

	projectDir = filepath.Join(workDir, projectName)
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// case not in go mod and go path
	err = g.Generate("./test.proto", projectDir, []string{src})
	assert.Nil(t, err)
	_, err = execx.Run("go test "+projectName, projectDir)
	if err != nil {
		assert.Contains(t, err.Error(), "not in GOROOT")
	}

	// invalid directory
	projectDir = filepath.Join(t.TempDir(), ".....")
	err = g.Generate("./test.proto", projectDir, nil)
	assert.NotNil(t, err)
}
