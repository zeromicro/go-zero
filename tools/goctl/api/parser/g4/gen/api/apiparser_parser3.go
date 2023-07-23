package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 3
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

// ITypeBlockBodyContext is an interface to support dynamic dispatch.
type ITypeBlockBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeBlockBodyContext differentiates from other interfaces.
	IsTypeBlockBodyContext()
}

type TypeBlockBodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeBlockBodyContext() *TypeBlockBodyContext {
	var p = new(TypeBlockBodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeBlockBody
	return p
}

func (*TypeBlockBodyContext) IsTypeBlockBodyContext() {}

func NewTypeBlockBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeBlockBodyContext {
	var p = new(TypeBlockBodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_typeBlockBody

	return p
}

func (s *TypeBlockBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeBlockBodyContext) TypeBlockStruct() ITypeBlockStructContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeBlockStructContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeBlockStructContext)
}

func (s *TypeBlockBodyContext) TypeBlockAlias() ITypeBlockAliasContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeBlockAliasContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeBlockAliasContext)
}

func (s *TypeBlockBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeBlockBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeBlockBodyContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitTypeBlockBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) TypeBlockBody() (localctx ITypeBlockBodyContext) {
	localctx = NewTypeBlockBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, ApiParserParserRULE_typeBlockBody)

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

	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.TypeBlockStruct()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.TypeBlockAlias()
		}

	}

	return localctx
}

// ITypeStructContext is an interface to support dynamic dispatch.
type ITypeStructContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetStructName returns the structName token.
	GetStructName() antlr.Token

	// GetStructToken returns the structToken token.
	GetStructToken() antlr.Token

	// GetLbrace returns the lbrace token.
	GetLbrace() antlr.Token

	// GetRbrace returns the rbrace token.
	GetRbrace() antlr.Token

	// SetStructName sets the structName token.
	SetStructName(antlr.Token)

	// SetStructToken sets the structToken token.
	SetStructToken(antlr.Token)

	// SetLbrace sets the lbrace token.
	SetLbrace(antlr.Token)

	// SetRbrace sets the rbrace token.
	SetRbrace(antlr.Token)

	// IsTypeStructContext differentiates from other interfaces.
	IsTypeStructContext()
}

type TypeStructContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	structName  antlr.Token
	structToken antlr.Token
	lbrace      antlr.Token
	rbrace      antlr.Token
}

func NewEmptyTypeStructContext() *TypeStructContext {
	var p = new(TypeStructContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeStruct
	return p
}

func (*TypeStructContext) IsTypeStructContext() {}

func NewTypeStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeStructContext {
	var p = new(TypeStructContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_typeStruct

	return p
}

func (s *TypeStructContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeStructContext) GetStructName() antlr.Token { return s.structName }

func (s *TypeStructContext) GetStructToken() antlr.Token { return s.structToken }

func (s *TypeStructContext) GetLbrace() antlr.Token { return s.lbrace }

func (s *TypeStructContext) GetRbrace() antlr.Token { return s.rbrace }

func (s *TypeStructContext) SetStructName(v antlr.Token) { s.structName = v }

func (s *TypeStructContext) SetStructToken(v antlr.Token) { s.structToken = v }

func (s *TypeStructContext) SetLbrace(v antlr.Token) { s.lbrace = v }

func (s *TypeStructContext) SetRbrace(v antlr.Token) { s.rbrace = v }

func (s *TypeStructContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ApiParserParserID)
}

func (s *TypeStructContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, i)
}

func (s *TypeStructContext) AllField() []IFieldContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldContext)(nil)).Elem())
	var tst = make([]IFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldContext)
		}
	}

	return tst
}

func (s *TypeStructContext) Field(i int) IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *TypeStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeStructContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeStructContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitTypeStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) TypeStruct() (localctx ITypeStructContext) {
	localctx = NewTypeStructContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, ApiParserParserRULE_typeStruct)
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
		p.SetState(156)

		var _m = p.Match(ApiParserParserID)

		localctx.(*TypeStructContext).structName = _m
	}
	p.SetState(158)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserID {
		{
			p.SetState(157)

			var _m = p.Match(ApiParserParserID)

			localctx.(*TypeStructContext).structToken = _m
		}

	}
	{
		p.SetState(160)

		var _m = p.Match(ApiParserParserT__3)

		localctx.(*TypeStructContext).lbrace = _m
	}
	p.SetState(164)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(161)
				p.Field()
			}

		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
	}
	{
		p.SetState(167)

		var _m = p.Match(ApiParserParserT__4)

		localctx.(*TypeStructContext).rbrace = _m
	}

	return localctx
}

// ITypeAliasContext is an interface to support dynamic dispatch.
type ITypeAliasContext interface {
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

	// IsTypeAliasContext differentiates from other interfaces.
	IsTypeAliasContext()
}

type TypeAliasContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	alias  antlr.Token
	assign antlr.Token
}

func NewEmptyTypeAliasContext() *TypeAliasContext {
	var p = new(TypeAliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeAlias
	return p
}

func (*TypeAliasContext) IsTypeAliasContext() {}

func NewTypeAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAliasContext {
	var p = new(TypeAliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_typeAlias

	return p
}

func (s *TypeAliasContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeAliasContext) GetAlias() antlr.Token { return s.alias }

func (s *TypeAliasContext) GetAssign() antlr.Token { return s.assign }

func (s *TypeAliasContext) SetAlias(v antlr.Token) { s.alias = v }

func (s *TypeAliasContext) SetAssign(v antlr.Token) { s.assign = v }

func (s *TypeAliasContext) DataType() IDataTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *TypeAliasContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *TypeAliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeAliasContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitTypeAlias(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) TypeAlias() (localctx ITypeAliasContext) {
	localctx = NewTypeAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, ApiParserParserRULE_typeAlias)
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
		p.SetState(170)

		var _m = p.Match(ApiParserParserID)

		localctx.(*TypeAliasContext).alias = _m
	}
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__0 {
		{
			p.SetState(171)

			var _m = p.Match(ApiParserParserT__0)

			localctx.(*TypeAliasContext).assign = _m
		}

	}
	{
		p.SetState(174)
		p.DataType()
	}

	return localctx
}

// ITypeBlockStructContext is an interface to support dynamic dispatch.
type ITypeBlockStructContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetStructName returns the structName token.
	GetStructName() antlr.Token

	// GetStructToken returns the structToken token.
	GetStructToken() antlr.Token

	// GetLbrace returns the lbrace token.
	GetLbrace() antlr.Token

	// GetRbrace returns the rbrace token.
	GetRbrace() antlr.Token

	// SetStructName sets the structName token.
	SetStructName(antlr.Token)

	// SetStructToken sets the structToken token.
	SetStructToken(antlr.Token)

	// SetLbrace sets the lbrace token.
	SetLbrace(antlr.Token)

	// SetRbrace sets the rbrace token.
	SetRbrace(antlr.Token)

	// IsTypeBlockStructContext differentiates from other interfaces.
	IsTypeBlockStructContext()
}

type TypeBlockStructContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	structName  antlr.Token
	structToken antlr.Token
	lbrace      antlr.Token
	rbrace      antlr.Token
}

func NewEmptyTypeBlockStructContext() *TypeBlockStructContext {
	var p = new(TypeBlockStructContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_typeBlockStruct
	return p
}

func (*TypeBlockStructContext) IsTypeBlockStructContext() {}

func NewTypeBlockStructContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeBlockStructContext {
	var p = new(TypeBlockStructContext)

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldContext)(nil)).Elem())
	var tst = make([]IFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldContext)
		}
	}

	return tst
}

func (s *TypeBlockStructContext) Field(i int) IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), i)

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

func (s *TypeBlockStructContext) Accept(visitor antlr.ParseTreeVisitor) any {
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
		p.SetState(177)

		var _m = p.Match(ApiParserParserID)

		localctx.(*TypeBlockStructContext).structName = _m
	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserID {
		{
			p.SetState(178)

			var _m = p.Match(ApiParserParserID)

			localctx.(*TypeBlockStructContext).structToken = _m
		}

	}
	{
		p.SetState(181)

		var _m = p.Match(ApiParserParserT__3)

		localctx.(*TypeBlockStructContext).lbrace = _m
	}
	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(182)
				p.Field()
			}

		}
		p.SetState(187)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())
	}
	{
		p.SetState(188)

		var _m = p.Match(ApiParserParserT__4)

		localctx.(*TypeBlockStructContext).rbrace = _m
	}

	return localctx
}
