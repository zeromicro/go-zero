package zipx

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func Unpacking(name, destPath string, mapper func(f *zip.File) bool) error {
	r, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer r.Close()

	destAbsPath, err := filepath.Abs(destPath)
	if err != nil {
		return err
	}

	for _, file := range r.File {
		ok := mapper(file)
		if ok {
			err = fileCopy(file, destAbsPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func fileCopy(file *zip.File, destPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Ensure the file path does not contain directory traversal elements
	if strings.Contains(file.Name, "..") {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	abs, err := filepath.Abs(file.Name)
	if err != nil {
		return err
	}

	filename := filepath.Join(destPath, filepath.Base(abs))
	dir := filepath.Dir(filename)
	err = pathx.MkdirIfNotExist(dir)
	if err != nil {
		return err
	}

	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	return err
}
