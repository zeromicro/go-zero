package gogen

import (
	"fmt"
	goformat "go/format"
	"io"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	goctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/project"
)

func getParentPackage(dir string) (string, error) {
	p, err := project.Prepare(dir, false)
	if err != nil {
		return "", err
	}

	if len(p.GoMod.Path) > 0 {
		goModePath := filepath.Clean(filepath.Dir(p.GoMod.Path))
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		parent := filepath.Clean(goctlutil.JoinPackages(p.GoMod.Module, absPath[len(goModePath):]))
		parent = strings.ReplaceAll(parent, "\\", "/")
		return parent, nil
	}

	return p.Package, nil
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func writeProperty(writer io.Writer, name, tp, tag, comment string, indent int) error {
	writeIndent(writer, indent)
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp, tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp, tag)
	}
	return err
}

func getAuths(api *spec.ApiSpec) []string {
	authNames := collection.NewSet()
	for _, g := range api.Service.Groups {
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "jwt"); ok {
			authNames.Add(value)
		}
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "signature"); ok {
			authNames.Add(value)
		}
	}
	return authNames.KeysStr()
}

func formatCode(code string) string {
	ret, err := goformat.Source([]byte(code))
	if err != nil {
		return code
	}

	return string(ret)
}
