package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 4
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

type TypeBlockStructContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	structName  antlr.Token
	structToken antlr.Token
	lbrace      antlr.Token
	rbrace      antlr.Token
}

func NewEmptyTypeBlockStructContext() *TypeBlockStructContext {
	p := new(TypeBlockStructContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeBlockStruct
	return p
}

func (*TypeBlockStructContext) IsTypeBlockStructContext() {}

func NewTypeBlockStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeBlockStructContext {
	p := new(TypeBlockStructContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_typeBlockStruct

	return p
}

func (s *TypeBlockStructContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeBlockStructContext) GetStructName() antlr.Token { return s.structName }

func (s *TypeBlockStructContext) GetStructToken() antlr.Token { return s.structToken }

func (s *TypeBlockStructContext) GetLbrace() antlr.Token { return s.lbrace }

func (s *TypeBlockStructContext) GetRbrace() antlr.Token { return s.rbrace }

func (s *TypeBlockStructContext) SetStructName(v antlr.Token) { s.structName = v }

func (s *TypeBlockStructContext) SetStructToken(v antlr.Token) { s.structToken = v }

func (s *TypeBlockStructContext) SetLbrace(v antlr.Token) { s.lbrace = v }

func (s *TypeBlockStructContext) SetRbrace(v antlr.Token) { s.rbrace = v }

func (s *TypeBlockStructContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ApiParserParserID)
}

func (s *TypeBlockStructContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, i)
}

func (s *TypeBlockStructContext) AllField() []IFieldContext {
	ts := s.GetTypedRuleContexts(reflect.TypeOf((*IFieldContext)(nil)).Elem())
	tst := make([]IFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldContext)
		}
	}

	return tst
}

func (s *TypeBlockStructContext) Field(i int) IFieldContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *TypeBlockStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeBlockStructContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeBlockStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitTypeBlockStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) TypeBlockStruct() (localctx ITypeBlockStructContext) {
	localctx = NewTypeBlockStructContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, ApiParserParserRULE_typeBlockStruct)
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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	checkKeyword(p)
	{
		p.SetState(175)

		_m := p.Match(ApiParserParserID)

		localctx.(*TypeBlockStructContext).structName = _m
	}
	p.SetState(177)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserID {
		{
			p.SetState(176)

			_m := p.Match(ApiParserParserID)

			localctx.(*TypeBlockStructContext).structToken = _m
		}
	}
	{
		p.SetState(179)

		_m := p.Match(ApiParserParserT__3)

		localctx.(*TypeBlockStructContext).lbrace = _m
	}
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(180)
				p.Field()
			}
		}
		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())
	}
	{
		p.SetState(186)

		_m := p.Match(ApiParserParserT__4)

		localctx.(*TypeBlockStructContext).rbrace = _m
	}

	return localctx
}

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
	p := new(TypeBlockAliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeBlockAlias
	return p
}

func (*TypeBlockAliasContext) IsTypeBlockAliasContext() {}

func NewTypeBlockAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeBlockAliasContext {
	p := new(TypeBlockAliasContext)

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
	t := s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

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

func (s *TypeBlockAliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
		p.SetState(189)

		_m := p.Match(ApiParserParserID)

		localctx.(*TypeBlockAliasContext).alias = _m
	}
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__0 {
		{
			p.SetState(190)

			_m := p.Match(ApiParserParserT__0)

			localctx.(*TypeBlockAliasContext).assign = _m
		}
	}
	{
		p.SetState(193)
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
	p := new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	p := new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) NormalField() INormalFieldContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*INormalFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INormalFieldContext)
}

func (s *FieldContext) AnonymousFiled() IAnonymousFiledContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IAnonymousFiledContext)(nil)).Elem(), 0)

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

func (s *FieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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

	p.SetState(198)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(195)

		if !(isNormal(p)) {
			panic(antlr.NewFailedPredicateException(p, "isNormal(p)", ""))
		}
		{
			p.SetState(196)
			p.NormalField()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(197)
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
	p := new(NormalFieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_normalField
	return p
}

func (*NormalFieldContext) IsNormalFieldContext() {}

func NewNormalFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NormalFieldContext {
	p := new(NormalFieldContext)

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
	t := s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

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

func (s *NormalFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
		p.SetState(201)

		_m := p.Match(ApiParserParserID)

		localctx.(*NormalFieldContext).fieldName = _m
	}
	{
		p.SetState(202)
		p.DataType()
	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(203)

			_m := p.Match(ApiParserParserRAW_STRING)

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
