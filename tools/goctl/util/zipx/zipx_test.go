package zipx

import (
	"archive/zip"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestZip(t *testing.T, files map[string]struct {
	content string
	mode    os.FileMode
}) string {
	t.Helper()

	zipPath := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(zipPath)
	require.NoError(t, err)
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for name, file := range files {
		header := &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}
		header.SetMode(file.mode)
		writer, err := w.CreateHeader(header)
		require.NoError(t, err)
		_, err = writer.Write([]byte(file.content))
		require.NoError(t, err)
	}

	return zipPath
}

func TestUnpacking(t *testing.T) {
	dest := t.TempDir()
	zipPath := createTestZip(t, map[string]struct {
		content string
		mode    os.FileMode
	}{
		"hello.txt": {content: "hello world", mode: 0644},
		"skip.txt":  {content: "should be skipped", mode: 0644},
	})

	err := Unpacking(zipPath, dest, func(f *zip.File) bool {
		return f.Name == "hello.txt"
	})
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dest, "hello.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(content))

	_, err = os.Stat(filepath.Join(dest, "skip.txt"))
	assert.True(t, os.IsNotExist(err))
}

func TestUnpackingPreservesExecutablePermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions not applicable on Windows")
	}

	dest := t.TempDir()
	zipPath := createTestZip(t, map[string]struct {
		content string
		mode    os.FileMode
	}{
		"bin/mybinary": {content: "binary content", mode: 0755},
		"readme.txt":   {content: "readme", mode: 0644},
	})

	err := Unpacking(zipPath, dest, func(f *zip.File) bool {
		return true
	})
	require.NoError(t, err)

	info, err := os.Stat(filepath.Join(dest, "mybinary"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0111, "executable bit should be set")

	info, err = os.Stat(filepath.Join(dest, "readme.txt"))
	require.NoError(t, err)
	assert.Zero(t, info.Mode()&0111, "executable bit should not be set for regular files")
}

func TestUnpackingInvalidZip(t *testing.T) {
	err := Unpacking("/nonexistent/path.zip", t.TempDir(), func(f *zip.File) bool {
		return true
	})
	assert.Error(t, err)
}

func TestUnpackingAllFilesFiltered(t *testing.T) {
	dest := t.TempDir()
	zipPath := createTestZip(t, map[string]struct {
		content string
		mode    os.FileMode
	}{
		"a.txt": {content: "a", mode: 0644},
	})

	err := Unpacking(zipPath, dest, func(f *zip.File) bool {
		return false
	})
	require.NoError(t, err)

	entries, err := os.ReadDir(dest)
	assert.NoError(t, err)
	assert.Empty(t, entries)
}
