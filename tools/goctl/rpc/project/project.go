package project

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

var (
	ErrNotFound  = errors.New("not found")
	emptyProject = Project{}
)

type (
	Project struct {
		Path           string
		Name           string
		IsGoModProject bool
		GoPath         string
		GoMod          string
		LibraryPath    string
	}
)

func Info() (Project, error) {
	_, err := exec.LookPath("go")
	if err != nil {
		return emptyProject, err
	}
	stdout, err := execx.Run("go", "env")
	if err != nil {
		return emptyProject, err
	}
	kvs := strings.Split(stdout, "\n")
	var (
		goMod      string
		goPath     string
		goModCache string
	)
	for _, item := range kvs {
		signFlagIndex := strings.Index(item, "=")
		if signFlagIndex < 0 {
			continue
		}
		key := item[:signFlagIndex]
		value := item[signFlagIndex+1:]
		switch key {
		case "GOMOD":
			goMod = value
		case "GOMODCACHE":
			goModCache = value
		case "GOPATH":
			goPath = value
		default:
			continue
		}
	}
	if len(goMod) != 0 {
		path := filepath.Dir(goMod)
		return Project{
			Path:           filepath.Dir(goMod),
			Name:           filepath.Base(path),
			IsGoModProject: true,
			GoMod:          goMod,
			LibraryPath:    goModCache,
		}, nil
	}
	if len(goPath) != 0 {
		current, err := filepath.Abs(".")
		if err != nil {
			return emptyProject, err
		}
		src := filepath.Join(goPath, "src")
		index := strings.Index(current, src)
		if index < 0 {
			return emptyProject, fmt.Errorf("%s: expected project is in gopath", current)
		}
		project := strings.TrimPrefix(current, src+string(filepath.Separator))
		dir := filepath.Dir(project)
		if dir == "." {
			return Project{
				Path:        filepath.Join(src, project),
				Name:        project,
				GoPath:      goPath,
				LibraryPath: src,
			}, nil
		}
	}
	return emptyProject, ErrNotFound
}
