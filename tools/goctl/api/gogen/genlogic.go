package gogen

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const logicTemplate = `package logic

import (
	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext) {{.logic}} {
	return {{.logic}}{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	// todo: add your logic here and delete this line

	{{.returnString}}
}
`

func genLogic(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByRoute(dir, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRoute(dir string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	goFile = goFile + ".go"
	fp, created, err := util.MaybeCreateFile(dir, getLogicFolderPath(group, route), goFile)
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

	imports := genLogicImports(route, parentPkg)
	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseType.Name) > 0 {
		resp := strings.Title(route.ResponseType.Name)
		responseString = "(*types." + resp + ", error)"
		returnString = fmt.Sprintf("return &types.%s{}, nil", resp)
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestType.Name) > 0 {
		requestString = "req " + "types." + strings.Title(route.RequestType.Name)
	}

	text, err := ctlutil.LoadTemplate(category, logicTemplateFile, logicTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("logicTemplate").Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(fp, map[string]string{
		"imports":      imports,
		"logic":        strings.Title(logic),
		"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
		"responseType": responseString,
		"returnString": returnString,
		"request":      requestString,
	})
	if err != nil {
		return err
	}

	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}

func getLogicFolderPath(group spec.Group, route spec.Route) string {
	folder, ok := util.GetAnnotationValue(route.Annotations, "server", groupProperty)
	if !ok {
		folder, ok = util.GetAnnotationValue(group.Annotations, "server", groupProperty)
		if !ok {
			return logicDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func genLogicImports(route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", ctlutil.JoinPackages(parentPkg, contextDir)))
	if len(route.ResponseType.Name) > 0 || len(route.RequestType.Name) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceUrl))
	return strings.Join(imports, "\n\t")
}
