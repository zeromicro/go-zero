package gen

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

func GenerateDockerfile(goFile string, args ...string) error {
	projPath, err := getFilePath(filepath.Dir(goFile))
	if err != nil {
		return err
	}

	pos := strings.IndexByte(projPath, '/')
	if pos >= 0 {
		projPath = projPath[pos+1:]
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
		"goRelPath":   projPath,
		"goFile":      goFile,
		"exeFile":     util.FileNameWithoutExt(goFile),
		"argument":    builder.String(),
	})
}
