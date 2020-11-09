package docker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
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

	return generateDockerfile(goFile, "-f", "etc/"+cfg)
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

func generateDockerfile(goFile string, args ...string) error {
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

	text, err := ctlutil.LoadTemplate(category, dockerTemplateFile, dockerTemplate)
	if err != nil {
		return err
	}

	var builder strings.Builder
	for _, arg := range args {
		builder.WriteString(`, "` + arg + `"`)
	}

	t := template.Must(template.New("dockerfile").Parse(text))
	return t.Execute(out, map[string]string{
		"goRelPath": projPath,
		"goFile":    goFile,
		"exeFile":   util.FileNameWithoutExt(filepath.Base(goFile)),
		"argument":  builder.String(),
	})
}

func getFilePath(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projPath, ok := util.FindGoModPath(filepath.Join(wd, file))
	if !ok {
		projPath, err = util.PathFromGoSrc()
		if err != nil {
			return "", errors.New("no go.mod found, or not in GOPATH")
		}
	}

	return projPath, nil
}
