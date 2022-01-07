package api

import (
	"reflect"

	"github.com/zeromicro/antlr"
)

// Part 1
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.

func (s *SyntaxLitContext) GetParser() antlr.Parser { return s.parser }

func (s *SyntaxLitContext) GetSyntaxToken() antlr.Token { return s.syntaxToken }

func (s *SyntaxLitContext) GetAssign() antlr.Token { return s.assign }

func (s *SyntaxLitContext) GetVersion() antlr.Token { return s.version }

func (s *SyntaxLitContext) SetSyntaxToken(v antlr.Token) { s.syntaxToken = v }

func (s *SyntaxLitContext) SetAssign(v antlr.Token) { s.assign = v }

func (s *SyntaxLitContext) SetVersion(v antlr.Token) { s.version = v }

func (s *SyntaxLitContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *SyntaxLitContext) STRING() antlr.TerminalNode {
	return s.GetToken(ApiParserParserSTRING, 0)
}

func (s *SyntaxLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SyntaxLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SyntaxLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitSyntaxLit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) SyntaxLit() (localctx ISyntaxLitContext) {
	localctx = NewSyntaxLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, ApiParserParserRULE_syntaxLit)

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
	match(p, "syntax")
	{
		p.SetState(88)

		_m := p.Match(ApiParserParserID)

		localctx.(*SyntaxLitContext).syntaxToken = _m
	}
	{
		p.SetState(89)

		_m := p.Match(ApiParserParserT__0)

		localctx.(*SyntaxLitContext).assign = _m
	}
	checkVersion(p)
	{
		p.SetState(91)

		_m := p.Match(ApiParserParserSTRING)

		localctx.(*SyntaxLitContext).version = _m
	}

	return localctx
}

// IImportSpecContext is an interface to support dynamic dispatch.
type IImportSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportSpecContext differentiates from other interfaces.
	IsImportSpecContext()
}

type ImportSpecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportSpecContext() *ImportSpecContext {
	p := new(ImportSpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_importSpec
	return p
}

func (*ImportSpecContext) IsImportSpecContext() {}

func NewImportSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportSpecContext {
	p := new(ImportSpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_importSpec

	return p
}

func (s *ImportSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportSpecContext) ImportLit() IImportLitContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportLitContext)
}

func (s *ImportSpecContext) ImportBlock() IImportBlockContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportBlockContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportBlockContext)
}

func (s *ImportSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportSpecContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitImportSpec(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ImportSpec() (localctx IImportSpecContext) {
	localctx = NewImportSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ApiParserParserRULE_importSpec)

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

	p.SetState(95)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(93)
			p.ImportLit()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(94)
			p.ImportBlock()
		}

	}

	return localctx
}

// IImportLitContext is an interface to support dynamic dispatch.
type IImportLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetImportToken returns the importToken token.
	GetImportToken() antlr.Token

	// SetImportToken sets the importToken token.
	SetImportToken(antlr.Token)

	// IsImportLitContext differentiates from other interfaces.
	IsImportLitContext()
}

type ImportLitContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	importToken antlr.Token
}

func NewEmptyImportLitContext() *ImportLitContext {
	p := new(ImportLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_importLit
	return p
}

func (*ImportLitContext) IsImportLitContext() {}

func NewImportLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportLitContext {
	p := new(ImportLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_importLit

	return p
}

func (s *ImportLitContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportLitContext) GetImportToken() antlr.Token { return s.importToken }

func (s *ImportLitContext) SetImportToken(v antlr.Token) { s.importToken = v }

func (s *ImportLitContext) ImportValue() IImportValueContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportValueContext)
}

func (s *ImportLitContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *ImportLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitImportLit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ImportLit() (localctx IImportLitContext) {
	localctx = NewImportLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ApiParserParserRULE_importLit)

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
	match(p, "import")
	{
		p.SetState(98)

		_m := p.Match(ApiParserParserID)

		localctx.(*ImportLitContext).importToken = _m
	}
	{
		p.SetState(99)
		p.ImportValue()
	}

	return localctx
}

// IImportBlockContext is an interface to support dynamic dispatch.
type IImportBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetImportToken returns the importToken token.
	GetImportToken() antlr.Token

	// SetImportToken sets the importToken token.
	SetImportToken(antlr.Token)

	// IsImportBlockContext differentiates from other interfaces.
	IsImportBlockContext()
}

type ImportBlockContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	importToken antlr.Token
}

func NewEmptyImportBlockContext() *ImportBlockContext {
	p := new(ImportBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_importBlock
	return p
}

func (*ImportBlockContext) IsImportBlockContext() {}

func NewImportBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportBlockContext {
	p := new(ImportBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_importBlock

	return p
}

func (s *ImportBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportBlockContext) GetImportToken() antlr.Token { return s.importToken }

func (s *ImportBlockContext) SetImportToken(v antlr.Token) { s.importToken = v }

func (s *ImportBlockContext) ID() antlr.TerminalNode {
	return s.GetToken(ApiParserParserID, 0)
}

func (s *ImportBlockContext) AllImportBlockValue() []IImportBlockValueContext {
	ts := s.GetTypedRuleContexts(reflect.TypeOf((*IImportBlockValueContext)(nil)).Elem())
	tst := make([]IImportBlockValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IImportBlockValueContext)
		}
	}

	return tst
}

func (s *ImportBlockContext) ImportBlockValue(i int) IImportBlockValueContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportBlockValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IImportBlockValueContext)
}

func (s *ImportBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitImportBlock(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ImportBlock() (localctx IImportBlockContext) {
	localctx = NewImportBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, ApiParserParserRULE_importBlock)
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
	match(p, "import")
	{
		p.SetState(102)

		_m := p.Match(ApiParserParserID)

		localctx.(*ImportBlockContext).importToken = _m
	}
	{
		p.SetState(103)
		p.Match(ApiParserParserT__1)
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)

	for ok := true; ok; ok = _la == ApiParserParserSTRING {
		{
			p.SetState(104)
			p.ImportBlockValue()
		}

		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(109)
		p.Match(ApiParserParserT__2)
	}

	return localctx
}

// IImportBlockValueContext is an interface to support dynamic dispatch.
type IImportBlockValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportBlockValueContext differentiates from other interfaces.
	IsImportBlockValueContext()
}

type ImportBlockValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportBlockValueContext() *ImportBlockValueContext {
	p := new(ImportBlockValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_importBlockValue
	return p
}

func (*ImportBlockValueContext) IsImportBlockValueContext() {}

func NewImportBlockValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportBlockValueContext {
	p := new(ImportBlockValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_importBlockValue

	return p
}

func (s *ImportBlockValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportBlockValueContext) ImportValue() IImportValueContext {
	t := s.GetTypedRuleContext(reflect.TypeOf((*IImportValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportValueContext)
}

func (s *ImportBlockValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportBlockValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportBlockValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case ApiParserVisitor:
		return t.VisitImportBlockValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *ApiParserParser) ImportBlockValue() (localctx IImportBlockValueContext) {
	localctx = NewImportBlockValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, ApiParserParserRULE_importBlockValue)

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
		p.SetState(111)
		p.ImportValue()
	}

	return localctx
}

// IImportValueContext is an interface to support dynamic dispatch.
type IImportValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportValueContext differentiates from other interfaces.
	IsImportValueContext()
}

type ImportValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportValueContext() *ImportValueContext {
	p := new(ImportValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = ApiParserParserRULE_importValue
	return p
}

func (*ImportValueContext) IsImportValueContext() {}

func NewImportValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportValueContext {
	p := new(ImportValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = ApiParserParserRULE_importValue

	return p
}
