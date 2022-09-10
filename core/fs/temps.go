package fs

import (
	"os"

	"github.com/zeromicro/go-zero/core/hash"
)

// TempFileWithText creates the temporary file with the given content,
// and returns the opened *os.File instance.
// The file is kept as open, the caller should close the file handle,
// and remove the file by name.
func TempFileWithText(text string) (*os.File, error) {
	tmpfile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(text)))
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(tmpfile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return nil, err
	}

	return tmpfile, nil
}

// TempFilenameWithText creates the file with the given content,
// and returns the filename (full path).
// The caller should remove the file after use.
func TempFilenameWithText(text string) (string, error) {
	tmpfile, err := TempFileWithText(text)
	if err != nil {
		return "", err
	}

	filename := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		return "", err
	}

	return filename, nil
}
