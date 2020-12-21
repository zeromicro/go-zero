package ast

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/g4gen"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const serverAnnotationName = "server"

type (
	ApiVisitor struct {
		parser.BaseApiParserVisitor
		serviceGroup *spec.Group
		apiSpec      spec.ApiSpec
	}
	kv struct {
		key   string
		value string
	}
)

func NewApiVisitor() *ApiVisitor {
	return &ApiVisitor{}
}

func (v *ApiVisitor) VisitApi(ctx *parser.ApiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitBody(ctx *parser.BodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitSyntaxLit(ctx *parser.SyntaxLitContext) interface{} {
	version := v.getTokenText(ctx.GetVersion(), true)

	return &spec.ApiSyntax{Version: version}
}

func (v *ApiVisitor) VisitImportSpec(ctx *parser.ImportSpecContext) interface{} {
	iImportLitContext := ctx.ImportLit()
	iImportLitGroupContext := ctx.ImportLitGroup()
	var list []string
	if iImportLitContext != nil {
		importLitContext, ok := iImportLitContext.(*parser.ImportLitContext)
		if ok {
			result := v.VisitImportLit(importLitContext)
			importValue, ok := result.(*spec.ApiImport)
			if ok {
				list = append(list, importValue.List...)
			}
		}
	}

	if iImportLitGroupContext != nil {
		importGroupContext, ok := iImportLitGroupContext.(*parser.ImportLitGroupContext)
		if ok {
			result := v.VisitImportLitGroup(importGroupContext)
			importValue, ok := result.(*spec.ApiImport)
			if ok {
				list = append(list, importValue.List...)
			}
		}
	}
	return &spec.ApiImport{List: list}
}

func (v *ApiVisitor) VisitImportLit(ctx *parser.ImportLitContext) interface{} {
	importPath := v.getTokenText(ctx.GetImportPath(), true)

	return &spec.ApiImport{
		List: []string{importPath},
	}
}

func (v *ApiVisitor) VisitImportLitGroup(ctx *parser.ImportLitGroupContext) interface{} {
	nodes := ctx.AllIMPORT_PATH()
	var list []string
	for _, node := range nodes {
		importPath := v.getNodeText(node, true)

		list = append(list, importPath)
	}
	return &spec.ApiImport{List: list}
}

func (v *ApiVisitor) VisitInfoBlock(ctx *parser.InfoBlockContext) interface{} {
	var info spec.Info
	info.Proterties = make(map[string]string)
	iKvLitContexts := ctx.AllKvLit()
	for _, each := range iKvLitContexts {
		kvLitContext, ok := each.(*parser.KvLitContext)
		if !ok {
			continue
		}

		r := v.VisitKvLit(kvLitContext)
		kv, ok := r.(*kv)
		if !ok {
			continue
		}
		info.Proterties[kv.key] = kv.value

	}
	return &info
}

func (v *ApiVisitor) VisitTypeBlock(ctx *parser.TypeBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeLit(ctx *parser.TypeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeGroup(ctx *parser.TypeGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeSpec(ctx *parser.TypeSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeAlias(ctx *parser.TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeStruct(ctx *parser.TypeStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitTypeField(ctx *parser.TypeFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitFiled(ctx *parser.FiledContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitInnerStruct(ctx *parser.InnerStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitDataType(ctx *parser.DataTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitMapType(ctx *parser.MapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitArrayType(ctx *parser.ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitPointer(ctx *parser.PointerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServiceBlock(ctx *parser.ServiceBlockContext) interface{} {
	if v.serviceGroup == nil {
		v.serviceGroup = new(spec.Group)
		v.apiSpec.Service.Groups = append(v.apiSpec.Service.Groups, *v.serviceGroup)
	}
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServerMeta(ctx *parser.ServerMetaContext) interface{} {
	if v.serviceGroup == nil {
		v.serviceGroup = new(spec.Group)
		v.apiSpec.Service.Groups = append(v.apiSpec.Service.Groups, *v.serviceGroup)
	}
	v.serviceGroup.Annotation.Name = serverAnnotationName
	v.serviceGroup.Annotation.Properties = make(map[string]string, 0)
	annos := ctx.AllAnnotation()
	for _, anno := range annos {
		anno.Accept(v)
	}

	return v.serviceGroup.Annotation
}

func (v *ApiVisitor) VisitAnnotation(ctx *parser.AnnotationContext) interface{} {
	key := v.getTokenText(ctx.GetKey(), true)

	if len(key) == 0 || ctx.GetValue() == nil {
		panic(errors.New("empty annotation key or value"))
	}

	v.serviceGroup.Annotation.Properties[key] = ctx.GetValue().GetText()
	return nil
}

func (v *ApiVisitor) VisitAnnotationKeyValue(ctx *parser.AnnotationKeyValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServiceBody(ctx *parser.ServiceBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServiceName(ctx *parser.ServiceNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServiceRoute(ctx *parser.ServiceRouteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitRouteDoc(ctx *parser.RouteDocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitDoc(ctx *parser.DocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitLineDoc(ctx *parser.LineDocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitRouteHandler(ctx *parser.RouteHandlerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitRoutePath(ctx *parser.RoutePathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitPath(ctx *parser.PathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitRequest(ctx *parser.RequestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitReply(ctx *parser.ReplyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitKvLit(ctx *parser.KvLitContext) interface{} {
	key := v.getTokenText(ctx.GetKey(), false)
	value := v.getTokenText(ctx.GetValue(), true)

	return &kv{
		key:   key,
		value: value,
	}
}

func (v *ApiVisitor) getTokenInt(token antlr.Token) (int64, error) {
	text := v.getTokenText(token, true)
	if len(text) == 0 {
		return 0, nil
	}

	vInt, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, err
	}

	return vInt, nil
}

func (v *ApiVisitor) getTokenText(token antlr.Token, trimQuote bool) string {
	if token == nil {
		return ""
	}

	text := token.GetText()
	if trimQuote {
		text = v.trimQuote(text)
	}
	return text
}

func (v *ApiVisitor) getNodeInt(node antlr.TerminalNode) (int64, error) {
	text := v.getNodeText(node, true)
	if len(text) == 0 {
		return 0, nil
	}

	vInt, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, err
	}

	return vInt, nil

}

func (v *ApiVisitor) getNodeText(node antlr.TerminalNode, trimQuote bool) string {
	if node == nil {
		return ""
	}

	text := node.GetText()
	if trimQuote {
		text = v.trimQuote(text)
	}
	return text
}

func (v *ApiVisitor) trimQuote(text string) string {
	text = strings.ReplaceAll(text, `"`, "")
	text = strings.ReplaceAll(text, `'`, "")
	text = strings.ReplaceAll(text, "`", "")
	return text
}
