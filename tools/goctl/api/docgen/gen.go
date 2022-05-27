package docgen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringDir describes a directory.
	VarStringDir string
	// VarStringOutput describes an output directory.
	VarStringOutput string
)

// DocCommand generate Markdown doc file
func DocCommand(_ *cobra.Command, _ []string) error {
	dir := VarStringDir
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	outputDir := VarStringOutput
	if len(outputDir) == 0 {
		var err error
		outputDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	if !pathx.FileExists(dir) {
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
			return fmt.Errorf("parse file: %s, err: %w", p, err)
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
