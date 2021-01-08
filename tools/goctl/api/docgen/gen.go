package docgen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/urfave/cli"
)

var docDir = "doc"

func DocCommand(c *cli.Context) error {
	dir := c.String("dir")
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	files, err := filePathWalkDir(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("dir %s not exist", dir))
	}

	err = os.RemoveAll(dir + "/" + docDir + "/")
	if err != nil {
		return err
	}
	for _, f := range files {
		api, err := parser.Parse(f)
		if err != nil {
			return errors.New(fmt.Sprintf("parse file: %s, err: %s", f, err.Error()))
		}

		index := strings.Index(f, dir)
		if index < 0 {
			continue
		}
		dst := dir + "/" + docDir + f[index+len(dir):]
		index = strings.LastIndex(dst, "/")
		if index < 0 {
			continue
		}
		dir := dst[:index]
		genDoc(api, dir, strings.Replace(dst[index+1:], ".api", ".md", 1))
	}
	return nil
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".api") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
