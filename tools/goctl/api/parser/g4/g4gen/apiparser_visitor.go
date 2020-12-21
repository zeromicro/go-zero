// Code generated from /Users/anqiansong/goland/go/go-zero_kingxt/tools/goctl/api/parser/g4/ApiParser.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // ApiParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ApiParser.
type ApiParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ApiParser#api.
	VisitApi(ctx *ApiContext) interface{}

	// Visit a parse tree produced by ApiParser#body.
	VisitBody(ctx *BodyContext) interface{}

	// Visit a parse tree produced by ApiParser#syntaxLit.
	VisitSyntaxLit(ctx *SyntaxLitContext) interface{}

	// Visit a parse tree produced by ApiParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by ApiParser#importLit.
	VisitImportLit(ctx *ImportLitContext) interface{}

	// Visit a parse tree produced by ApiParser#importLitGroup.
	VisitImportLitGroup(ctx *ImportLitGroupContext) interface{}

	// Visit a parse tree produced by ApiParser#infoBlock.
	VisitInfoBlock(ctx *InfoBlockContext) interface{}

	// Visit a parse tree produced by ApiParser#typeBlock.
	VisitTypeBlock(ctx *TypeBlockContext) interface{}

	// Visit a parse tree produced by ApiParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) interface{}

	// Visit a parse tree produced by ApiParser#typeGroup.
	VisitTypeGroup(ctx *TypeGroupContext) interface{}

	// Visit a parse tree produced by ApiParser#typeSpec.
	VisitTypeSpec(ctx *TypeSpecContext) interface{}

	// Visit a parse tree produced by ApiParser#typeAlias.
	VisitTypeAlias(ctx *TypeAliasContext) interface{}

	// Visit a parse tree produced by ApiParser#typeStruct.
	VisitTypeStruct(ctx *TypeStructContext) interface{}

	// Visit a parse tree produced by ApiParser#typeField.
	VisitTypeField(ctx *TypeFieldContext) interface{}

	// Visit a parse tree produced by ApiParser#filed.
	VisitFiled(ctx *FiledContext) interface{}

	// Visit a parse tree produced by ApiParser#innerStruct.
	VisitInnerStruct(ctx *InnerStructContext) interface{}

	// Visit a parse tree produced by ApiParser#dataType.
	VisitDataType(ctx *DataTypeContext) interface{}

	// Visit a parse tree produced by ApiParser#mapType.
	VisitMapType(ctx *MapTypeContext) interface{}

	// Visit a parse tree produced by ApiParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by ApiParser#pointer.
	VisitPointer(ctx *PointerContext) interface{}

	// Visit a parse tree produced by ApiParser#serviceBlock.
	VisitServiceBlock(ctx *ServiceBlockContext) interface{}

	// Visit a parse tree produced by ApiParser#serverMeta.
	VisitServerMeta(ctx *ServerMetaContext) interface{}

	// Visit a parse tree produced by ApiParser#annotation.
	VisitAnnotation(ctx *AnnotationContext) interface{}

	// Visit a parse tree produced by ApiParser#annotationKeyValue.
	VisitAnnotationKeyValue(ctx *AnnotationKeyValueContext) interface{}

	// Visit a parse tree produced by ApiParser#serviceBody.
	VisitServiceBody(ctx *ServiceBodyContext) interface{}

	// Visit a parse tree produced by ApiParser#serviceName.
	VisitServiceName(ctx *ServiceNameContext) interface{}

	// Visit a parse tree produced by ApiParser#serviceRoute.
	VisitServiceRoute(ctx *ServiceRouteContext) interface{}

	// Visit a parse tree produced by ApiParser#routeDoc.
	VisitRouteDoc(ctx *RouteDocContext) interface{}

	// Visit a parse tree produced by ApiParser#doc.
	VisitDoc(ctx *DocContext) interface{}

	// Visit a parse tree produced by ApiParser#lineDoc.
	VisitLineDoc(ctx *LineDocContext) interface{}

	// Visit a parse tree produced by ApiParser#routeHandler.
	VisitRouteHandler(ctx *RouteHandlerContext) interface{}

	// Visit a parse tree produced by ApiParser#routePath.
	VisitRoutePath(ctx *RoutePathContext) interface{}

	// Visit a parse tree produced by ApiParser#path.
	VisitPath(ctx *PathContext) interface{}

	// Visit a parse tree produced by ApiParser#request.
	VisitRequest(ctx *RequestContext) interface{}

	// Visit a parse tree produced by ApiParser#reply.
	VisitReply(ctx *ReplyContext) interface{}

	// Visit a parse tree produced by ApiParser#kvLit.
	VisitKvLit(ctx *KvLitContext) interface{}
}
