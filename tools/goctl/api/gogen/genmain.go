package gogen

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const mainTemplate = `package main

import (
	"flag"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.json", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	server.Start()
}
`

func genMain(dir string, api *spec.ApiSpec) error {
	name := strings.ToLower(api.Service.Name)
	if strings.HasSuffix(name, "-api") {
		name = strings.ReplaceAll(name, "-api", "")
	}
	goFile := name + ".go"
	fp, created, err := util.MaybeCreateFile(dir, "", goFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	t := template.Must(template.New("mainTemplate").Parse(mainTemplate))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"importPackages": genMainImports(parentPkg),
		"serviceName":    api.Service.Name,
	})
	if err != nil {
		return nil
	}
	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}

func genMainImports(parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s/core/conf\"", vars.ProjectOpenSourceUrl),
		fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceUrl),
	}
	imports = append(imports, fmt.Sprintf("\"%s\"", path.Join(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", path.Join(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", path.Join(parentPkg, contextDir)))
	sort.Strings(imports)
	return strings.Join(imports, "\n\t")
}
