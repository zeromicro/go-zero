package api

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

func TestApi(t *testing.T) {
	fileSystem := os.DirFS("/Users/keson/workspace/go-zero/tools/goctl/example/api")
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext != ".api" {
			return nil
		}
		_, err = parser.Parse(path)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
