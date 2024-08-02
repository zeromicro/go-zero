package dartgen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

var (
	// VarStringDir describes the directory.
	VarStringDir string
	// VarStringAPI defines the API.
	VarStringAPI string
	// VarStringLegacy describes whether legacy.
	VarStringLegacy bool
	// VarStringHostname defines the hostname.
	VarStringHostname string
	// VarStringSchema defines the scheme.
	VarStringScheme string
)

// DartCommand create dart network request code
func DartCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir
	isLegacy := VarStringLegacy
	hostname := VarStringHostname
	scheme := VarStringScheme
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}
	if len(hostname) == 0 {
		fmt.Println("you could use '-hostname' flag to specify your server hostname")
		hostname = "go-zero.dev"
	}
	if len(scheme) == 0 {
		fmt.Println("you could use '-scheme' flag to specify your server scheme")
		scheme = "http"
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	api.Service = api.Service.JoinPrefix()
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	api.Info.Title = strings.Replace(apiFile, ".api", "", -1)
	logx.Must(genData(dir+"data/", api, isLegacy))
	logx.Must(genApi(dir+"api/", api, isLegacy))
	logx.Must(genVars(dir+"vars/", isLegacy, scheme, hostname))
	if err := formatDir(dir); err != nil {
		logx.Errorf("failed to format, %v", err)
	}
	return nil
}
