package docker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/gen"
	"github.com/urfave/cli"
)

const (
	etcDir  = "etc"
	yamlEtx = ".yaml"
)

func DockerCommand(c *cli.Context) error {
	goFile := c.String("go")
	if len(goFile) == 0 {
		return errors.New("-go can't be empty")
	}

	cfg, err := findConfig(goFile, etcDir)
	if err != nil {
		return err
	}

	return gen.GenerateDockerfile(goFile, "-f", "etc/"+cfg)
}

func findConfig(file, dir string) (string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == yamlEtx {
				files = append(files, f.Name())
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", errors.New("no yaml file")
	}

	name := strings.TrimSuffix(filepath.Base(file), ".go")
	for _, f := range files {
		if strings.Index(f, name) == 0 {
			return f, nil
		}
	}

	return files[0], nil
}
