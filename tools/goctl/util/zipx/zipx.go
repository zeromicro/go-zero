package zipx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func Unpacking(name, destPath string, mapper func(f *zip.File) bool) error {
	r, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		ok := mapper(file)
		if ok {
			err = fileCopy(file, destPath)
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
	filename := filepath.Join(destPath, filepath.Base(file.Name))
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
