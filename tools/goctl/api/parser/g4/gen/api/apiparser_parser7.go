package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 7
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

// IAtHandlerContext is an interface to support dynamic dispatch.
type IAtHandlerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAtHandlerContext differentiates from other interfaces.
	IsAtHandlerContext()
}

type AtHandlerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtHandlerContext() *AtHandlerContext {
	var p = new(AtHandlerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atHandler
	return p
}

func (*AtHandlerContext) IsAtHandlerContext() {}

func NewAtHandlerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtHandlerContext {
	var p = new(AtHandlerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_atHandler

	return p
}

func (s *AtHandlerContext) GetParser() antlr.Parser { return s.parser }

func (s *AtHandlerContext) ATHANDLER() antlr.TerminalNode {
	return s.GetToken(ApiParserParserATHANDLER, 0)
}

func (s *AtHandlerContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *AtHandlerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtHandlerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtHandlerContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitAtHandler(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) AtHandler() (localctx IAtHandlerContext) {
	localctx = NewAtHandlerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, ApiParserParserRULE_atHandler)

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
	{
		p.SetState(289)
		p.Match(ApiParserParserATHANDLER)
	}
	{
		p.SetState(290)
		p.Match(ApiParserParserID)
	}

	return localctx
}

// IRouteContext is an interface to support dynamic dispatch.
type IRouteContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetHttpMethod returns the httpMethod token.
	GetHttpMethod() antlr.Token

	// SetHttpMethod sets the httpMethod token.
	SetHttpMethod(antlr.Token)

	// GetRequest returns the request rule contexts.
	GetRequest() IBodyContext

	// GetResponse returns the response rule contexts.
	GetResponse() IReplybodyContext

	// SetRequest sets the request rule contexts.
	SetRequest(IBodyContext)

	// SetResponse sets the response rule contexts.
	SetResponse(IReplybodyContext)

	// IsRouteContext differentiates from other interfaces.
	IsRouteContext()
}

type RouteContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	httpMethod antlr.Token
	request    IBodyContext
	response   IReplybodyContext
}

func NewEmptyRouteContext() *RouteContext {
	var p = new(RouteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_route
	return p
}

func (*RouteContext) IsRouteContext() {}

func NewRouteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RouteContext {
	var p = new(RouteContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_route

	return p
}

func (s *RouteContext) GetParser() antlr.Parser { return s.parser }

func (s *RouteContext) GetHttpMethod() antlr.Token { return s.httpMethod }

func (s *RouteContext) SetHttpMethod(v antlr.Token) { s.httpMethod = v }

func (s *RouteContext) GetRequest() IBodyContext { return s.request }

func (s *RouteContext) GetResponse() IReplybodyContext { return s.response }

func (s *RouteContext) SetRequest(v IBodyContext) { s.request = v }

func (s *RouteContext) SetResponse(v IReplybodyContext) { s.response = v }

func (s *RouteContext) Path() IPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathContext)
}

func (s *RouteContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *RouteContext) Body() IBodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
}

func (s *RouteContext) Replybody() IReplybodyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReplybodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReplybodyContext)
}

func (s *RouteContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RouteContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RouteContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitRoute(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Route() (localctx IRouteContext) {
	localctx = NewRouteContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, ApiParserParserRULE_route)
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
	checkHTTPMethod(p)
	{
		p.SetState(293)

		var _m = p.Match(ApiParserParserID)

		localctx.(*RouteContext).httpMethod = _m
	}
	{
		p.SetState(294)
		p.Path()
	}
	p.SetState(296)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__1 {
		{
			p.SetState(295)

			var _x = p.Body()

			localctx.(*RouteContext).request = _x
		}

	}
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__9 {
		{
			p.SetState(298)

			var _x = p.Replybody()

			localctx.(*RouteContext).response = _x
		}

	}

	return localctx
}

// IBodyContext is an interface to support dynamic dispatch.
type IBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLp returns the lp token.
	GetLp() antlr.Token

	// GetRp returns the rp token.
	GetRp() antlr.Token

	// SetLp sets the lp token.
	SetLp(antlr.Token)

	// SetRp sets the rp token.
	SetRp(antlr.Token)

	// IsBodyContext differentiates from other interfaces.
	IsBodyContext()
}

type BodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	lp     antlr.Token
	rp     antlr.Token
}

func NewEmptyBodyContext() *BodyContext {
	var p = new(BodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_body
	return p
}

func (*BodyContext) IsBodyContext() {}

func NewBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyContext {
	var p = new(BodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_body

	return p
}

func (s *BodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyContext) GetLp() antlr.Token { return s.lp }

func (s *BodyContext) GetRp() antlr.Token { return s.rp }

func (s *BodyContext) SetLp(v antlr.Token) { s.lp = v }

func (s *BodyContext) SetRp(v antlr.Token) { s.rp = v }

func (s *BodyContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *BodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Body() (localctx IBodyContext) {
	localctx = NewBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, ApiParserParserRULE_body)
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
	{
		p.SetState(301)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*BodyContext).lp = _m
	}
	p.SetState(303)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserID {
		{
			p.SetState(302)
			p.Match(ApiParserParserID)
		}

	}
	{
		p.SetState(305)

		var _m = p.Match(ApiParserParserT__2)

		localctx.(*BodyContext).rp = _m
	}

	return localctx
}

// IReplybodyContext is an interface to support dynamic dispatch.
type IReplybodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetReturnToken returns the returnToken token.
	GetReturnToken() antlr.Token

	// GetLp returns the lp token.
	GetLp() antlr.Token

	// GetRp returns the rp token.
	GetRp() antlr.Token

	// SetReturnToken sets the returnToken token.
	SetReturnToken(antlr.Token)

	// SetLp sets the lp token.
	SetLp(antlr.Token)

	// SetRp sets the rp token.
	SetRp(antlr.Token)

	// IsReplybodyContext differentiates from other interfaces.
	IsReplybodyContext()
}

type ReplybodyContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	returnToken antlr.Token
	lp          antlr.Token
	rp          antlr.Token
}

func NewEmptyReplybodyContext() *ReplybodyContext {
	var p = new(ReplybodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_replybody
	return p
}

func (*ReplybodyContext) IsReplybodyContext() {}

func NewReplybodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplybodyContext {
	var p = new(ReplybodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_replybody

	return p
}

func (s *ReplybodyContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplybodyContext) GetReturnToken() antlr.Token { return s.returnToken }

func (s *ReplybodyContext) GetLp() antlr.Token { return s.lp }

func (s *ReplybodyContext) GetRp() antlr.Token { return s.rp }

func (s *ReplybodyContext) SetReturnToken(v antlr.Token) { s.returnToken = v }

func (s *ReplybodyContext) SetLp(v antlr.Token) { s.lp = v }

func (s *ReplybodyContext) SetRp(v antlr.Token) { s.rp = v }

func (s *ReplybodyContext) DataType() IDataTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *ReplybodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplybodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplybodyContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitReplybody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Replybody() (localctx IReplybodyContext) {
	localctx = NewReplybodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, ApiParserParserRULE_replybody)
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
	{
		p.SetState(307)

		var _m = p.Match(ApiParserParserT__9)

		localctx.(*ReplybodyContext).returnToken = _m
	}
	{
		p.SetState(308)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*ReplybodyContext).lp = _m
	}
	p.SetState(310)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ApiParserParserT__5)|(1<<ApiParserParserT__6)|(1<<ApiParserParserT__7)|(1<<ApiParserParserINTERFACE)|(1<<ApiParserParserID))) != 0 {
		{
			p.SetState(309)
			p.DataType()
		}

	}
	{
		p.SetState(312)

		var _m = p.Match(ApiParserParserT__2)

		localctx.(*ReplybodyContext).rp = _m
	}

	return localctx
}

// IKvLitContext is an interface to support dynamic dispatch.
type IKvLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetKey returns the key token.
	GetKey() antlr.Token

	// GetValue returns the value token.
	GetValue() antlr.Token

	// SetKey sets the key token.
	SetKey(antlr.Token)

	// SetValue sets the value token.
	SetValue(antlr.Token)

	// IsKvLitContext differentiates from other interfaces.
	IsKvLitContext()
}

type KvLitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	key    antlr.Token
	value  antlr.Token
}

func NewEmptyKvLitContext() *KvLitContext {
	var p = new(KvLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_kvLit
	return p
}

func (*KvLitContext) IsKvLitContext() {}

func NewKvLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *KvLitContext {
	var p = new(KvLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_kvLit

	return p
}

func (s *KvLitContext) GetParser() antlr.Parser { return s.parser }

func (s *KvLitContext) GetKey() antlr.Token { return s.key }

func (s *KvLitContext) GetValue() antlr.Token { return s.value }

func (s *KvLitContext) SetKey(v antlr.Token) { s.key = v }

func (s *KvLitContext) SetValue(v antlr.Token) { s.value = v }

func (s *KvLitContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *KvLitContext) LINE_VALUE() antlr.TerminalNode {
	return s.GetToken(ApiParserParserLINE_VALUE, 0)
}

func (s *KvLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *KvLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}
