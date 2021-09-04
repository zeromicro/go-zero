package docgen

import (
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

// DocCommand generate markdown doc file
func DocCommand(c *cli.Context) error {
	dir := c.String("dir")
	if len(dir) == 0 {
		errorx.Must(errors.New("missing -dir"))
	}

	outputDir := c.String("o")
	if len(outputDir) == 0 {
		var err error
		outputDir, err = os.Getwd()
		errorx.Must(err)
	}

	if !util.FileExists(dir) {
		errorx.Must(fmt.Errorf("dir %s not exsit", dir))
	}

	dir, err := filepath.Abs(dir)
	errorx.Must(err)

	files, err := filePathWalkDir(dir)
	errorx.Must(err)

	for _, path := range files {
		api, err := parser.Parse(path)
		errorx.Must(fmt.Errorf("parse file: %s, err: %s", path, err.Error()))
		err = genDoc(api, filepath.Dir(filepath.Join(outputDir, path[len(dir):])),
			strings.Replace(path[len(filepath.Dir(path)):], ".api", ".md", 1))
		errorx.Must(err)
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
