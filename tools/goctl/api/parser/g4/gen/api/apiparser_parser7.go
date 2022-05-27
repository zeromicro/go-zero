package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 7
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

// IServiceRouteContext is an interface to support dynamic dispatch.
type IServiceRouteContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsServiceRouteContext differentiates from other interfaces.
	IsServiceRouteContext()
}

type ServiceRouteContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyServiceRouteContext() *ServiceRouteContext {
	p := new(ServiceRouteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_serviceRoute
	return p
}

func (*ServiceRouteContext) IsServiceRouteContext() {}

func NewServiceRouteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ServiceRouteContext {
	p := new(ServiceRouteContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_serviceRoute

	return p
}

func (s *ServiceRouteContext) GetParser() antlr.Parser { return s.parser }

func (s *ServiceRouteContext) Route() IRouteContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IRouteContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRouteContext)
}

func (s *ServiceRouteContext) AtServer() IAtServerContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IAtServerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtServerContext)
}

func (s *ServiceRouteContext) AtHandler() IAtHandlerContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IAtHandlerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtHandlerContext)
}

func (s *ServiceRouteContext) AtDoc() IAtDocContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IAtDocContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtDocContext)
}

func (s *ServiceRouteContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ServiceRouteContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ServiceRouteContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitServiceRoute(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ServiceRoute() (localctx IServiceRouteContext) {
	localctx = NewServiceRouteContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, ApiParserParserRULE_serviceRoute)
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
	p.SetState(264)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserATDOC {
		{
			p.SetState(263)
			p.AtDoc()
		}
	}
	p.SetState(268)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ApiParserParserATSERVER:
		{
			p.SetState(266)
			p.AtServer()
		}

	case ApiParserParserATHANDLER:
		{
			p.SetState(267)
			p.AtHandler()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(270)
		p.Route()
	}

	return localctx
}

// IAtDocContext is an interface to support dynamic dispatch.
type IAtDocContext interface {
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

	// IsAtDocContext differentiates from other interfaces.
	IsAtDocContext()
}

type AtDocContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	lp     antlr.Token
	rp     antlr.Token
}

func NewEmptyAtDocContext() *AtDocContext {
	p := new(AtDocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atDoc
	return p
}

func (*AtDocContext) IsAtDocContext() {}

func NewAtDocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtDocContext {
	p := new(AtDocContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_atDoc

	return p
}

func (s *AtDocContext) GetParser() antlr.Parser { return s.parser }

func (s *AtDocContext) GetLp() antlr.Token { return s.lp }

func (s *AtDocContext) GetRp() antlr.Token { return s.rp }

func (s *AtDocContext) SetLp(v antlr.Token) { s.lp = v }

func (s *AtDocContext) SetRp(v antlr.Token) { s.rp = v }

func (s *AtDocContext) ATDOC() antlr.TerminalNode {
	return s.GetToken(ApiParserParserATDOC, 0)
}

func (s *AtDocContext) STRING() antlr.TerminalNode {
	return s.GetToken(ApiParserParserSTRING, 0)
}

func (s *AtDocContext) AllKvLit() []IKvLitContext {
	ts := s.GetTypedRuleContexts(reflect.TypeOf((*IKvLitContext)(nil)).Elem())
	tst := make([]IKvLitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IKvLitContext)
		}
	}

	return tst
}

func (s *AtDocContext) KvLit(i int) IKvLitContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IKvLitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IKvLitContext)
}

func (s *AtDocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtDocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtDocContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitAtDoc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) AtDoc() (localctx IAtDocContext) {
	localctx = NewAtDocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, ApiParserParserRULE_atDoc)
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
		p.SetState(272)
		p.Match(ApiParserParserATDOC)
	}
	p.SetState(274)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__1 {
		{
			p.SetState(273)

			_m := p.Match(ApiParserParserT__1)

			localctx.(*AtDocContext).lp = _m
		}
	}
	p.SetState(282)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ApiParserParserID:
		p.SetState(277)
		p.GetErrorHandler().Sync(p)

		for ok := true; ok; ok = _la == ApiParserParserID {
			{
				p.SetState(276)
				p.KvLit()
			}

			p.SetState(279)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case ApiParserParserSTRING:
		{
			p.SetState(281)
			p.Match(ApiParserParserSTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(285)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__2 {
		{
			p.SetState(284)

			_m := p.Match(ApiParserParserT__2)

			localctx.(*AtDocContext).rp = _m
		}
	}

	return localctx
}

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
	p := new(AtHandlerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atHandler
	return p
}

func (*AtHandlerContext) IsAtHandlerContext() {}

func NewAtHandlerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtHandlerContext {
	p := new(AtHandlerContext)

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

func (s *AtHandlerContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
		p.SetState(287)
		p.Match(ApiParserParserATHANDLER)
	}
	{
		p.SetState(288)
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
	p := new(RouteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_route
	return p
}

func (*RouteContext) IsRouteContext() {}

func NewRouteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RouteContext {
	p := new(RouteContext)

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
	t := s.GetTypedRuleContext(reflect.TypeOf((*IPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathContext)
}

func (s *RouteContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *RouteContext) Body() IBodyContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IBodyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
}

func (s *RouteContext) Replybody() IReplybodyContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IReplybodyContext)(nil)).Elem(), 0)

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

func (s *RouteContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
		p.SetState(291)

		_m := p.Match(ApiParserParserID)

		localctx.(*RouteContext).httpMethod = _m
	}
	{
		p.SetState(292)
		p.Path()
	}
	p.SetState(294)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__1 {
		{
			p.SetState(293)

			_x := p.Body()

			localctx.(*RouteContext).request = _x
		}
	}
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__9 {
		{
			p.SetState(296)

			_x := p.Replybody()

			localctx.(*RouteContext).response = _x
		}
	}

	return localctx
}
