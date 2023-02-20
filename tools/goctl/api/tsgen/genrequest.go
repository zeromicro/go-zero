package tsgen

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed request.ts
var requestTemplate string

func genRequest(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	filename := filepath.Join(abs, "gocliRequest.ts")
	if pathx.FileExists(filename) {
		return nil
	}

	return os.WriteFile(filename, []byte(requestTemplate), 0644)
}
