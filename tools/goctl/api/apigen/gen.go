package apigen

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed api.tpl
var apiTemplate string

// ApiCommand create api template file
func ApiCommand(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowAppHelpAndExit(c, 1)
	}

	apiFile := c.String("o")
	if len(apiFile) == 0 {
		return errors.New("missing -o")
	}

	fp, err := pathx.CreateIfNotExist(apiFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	home := c.String("home")
	remote := c.String("remote")
	branch := c.String("branch")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	text, err := pathx.LoadTemplate(category, apiTemplateFile, apiTemplate)
	if err != nil {
		return err
	}

	baseName := pathx.FileNameWithoutExt(filepath.Base(apiFile))
	if strings.HasSuffix(strings.ToLower(baseName), "-api") {
		baseName = baseName[:len(baseName)-4]
	} else if strings.HasSuffix(strings.ToLower(baseName), "api") {
		baseName = baseName[:len(baseName)-3]
	}

	t := template.Must(template.New("etcTemplate").Parse(text))
	if err := t.Execute(fp, map[string]string{
		"gitUser":     getGitName(),
		"gitEmail":    getGitEmail(),
		"serviceName": baseName + "-api",
	}); err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}
