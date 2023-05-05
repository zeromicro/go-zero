package docker

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	dockerfileName = "Dockerfile"
	etcDir         = "etc"
	yamlEtx        = ".yaml"
)

// Docker describes a dockerfile
type Docker struct {
	Chinese     bool
	GoMainFrom  string
	GoRelPath   string
	GoFile      string
	ExeFile     string
	BaseImage   string
	HasPort     bool
	Port        int
	Argument    string
	Version     string
	HasTimezone bool
	Timezone    string
}

// dockerCommand provides the entry for goctl docker
func dockerCommand(_ *cobra.Command, _ []string) (err error) {
	defer func() {
		if err == nil {
			fmt.Println(color.Green.Render("Done."))
		}
	}()

	goFile := varStringGo
	home := varStringHome
	version := varStringVersion
	remote := varStringRemote
	branch := varStringBranch
	timezone := varStringTZ
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(version) > 0 {
		version = version + "-"
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(goFile) > 0 && !pathx.FileExists(goFile) {
		return fmt.Errorf("file %q not found", goFile)
	}

	base := varStringBase
	port := varIntPort
	if _, err := os.Stat(etcDir); os.IsNotExist(err) {
		return generateDockerfile(goFile, base, port, version, timezone)
	}

	cfg, err := findConfig(goFile, etcDir)
	if err != nil {
		return err
	}

	if err := generateDockerfile(goFile, base, port, version, timezone, "-f", "etc/"+cfg); err != nil {
		return err
	}

	projDir, ok := pathx.FindProjectPath(goFile)
	if ok {
		fmt.Printf("Hint: run \"docker build ...\" command in dir:\n    %s\n", projDir)
	}

	return nil
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

func generateDockerfile(goFile, base string, port int, version, timezone string, args ...string) error {
	var projPath string
	var err error
	if len(goFile) > 0 {
		projPath, err = getFilePath(filepath.Dir(goFile))
		if err != nil {
			return err
		}
	}

	if len(projPath) == 0 {
		projPath = "."
	}

	out, err := pathx.CreateIfNotExist(dockerfileName)
	if err != nil {
		return err
	}
	defer out.Close()

	text, err := pathx.LoadTemplate(category, dockerTemplateFile, dockerTemplate)
	if err != nil {
		return err
	}

	var builder strings.Builder
	for _, arg := range args {
		builder.WriteString(`, "` + arg + `"`)
	}

	var exeName string
	if len(varExeName) > 0 {
		exeName = varExeName
	} else if len(goFile) > 0 {
		exeName = pathx.FileNameWithoutExt(filepath.Base(goFile))
	} else {
		absPath, err := filepath.Abs(projPath)
		if err != nil {
			return err
		}

		exeName = filepath.Base(absPath)
	}

	t := template.Must(template.New("dockerfile").Parse(text))
	return t.Execute(out, Docker{
		Chinese:     env.InChina(),
		GoMainFrom:  path.Join(projPath, goFile),
		GoRelPath:   projPath,
		GoFile:      goFile,
		ExeFile:     exeName,
		BaseImage:   base,
		HasPort:     port > 0,
		Port:        port,
		Argument:    builder.String(),
		Version:     version,
		HasTimezone: len(timezone) > 0,
		Timezone:    timezone,
	})
}

func getFilePath(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projPath, ok := pathx.FindGoModPath(filepath.Join(wd, file))
	if !ok {
		projPath, err = pathx.PathFromGoSrc()
		if err != nil {
			return "", errors.New("no go.mod found, or not in GOPATH")
		}

		// ignore project root directory for GOPATH mode
		pos := strings.IndexByte(projPath, os.PathSeparator)
		if pos >= 0 {
			projPath = projPath[pos+1:]
		}
	}

	return projPath, nil
}
