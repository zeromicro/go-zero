package gogen

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const typesFile = "types"

//go:embed types.tpl
var typesTemplate string

// BuildTypes gen types to string
func BuildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}

	return builder.String(), nil
}

func removeTypeFromDefault(tp spec.Type, group string, groupTypes map[string]map[string]spec.Type) map[string]map[string]spec.Type {
	switch val := tp.(type) {
	case spec.DefineStruct:
		typeName := util.Title(tp.Name())
		defaultGroups, ok := groupTypes[groupTypeDefault]
		if ok {
			delete(defaultGroups, typeName)
			types, ok := groupTypes[group]
			if !ok {
				types = make(map[string]spec.Type)
			}
			types[typeName] = tp
			groupTypes[group] = types
		}
		groupTypes[groupTypeDefault] = defaultGroups
	case spec.PointerType:
		groupTypes = removeTypeFromDefault(val.Type, group, groupTypes)
	case spec.ArrayType:
		groupTypes = removeTypeFromDefault(val.Value, group, groupTypes)
	}
	return groupTypes
}

func genTypesWithGroup(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	groupTypes := make(map[string]map[string]spec.Type)
	for _, v := range api.Types {
		types, ok := groupTypes[groupTypeDefault]
		if !ok {
			types = make(map[string]spec.Type)
		}
		types[util.Title(v.Name())] = v
		groupTypes[groupTypeDefault] = types
	}

	for _, v := range api.Service.Groups {
		group := v.GetAnnotation(groupProperty)
		if len(group) == 0 {
			continue
		}
		for _, v := range v.Routes {
			if v.RequestType != nil {
				groupTypes = removeTypeFromDefault(v.RequestType, group, groupTypes)
			}
			if v.ResponseType != nil {
				groupTypes = removeTypeFromDefault(v.ResponseType, group, groupTypes)
			}
		}
	}

	for group, typeGroup := range groupTypes {
		var types []spec.Type
		for _, v := range typeGroup {
			types = append(types, v)
		}
		sort.Slice(types, func(i, j int) bool {
			return types[i].Name() < types[j].Name()
		})

		if err := writeTypes(dir, rootPkg, group, cfg, types); err != nil {
			return err
		}
	}

	return nil
}

func containsBaseType(types []spec.Type) bool {
	for _, tp := range types {
		switch val := tp.(type) {
		case spec.DefineStruct:
			for _, member := range val.Members {
				if member.IsInline {
					return true
				}
				if isStructType(member.Type) {
					return true
				}
			}
		}
	}
	return false
}

func isStructType(tp spec.Type) bool {
	switch val := tp.(type) {
	case spec.DefineStruct:
		return true
	case spec.PointerType:
		return isStructType(val.Type)
	case spec.ArrayType:
		return isStructType(val.Value)
	case spec.MapType:
		return isStructType(val.Value)
	default:
		return false
	}
}

func writeTypes(dir, parentPkg, group string, cfg *config.Config, types []spec.Type) error {
	if len(types) == 0 {
		return nil
	}
	val, err := BuildTypes(types)
	if err != nil {
		return err
	}

	td := typesDir
	var baseFilename string
	if len(group) == 0 {
		baseFilename = typesFile + ".go"
	} else {
		td = path.Join(td, group)
		base := filepath.Base(group)
		baseFilename = base + ".go"
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, baseFilename)
	if err != nil {
		return err
	}

	filename := path.Join(dir, td, typeFilename)
	outputDir := filepath.Dir(filename)
	if err := pathx.MkdirIfNotExist(outputDir); err != nil {
		return err
	}
	_ = os.Remove(filename)

	dirOfBaseTypes := path.Join(dir, typesDir)
	shouldImportBasePkg := dirOfBaseTypes != outputDir && containsBaseType(types)
	var baseTypesPkg string
	if shouldImportBasePkg {
		baseTypesPkg = fmt.Sprintf("import . %q", pathx.JoinPackages(parentPkg, typesDir))
	}
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          td,
		filename:        typeFilename,
		templateName:    "typesTemplate",
		category:        category,
		templateFile:    typesTemplateFile,
		builtinTemplate: typesTemplate,
		data: map[string]any{
			"types":        val,
			"baseTypesPkg": baseTypesPkg,
			"containsTime": false,
			"version":      version.BuildVersion,
		},
	})
}

func genTypes(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	if env.UseExperimental() {
		return genTypesWithGroup(dir, rootPkg, cfg, api)
	}
	return writeTypes(dir, rootPkg, "", cfg, api.Types)
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	_, err := fmt.Fprintf(writer, "type %s struct {\n", util.Title(tp.Name()))
	if err != nil {
		return err
	}

	if err := writeMember(writer, structType.Members); err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "}")
	return err
}

func writeMember(writer io.Writer, members []spec.Member) error {
	for _, member := range members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
				return err
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1); err != nil {
			return err
		}
	}
	return nil
}
