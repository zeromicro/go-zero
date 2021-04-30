// Code generated from tools/goctl/api/parser/g4/ApiParser.g4 by ANTLR 4.9. DO NOT EDIT.

package api // ApiParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseApiParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseApiParserVisitor) VisitApi(ctx *ApiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitSpec(ctx *SpecContext) interface{} {
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

func (v *BaseApiParserVisitor) VisitImportBlock(ctx *ImportBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportBlockValue(ctx *ImportBlockValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportValue(ctx *ImportValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitInfoSpec(ctx *InfoSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeSpec(ctx *TypeSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeLit(ctx *TypeLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlock(ctx *TypeBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeLitBody(ctx *TypeLitBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockBody(ctx *TypeBlockBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeStruct(ctx *TypeStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeAlias(ctx *TypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockStruct(ctx *TypeBlockStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockAlias(ctx *TypeBlockAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitField(ctx *FieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitNormalField(ctx *NormalFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAnonymousFiled(ctx *AnonymousFiledContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitDataType(ctx *DataTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPointerType(ctx *PointerTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitMapType(ctx *MapTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceSpec(ctx *ServiceSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtServer(ctx *AtServerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceApi(ctx *ServiceApiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceRoute(ctx *ServiceRouteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtDoc(ctx *AtDocContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtHandler(ctx *AtHandlerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRoute(ctx *RouteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitBody(ctx *BodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitReplybody(ctx *ReplybodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitKvLit(ctx *KvLitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceName(ctx *ServiceNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPath(ctx *PathContext) interface{} {
	return v.VisitChildren(ctx)
}
