package gogen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
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

func getTypeName(tp spec.Type) string {
	if tp == nil {
		return ""
	}
	switch val := tp.(type) {
	case spec.DefineStruct:
		typeName := util.Title(tp.Name())
		return typeName
	case spec.PointerType:
		return getTypeName(val.Type)
	case spec.ArrayType:
		return getTypeName(val.Value)
	}
	return ""
}

func genTypesWithGroup(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	groupTypes := make(map[string]map[string]spec.Type)
	typesBelongToFiles := make(map[string]*collection.Set)

	for _, v := range api.Service.Groups {
		group := v.GetAnnotation(groupProperty)
		if len(group) == 0 {
			group = groupTypeDefault
		}
		// convert filepath to Identifier name spec.
		group = strings.TrimPrefix(group, "/")
		group = strings.TrimSuffix(group, "/")
		group = util.SafeString(group)
		for _, v := range v.Routes {
			requestTypeName := getTypeName(v.RequestType)
			responseTypeName := getTypeName(v.ResponseType)
			requestTypeFileSet, ok := typesBelongToFiles[requestTypeName]
			if !ok {
				requestTypeFileSet = collection.NewSet()
			}
			if len(requestTypeName) > 0 {
				requestTypeFileSet.AddStr(group)
				typesBelongToFiles[requestTypeName] = requestTypeFileSet
			}

			responseTypeFileSet, ok := typesBelongToFiles[responseTypeName]
			if !ok {
				responseTypeFileSet = collection.NewSet()
			}
			if len(responseTypeName) > 0 {
				responseTypeFileSet.AddStr(group)
				typesBelongToFiles[responseTypeName] = responseTypeFileSet
			}
		}
	}

	typesInOneFile := make(map[string]*collection.Set)
	for typeName, fileSet := range typesBelongToFiles {
		count := fileSet.Count()
		switch {
		case count == 0: // it means there has no structure type or no request/response body
			continue
		case count == 1: // it means a structure type used in only one group.
			groupName := fileSet.KeysStr()[0]
			typeSet, ok := typesInOneFile[groupName]
			if !ok {
				typeSet = collection.NewSet()
			}
			typeSet.AddStr(typeName)
			typesInOneFile[groupName] = typeSet
		default: // it means this type is used in multiple groups.
			continue
		}
	}

	for _, v := range api.Types {
		typeName := util.Title(v.Name())
		groupSet, ok := typesBelongToFiles[typeName]
		var typeCount int
		if !ok {
			typeCount = 0
		} else {
			typeCount = groupSet.Count()
		}

		if typeCount == 0 { // not belong to any group
			types, ok := groupTypes[groupTypeDefault]
			if !ok {
				types = make(map[string]spec.Type)
			}
			types[typeName] = v
			groupTypes[groupTypeDefault] = types
			continue
		}

		if typeCount == 1 { // belong to one group
			groupName := groupSet.KeysStr()[0]
			types, ok := groupTypes[groupName]
			if !ok {
				types = make(map[string]spec.Type)
			}
			types[typeName] = v
			groupTypes[groupName] = types
			continue
		}

		// belong to multiple groups
		types, ok := groupTypes[groupTypeDefault]
		if !ok {
			types = make(map[string]spec.Type)
		}
		types[typeName] = v
		groupTypes[groupTypeDefault] = types

	}

	for group, typeGroup := range groupTypes {
		var types []spec.Type
		for _, v := range typeGroup {
			types = append(types, v)
		}
		sort.Slice(types, func(i, j int) bool {
			return types[i].Name() < types[j].Name()
		})

		if err := writeTypes(dir, group, cfg, types); err != nil {
			return err
		}
	}

	return nil
}

func writeTypes(dir, baseFilename string, cfg *config.Config, types []spec.Type) error {
	if len(types) == 0 {
		return nil
	}
	val, err := BuildTypes(types)
	if err != nil {
		return err
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, baseFilename)
	if err != nil {
		return err
	}

	typeFilename = typeFilename + ".go"
	filename := path.Join(dir, typesDir, typeFilename)
	_ = os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          typesDir,
		filename:        typeFilename,
		templateName:    "typesTemplate",
		category:        category,
		templateFile:    typesTemplateFile,
		builtinTemplate: typesTemplate,
		data: map[string]any{
			"types":        val,
			"containsTime": false,
			"version":      version.BuildVersion,
		},
	})
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	if VarBoolTypeGroup {
		return genTypesWithGroup(dir, cfg, api)
	}
	return writeTypes(dir, typesFile, cfg, api.Types)
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
