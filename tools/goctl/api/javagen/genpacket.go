package javagen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const packetTemplate = `package com.xhb.logic.http.packet.{{.packet}};

import com.google.gson.Gson;
import com.xhb.commons.JSON;
import com.xhb.commons.JsonMarshal;
import com.xhb.core.network.HttpRequestClient;
import com.xhb.core.packet.HttpRequestPacket;
import com.xhb.core.response.HttpResponseData;
import com.xhb.logic.http.DeProguardable;
{{if not .HasRequestBody}}
import com.xhb.logic.http.request.EmptyRequest;
{{end}}
{{.import}}

import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;
import org.json.JSONObject;

public class {{.packetName}} extends HttpRequestPacket<{{.packetName}}.{{.packetName}}Response> {

	{{.paramsDeclaration}}

	public {{.packetName}}({{.params}}{{if .HasRequestBody}}, {{.requestType}} request{{end}}) {
		{{if .HasRequestBody}}super(request);{{else}}super(EmptyRequest.instance);{{end}}
		{{if .HasRequestBody}}this.request = request;{{end}}{{.paramsSet}}
    }

	@Override
    public HttpRequestClient.Method requestMethod() {
        return HttpRequestClient.Method.{{.method}};
    }

	@Override
    public String requestUri() {
        return {{.uri}};
    }

	@Override
    public {{.packetName}}Response newInstanceFrom(JSON json) {
        return new {{.packetName}}Response(json);
    }

	public static class {{.packetName}}Response extends HttpResponseData {

		private {{.responseType}} responseData;

        {{.packetName}}Response(@NotNull JSON json) {
            super(json);
            JSONObject jsonObject = json.asObject();
			if (JsonParser.hasKey(jsonObject, "data")) {
				Gson gson = new Gson();
				JSONObject dataJson = JsonParser.getJSONObject(jsonObject, "data");
				responseData = gson.fromJson(dataJson.toString(), {{.responseType}}.class);
			}
        }

		public {{.responseType}} get{{.responseType}} () {
            return responseData;
        }
    }

	{{.types}}
}
`

func genPacket(dir, packetName string, api *spec.ApiSpec) error {
	for _, route := range api.Service.Routes() {
		if err := createWith(dir, api, route, packetName); err != nil {
			return err
		}
	}

	return nil
}

func createWith(dir string, api *spec.ApiSpec, route spec.Route, packetName string) error {
	packet, ok := apiutil.GetAnnotationValue(route.Annotations, "server", "handler")
	packet = strings.Replace(packet, "Handler", "Packet", 1)
	if !ok {
		return fmt.Errorf("missing packet annotation for %q", route.Path)
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

	var builder strings.Builder
	var first bool
	tps := apiutil.GetLocalTypes(api, route)

	for _, tp := range tps {
		if first {
			first = false
		} else {
			fmt.Fprintln(&builder)
		}

		if err := genType(&builder, tp, api.Types); err != nil {
			return err
		}
	}
	types := builder.String()
	writeIndent(&builder, 1)

	params := paramsForRoute(route)
	paramsDeclaration := declarationForRoute(route)
	paramsSet := paramsSet(route)

	t := template.Must(template.New("packetTemplate").Parse(packetTemplate))
	var tmplBytes bytes.Buffer
	err = t.Execute(&tmplBytes, map[string]interface{}{
		"packetName":        packet,
		"method":            strings.ToUpper(route.Method),
		"uri":               processUri(route),
		"types":             strings.TrimSpace(types),
		"responseType":      stringx.TakeOne(util.Title(route.ResponseType.Name), "Object"),
		"params":            params,
		"paramsDeclaration": strings.TrimSpace(paramsDeclaration),
		"paramsSet":         paramsSet,
		"packet":            packetName,
		"requestType":       util.Title(route.RequestType.Name),
		"HasRequestBody":    len(route.RequestType.GetBodyMembers()) > 0,
		"import":            getImports(api, route, packetName),
	})
	if err != nil {
		return err
	}
	formatFile(&tmplBytes, fp)
	return nil
}

func getImports(api *spec.ApiSpec, route spec.Route, packetName string) string {
	var builder strings.Builder
	allTypes := apiutil.GetAllTypes(api, route)
	sharedTypes := apiutil.GetSharedTypes(api)
	for _, at := range allTypes {
		for _, item := range sharedTypes {
			if item.Name == at.Name {
				fmt.Fprintf(&builder, "import com.xhb.logic.http.packet.%s.model.%s;\n", packetName, item.Name)
				break
			}
		}
	}
	return builder.String()
}

func formatFile(tmplBytes *bytes.Buffer, file *os.File) {
	scanner := bufio.NewScanner(tmplBytes)
	builder := bufio.NewWriter(file)
	defer builder.Flush()
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
	result := builder.String()
	if strings.HasSuffix(result, " + \"") {
		result = result[:len(result)-4]
	}
	if strings.HasPrefix(result, "/") {
		result = "\"" + result
	}
	return result
}

func genType(writer io.Writer, tp spec.Type, types []spec.Type) error {
	if len(tp.GetBodyMembers()) == 0 {
		return nil
	}

	writeIndent(writer, 1)
	fmt.Fprintf(writer, "static class %s implements DeProguardable {\n", util.Title(tp.Name))
	var members []spec.Member
	err := writeMembers(writer, types, tp.Members, &members, 2)
	if err != nil {
		return err
	}

	writeNewline(writer)
	writeIndent(writer, 1)
	genGetSet(writer, members, 2)
	writeIndent(writer, 1)
	fmt.Fprintln(writer, "}")

	return nil
}
