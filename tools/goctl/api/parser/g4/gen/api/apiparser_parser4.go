package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 4
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

// ITypeBlockAliasContext is an interface to support dynamic dispatch.
type ITypeBlockAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetAlias returns the alias token.
	GetAlias() antlr.Token

	// GetAssign returns the assign token.
	GetAssign() antlr.Token

	// SetAlias sets the alias token.
	SetAlias(antlr.Token)

	// SetAssign sets the assign token.
	SetAssign(antlr.Token)

	// IsTypeBlockAliasContext differentiates from other interfaces.
	IsTypeBlockAliasContext()
}

type TypeBlockAliasContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	alias  antlr.Token
	assign antlr.Token
}

func NewEmptyTypeBlockAliasContext() *TypeBlockAliasContext {
	var p = new(TypeBlockAliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeBlockAlias
	return p
}

func (*TypeBlockAliasContext) IsTypeBlockAliasContext() {}

func NewTypeBlockAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeBlockAliasContext {
	var p = new(TypeBlockAliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_typeBlockAlias

	return p
}

func (s *TypeBlockAliasContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeBlockAliasContext) GetAlias() antlr.Token { return s.alias }

func (s *TypeBlockAliasContext) GetAssign() antlr.Token { return s.assign }

func (s *TypeBlockAliasContext) SetAlias(v antlr.Token) { s.alias = v }

func (s *TypeBlockAliasContext) SetAssign(v antlr.Token) { s.assign = v }

func (s *TypeBlockAliasContext) DataType() IDataTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *TypeBlockAliasContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *TypeBlockAliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeBlockAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeBlockAliasContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitTypeBlockAlias(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) TypeBlockAlias() (localctx ITypeBlockAliasContext) {
	localctx = NewTypeBlockAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, ApiParserParserRULE_typeBlockAlias)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	checkKeyword(p)
	{
		p.SetState(191)

		var _m = p.Match(ApiParserParserID)

		localctx.(*TypeBlockAliasContext).alias = _m
	}
	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__0 {
		{
			p.SetState(192)

			var _m = p.Match(ApiParserParserT__0)

			localctx.(*TypeBlockAliasContext).assign = _m
		}

	}
	{
		p.SetState(195)
		p.DataType()
	}

	return localctx
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) NormalField() INormalFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INormalFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INormalFieldContext)
}

func (s *FieldContext) AnonymousFiled() IAnonymousFiledContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnonymousFiledContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAnonymousFiledContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, ApiParserParserRULE_field)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(200)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(197)

		if !(isNormal(p)) {
			panic(antlr.NewFailedPredicateException(p, "isNormal(p)", ""))
		}
		{
			p.SetState(198)
			p.NormalField()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(199)
			p.AnonymousFiled()
		}

	}

	return localctx
}

// INormalFieldContext is an interface to support dynamic dispatch.
type INormalFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetFieldName returns the fieldName token.
	GetFieldName() antlr.Token

	// GetTag returns the tag token.
	GetTag() antlr.Token

	// SetFieldName sets the fieldName token.
	SetFieldName(antlr.Token)

	// SetTag sets the tag token.
	SetTag(antlr.Token)

	// IsNormalFieldContext differentiates from other interfaces.
	IsNormalFieldContext()
}

type NormalFieldContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	fieldName antlr.Token
	tag       antlr.Token
}

func NewEmptyNormalFieldContext() *NormalFieldContext {
	var p = new(NormalFieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_normalField
	return p
}

func (*NormalFieldContext) IsNormalFieldContext() {}

func NewNormalFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NormalFieldContext {
	var p = new(NormalFieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_normalField

	return p
}

func (s *NormalFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *NormalFieldContext) GetFieldName() antlr.Token { return s.fieldName }

func (s *NormalFieldContext) GetTag() antlr.Token { return s.tag }

func (s *NormalFieldContext) SetFieldName(v antlr.Token) { s.fieldName = v }

func (s *NormalFieldContext) SetTag(v antlr.Token) { s.tag = v }

func (s *NormalFieldContext) DataType() IDataTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *NormalFieldContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *NormalFieldContext) RAW_STRING() antlr.TerminalNode {
	return s.GetToken(ApiParserParserRAW_STRING, 0)
}

func (s *NormalFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NormalFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NormalFieldContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitNormalField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) NormalField() (localctx INormalFieldContext) {
	localctx = NewNormalFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, ApiParserParserRULE_normalField)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	checkKeyword(p)
	{
		p.SetState(203)

		var _m = p.Match(ApiParserParserID)

		localctx.(*NormalFieldContext).fieldName = _m
	}
	{
		p.SetState(204)
		p.DataType()
	}
	p.SetState(206)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(205)

			var _m = p.Match(ApiParserParserRAW_STRING)

			localctx.(*NormalFieldContext).tag = _m
		}

	}

	return localctx
}

// IAnonymousFiledContext is an interface to support dynamic dispatch.
type IAnonymousFiledContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetStar returns the star token.
	GetStar() antlr.Token

	// SetStar sets the star token.
	SetStar(antlr.Token)

	// IsAnonymousFiledContext differentiates from other interfaces.
	IsAnonymousFiledContext()
}

type AnonymousFiledContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	star   antlr.Token
}

func NewEmptyAnonymousFiledContext() *AnonymousFiledContext {
	var p = new(AnonymousFiledContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_anonymousFiled
	return p
}

func (*AnonymousFiledContext) IsAnonymousFiledContext() {}

func NewAnonymousFiledContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnonymousFiledContext {
	var p = new(AnonymousFiledContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_anonymousFiled

	return p
}

func (s *AnonymousFiledContext) GetParser() antlr.Parser { return s.parser }

func (s *AnonymousFiledContext) GetStar() antlr.Token { return s.star }

func (s *AnonymousFiledContext) SetStar(v antlr.Token) { s.star = v }

func (s *AnonymousFiledContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *AnonymousFiledContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnonymousFiledContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnonymousFiledContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitAnonymousFiled(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) AnonymousFiled() (localctx IAnonymousFiledContext) {
	localctx = NewAnonymousFiledContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, ApiParserParserRULE_anonymousFiled)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(209)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__5 {
		{
			p.SetState(208)

			var _m = p.Match(ApiParserParserT__5)

			localctx.(*AnonymousFiledContext).star = _m
		}

	}
	{
		p.SetState(211)
		p.Match(ApiParserParserID)
	}

	return localctx
}

// IDataTypeContext is an interface to support dynamic dispatch.
type IDataTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetInter returns the inter token.
	GetInter() antlr.Token

	// GetTime returns the time token.
	GetTime() antlr.Token

	// SetInter sets the inter token.
	SetInter(antlr.Token)

	// SetTime sets the time token.
	SetTime(antlr.Token)

	// IsDataTypeContext differentiates from other interfaces.
	IsDataTypeContext()
}

type DataTypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	inter  antlr.Token
	time   antlr.Token
}

func NewEmptyDataTypeContext() *DataTypeContext {
	var p = new(DataTypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_dataType
	return p
}

func (*DataTypeContext) IsDataTypeContext() {}

func NewDataTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataTypeContext {
	var p = new(DataTypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_dataType

	return p
}

func (s *DataTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *DataTypeContext) GetInter() antlr.Token { return s.inter }

func (s *DataTypeContext) GetTime() antlr.Token { return s.time }

func (s *DataTypeContext) SetInter(v antlr.Token) { s.inter = v }

func (s *DataTypeContext) SetTime(v antlr.Token) { s.time = v }

func (s *DataTypeContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *DataTypeContext) MapType() IMapTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMapTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMapTypeContext)
}

func (s *DataTypeContext) ArrayType() IArrayTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArrayTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArrayTypeContext)
}

func (s *DataTypeContext) INTERFACE() antlr.TerminalNode {
	return s.GetToken(ApiParserParserINTERFACE, 0)
}

func (s *DataTypeContext) PointerType() IPointerTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPointerTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPointerTypeContext)
}

func (s *DataTypeContext) TypeStruct() ITypeStructContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeStructContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeStructContext)
}

func (s *DataTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DataTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DataTypeContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitDataType(s)

	default:
		return t.VisitChildren(s)
	}
}
