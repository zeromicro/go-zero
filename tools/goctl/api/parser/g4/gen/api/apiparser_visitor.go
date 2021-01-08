// Code generated from tools/goctl/api/parser/g4/ApiParser.g4 by ANTLR 4.9. DO NOT EDIT.

package api // ApiParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ApiParserParser.
type ApiParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ApiParserParser#api.
	VisitApi(ctx *ApiContext) interface{}

	// Visit a parse tree produced by ApiParserParser#spec.
	VisitSpec(ctx *SpecContext) interface{}

	// Visit a parse tree produced by ApiParserParser#syntaxLit.
	VisitSyntaxLit(ctx *SyntaxLitContext) interface{}

	// Visit a parse tree produced by ApiParserParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by ApiParserParser#importLit.
	VisitImportLit(ctx *ImportLitContext) interface{}

	// Visit a parse tree produced by ApiParserParser#importBlock.
	VisitImportBlock(ctx *ImportBlockContext) interface{}

	// Visit a parse tree produced by ApiParserParser#importBlockValue.
	VisitImportBlockValue(ctx *ImportBlockValueContext) interface{}

	// Visit a parse tree produced by ApiParserParser#importValue.
	VisitImportValue(ctx *ImportValueContext) interface{}

	// Visit a parse tree produced by ApiParserParser#infoSpec.
	VisitInfoSpec(ctx *InfoSpecContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeSpec.
	VisitTypeSpec(ctx *TypeSpecContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeBlock.
	VisitTypeBlock(ctx *TypeBlockContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeLitBody.
	VisitTypeLitBody(ctx *TypeLitBodyContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeBlockBody.
	VisitTypeBlockBody(ctx *TypeBlockBodyContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeStruct.
	VisitTypeStruct(ctx *TypeStructContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeAlias.
	VisitTypeAlias(ctx *TypeAliasContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeBlockStruct.
	VisitTypeBlockStruct(ctx *TypeBlockStructContext) interface{}

	// Visit a parse tree produced by ApiParserParser#typeBlockAlias.
	VisitTypeBlockAlias(ctx *TypeBlockAliasContext) interface{}

	// Visit a parse tree produced by ApiParserParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by ApiParserParser#normalField.
	VisitNormalField(ctx *NormalFieldContext) interface{}

	// Visit a parse tree produced by ApiParserParser#anonymousFiled.
	VisitAnonymousFiled(ctx *AnonymousFiledContext) interface{}

	// Visit a parse tree produced by ApiParserParser#dataType.
	VisitDataType(ctx *DataTypeContext) interface{}

	// Visit a parse tree produced by ApiParserParser#pointerType.
	VisitPointerType(ctx *PointerTypeContext) interface{}

	// Visit a parse tree produced by ApiParserParser#mapType.
	VisitMapType(ctx *MapTypeContext) interface{}

	// Visit a parse tree produced by ApiParserParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by ApiParserParser#serviceSpec.
	VisitServiceSpec(ctx *ServiceSpecContext) interface{}

	// Visit a parse tree produced by ApiParserParser#atServer.
	VisitAtServer(ctx *AtServerContext) interface{}

	// Visit a parse tree produced by ApiParserParser#serviceApi.
	VisitServiceApi(ctx *ServiceApiContext) interface{}

	// Visit a parse tree produced by ApiParserParser#serviceRoute.
	VisitServiceRoute(ctx *ServiceRouteContext) interface{}

	// Visit a parse tree produced by ApiParserParser#atDoc.
	VisitAtDoc(ctx *AtDocContext) interface{}

	// Visit a parse tree produced by ApiParserParser#atHandler.
	VisitAtHandler(ctx *AtHandlerContext) interface{}

	// Visit a parse tree produced by ApiParserParser#route.
	VisitRoute(ctx *RouteContext) interface{}

	// Visit a parse tree produced by ApiParserParser#body.
	VisitBody(ctx *BodyContext) interface{}

	// Visit a parse tree produced by ApiParserParser#replybody.
	VisitReplybody(ctx *ReplybodyContext) interface{}

	// Visit a parse tree produced by ApiParserParser#kvLit.
	VisitKvLit(ctx *KvLitContext) interface{}

	// Visit a parse tree produced by ApiParserParser#serviceName.
	VisitServiceName(ctx *ServiceNameContext) interface{}

	// Visit a parse tree produced by ApiParserParser#path.
	VisitPath(ctx *PathContext) interface{}
}
