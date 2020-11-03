package ctx

import (
	"errors"
)

var moduleCheckErr = errors.New("the work directory must be found in the go mod or the $GOPATH")

type (
	ProjectContext struct {
		WorkDir string
		// Name is the root name of project
		// eg: go-zero、greet
		Name string
		// Path is the module path,which is module value if it's go mod project,
		// or else is the root name of project
		// eg: github.com/tal-tech/go-zero、greet
		Path string
		// Dir is the path of project
		// eg: /Users/keson/goland/go/go-zero、/Users/keson/go/src/greet
		Dir string
	}
)

// Background checks the project which module belongs to,and returns the path and module.
// workDir parameter is the directory of the source of generating code,
// where can be found the project path and the project module,
func Background(workDir string) (*ProjectContext, error) {
	isGoMod, err := IsGoMod(workDir)
	if err != nil {
		return nil, err
	}

	if isGoMod {
		return projectFromGoMod(workDir)
	}
	return projectFromGoPath(workDir)
}
