package docgen

import (
	"errors"
	"fmt"
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
		return errors.New("missing -dir")
	}

	outputDir := c.String("o")
	if len(outputDir) == 0 {
		var err error
		outputDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	if !util.FileExists(dir) {
		return fmt.Errorf("dir %s not exsit", dir)
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	files, err := filePathWalkDir(dir)
	if err != nil {
		return err
	}

	for _, p := range files {
		api, err := parser.Parse(p)
		if err != nil {
			return fmt.Errorf("parse file: %s, err: %s", p, err.Error())
		}

		api.Service = api.Service.JoinPrefix()
		err = genDoc(api, filepath.Dir(filepath.Join(outputDir, p[len(dir):])),
			strings.Replace(p[len(filepath.Dir(p)):], ".api", ".md", 1))
		if err != nil {
			return err
		}
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
