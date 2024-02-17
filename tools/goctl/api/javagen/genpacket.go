package javagen

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

//go:embed packet.tpl
var packetTemplate string

func genPacket(dir, packetName string, api *spec.ApiSpec) error {
	for _, route := range api.Service.Routes() {
		if err := createWith(dir, api, route, packetName); err != nil {
			return err
		}
	}

	return nil
}

func createWith(dir string, api *spec.ApiSpec, route spec.Route, packetName string) error {
	packet := route.Handler
	packet = strings.Replace(packet, "Handler", "Packet", 1)
	packet = strings.Title(packet)
	if !strings.HasSuffix(packet, "Packet") {
		packet += "Packet"
	}

	javaFile := packet + ".java"
	fp, created, err := apiutil.MaybeCreateFile(dir, "", javaFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	hasRequestBody := false
	if route.RequestType != nil {
		if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
			hasRequestBody = len(defineStruct.GetBodyMembers()) > 0 || len(defineStruct.GetFormMembers()) > 0
		}
	}

	params := strings.TrimSpace(paramsForRoute(route))
	if len(params) > 0 && hasRequestBody {
		params += ", "
	}
	paramsDeclaration := declarationForRoute(route)
	paramsSetter := paramsSet(route)
	imports := getImports(api, packetName)

	if len(route.ResponseTypeName()) == 0 {
		imports += fmt.Sprintf("\v%s", "import com.xhb.core.response.EmptyResponse;")
	}

	t := template.Must(template.New("packetTemplate").Parse(packetTemplate))
	var tmplBytes bytes.Buffer
	err = t.Execute(&tmplBytes, map[string]any{
		"packetName":        packet,
		"method":            strings.ToUpper(route.Method),
		"uri":               processUri(route),
		"responseType":      stringx.TakeOne(util.Title(route.ResponseTypeName()), "EmptyResponse"),
		"params":            params,
		"paramsDeclaration": strings.TrimSpace(paramsDeclaration),
		"paramsSetter":      paramsSetter,
		"packet":            packetName,
		"requestType":       util.Title(route.RequestTypeName()),
		"HasRequestBody":    hasRequestBody,
		"imports":           imports,
		"doc":               doc(route),
	})
	if err != nil {
		return err
	}

	_, err = fp.WriteString(formatSource(tmplBytes.String()))
	return err
}

func doc(route spec.Route) string {
	comment := route.JoinedDoc()
	if len(comment) > 0 {
		formatter := `
/*
    %s	
*/`
		return fmt.Sprintf(formatter, comment)
	}
	return ""
}

func getImports(api *spec.ApiSpec, packetName string) string {
	var builder strings.Builder
	allTypes := api.Types
	if len(allTypes) > 0 {
		fmt.Fprintf(&builder, "import com.xhb.logic.http.packet.%s.model.*;\n", packetName)
	}

	return builder.String()
}

func paramsSet(route spec.Route) string {
	path := route.Path
	cops := strings.Split(path, "/")
	var builder strings.Builder
	for _, cop := range cops {
		if len(cop) == 0 {
			continue
		}
		if strings.HasPrefix(cop, ":") {
			param := cop[1:]
			builder.WriteString("\n")
			builder.WriteString(fmt.Sprintf("\t\tthis.%s = %s;", param, param))
		}
	}
	result := builder.String()
	return result
}

func paramsForRoute(route spec.Route) string {
	path := route.Path
	cops := strings.Split(path, "/")
	var builder strings.Builder
	for _, cop := range cops {
		if len(cop) == 0 {
			continue
		}
		if strings.HasPrefix(cop, ":") {
			builder.WriteString(fmt.Sprintf("String %s, ", cop[1:]))
		}
	}
	return strings.TrimSuffix(builder.String(), ", ")
}

func declarationForRoute(route spec.Route) string {
	path := route.Path
	cops := strings.Split(path, "/")
	var builder strings.Builder
	writeIndent(&builder, 1)
	for _, cop := range cops {
		if len(cop) == 0 {
			continue
		}
		if strings.HasPrefix(cop, ":") {
			writeIndent(&builder, 1)
			builder.WriteString(fmt.Sprintf("private String %s;\n", cop[1:]))
		}
	}
	result := strings.TrimSpace(builder.String())
	if len(result) > 0 {
		result = "\n" + result
	}
	return result
}

func processUri(route spec.Route) string {
	path := route.Path

	var builder strings.Builder
	cops := strings.Split(path, "/")
	for index, cop := range cops {
		if len(cop) == 0 {
			continue
		}
		if strings.HasPrefix(cop, ":") {
			builder.WriteString("/\" + " + cop[1:] + " + \"")
		} else {
			builder.WriteString("/" + cop)
			if index == len(cops)-1 {
				builder.WriteString("\"")
			}
		}
	}
	result := strings.TrimSuffix(builder.String(), " + \"")
	if strings.HasPrefix(result, "/") {
		result = strings.TrimPrefix(result, "/")
		result = "\"" + result
	}
	return result + formString(route)
}

func formString(route spec.Route) string {
	var keyValues []string
	if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
		forms := defineStruct.GetFormMembers()
		for _, item := range forms {
			name, err := item.GetPropertyName()
			if err != nil {
				panic(err)
			}

			strcat := "?"
			if len(keyValues) > 0 {
				strcat = "&"
			}
			if item.Type.Name() == "bool" {
				name = strings.TrimPrefix(name, "Is")
				name = strings.TrimPrefix(name, "is")
				keyValues = append(keyValues, fmt.Sprintf(`"%s%s=" + request.is%s()`, strcat, name, strings.Title(name)))
			} else {
				keyValues = append(keyValues, fmt.Sprintf(`"%s%s=" + request.get%s()`, strcat, name, strings.Title(name)))
			}
		}
		if len(keyValues) > 0 {
			return " + " + strings.Join(keyValues, " + ")
		}
	}
	return ""
}
