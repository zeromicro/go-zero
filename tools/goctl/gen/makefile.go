package gen

import (
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func GenerateMakefile(goFile, namespace string) error {
	relPath, err := util.PathFromGoSrc()
	if err != nil {
		return err
	}

	movePath, err := getMovePath()
	if err != nil {
		return err
	}

	out, err := util.CreateIfNotExist("Makefile")
	if err != nil {
		return err
	}
	defer out.Close()

	t := template.Must(template.New("makefile").Parse(makefileTemplate))
	return t.Execute(out, map[string]string{
		"rootRelPath": movePath,
		"relPath":     relPath,
		"exeFile":     util.FileNameWithoutExt(goFile),
		"namespace":   namespace,
	})
}

func getMovePath() (string, error) {
	relPath, err := util.PathFromGoSrc()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for range strings.Split(relPath, "/") {
		builder.WriteString("../")
	}

	if move := builder.String(); len(move) == 0 {
		return ".", nil
	} else {
		return move, nil
	}
}
