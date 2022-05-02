package docker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	dockerfileName = "Dockerfile"
	etcDir         = "etc"
	yamlEtx        = ".yaml"
	cstOffset      = 60 * 60 * 8 // 8 hours offset for Chinese Standard Time
)

// Docker describes a dockerfile
type Docker struct {
	Chinese     bool
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
			fmt.Println(aurora.Green("Done."))
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

	if len(goFile) == 0 {
		return errors.New("-go can't be empty")
	}

	if !pathx.FileExists(goFile) {
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
	projPath, err := getFilePath(filepath.Dir(goFile))
	if err != nil {
		return err
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

	_, offset := time.Now().Zone()
	t := template.Must(template.New("dockerfile").Parse(text))
	return t.Execute(out, Docker{
		Chinese:     offset == cstOffset,
		GoRelPath:   projPath,
		GoFile:      goFile,
		ExeFile:     pathx.FileNameWithoutExt(filepath.Base(goFile)),
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
