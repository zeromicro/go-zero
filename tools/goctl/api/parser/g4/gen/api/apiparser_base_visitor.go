package api // ApiParser
import "github.com/zeromicro/antlr"

type BaseApiParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseApiParserVisitor) VisitApi(ctx *ApiContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitSpec(ctx *SpecContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitSyntaxLit(ctx *SyntaxLitContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportSpec(ctx *ImportSpecContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportLit(ctx *ImportLitContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportBlock(ctx *ImportBlockContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportBlockValue(ctx *ImportBlockValueContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitImportValue(ctx *ImportValueContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitInfoSpec(ctx *InfoSpecContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeSpec(ctx *TypeSpecContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeLit(ctx *TypeLitContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlock(ctx *TypeBlockContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeLitBody(ctx *TypeLitBodyContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockBody(ctx *TypeBlockBodyContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeStruct(ctx *TypeStructContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeAlias(ctx *TypeAliasContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockStruct(ctx *TypeBlockStructContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitTypeBlockAlias(ctx *TypeBlockAliasContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitField(ctx *FieldContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitNormalField(ctx *NormalFieldContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAnonymousFiled(ctx *AnonymousFiledContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitDataType(ctx *DataTypeContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPointerType(ctx *PointerTypeContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitMapType(ctx *MapTypeContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitArrayType(ctx *ArrayTypeContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceSpec(ctx *ServiceSpecContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtServer(ctx *AtServerContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceApi(ctx *ServiceApiContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceRoute(ctx *ServiceRouteContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtDoc(ctx *AtDocContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitAtHandler(ctx *AtHandlerContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitRoute(ctx *RouteContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitBody(ctx *BodyContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitReplybody(ctx *ReplybodyContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitKvLit(ctx *KvLitContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitServiceName(ctx *ServiceNameContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPath(ctx *PathContext) any {
	return v.VisitChildren(ctx)
}

func (v *BaseApiParserVisitor) VisitPathItem(ctx *PathItemContext) any {
	return v.VisitChildren(ctx)
}
