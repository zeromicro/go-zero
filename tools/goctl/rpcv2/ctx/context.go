package ctx

import (
	"fmt"
)

type (
	ProjectContext struct {
		WorkDir string
		// the project name
		// eg: go-zero、greet
		Name string
		// the module path,
		// eg:github.com/tal-tech/go-zero、greet
		Path string
		// the path of current project
		// eg:/Users/keson/goland/go/go-zero、/Users/keson/go/src/greet
		Dir string
	}
)

// workDir is the directory of the source of generating code,
// where can be found the project path and the project dir,
func Background(workDir string) (ret ProjectContext, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("the project must using go mod or in $GOPATH/src")
		}
	}()
	isGoMod, err := IsGoMod(workDir)
	if err != nil {
		ret = ProjectContext{}
		return
	}

	if isGoMod {
		ret, err = GOMOD(workDir)
		return
	}
	ret, err = GOPATH(workDir)
	return
}
