package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func MaybeCreateFile(dir, subdir, file string) (fp *os.File, created bool, err error) {
	logx.Must(util.MkdirIfNotExist(path.Join(dir, subdir)))
	fpath := path.Join(dir, subdir, file)
	if util.FileExists(fpath) {
		fmt.Printf("%s exists, ignored generation\n", fpath)
		return nil, false, nil
	}

	fp, err = util.CreateIfNotExist(fpath)
	created = err == nil
	return
}

func ClearAndOpenFile(fpath string) (*os.File, error) {
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString("")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func WrapErr(err error, message string) error {
	return errors.New(message + ", " + err.Error())
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ComponentName(api *spec.ApiSpec) string {
	name := api.Service.Name
	if strings.HasSuffix(name, "-api") {
		return name[:len(name)-4] + "Components"
	}
	return name + "Components"
}
