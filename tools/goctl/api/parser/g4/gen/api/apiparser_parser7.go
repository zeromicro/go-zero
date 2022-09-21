package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 7
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

// IAtRespDocContext is an interface to support dynamic dispatch.
type IAtRespDocContext interface {
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

	// GetCode returns the code rule contexts.
	GetCode() IRespDocCodeContext

	// SetCode sets the code rule contexts.
	SetCode(IRespDocCodeContext)

	// IsAtRespDocContext differentiates from other interfaces.
	IsAtRespDocContext()
}

type AtRespDocContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	code   IRespDocCodeContext
	lp     antlr.Token
	rp     antlr.Token
}

func NewEmptyAtRespDocContext() *AtRespDocContext {
	var p = new(AtRespDocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_atRespDoc
	return p
}

func (*AtRespDocContext) IsAtRespDocContext() {}

func NewAtRespDocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtRespDocContext {
	var p = new(AtRespDocContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_atRespDoc

	return p
}

func (s *AtRespDocContext) GetParser() antlr.Parser { return s.parser }

func (s *AtRespDocContext) GetLp() antlr.Token { return s.lp }

func (s *AtRespDocContext) GetRp() antlr.Token { return s.rp }

func (s *AtRespDocContext) SetLp(v antlr.Token) { s.lp = v }

func (s *AtRespDocContext) SetRp(v antlr.Token) { s.rp = v }

func (s *AtRespDocContext) GetCode() IRespDocCodeContext { return s.code }

func (s *AtRespDocContext) SetCode(v IRespDocCodeContext) { s.code = v }

func (s *AtRespDocContext) ATRESPDOC() antlr.TerminalNode {
	return s.GetToken(ApiParserParserATRESPDOC, 0)
}

func (s *AtRespDocContext) RespDocCode() IRespDocCodeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRespDocCodeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRespDocCodeContext)
}

func (s *AtRespDocContext) DataType() IDataTypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDataTypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *AtRespDocContext) AllRespDocKvLit() []IRespDocKvLitContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRespDocKvLitContext)(nil)).Elem())
	var tst = make([]IRespDocKvLitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRespDocKvLitContext)
		}
	}

	return tst
}

func (s *AtRespDocContext) RespDocKvLit(i int) IRespDocKvLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRespDocKvLitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRespDocKvLitContext)
}

func (s *AtRespDocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtRespDocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AtRespDocContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitAtRespDoc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) AtRespDoc() (localctx IAtRespDocContext) {
	localctx = NewAtRespDocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, ApiParserParserRULE_atRespDoc)
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
		p.SetState(304)
		p.Match(ApiParserParserATRESPDOC)
	}
	{
		p.SetState(305)
		p.Match(ApiParserParserT__9)
	}
	{
		p.SetState(306)

		var _x = p.RespDocCode()

		localctx.(*AtRespDocContext).code = _x
	}
	{
		p.SetState(307)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*AtRespDocContext).lp = _m
	}
	p.SetState(314)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.SetState(309)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == ApiParserParserID || _la == ApiParserParserLetterOrDigit {
			{
				p.SetState(308)
				p.RespDocKvLit()
			}

			p.SetState(311)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		{
			p.SetState(313)
			p.DataType()
		}

	}
	{
		p.SetState(316)

		var _m = p.Match(ApiParserParserT__2)

		localctx.(*AtRespDocContext).rp = _m
	}

	return localctx
}

// IRespDocCodeContext is an interface to support dynamic dispatch.
type IRespDocCodeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRespDocCodeContext differentiates from other interfaces.
	IsRespDocCodeContext()
}

type RespDocCodeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRespDocCodeContext() *RespDocCodeContext {
	var p = new(RespDocCodeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_respDocCode
	return p
}

func (*RespDocCodeContext) IsRespDocCodeContext() {}

func NewRespDocCodeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RespDocCodeContext {
	var p = new(RespDocCodeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_respDocCode

	return p
}

func (s *RespDocCodeContext) GetParser() antlr.Parser { return s.parser }

func (s *RespDocCodeContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ApiParserParserID)
}

func (s *RespDocCodeContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, i)
}

func (s *RespDocCodeContext) AllLetterOrDigit() []antlr.TerminalNode {
	return s.GetTokens(ApiParserParserLetterOrDigit)
}

func (s *RespDocCodeContext) LetterOrDigit(i int) antlr.TerminalNode {
	return s.GetToken(ApiParserParserLetterOrDigit, i)
}

func (s *RespDocCodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RespDocCodeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RespDocCodeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitRespDocCode(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) RespDocCode() (localctx IRespDocCodeContext) {
	localctx = NewRespDocCodeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, ApiParserParserRULE_respDocCode)
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
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ApiParserParserID || _la == ApiParserParserLetterOrDigit {
		{
			p.SetState(318)
			_la = p.GetTokenStream().LA(1)

			if !(_la == ApiParserParserID || _la == ApiParserParserLetterOrDigit) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(321)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 66, ApiParserParserRULE_route)
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
		p.SetState(324)

		var _m = p.Match(ApiParserParserID)

		localctx.(*RouteContext).httpMethod = _m
	}
	{
		p.SetState(325)
		p.Path()
	}
	p.SetState(327)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__1 {
		{
			p.SetState(326)

			var _x = p.Body()

			localctx.(*RouteContext).request = _x
		}

	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserT__10 {
		{
			p.SetState(329)

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

func (s *BodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitBody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Body() (localctx IBodyContext) {
	localctx = NewBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, ApiParserParserRULE_body)
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
		p.SetState(332)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*BodyContext).lp = _m
	}
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == ApiParserParserID {
		{
			p.SetState(333)
			p.Match(ApiParserParserID)
		}

	}
	{
		p.SetState(336)

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

func (s *ReplybodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitReplybody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) Replybody() (localctx IReplybodyContext) {
	localctx = NewReplybodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, ApiParserParserRULE_replybody)
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
		p.SetState(338)

		var _m = p.Match(ApiParserParserT__10)

		localctx.(*ReplybodyContext).returnToken = _m
	}
	{
		p.SetState(339)

		var _m = p.Match(ApiParserParserT__1)

		localctx.(*ReplybodyContext).lp = _m
	}
	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<ApiParserParserT__5)|(1<<ApiParserParserT__6)|(1<<ApiParserParserT__7)|(1<<ApiParserParserINTERFACE)|(1<<ApiParserParserID))) != 0 {
		{
			p.SetState(340)
			p.DataType()
		}

	}
	{
		p.SetState(343)

		var _m = p.Match(ApiParserParserT__2)

		localctx.(*ReplybodyContext).rp = _m
	}

	return localctx
}
