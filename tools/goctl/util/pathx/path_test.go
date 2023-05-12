package pathx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLink(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-zero")
	assert.Nil(t, err)
	symLink := filepath.Join(dir, "test")
	pwd, err := os.Getwd()
	assertError(err, t)

	err = os.Symlink(pwd, symLink)
	assertError(err, t)

	t.Run("linked", func(t *testing.T) {
		ret, err := ReadLink(symLink)
		assert.Nil(t, err)
		assert.Equal(t, pwd, ret)
	})

	t.Run("unlink", func(t *testing.T) {
		ret, err := ReadLink(pwd)
		assert.Nil(t, err)
		assert.Equal(t, pwd, ret)
	})
}

func assertError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
