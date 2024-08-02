package pathx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
)

func TestGetTemplateDir(t *testing.T) {
	category := "foo"
	t.Run("before_have_templates", func(t *testing.T) {
		home := t.TempDir()
		RegisterGoctlHome("")
		RegisterGoctlHome(home)
		v := version.GetGoctlVersion()
		dir := filepath.Join(home, v, category)
		err := MkdirIfNotExist(dir)
		if err != nil {
			return
		}
		tempFile := filepath.Join(dir, "bar.txt")
		err = os.WriteFile(tempFile, []byte("foo"), os.ModePerm)
		if err != nil {
			return
		}
		templateDir, err := GetTemplateDir(category)
		if err != nil {
			return
		}
		assert.Equal(t, dir, templateDir)
		RegisterGoctlHome("")
	})

	t.Run("before_has_no_template", func(t *testing.T) {
		home := t.TempDir()
		RegisterGoctlHome("")
		RegisterGoctlHome(home)
		dir := filepath.Join(home, category)
		err := MkdirIfNotExist(dir)
		if err != nil {
			return
		}
		templateDir, err := GetTemplateDir(category)
		if err != nil {
			return
		}
		assert.Equal(t, dir, templateDir)
	})

	t.Run("default", func(t *testing.T) {
		RegisterGoctlHome("")
		dir, err := GetTemplateDir(category)
		if err != nil {
			return
		}
		assert.Contains(t, dir, version.BuildVersion)
	})
}

func TestGetGitHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	actual, err := GetGitHome()
	if err != nil {
		return
	}

	expected := filepath.Join(homeDir, goctlDir, gitDir)
	assert.Equal(t, expected, actual)
}

func TestGetGoctlHome(t *testing.T) {
	t.Run("goctl_is_file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "a.tmp")
		backupTempFile := tmpFile + ".old"
		err := os.WriteFile(tmpFile, nil, 0o666)
		if err != nil {
			return
		}
		RegisterGoctlHome(tmpFile)
		home, err := GetGoctlHome()
		if err != nil {
			return
		}
		info, err := os.Stat(home)
		assert.Nil(t, err)
		assert.True(t, info.IsDir())

		_, err = os.Stat(backupTempFile)
		assert.Nil(t, err)
	})

	t.Run("goctl_is_dir", func(t *testing.T) {
		RegisterGoctlHome("")
		dir := t.TempDir()
		RegisterGoctlHome(dir)
		home, err := GetGoctlHome()
		assert.Nil(t, err)
		assert.Equal(t, dir, home)
	})
}
