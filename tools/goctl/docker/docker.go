package docker

import (
	"fmt"
	"os"
	"text/template"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/suyuan32/knife/core/io/filex"

	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	dockerfileName = "Dockerfile"
	yamlEtx        = ".yaml"
)

// Docker describes a dockerfile
type Docker struct {
	Chinese     bool
	GoMainFrom  string
	GoRelPath   string
	GoFile      string
	ServiceName string
	ServiceType string
	BaseImage   string
	HasPort     bool
	Port        int
	Image       string
	HasTimezone bool
	Timezone    string
}

type GenContext struct {
	Home        string
	Image       string
	Remote      string
	Branch      string
	TimeZone    string
	Base        string
	Port        int
	ServiceName string
	ServiceType string
}

// dockerCommand provides the entry for goctl docker
func dockerCommand(_ *cobra.Command, _ []string) (err error) {
	defer func() {
		if err == nil {
			fmt.Println(aurora.Green("Done."))
		}
	}()

	home := varStringHome
	remote := varStringRemote
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, varStringBranch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	g := &GenContext{
		Home:        home,
		Image:       varStringImage,
		Remote:      remote,
		Branch:      varStringBranch,
		TimeZone:    varStringTZ,
		Base:        varStringBase,
		Port:        varIntPort,
		ServiceType: varServiceType,
		ServiceName: varServiceName,
	}

	if err := generateDockerfile(g); err != nil {
		return err
	}

	return nil
}

func generateDockerfile(g *GenContext) error {
	var projPath string
	var err error

	if len(projPath) == 0 {
		projPath = "."
	}

	if filex.Exist(dockerfileName) == nil {
		err = os.Remove(dockerfileName)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(dockerfileName)
	if err != nil {
		return err
	}

	text, err := pathx.LoadTemplate(category, dockerTemplateFile, dockerTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("dockerfile").Parse(text))
	return t.Execute(out, Docker{
		Chinese:     env.InChina(),
		GoRelPath:   projPath,
		ServiceName: g.ServiceName,
		ServiceType: g.ServiceType,
		BaseImage:   g.Base,
		HasPort:     g.Port > 0,
		Port:        g.Port,
		Image:       g.Image,
		HasTimezone: len(g.TimeZone) > 0,
		Timezone:    g.TimeZone,
	})
}
