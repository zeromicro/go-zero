package mongo

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/mongo/generate"
	file "github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringSliceType describes a golang data structure name for mongo.
	VarStringSliceType []string
	// VarStringDir describes an output directory.
	VarStringDir string
	//VarBoolAutoPackage Default false;
	//  - Whether to automatically use the name of the directory where the build file is located as the package name
	//  - It is mainly used to be compatible with previous versions,If this parameter is set to false, the package name is set to model by default
	VarBoolAutoPackage bool
	//VarStringPackage Manually specify the package name, which takes precedence over VarBoolPkgAuto
	VarStringPackage string
	// VarBoolCache describes whether cache is enabled.
	VarBoolCache bool
	// VarBoolEasy  describes whether to generate Collection Name in the code for easy declare.
	VarBoolEasy bool
	// VarStringStyle describes the style.
	VarStringStyle string
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
)

// Action provides the entry for goctl mongo code generation.
func Action(_ *cobra.Command, _ []string) error {
	tp := VarStringSliceType
	c := VarBoolCache
	easy := VarBoolEasy
	o := strings.TrimSpace(VarStringDir)
	s := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch

	if len(remote) > 0 {
		repo, _ := file.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(tp) == 0 {
		return errors.New("missing type")
	}

	cfg, err := config.NewConfig(s)
	if err != nil {
		return err
	}

	a, err := filepath.Abs(o)
	if err != nil {
		return err
	}
	if err = pathx.MkdirIfNotExist(a); err != nil {
		return err
	}

	var pkg string
	if VarStringPackage != "" {
		pkg = VarStringPackage
	} else if VarBoolAutoPackage {
		pkg = file.SafeString(filepath.Base(a))
	} else {
		pkg = "model"
	}

	return generate.Do(&generate.Context{
		Types:       tp,
		Cache:       c,
		Easy:        easy,
		Output:      a,
		Cfg:         cfg,
		PackageName: pkg,
	})
}
