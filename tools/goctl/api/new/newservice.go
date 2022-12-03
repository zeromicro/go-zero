package new

import (
	_ "embed"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed api.tpl
var apiTemplate string

var (
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarBoolErrorTranslate describes whether to translate error
	VarBoolErrorTranslate bool
	// VarBoolUseCasbin describe whether to use Casbin
	VarBoolUseCasbin bool
	// VarBoolUseI18n describe whether to use i18n
	VarBoolUseI18n bool
	// VarStringGoZeroVersion describe the version of Go Zero
	VarStringGoZeroVersion string
	// VarStringToolVersion describe the version of Simple Admin Tools
	VarStringToolVersion string
	// VarModuleName describe the module name
	VarModuleName string
	// VarIntServicePort describe the service port exposed
	VarIntServicePort int
)

// CreateServiceCommand fast create service
func CreateServiceCommand(args []string) error {
	dirName := args[0]
	if len(VarStringStyle) == 0 {
		VarStringStyle = conf.DefaultFormat
	}
	if strings.Contains(dirName, "-") {
		return errors.New("api new command service name not support strikethrough, because this will used by function name")
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	dirName = filepath.Base(filepath.Clean(abs))
	filename := dirName + ".api"
	apiFilePath := filepath.Join(abs, "desc", filename)

	err = pathx.MkdirIfNotExist(filepath.Join(abs, "desc"))
	if err != nil {
		return err
	}

	fp, err := os.Create(apiFilePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	if len(VarStringRemote) > 0 {
		repo, _ := util.CloneIntoGitHome(VarStringRemote, VarStringBranch)
		if len(repo) > 0 {
			VarStringHome = repo
		}
	}

	if len(VarStringHome) > 0 {
		pathx.RegisterGoctlHome(VarStringHome)
	}

	text, err := pathx.LoadTemplate(category, apiTemplateFile, apiTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("template").Parse(text))
	if err := t.Execute(fp, map[string]string{
		"name": dirName,
	}); err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(abs, "desc", "base.api"), []byte(baseApiTmpl), os.ModePerm)
	if err != nil {
		return err
	}

	allApiFile, err := os.Create(filepath.Join(abs, "desc", "all.api"))
	if err != nil {
		return err
	}
	defer allApiFile.Close()

	allTpl := template.Must(template.New("allApuTemplate").Parse(allApiTmpl))
	if err := allTpl.Execute(allApiFile, map[string]string{
		"name": dirName,
	}); err != nil {
		return err
	}

	genCtx := &gogen.GenContext{
		GoZeroVersion: VarStringGoZeroVersion,
		ToolVersion:   VarStringToolVersion,
		UseCasbin:     VarBoolUseCasbin,
		UseI18n:       VarBoolUseI18n,
		TransErr:      VarBoolErrorTranslate,
		ModuleName:    VarModuleName,
		Port:          VarIntServicePort,
	}

	err = gogen.DoGenProject(apiFilePath, abs, VarStringStyle, genCtx)
	return err
}
