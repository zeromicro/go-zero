package ast

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/g4gen"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type (
	ApiVisitor struct {
		parser.BaseApiParserVisitor

		serviceGroup *spec.Group
	}
	VisitResult struct {
		v   interface{}
		err error
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
	version, err := v.getText(ctx.GetVersion(), true)
	syntax := &spec.ApiSyntax{Version: version}
	return &VisitResult{
		err: err,
		v:   syntax,
	}
}

func (v *ApiVisitor) VisitImportSpec(ctx *parser.ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitImportLit(ctx *parser.ImportLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitImportLitGroup(ctx *parser.ImportLitGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitInfoBlock(ctx *parser.InfoBlockContext) interface{} {
	return v.VisitChildren(ctx)
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
	v.serviceGroup = new(spec.Group)
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitServerMeta(ctx *parser.ServerMetaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ApiVisitor) VisitAnnotation(ctx *parser.AnnotationContext) interface{} {
	return v.VisitChildren(ctx)
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

func (v *VisitResult) Result() (interface{}, error) {
	return v.v, v.err
}

func (v *ApiVisitor) getInt(token antlr.Token) (int64, error) {
	text, err := v.getText(token, true)
	if err != nil {
		return 0, err
	}

	vInt, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, err
	}

	return vInt, nil
}

func (v *ApiVisitor) getText(token antlr.Token, trimQuote bool) (string, error) {
	if token == nil {
		return "", nil
	}
	text := token.GetText()
	if trimQuote {
		text = v.trimQuote(text)
	}
	return text, nil
}

func (v *ApiVisitor) trimQuote(text string) string {
	text = strings.ReplaceAll(text, `"`, "")
	text = strings.ReplaceAll(text, `'`, "")
	text = strings.ReplaceAll(text, "`", "")
	return text
}
