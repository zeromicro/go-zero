package api // ApiParser

import "github.com/zeromicro/antlr"

// ApiParserVisitor is a complete Visitor for a parse tree produced by ApiParserParser.
type ApiParserVisitor interface {
	antlr.ParseTreeVisitor

	// VisitApi is a parse tree produced by ApiParserParser#api.
	VisitApi(ctx *ApiContext) interface{}

	// VisitSpec is a parse tree produced by ApiParserParser#spec.
	VisitSpec(ctx *SpecContext) interface{}

	// VisitSyntaxLit is a parse tree produced by ApiParserParser#syntaxLit.
	VisitSyntaxLit(ctx *SyntaxLitContext) interface{}

	// VisitImportSpec is a parse tree produced by ApiParserParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// VisitImportLit is a parse tree produced by ApiParserParser#importLit.
	VisitImportLit(ctx *ImportLitContext) interface{}

	// VisitImportBlock is a parse tree produced by ApiParserParser#importBlock.
	VisitImportBlock(ctx *ImportBlockContext) interface{}

	// VisitImportBlockValue is a parse tree produced by ApiParserParser#importBlockValue.
	VisitImportBlockValue(ctx *ImportBlockValueContext) interface{}

	// VisitImportValue is a parse tree produced by ApiParserParser#importValue.
	VisitImportValue(ctx *ImportValueContext) interface{}

	// VisitInfoSpec is a parse tree produced by ApiParserParser#infoSpec.
	VisitInfoSpec(ctx *InfoSpecContext) interface{}

	// VisitTypeSpec is a parse tree produced by ApiParserParser#typeSpec.
	VisitTypeSpec(ctx *TypeSpecContext) interface{}

	// VisitTypeLit is a parse tree produced by ApiParserParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) interface{}

	// VisitTypeBlock is a parse tree produced by ApiParserParser#typeBlock.
	VisitTypeBlock(ctx *TypeBlockContext) interface{}

	// VisitTypeLitBody is a parse tree produced by ApiParserParser#typeLitBody.
	VisitTypeLitBody(ctx *TypeLitBodyContext) interface{}

	// VisitTypeBlockBody is a parse tree produced by ApiParserParser#typeBlockBody.
	VisitTypeBlockBody(ctx *TypeBlockBodyContext) interface{}

	// VisitTypeStruct is a parse tree produced by ApiParserParser#typeStruct.
	VisitTypeStruct(ctx *TypeStructContext) interface{}

	// VisitTypeAlias is a parse tree produced by ApiParserParser#typeAlias.
	VisitTypeAlias(ctx *TypeAliasContext) interface{}

	// VisitTypeBlockStruct is a parse tree produced by ApiParserParser#typeBlockStruct.
	VisitTypeBlockStruct(ctx *TypeBlockStructContext) interface{}

	// VisitTypeBlockAlias is a parse tree produced by ApiParserParser#typeBlockAlias.
	VisitTypeBlockAlias(ctx *TypeBlockAliasContext) interface{}

	// VisitField is a parse tree produced by ApiParserParser#field.
	VisitField(ctx *FieldContext) interface{}

	// VisitNormalField is a parse tree produced by ApiParserParser#normalField.
	VisitNormalField(ctx *NormalFieldContext) interface{}

	// VisitAnonymousFiled is a parse tree produced by ApiParserParser#anonymousFiled.
	VisitAnonymousFiled(ctx *AnonymousFiledContext) interface{}

	// VisitDataType is a parse tree produced by ApiParserParser#dataType.
	VisitDataType(ctx *DataTypeContext) interface{}

	// VisitPointerType is a parse tree produced by ApiParserParser#pointerType.
	VisitPointerType(ctx *PointerTypeContext) interface{}

	// VisitMapType is a parse tree produced by ApiParserParser#mapType.
	VisitMapType(ctx *MapTypeContext) interface{}

	// VisitArrayType is a parse tree produced by ApiParserParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// VisitServiceSpec is a parse tree produced by ApiParserParser#serviceSpec.
	VisitServiceSpec(ctx *ServiceSpecContext) interface{}

	// VisitAtServer is a parse tree produced by ApiParserParser#atServer.
	VisitAtServer(ctx *AtServerContext) interface{}

	// VisitServiceApi is a parse tree produced by ApiParserParser#serviceApi.
	VisitServiceApi(ctx *ServiceApiContext) interface{}

	// VisitServiceRoute is a parse tree produced by ApiParserParser#serviceRoute.
	VisitServiceRoute(ctx *ServiceRouteContext) interface{}

	// VisitAtDoc is a parse tree produced by ApiParserParser#atDoc.
	VisitAtDoc(ctx *AtDocContext) interface{}

	// VisitAtHandler is a parse tree produced by ApiParserParser#atHandler.
	VisitAtHandler(ctx *AtHandlerContext) interface{}

	// VisitRoute is a parse tree produced by ApiParserParser#route.
	VisitRoute(ctx *RouteContext) interface{}

	// VisitBody is a parse tree produced by ApiParserParser#body.
	VisitBody(ctx *BodyContext) interface{}

	// VisitReplybody is a parse tree produced by ApiParserParser#replybody.
	VisitReplybody(ctx *ReplybodyContext) interface{}

	// VisitKvLit is a parse tree produced by ApiParserParser#kvLit.
	VisitKvLit(ctx *KvLitContext) interface{}

	// VisitServiceName is a parse tree produced by ApiParserParser#serviceName.
	VisitServiceName(ctx *ServiceNameContext) interface{}

	// VisitPath is a parse tree produced by ApiParserParser#path.
	VisitPath(ctx *PathContext) interface{}
}
