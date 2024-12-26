package javagen

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	httpResponseData = "import com.xhb.core.response.HttpResponseData;"
	httpData         = "import com.xhb.core.packet.HttpData;"
)

var (
	//go:embed component.tpl
	componentTemplate string
	//go:embed getset.tpl
	getSetTemplate string
	//go:embed bool.tpl
	boolTemplate string
)

type componentsContext struct {
	api           *spec.ApiSpec
	requestTypes  []spec.Type
	responseTypes []spec.Type
	imports       []string
	members       []spec.Member
}

func genComponents(dir, packetName string, api *spec.ApiSpec) error {
	types := api.Types
	if len(types) == 0 {
		return nil
	}

	var requestTypes []spec.Type
	var responseTypes []spec.Type
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if route.RequestType != nil {
				requestTypes = append(requestTypes, route.RequestType)
			}
			if route.ResponseType != nil {
				responseTypes = append(responseTypes, route.ResponseType)
			}
		}
	}

	context := componentsContext{api: api, requestTypes: requestTypes, responseTypes: responseTypes}
	for _, ty := range types {
		if err := context.createComponent(dir, packetName, ty); err != nil {
			return err
		}
	}

	return nil
}

func (c *componentsContext) createComponent(dir, packetName string, ty spec.Type) error {
	defineStruct, done, err := c.checkStruct(ty)
	if done {
		return err
	}

	modelFile := util.Title(ty.Name()) + ".java"
	filename := path.Join(dir, modelDir, modelFile)
	if err := pathx.RemoveOrQuit(filename); err != nil {
		return err
	}

	propertiesString, err := c.buildProperties(defineStruct)
	if err != nil {
		return err
	}

	getSetString, err := c.buildGetterSetter(defineStruct)
	if err != nil {
		return err
	}

	superClassName := "HttpData"
	for _, item := range c.responseTypes {
		if item.Name() == defineStruct.Name() {
			superClassName = "HttpResponseData"
			if !stringx.Contains(c.imports, httpResponseData) {
				c.imports = append(c.imports, httpResponseData)
			}
			break
		}
	}
	if superClassName == "HttpData" && !stringx.Contains(c.imports, httpData) {
		c.imports = append(c.imports, httpData)
	}

	params, constructorSetter, err := c.buildConstructor()
	if err != nil {
		return err
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, modelDir, modelFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	buffer := new(bytes.Buffer)
	t := template.Must(template.New("componentType").Parse(componentTemplate))
	err = t.Execute(buffer, map[string]any{
		"properties":        propertiesString,
		"params":            params,
		"constructorSetter": constructorSetter,
		"getSet":            getSetString,
		"packet":            packetName,
		"imports":           strings.Join(c.imports, "\n"),
		"className":         util.Title(defineStruct.Name()),
		"superClassName":    superClassName,
		"HasProperty":       len(strings.TrimSpace(propertiesString)) > 0,
		"version":           version.BuildVersion,
	})
	if err != nil {
		return err
	}

	_, err = fp.WriteString(formatSource(buffer.String()))
	return err
}

func (c *componentsContext) checkStruct(ty spec.Type) (spec.DefineStruct, bool, error) {
	defineStruct, ok := ty.(spec.DefineStruct)
	if !ok {
		return spec.DefineStruct{}, true, errors.New("unsupported type %s" + ty.Name())
	}

	for _, item := range c.requestTypes {
		if item.Name() == defineStruct.Name() {
			if len(defineStruct.GetFormMembers())+len(defineStruct.GetBodyMembers()) == 0 {
				return spec.DefineStruct{}, true, nil
			}
		}
	}
	return defineStruct, false, nil
}

func (c *componentsContext) buildProperties(defineStruct spec.DefineStruct) (string, error) {
	var builder strings.Builder
	if err := c.writeType(&builder, defineStruct); err != nil {
		return "", apiutil.WrapErr(err, "Type "+defineStruct.Name()+" generate error")
	}

	return builder.String(), nil
}

func (c *componentsContext) buildGetterSetter(defineStruct spec.DefineStruct) (string, error) {
	var builder strings.Builder
	if err := c.genGetSet(&builder, 1); err != nil {
		return "", apiutil.WrapErr(err, "Type "+defineStruct.Name()+" get or set generate error")
	}

	return builder.String(), nil
}

func (c *componentsContext) writeType(writer io.Writer, defineStruct spec.DefineStruct) error {
	c.members = make([]spec.Member, 0)
	err := c.writeMembers(writer, defineStruct, 1)
	if err != nil {
		return err
	}

	return nil
}

func (c *componentsContext) writeMembers(writer io.Writer, tp spec.Type, indent int) error {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return c.writeMembers(writer, pointType.Type, indent)
		}
		return fmt.Errorf("type %s not supported", tp.Name())
	}

	for _, member := range definedType.Members {
		if member.IsInline {
			err := c.writeMembers(writer, member.Type, indent)
			if err != nil {
				return err
			}
			continue
		}

		if member.IsBodyMember() || member.IsFormMember() {
			if err := writeProperty(writer, member, indent); err != nil {
				return err
			}

			c.members = append(c.members, member)
		}
	}

	return nil
}

func (c *componentsContext) buildConstructor() (string, string, error) {
	var params strings.Builder
	var constructorSetter strings.Builder
	for index, member := range c.members {
		tp, err := specTypeToJava(member.Type)
		if err != nil {
			return "", "", err
		}

		params.WriteString(fmt.Sprintf("%s %s", tp, util.Untitle(member.Name)))
		pn, err := member.GetPropertyName()
		if err != nil {
			return "", "", err
		}

		if index != len(c.members)-1 {
			params.WriteString(", ")
		}

		writeIndent(&constructorSetter, 2)
		constructorSetter.WriteString(fmt.Sprintf("this.%s = %s;", pn, util.Untitle(member.Name)))
		if index != len(c.members)-1 {
			constructorSetter.WriteString(pathx.NL)
		}
	}
	return params.String(), constructorSetter.String(), nil
}

func (c *componentsContext) genGetSet(writer io.Writer, indent int) error {
	members := c.members
	for _, member := range members {
		javaType, err := specTypeToJava(member.Type)
		if err != nil {
			return nil
		}

		property := util.Title(member.Name)
		templateStr := getSetTemplate
		if javaType == "boolean" {
			templateStr = boolTemplate
			property = strings.TrimPrefix(property, "Is")
			property = strings.TrimPrefix(property, "is")
		}
		t := template.Must(template.New(templateStr).Parse(getSetTemplate))
		var tmplBytes bytes.Buffer

		tyString := javaType
		decorator := ""
		javaPrimitiveType := []string{"int", "long", "boolean", "float", "double", "short"}
		if !stringx.Contains(javaPrimitiveType, javaType) {
			if member.IsOptional() || member.IsOmitEmpty() {
				decorator = "@Nullable "
			} else {
				decorator = "@NotNull "
			}
			tyString = decorator + tyString
		}

		tagName, err := member.GetPropertyName()
		if err != nil {
			return err
		}

		err = t.Execute(&tmplBytes, map[string]string{
			"property":      property,
			"propertyValue": util.Untitle(member.Name),
			"tagValue":      tagName,
			"type":          tyString,
			"decorator":     decorator,
			"returnType":    javaType,
			"indent":        indentString(indent),
		})
		if err != nil {
			return err
		}

		r := tmplBytes.String()
		r = strings.Replace(r, " boolean get", " boolean is", 1)
		writer.Write([]byte(r))
	}
	return nil
}

func formatSource(source string) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(source))
	preIsBreakLine := false
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" && preIsBreakLine {
			continue
		}
		preIsBreakLine = text == ""
		builder.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return builder.String()
}
