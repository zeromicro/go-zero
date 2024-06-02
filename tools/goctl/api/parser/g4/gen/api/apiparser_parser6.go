package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 6
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

type AtServerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	lp     antlr.Token
	rp     antlr.Token
}

func NewEmptyAtServerContext() *AtServerContext {
	var p = new(AtServerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atServer
	return p
}

func (*AtServerContext) IsAtServerContext() {}

func NewAtServerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtServerContext {
	var p = new(AtServerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_atServer

	return p
}

func (s *AtServerContext) GetParser() antlr.Parser { return s.parser }

func (s *AtServerContext) GetLp() antlr.Token { return s.lp }

func (s *AtServerContext) GetRp() antlr.Token { return s.rp }

func (s *AtServerContext) SetLp(v antlr.Token) { s.lp = v }

func (s *AtServerContext) SetRp(v antlr.Token) { s.rp = v }

func (s *AtServerContext) ATSERVER() antlr.TerminalNode {
	return s.GetToken(ApiParserParserATSERVER, 0)
}

func (s *AtServerContext) AllKvLit() []IKvLitContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IKvLitContext)(nil)).Elem())
	var tst = make([]IKvLitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IKvLitContext)
		}
	}

	return tst
}

func (s *AtServerContext) KvLit(i int) IKvLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IKvLitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IKvLitContext)
}

func (s *AtServerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtServerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtServerContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitAtServer(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) AtServer() (localctx IAtServerContext) {
	localctx = NewAtServerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, ApiParserParserRULE_atServer)
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
		p.SetState(244)
		p.Match(ApiParserParserATSERVER)
	}
	{
		p.SetState(245)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*AtServerContext).lp = _m
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ApiParserParserID {
		{
			p.SetState(246)
			p.KvLit()
		}

		p.SetState(249)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(251)

		var _m = p.Match(ApiParserParserT__2)

		localctx.(*AtServerContext).rp = _m
	}

	return localctx
}

// IServiceApiContext is an interface to support dynamic dispatch.
type IServiceApiContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetServiceToken returns the serviceToken token.
	GetServiceToken() antlr.Token

	// GetLbrace returns the lbrace token.
	GetLbrace() antlr.Token

	// GetRbrace returns the rbrace token.
	GetRbrace() antlr.Token

	// SetServiceToken sets the serviceToken token.
	SetServiceToken(antlr.Token)

	// SetLbrace sets the lbrace token.
	SetLbrace(antlr.Token)

	// SetRbrace sets the rbrace token.
	SetRbrace(antlr.Token)

	// IsServiceApiContext differentiates from other interfaces.
	IsServiceApiContext()
}

type ServiceApiContext struct {
	*antlr.BaseParserRuleContext
	parser       antlr.Parser
	serviceToken antlr.Token
	lbrace       antlr.Token
	rbrace       antlr.Token
}

func NewEmptyServiceApiContext() *ServiceApiContext {
	var p = new(ServiceApiContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_serviceApi
	return p
}

func (*ServiceApiContext) IsServiceApiContext() {}

func NewServiceApiContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ServiceApiContext {
	var p = new(ServiceApiContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_serviceApi

	return p
}

func (s *ServiceApiContext) GetParser() antlr.Parser { return s.parser }

func (s *ServiceApiContext) GetServiceToken() antlr.Token { return s.serviceToken }

func (s *ServiceApiContext) GetLbrace() antlr.Token { return s.lbrace }

func (s *ServiceApiContext) GetRbrace() antlr.Token { return s.rbrace }

func (s *ServiceApiContext) SetServiceToken(v antlr.Token) { s.serviceToken = v }

func (s *ServiceApiContext) SetLbrace(v antlr.Token) { s.lbrace = v }

func (s *ServiceApiContext) SetRbrace(v antlr.Token) { s.rbrace = v }

func (s *ServiceApiContext) ServiceName() IServiceNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IServiceNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IServiceNameContext)
}

func (s *ServiceApiContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *ServiceApiContext) AllServiceRoute() []IServiceRouteContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IServiceRouteContext)(nil)).Elem())
	var tst = make([]IServiceRouteContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IServiceRouteContext)
		}
	}

	return tst
}

func (s *ServiceApiContext) ServiceRoute(i int) IServiceRouteContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IServiceRouteContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IServiceRouteContext)
}

func (s *ServiceApiContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ServiceApiContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ServiceApiContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitServiceApi(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ServiceApi() (localctx IServiceApiContext) {
	localctx = NewServiceApiContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, ApiParserParserRULE_serviceApi)
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
	match(p, "service")
	{
		p.SetState(254)

		var _m = p.Match(ApiParserParserID)

		localctx.(*ServiceApiContext).serviceToken = _m
	}
	{
		p.SetState(255)
		p.ServiceName()
	}
	{
		p.SetState(256)

		var _m = p.Match(ApiParserParserT__3)

		localctx.(*ServiceApiContext).lbrace = _m
	}
	p.SetState(260)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ApiParserParserATDOC)|(1<<ApiParserParserATHANDLER)|(1<<ApiParserParserATSERVER))) != 0 {
		{
			p.SetState(257)
			p.ServiceRoute()
		}

		p.SetState(262)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(263)

		var _m = p.Match(ApiParserParserT__4)

		localctx.(*ServiceApiContext).rbrace = _m
	}

	return localctx
}

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
	var p = new(ServiceRouteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_serviceRoute
	return p
}

func (*ServiceRouteContext) IsServiceRouteContext() {}

func NewServiceRouteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ServiceRouteContext {
	var p = new(ServiceRouteContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_serviceRoute

	return p
}

func (s *ServiceRouteContext) GetParser() antlr.Parser { return s.parser }

func (s *ServiceRouteContext) Route() IRouteContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRouteContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRouteContext)
}

func (s *ServiceRouteContext) AtServer() IAtServerContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtServerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtServerContext)
}

func (s *ServiceRouteContext) AtHandler() IAtHandlerContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtHandlerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtHandlerContext)
}

func (s *ServiceRouteContext) AtDoc() IAtDocContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtDocContext)(nil)).Elem(), 0)

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

func (s *ServiceRouteContext) Accept(visitor antlr.ParseTreeVisitor) any {
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
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserATDOC {
		{
			p.SetState(265)
			p.AtDoc()
		}

	}
	p.SetState(270)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ApiParserParserATSERVER:
		{
			p.SetState(268)
			p.AtServer()
		}

	case ApiParserParserATHANDLER:
		{
			p.SetState(269)
			p.AtHandler()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(272)
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
	var p = new(AtDocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atDoc
	return p
}

func (*AtDocContext) IsAtDocContext() {}

func NewAtDocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtDocContext {
	var p = new(AtDocContext)

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IKvLitContext)(nil)).Elem())
	var tst = make([]IKvLitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IKvLitContext)
		}
	}

	return tst
}

func (s *AtDocContext) KvLit(i int) IKvLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IKvLitContext)(nil)).Elem(), i)

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

func (s *AtDocContext) Accept(visitor antlr.ParseTreeVisitor) any {
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
		p.SetState(274)
		p.Match(ApiParserParserATDOC)
	}
	p.SetState(276)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__1 {
		{
			p.SetState(275)

			var _m = p.Match(ApiParserParserT__1)

			localctx.(*AtDocContext).lp = _m
		}

	}
	p.SetState(284)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case ApiParserParserID:
		p.SetState(279)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == ApiParserParserID {
			{
				p.SetState(278)
				p.KvLit()
			}

			p.SetState(281)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case ApiParserParserSTRING:
		{
			p.SetState(283)
			p.Match(ApiParserParserSTRING)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__2 {
		{
			p.SetState(286)

			var _m = p.Match(ApiParserParserT__2)

			localctx.(*AtDocContext).rp = _m
		}

	}

	return localctx
}
