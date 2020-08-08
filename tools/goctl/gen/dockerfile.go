package gen

import (
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

func GenerateDockerfile(goFile string, args ...string) error {
	relPath, err := util.PathFromGoSrc()
	if err != nil {
		return err
	}

	out, err := util.CreateIfNotExist("Dockerfile")
	if err != nil {
		return err
	}
	defer out.Close()

	var builder strings.Builder
	for _, arg := range args {
		builder.WriteString(`, "` + arg + `"`)
	}

	t := template.Must(template.New("dockerfile").Parse(dockerTemplate))
	return t.Execute(out, map[string]string{
		"projectName": vars.ProjectName,
		"goRelPath":   relPath,
		"goFile":      goFile,
		"exeFile":     util.FileNameWithoutExt(goFile),
		"argument":    builder.String(),
	})
}
