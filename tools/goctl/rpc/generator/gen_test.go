package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

func TestRpcGenerate(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	err := dispatcher.Prepare()
	if err != nil {
		return
	}
	projectName := stringx.Rand()

	// case in go path
	src := filepath.Join(build.Default.GOPATH, "src")
	_, err = os.Stat(src)
	if err == nil {
		g := NewRpcGenerator(dispatcher, namingLower)
		projectDir := filepath.Join(src, projectName)
		err = g.Generate("./test.proto", projectDir, []string{src})
		assert.Nil(t, err)
		_, err = execx.Run("go test "+projectName, projectDir)
		assert.Nil(t, err)
		return
	}

	// case go mod

}
