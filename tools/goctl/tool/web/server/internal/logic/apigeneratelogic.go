package logic

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	sortedmap "github.com/zeromicro/go-zero/tools/goctl/pkg/collection"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/logic/sortmap"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	typex "github.com/zeromicro/go-zero/tools/goctl/util/types"
	"github.com/zeromicro/go-zero/tools/goctl/util/writer"
)

const (
	indent          = "  "
	applicationJSON = "application/json"
)

var (
	//go:embed tpl/api.api
	apiTemplate string
	//go:embed tpl/field.tpl
	filedTemplate string

	errMissingServiceName = errors.New("missing service name")
	errMissingRouteGroups = errors.New("missing route groups")
)

type ApiGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiGenerateLogic {
	return &ApiGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiGenerateLogic) ApiGenerate(req *types.APIGenerateRequest) (resp *types.APIGenerateResponse, err error) {
	if err := l.validateAPIGenerateRequest(req); err != nil {
		return nil, err
	}
	mergedReq := l.mergeGroup(req)
	var data []KV
	for _, group := range mergedReq.List {
		var groupData = KV{}
		var hasServer bool
		var server = KV{}
		if group.Jwt {
			hasServer = true
			server["jwt"] = group.Jwt
		}
		if len(group.Prefix) > 0 {
			hasServer = true
			server["prefix"] = group.Prefix
		}
		if len(group.Group) > 0 {
			hasServer = true
			server["group"] = group.Group
		}
		if group.Timeout > 0 {
			hasServer = true
			server["timeout"] = fmt.Sprintf("%dms", group.Timeout)
		}
		if len(group.Middleware) > 0 {
			hasServer = true
			server["middleware"] = group.Middleware
		}
		if group.MaxBytes > 0 {
			hasServer = true
			server["maxBytes"] = group.MaxBytes
		}

		if hasServer {
			groupData["server"] = server
		}

		var routesData []KV
		for _, route := range group.Routes {
			var request, response string
			if len(route.RequestBody) > 0 {
				request = l.generateTypeName(route, true)
			}
			if !util.IsEmptyStringOrWhiteSpace(route.ResponseBody) {
				response = l.generateTypeName(route, false)
			}
			routesData = append(routesData, KV{
				"handlerName": l.generateHandlerName(route),
				"method":      strings.ToLower(route.Method),
				"path":        route.Path,
				"request":     request,
				"response":    response,
			})
		}
		var service = KV{
			"name":   req.Name,
			"routes": routesData,
		}
		groupData["service"] = service
		data = append(data, groupData)
	}

	t, err := template.New("api").Funcs(map[string]any{
		"lessThan": func(idx int, length int) bool {
			return idx < length-1
		},
	}).Parse(apiTemplate)
	if err != nil {
		return nil, err
	}

	tps, err := l.generateTypes(mergedReq.List)
	if err != nil {
		return nil, err
	}

	var typeString string
	if len(tps) > 0 {
		typeString = strings.Join(tps, "\n\n")
	}

	w := bytes.NewBuffer(nil)
	err = t.Execute(w, map[string]any{
		"types":  typeString,
		"groups": data,
	})
	if err != nil {
		return nil, err
	}

	formatWriter := bytes.NewBuffer(nil)
	err = format.Source(w.Bytes(), formatWriter)
	if err != nil {
		return nil, err
	}

	return &types.APIGenerateResponse{
		API: formatWriter.String(),
	}, nil
}

func (l *ApiGenerateLogic) generateTypes(groups []*types.APIRouteGroup) ([]string, error) {
	var resp []string
	for _, group := range groups {
		var groupTypes []string
		for _, route := range group.Routes {
			tp, err := l.generateType(route)
			if err != nil {
				return nil, err
			}
			if len(tp) > 0 {
				groupTypes = append(groupTypes, tp...)
			}
		}
		if len(groupTypes) > 0 {
			resp = append(resp, fmt.Sprintf(`type(
%s
)`, strings.Join(groupTypes, "\n\n")))
		}
	}
	return resp, nil
}

func (l *ApiGenerateLogic) generateType(route *types.APIRoute) ([]string, error) {
	var requestTypes []string
	if len(route.RequestBody) > 0 {
		postJson := strings.EqualFold(route.Method, http.MethodPost) && route.ContentType == applicationJSON
		requestType, err := l.generateRequestType(l.generateTypeName(route, true), postJson, route.RequestBody)
		if err != nil {
			return nil, err
		}
		if len(requestType) > 0 {
			requestTypes = append(requestTypes, requestType)
		}
	}

	responseType, err := l.generateResponseType(l.generateTypeName(route, false), route.ResponseBody)
	if err != nil {
		return nil, err
	}
	if len(responseType) > 0 {
		requestTypes = append(requestTypes, responseType)
	}
	return requestTypes, nil
}

func (l *ApiGenerateLogic) generateRequestType(typeName string, json bool, form []*types.FormItem) (string, error) {
	t, err := template.New("field").Funcs(map[string]any{
		"camel": func(s string) string {
			x := stringx.From(s)
			return strings.Title(x.ToCamel())
		},
	}).Parse(filedTemplate)
	if err != nil {
		return "", err
	}

	w := writer.New("")
	fieldWriter := bytes.NewBuffer(nil)
	for _, item := range form {
		fieldWriter.Reset()
		var rangeValue, enumValue string
		if item.CheckEnum == "range" &&
			item.LowerBound != item.UpperBound {
			rangeValue = fmt.Sprintf("range=%s", formatRange(item.LowerBound, item.UpperBound))
		}
		if item.CheckEnum == "enum" {
			enumValue = item.EnumValue
		}
		err = t.Execute(fieldWriter, map[string]any{
			"name":         item.Name,
			"type":         item.Type,
			"json":         json,
			"optional":     item.Optional,
			"defaultValue": item.DefaultValue,
			"checkEnum":    item.CheckEnum == "enum",
			"enumValue":    enumValue,
			"rangeValue":   rangeValue,
		})
		if err != nil {
			return "", err
		}
		w.WriteStringln(fieldWriter.String())
	}

	return fmt.Sprintf(`%s {
%s
}`, typeName, w.String()), nil
}

func formatRange(lowerBound, upperBound int64) string {
	if lowerBound == -1 {
		return fmt.Sprintf("[:%d]", upperBound)
	}
	if upperBound == -1 {
		return fmt.Sprintf("[%d:]", lowerBound)
	}
	return fmt.Sprintf("[%d:%d]", lowerBound, upperBound)
}

func (l *ApiGenerateLogic) generateResponseType(typeName, s string) (string, error) {
	if util.IsEmptyStringOrWhiteSpace(s) {
		return "", nil
	}
	var v any
	err := json.Unmarshal([]byte(s), &v)
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		return "", fmt.Errorf("invalid json, offset: %d, msg: %s", jsonSyntaxErr.Offset, jsonSyntaxErr.Error())
	}
	tps, _, err := json2APIType(json2APITypeReq{
		root:        true,
		indentCount: 1,
		typeName:    typeName,
		v:           v,
	})
	return tps, err
}

func json2APIType(req json2APITypeReq) (tp string, externalTypes []string, err error) {
	typeName := strcase.ToCamel(req.parentTypeName) + strcase.ToCamel(req.typeName)
	kv, ok := req.v.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("input must be object, got %T", req.v)
	}
	sm := sortmap.From(kv)

	w := writer.New(getIdent(req.indentCount))
	w.WriteWithIndentStringf("%s {\n", typeName)
	if len(kv) == 0 {
		w.UndoNewLine()
		w.Writef("}")
		return w.String(), nil, nil
	}

	var externalTypeList []string
	memberWriter := writer.New(getIdent(req.indentCount + 1))
	err = sm.Range(func(_ int, key string, value any) error {
		result, err := convertGoctlAPIMemberType(req.indentCount+1, typeName, key, value)
		if err != nil {
			return err
		}
		externalTypeList = append(externalTypeList, result.ExternalTypeExpr...)
		if result.IsStruct {
			externalTypeList = append(externalTypeList, result.TypeExpr)
		}
		if result.IsArray {
			memberWriter.WriteWithIndentStringf("%s []%s `json:\"%s\"`\n", strcase.ToCamel(key), result.TypeName, key)
		} else {
			memberWriter.WriteWithIndentStringf("%s %s `json:\"%s\"`\n", strcase.ToCamel(key), result.TypeName, key)
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	w.Writef(memberWriter.String())
	w.WriteWithIndentStringf("}")

	if req.root {
		w.NewLine()
		w.WriteStringln(strings.Join(externalTypeList, "\n\n"))
	}

	return w.String(), externalTypeList, nil
}

func convertGoctlAPIMemberType(indentCount int, parent, key string, value any) (*goctlAPIMemberResult, error) {
	resp := new(goctlAPIMemberResult)
	switch {
	case typex.IsInteger(value):
		resp.TypeExpr = "int64"
		resp.TypeName = "int64"
		return resp, nil
	case typex.IsFloat(value):
		resp.TypeExpr = "double"
		resp.TypeName = "double"
		return resp, nil
	case typex.IsBool(value):
		resp.TypeExpr = "bool"
		resp.TypeName = "bool"
		return resp, nil
	case typex.IsTime(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	case typex.IsString(value):
		resp.TypeExpr = "string"
		resp.TypeName = "string"
		return resp, nil
	default:
		_, ok := value.(map[string]any)
		if ok {
			tp, externalTypes, err := json2APIType(json2APITypeReq{
				indentCount:    indentCount,
				parentTypeName: parent,
				typeName:       key,
				v:              value,
			})
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = "*" + strcase.ToCamel(parent) + strcase.ToCamel(key)
			resp.IsStruct = true
			resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, externalTypes...)
			return resp, nil
		}
		list, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("unsupport type, got %T", value)
		}
		if len(list) == 0 {
			resp.TypeExpr = "interface{}"
			resp.TypeName = "interface{}"
			resp.IsArray = true
			return resp, nil
		}
		first := list[0]
		_, ok = first.(map[string]any)
		if ok {
			var memberSet = make(map[string]any)
			for _, v := range list {
				m, ok := v.(map[string]any)
				if !ok {
					continue
				}
				for k, v := range m {
					memberSet[k] = v
				}
			}
			tp, externalTypes, err := json2APIType(json2APITypeReq{
				indentCount:    indentCount,
				parentTypeName: parent,
				typeName:       key,
				v:              memberSet,
			})
			if err != nil {
				return nil, err
			}
			resp.TypeExpr = tp
			resp.TypeName = "*" + strcase.ToCamel(parent) + strcase.ToCamel(key)
			resp.IsStruct = true
			resp.IsArray = true
			resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, externalTypes...)
			return resp, nil
		}
		result, err := convertGoctlAPIMemberType(indentCount, parent, key, first)
		if err != nil {
			return nil, err
		}
		resp.TypeExpr = result.TypeExpr
		resp.TypeName = result.TypeName
		resp.IsStruct = result.IsStruct
		resp.IsArray = true
		resp.ExternalTypeExpr = append(resp.ExternalTypeExpr, result.ExternalTypeExpr...)
		return resp, nil
	}
}

func getIdent(c int) string {
	var list []string
	for i := 0; i < c; i++ {
		list = append(list, indent)
	}
	return strings.Join(list, "")
}

func (l *ApiGenerateLogic) generateTypeName(route *types.APIRoute, req bool) string {
	handlerName := l.generateHandlerName(route)
	if req {
		return handlerName + "Request"
	}
	return handlerName + "Response"
}

func (l *ApiGenerateLogic) generateHandlerName(route *types.APIRoute) string {
	if len(route.Handler) > 0 {
		return strings.Title(route.Handler)
	}
	if route.Path == "/" {
		return "Default"
	}

	r := strings.NewReplacer("/", "_", ":", "by_")
	formatedPath := r.Replace(route.Path)
	s := stringx.From(route.Method + "_" + formatedPath)
	return strings.Title(s.ToCamel())
}

func (l *ApiGenerateLogic) mergeGroup(req *types.APIGenerateRequest) *types.APIGenerateRequest {
	routeGroup := sortedmap.New()
	for _, group := range req.List {
		middlewares := strings.Split(group.Middleware, ",")
		sort.Strings(middlewares)
		middleware := strings.Join(middlewares, ", ")
		routeGroupStruct := RouteGroup{
			Jwt:        group.Jwt,
			Prefix:     group.Prefix,
			Group:      group.Group,
			Timeout:    group.Timeout,
			Middleware: middleware,
			MaxBytes:   group.MaxBytes,
		}
		val, ok := routeGroup.Get(routeGroupStruct)
		if ok {
			existGroup := val.(*types.APIRouteGroup)
			existGroup.Routes = l.appendAndMergeRoute(existGroup.Routes, group.Routes)
			routeGroup.SetKV(routeGroupStruct, existGroup)
		} else {
			routeGroup.SetKV(routeGroupStruct, &types.APIRouteGroup{
				Jwt:        group.Jwt,
				Prefix:     group.Prefix,
				Group:      group.Group,
				Timeout:    group.Timeout,
				Middleware: middleware,
				MaxBytes:   group.MaxBytes,
				Routes:     group.Routes,
			})
		}
	}

	var resp = types.APIGenerateRequest{
		Name: req.Name,
	}
	routeGroup.Range(func(key, value any) {
		resp.List = append(resp.List, value.(*types.APIRouteGroup))
	})

	return &resp
}

func (l *ApiGenerateLogic) appendAndMergeRoute(header, tailer []*types.APIRoute) []*types.APIRoute {
	routeMap := sortedmap.New()
	for _, route := range header {
		key := l.convertRoute(route)
		routeMap.SetKV(key, route)
	}

	for _, route := range tailer {
		key := l.convertRoute(route)
		ok := routeMap.HasKey(key)
		if !ok {
			routeMap.SetKV(route, route)
		}
	}

	var list []*types.APIRoute
	routeMap.Range(func(key, value any) {
		list = append(list, value.(*types.APIRoute))
	})

	return list
}

func (l *ApiGenerateLogic) convertRoute(route *types.APIRoute) APIRoute {
	var requestBody []types.FormItem
	for _, v := range route.RequestBody {
		requestBody = append(requestBody, *v)
	}
	return APIRoute{ // ignore request & response
		Handler:     route.Handler,
		Method:      route.Method,
		Path:        route.Path,
		ContentType: route.ContentType,
	}
}

func (l *ApiGenerateLogic) validateAPIGenerateRequest(req *types.APIGenerateRequest) error {
	if util.IsEmptyStringOrWhiteSpace(req.Name) {
		return errMissingServiceName
	}
	if len(req.List) == 0 {
		return errMissingRouteGroups
	}

	var err []string
	for idx, group := range req.List {
		if len(group.Routes) == 0 {
			if len(group.Group) > 0 {
				err = append(err, fmt.Sprintf("group %q: missing routes", group.Group))
			} else {
				err = append(err, fmt.Sprintf("group%d: missing routes", idx+1))
			}
		}
	}

	var (
		handlerDuplicateCheck = make(map[string]lang.PlaceholderType)
		routeDuplicateCheck   = make(map[string]lang.PlaceholderType)
	)

	for idx, group := range req.List {
		for _, route := range group.Routes {
			var handlerUniqueValue, routeUniqueValue string
			if len(group.Group) > 0 {
				handlerUniqueValue = fmt.Sprintf("%s/%s", group.Group, route.Handler)
				routeUniqueValue = fmt.Sprintf("%s/%s:%s/%s", group.Group, route.Method, group.Prefix, route.Path)
			} else {
				handlerUniqueValue = fmt.Sprintf("group[%d]/%s", idx, route.Handler)
				routeUniqueValue = fmt.Sprintf("group[%d]/%s:%s/%s", idx, route.Method, group.Prefix, route.Path)
			}

			if _, ok := handlerDuplicateCheck[handlerUniqueValue]; ok {
				err = append(err, fmt.Sprintf("duplicate handler: %q", handlerUniqueValue))
			}
			if _, ok := routeDuplicateCheck[routeUniqueValue]; ok {
				err = append(err, fmt.Sprintf("duplicate route: %q", routeUniqueValue))
			}

			handlerDuplicateCheck[handlerUniqueValue] = lang.Placeholder
			routeDuplicateCheck[routeUniqueValue] = lang.Placeholder
		}
	}

	if len(err) == 0 {
		return nil
	}

	return errors.New(strings.Join(err, "\n"))
}
