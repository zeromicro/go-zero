// Code generated from C:/Users/keson/GolandProjects/go-zero/tools/goctl/api/parser/g4\ApiParser.g4 by ANTLR 4.9. DO NOT EDIT.

package api // ApiParser
import "github.com/zeromicro/antlr"

// A complete Visitor for a parse tree produced by ApiParserParser.
type ApiParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ApiParserParser#api.
	VisitApi(ctx *ApiContext) any

	// Visit a parse tree produced by ApiParserParser#spec.
	VisitSpec(ctx *SpecContext) any

	// Visit a parse tree produced by ApiParserParser#syntaxLit.
	VisitSyntaxLit(ctx *SyntaxLitContext) any

	// Visit a parse tree produced by ApiParserParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) any

	// Visit a parse tree produced by ApiParserParser#importLit.
	VisitImportLit(ctx *ImportLitContext) any

	// Visit a parse tree produced by ApiParserParser#importBlock.
	VisitImportBlock(ctx *ImportBlockContext) any

	// Visit a parse tree produced by ApiParserParser#importBlockValue.
	VisitImportBlockValue(ctx *ImportBlockValueContext) any

	// Visit a parse tree produced by ApiParserParser#importValue.
	VisitImportValue(ctx *ImportValueContext) any

	// Visit a parse tree produced by ApiParserParser#infoSpec.
	VisitInfoSpec(ctx *InfoSpecContext) any

	// Visit a parse tree produced by ApiParserParser#typeSpec.
	VisitTypeSpec(ctx *TypeSpecContext) any

	// Visit a parse tree produced by ApiParserParser#typeLit.
	VisitTypeLit(ctx *TypeLitContext) any

	// Visit a parse tree produced by ApiParserParser#typeBlock.
	VisitTypeBlock(ctx *TypeBlockContext) any

	// Visit a parse tree produced by ApiParserParser#typeLitBody.
	VisitTypeLitBody(ctx *TypeLitBodyContext) any

	// Visit a parse tree produced by ApiParserParser#typeBlockBody.
	VisitTypeBlockBody(ctx *TypeBlockBodyContext) any

	// Visit a parse tree produced by ApiParserParser#typeStruct.
	VisitTypeStruct(ctx *TypeStructContext) any

	// Visit a parse tree produced by ApiParserParser#typeAlias.
	VisitTypeAlias(ctx *TypeAliasContext) any

	// Visit a parse tree produced by ApiParserParser#typeBlockStruct.
	VisitTypeBlockStruct(ctx *TypeBlockStructContext) any

	// Visit a parse tree produced by ApiParserParser#typeBlockAlias.
	VisitTypeBlockAlias(ctx *TypeBlockAliasContext) any

	// Visit a parse tree produced by ApiParserParser#field.
	VisitField(ctx *FieldContext) any

	// Visit a parse tree produced by ApiParserParser#normalField.
	VisitNormalField(ctx *NormalFieldContext) any

	// Visit a parse tree produced by ApiParserParser#anonymousFiled.
	VisitAnonymousFiled(ctx *AnonymousFiledContext) any

	// Visit a parse tree produced by ApiParserParser#dataType.
	VisitDataType(ctx *DataTypeContext) any

	// Visit a parse tree produced by ApiParserParser#pointerType.
	VisitPointerType(ctx *PointerTypeContext) any

	// Visit a parse tree produced by ApiParserParser#mapType.
	VisitMapType(ctx *MapTypeContext) any

	// Visit a parse tree produced by ApiParserParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) any

	// Visit a parse tree produced by ApiParserParser#serviceSpec.
	VisitServiceSpec(ctx *ServiceSpecContext) any

	// Visit a parse tree produced by ApiParserParser#atServer.
	VisitAtServer(ctx *AtServerContext) any

	// Visit a parse tree produced by ApiParserParser#serviceApi.
	VisitServiceApi(ctx *ServiceApiContext) any

	// Visit a parse tree produced by ApiParserParser#serviceRoute.
	VisitServiceRoute(ctx *ServiceRouteContext) any

	// Visit a parse tree produced by ApiParserParser#atDoc.
	VisitAtDoc(ctx *AtDocContext) any

	// Visit a parse tree produced by ApiParserParser#atHandler.
	VisitAtHandler(ctx *AtHandlerContext) any

	// Visit a parse tree produced by ApiParserParser#route.
	VisitRoute(ctx *RouteContext) any

	// Visit a parse tree produced by ApiParserParser#body.
	VisitBody(ctx *BodyContext) any

	// Visit a parse tree produced by ApiParserParser#replybody.
	VisitReplybody(ctx *ReplybodyContext) any

	// Visit a parse tree produced by ApiParserParser#kvLit.
	VisitKvLit(ctx *KvLitContext) any

	// Visit a parse tree produced by ApiParserParser#serviceName.
	VisitServiceName(ctx *ServiceNameContext) any

	// Visit a parse tree produced by ApiParserParser#path.
	VisitPath(ctx *PathContext) any

	// Visit a parse tree produced by ApiParserParser#pathItem.
	VisitPathItem(ctx *PathItemContext) any
}
