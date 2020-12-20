// Code generated from /Users/anqiansong/goland/go/go-zero_kingxt/tools/goctl/api/parser/g4/ApiParser.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // ApiParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseApiParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseApiParserVisitor) VisitApi(ctx *ApiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitBody(ctx *BodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitSyntaxLit(ctx *SyntaxLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportLit(ctx *ImportLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportLitGroup(ctx *ImportLitGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitInfoBlock(ctx *InfoBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlock(ctx *TypeBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeLit(ctx *TypeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeGroup(ctx *TypeGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeSpec(ctx *TypeSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeAlias(ctx *TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeStruct(ctx *TypeStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeField(ctx *TypeFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitFiled(ctx *FiledContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitInnerStruct(ctx *InnerStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitDataType(ctx *DataTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitMapType(ctx *MapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPointer(ctx *PointerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceBlock(ctx *ServiceBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServerMeta(ctx *ServerMetaContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitIdValue(ctx *IdValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceBody(ctx *ServiceBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceName(ctx *ServiceNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceRoute(ctx *ServiceRouteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRouteDoc(ctx *RouteDocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitDoc(ctx *DocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitLineDoc(ctx *LineDocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRouteHandler(ctx *RouteHandlerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRoutePath(ctx *RoutePathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPath(ctx *PathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRequest(ctx *RequestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitReply(ctx *ReplyContext) interface{} {
	return v.VisitChildren(ctx)
}
