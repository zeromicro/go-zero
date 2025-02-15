package gen

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen/template"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func BuildTagMembers(tp spec.Type, tagKey string, indent int) ([]*template.ComponentMemberTemplateData, error) {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		if pointType, ok := tp.(spec.PointerType); ok {
			return BuildTagMembers(pointType.Type, tagKey, indent)
		}

		if nestedType, ok := tp.(spec.NestedStruct); ok {
			return BuildTagMembers(spec.DefineStruct(nestedType), tagKey, indent)
		}

		return nil, fmt.Errorf("type %s not supported", tp.Name())
	}

	result := []*template.ComponentMemberTemplateData{}
	members := definedType.GetTagMembers(tagKey)
	for _, member := range members {
		if member.IsInline {
			if r, err := BuildTagMembers(member.Type, tagKey, indent); err != nil {
				return nil, err
			} else {
				result = append(result, r...)
			}
			continue
		}

		if p, err := BuildProperty(member, indent); err != nil {
			return nil, apiutil.WrapErr(err, " type "+tp.Name())
		} else {
			result = append(result, p)
		}
	}
	return result, nil
}

func BuildMembers(tp spec.Type, isParam bool, indent int) ([]*template.ComponentMemberTemplateData, error) {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		if pointType, ok := tp.(spec.PointerType); ok {
			return BuildMembers(pointType.Type, isParam, indent)
		}

		if nestedType, ok := tp.(spec.NestedStruct); ok {
			return BuildMembers(spec.DefineStruct(nestedType), isParam, indent)
		}

		return nil, fmt.Errorf("type %s not supported", tp.Name())
	}

	result := []*template.ComponentMemberTemplateData{}
	members := definedType.GetBodyMembers()
	if isParam {
		members = definedType.GetNonBodyMembers()
	}
	for _, member := range members {
		if member.IsInline {
			if ms, err := BuildMembers(member.Type, isParam, indent); err != nil {
				return nil, err
			} else {
				result = append(result, ms...)
			}
			continue
		}

		if p, err := BuildProperty(member, indent); err != nil {
			return nil, apiutil.WrapErr(err, " type "+tp.Name())
		} else {
			result = append(result, p)
		}
	}
	return result, nil
}

func BuildSubTypes(tp spec.Type, indent int) ([]*template.ComponentTypeTemplateData, error) {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		return nil, errors.New("no members of type " + tp.Name())
	}

	result := []*template.ComponentTypeTemplateData{}
	members := definedType.GetNonBodyMembers()
	if len(members) == 0 {
		return result, nil
	}

	if len(definedType.GetTagMembers(formTagKey)) > 0 {
		ctt := &template.ComponentTypeTemplateData{
			TypeName: fmt.Sprintf("%sParams", util.Title(tp.Name())),
		}

		if ms, err := BuildTagMembers(tp, formTagKey, indent); err != nil {
			return nil, err
		} else {
			ctt.Members = ms
		}
		result = append(result, ctt)
	}

	if len(definedType.GetTagMembers(headerTagKey)) > 0 {
		ctt := &template.ComponentTypeTemplateData{
			TypeName: fmt.Sprintf("%sHeaders", util.Title(tp.Name())),
		}
		if ms, err := BuildTagMembers(tp, headerTagKey, indent); err != nil {
			return nil, err
		} else {
			ctt.Members = ms
		}

		result = append(result, ctt)
	}

	return result, nil
}

func BuildProperty(member spec.Member, indent int) (*template.ComponentMemberTemplateData, error) {
	ty, err := GenTsType(member, indent)
	if err != nil {
		return nil, err
	}

	result := &template.ComponentMemberTemplateData{}

	if IsOptionalOrOmitEmpty(member) {
		result.OptionalTag = "?"
	}

	name, err := member.GetPropertyName()
	if err != nil {
		return nil, err
	}
	if strings.Contains(name, "-") {
		name = fmt.Sprintf("\"%s\"", name)
	}
	result.PropertyName = name

	comment := member.GetComment()
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		result.Comment = strings.TrimSpace(comment)
	}
	// result.Docs = member.Docs
	result.PropertyType = ty

	return result, nil
}

// BuildTypes generates the typescript code for the types.
func BuildTypes(types []spec.Type) ([]template.ComponentTypeTemplateData, error) {
	result := []template.ComponentTypeTemplateData{}
	for _, tp := range types {
		members, err := BuildMembers(tp, false, 0)
		if err != nil {
			return nil, err
		}

		subTypes, err := BuildSubTypes(tp, 0)
		if err != nil {
			return nil, err
		}
		ctt := template.ComponentTypeTemplateData{
			TypeName: util.Title(tp.Name()),
			Members:  members,
			SubTypes: subTypes,
		}
		result = append(result, ctt)
	}

	return result, nil
}

func GenComponents(dir string, api *spec.ApiSpec) error {
	types := api.Types
	if len(types) == 0 {
		return nil
	}

	val, err := BuildTypes(types)
	if err != nil {
		return err
	}

	componentName := apiutil.ComponentName(api)
	data := template.ComponentTemplateData{
		Version: version.BuildVersion,
		Types:   val,
	}
	return template.GenTsFile(dir, componentName, template.Components, data)
}
