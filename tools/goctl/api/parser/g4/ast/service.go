package ast

import (
	"fmt"
	"sort"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
)

// Service describes service for api syntax
type Service struct {
	AtServer   *AtServer
	ServiceApi *ServiceApi
}

// KV defines a slice for KvExpr
type KV []*KvExpr

// AtServer describes server metadata for api syntax
type AtServer struct {
	AtServerToken Expr
	Lp            Expr
	Rp            Expr
	Kv            KV
}

// ServiceApi describes service ast for api syntax
type ServiceApi struct {
	ServiceToken Expr
	Name         Expr
	Lbrace       Expr
	Rbrace       Expr
	ServiceRoute []*ServiceRoute
}

// ServiceRoute describes service route ast for api syntax
type ServiceRoute struct {
	AtDoc     *AtDoc
	AtServer  *AtServer
	AtHandler *AtHandler
	Route     *Route
}

// AtDoc describes service comments ast for api syntax
type AtDoc struct {
	AtDocToken Expr
	Lp         Expr
	Rp         Expr
	LineDoc    Expr
	Kv         []*KvExpr
}

// AtHandler describes service handler ast for api syntax
type AtHandler struct {
	AtHandlerToken Expr
	Name           Expr
	DocExpr        []Expr
	CommentExpr    Expr
}

// Route describes route ast for api syntax
type Route struct {
	Method      Expr
	Path        Expr
	Req         *Body
	ReturnToken Expr
	Reply       *Body
	DocExpr     []Expr
	CommentExpr Expr
}

// Body describes request,response body ast for api syntax
type Body struct {
	ReturnExpr Expr
	Lp         Expr
	Rp         Expr
	Name       DataType
}

// VisitServiceSpec implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitServiceSpec(ctx *api.ServiceSpecContext) any {
	var serviceSpec Service
	if ctx.AtServer() != nil {
		serviceSpec.AtServer = ctx.AtServer().Accept(v).(*AtServer)
	}

	serviceSpec.ServiceApi = ctx.ServiceApi().Accept(v).(*ServiceApi)
	return &serviceSpec
}

// VisitAtServer implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitAtServer(ctx *api.AtServerContext) any {
	var atServer AtServer
	atServer.AtServerToken = v.newExprWithTerminalNode(ctx.ATSERVER())
	atServer.Lp = v.newExprWithToken(ctx.GetLp())
	atServer.Rp = v.newExprWithToken(ctx.GetRp())

	for _, each := range ctx.AllKvLit() {
		atServer.Kv = append(atServer.Kv, each.Accept(v).(*KvExpr))
	}

	return &atServer
}

// VisitServiceApi implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitServiceApi(ctx *api.ServiceApiContext) any {
	var serviceApi ServiceApi
	serviceApi.ServiceToken = v.newExprWithToken(ctx.GetServiceToken())
	serviceName := ctx.ServiceName()
	serviceApi.Name = v.newExprWithText(serviceName.GetText(), serviceName.GetStart().GetLine(), serviceName.GetStart().GetColumn(), serviceName.GetStart().GetStart(), serviceName.GetStop().GetStop())
	serviceApi.Lbrace = v.newExprWithToken(ctx.GetLbrace())
	serviceApi.Rbrace = v.newExprWithToken(ctx.GetRbrace())

	for _, each := range ctx.AllServiceRoute() {
		serviceApi.ServiceRoute = append(serviceApi.ServiceRoute, each.Accept(v).(*ServiceRoute))
	}

	return &serviceApi
}

// VisitServiceRoute implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitServiceRoute(ctx *api.ServiceRouteContext) any {
	var serviceRoute ServiceRoute
	if ctx.AtDoc() != nil {
		serviceRoute.AtDoc = ctx.AtDoc().Accept(v).(*AtDoc)
	}

	if ctx.AtServer() != nil {
		serviceRoute.AtServer = ctx.AtServer().Accept(v).(*AtServer)
	} else if ctx.AtHandler() != nil {
		serviceRoute.AtHandler = ctx.AtHandler().Accept(v).(*AtHandler)
	}

	serviceRoute.Route = ctx.Route().Accept(v).(*Route)
	return &serviceRoute
}

// VisitAtDoc implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitAtDoc(ctx *api.AtDocContext) any {
	var atDoc AtDoc
	atDoc.AtDocToken = v.newExprWithTerminalNode(ctx.ATDOC())

	if ctx.STRING() != nil {
		atDoc.LineDoc = v.newExprWithTerminalNode(ctx.STRING())
	} else {
		for _, each := range ctx.AllKvLit() {
			atDoc.Kv = append(atDoc.Kv, each.Accept(v).(*KvExpr))
		}
	}
	atDoc.Lp = v.newExprWithToken(ctx.GetLp())
	atDoc.Rp = v.newExprWithToken(ctx.GetRp())

	if ctx.GetLp() != nil {
		if ctx.GetRp() == nil {
			v.panic(atDoc.Lp, "mismatched ')'")
		}
	}

	if ctx.GetRp() != nil {
		if ctx.GetLp() == nil {
			v.panic(atDoc.Rp, "mismatched '('")
		}
	}

	return &atDoc
}

// VisitAtHandler implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitAtHandler(ctx *api.AtHandlerContext) any {
	var atHandler AtHandler
	astHandlerExpr := v.newExprWithTerminalNode(ctx.ATHANDLER())
	atHandler.AtHandlerToken = astHandlerExpr
	atHandler.Name = v.newExprWithTerminalNode(ctx.ID())
	atHandler.DocExpr = v.getDoc(ctx)
	atHandler.CommentExpr = v.getComment(ctx)
	return &atHandler
}

// VisitRoute implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitRoute(ctx *api.RouteContext) any {
	var route Route
	path := ctx.Path()
	methodExpr := v.newExprWithToken(ctx.GetHttpMethod())
	route.Method = methodExpr
	route.Path = v.newExprWithText(path.GetText(), path.GetStart().GetLine(), path.GetStart().GetColumn(), path.GetStart().GetStart(), path.GetStop().GetStop())

	if ctx.GetRequest() != nil {
		req := ctx.GetRequest().Accept(v)
		if req != nil {
			route.Req = req.(*Body)
		}
	}

	if ctx.GetResponse() != nil {
		reply := ctx.GetResponse().Accept(v)
		if reply != nil {
			resp := reply.(*Body)
			route.ReturnToken = resp.ReturnExpr
			resp.ReturnExpr = nil
			route.Reply = resp
		}
	}

	route.DocExpr = v.getDoc(ctx)
	route.CommentExpr = v.getComment(ctx)
	return &route
}

// VisitBody implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitBody(ctx *api.BodyContext) any {
	if ctx.ID() == nil {
		if v.debug {
			msg := fmt.Sprintf(
				`%s line %d:  expr "()" is deprecated, if there has no request body, please omit it`,
				v.prefix,
				ctx.GetStart().GetLine(),
			)
			v.log.Warning(msg)
		}
		return nil
	}

	idExpr := v.newExprWithTerminalNode(ctx.ID())
	if api.IsGolangKeyWord(idExpr.Text()) {
		v.panic(idExpr, fmt.Sprintf("expecting 'ID', but found golang keyword '%s'", idExpr.Text()))
	}

	return &Body{
		Lp:   v.newExprWithToken(ctx.GetLp()),
		Rp:   v.newExprWithToken(ctx.GetRp()),
		Name: &Literal{Literal: idExpr},
	}
}

// VisitReplybody implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitReplybody(ctx *api.ReplybodyContext) any {
	if ctx.DataType() == nil {
		if v.debug {
			msg := fmt.Sprintf(
				`%s line %d:  expr "returns ()" or "()" is deprecated, if there has no response body, please omit it`,
				v.prefix,
				ctx.GetStart().GetLine(),
			)
			v.log.Warning(msg)
		}
		return nil
	}

	var returnExpr Expr
	if ctx.GetReturnToken() != nil {
		returnExpr = v.newExprWithToken(ctx.GetReturnToken())
		if ctx.GetReturnToken().GetText() != "returns" {
			v.panic(returnExpr, fmt.Sprintf("expecting returns, found input '%s'", ctx.GetReturnToken().GetText()))
		}
	}

	dt := ctx.DataType().Accept(v).(DataType)
	if dt == nil {
		return nil
	}

	switch dataType := dt.(type) {
	case *Array:
		lit := dataType.Literal
		switch lit.(type) {
		case *Literal, *Pointer:
			if api.IsGolangKeyWord(lit.Expr().Text()) {
				v.panic(lit.Expr(), fmt.Sprintf("expecting 'ID', but found golang keyword '%s'", lit.Expr().Text()))
			}
		default:
			v.panic(dt.Expr(), fmt.Sprintf("unsupport %s", dt.Expr().Text()))
		}
	case *Literal:
		lit := dataType.Literal.Text()
		if api.IsGolangKeyWord(lit) {
			v.panic(dataType.Literal, fmt.Sprintf("expecting 'ID', but found golang keyword '%s'", lit))
		}
	default:
		v.panic(dt.Expr(), fmt.Sprintf("unsupport %s", dt.Expr().Text()))
	}

	return &Body{
		ReturnExpr: returnExpr,
		Lp:         v.newExprWithToken(ctx.GetLp()),
		Rp:         v.newExprWithToken(ctx.GetRp()),
		Name:       dt,
	}
}

// Format provides a formatter for api command, now nothing to do
func (b *Body) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two Body are equal
func (b *Body) Equal(v any) bool {
	if v == nil {
		return false
	}

	body, ok := v.(*Body)
	if !ok {
		return false
	}

	if !b.Lp.Equal(body.Lp) {
		return false
	}

	if !b.Rp.Equal(body.Rp) {
		return false
	}

	return b.Name.Equal(body.Name)
}

// Format provides a formatter for api command, now nothing to do
func (r *Route) Format() error {
	// todo
	return nil
}

// Doc returns the document of Route, like // some text
func (r *Route) Doc() []Expr {
	return r.DocExpr
}

// Comment returns the comment of Route, like // some text
func (r *Route) Comment() Expr {
	return r.CommentExpr
}

// Equal compares whether the element literals in two Route are equal
func (r *Route) Equal(v any) bool {
	if v == nil {
		return false
	}

	route, ok := v.(*Route)
	if !ok {
		return false
	}

	if !r.Method.Equal(route.Method) {
		return false
	}

	if !r.Path.Equal(route.Path) {
		return false
	}

	if r.Req != nil {
		if !r.Req.Equal(route.Req) {
			return false
		}
	}

	if r.ReturnToken != nil {
		if !r.ReturnToken.Equal(route.ReturnToken) {
			return false
		}
	}

	if r.Reply != nil {
		if !r.Reply.Equal(route.Reply) {
			return false
		}
	}

	return EqualDoc(r, route)
}

// Doc returns the document of AtHandler, like // some text
func (a *AtHandler) Doc() []Expr {
	return a.DocExpr
}

// Comment returns the comment of AtHandler, like // some text
func (a *AtHandler) Comment() Expr {
	return a.CommentExpr
}

// Format provides a formatter for api command, now nothing to do
func (a *AtHandler) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two AtHandler are equal
func (a *AtHandler) Equal(v any) bool {
	if v == nil {
		return false
	}

	atHandler, ok := v.(*AtHandler)
	if !ok {
		return false
	}

	if !a.AtHandlerToken.Equal(atHandler.AtHandlerToken) {
		return false
	}

	if !a.Name.Equal(atHandler.Name) {
		return false
	}

	return EqualDoc(a, atHandler)
}

// Format provides a formatter for api command, now nothing to do
func (a *AtDoc) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two AtDoc are equal
func (a *AtDoc) Equal(v any) bool {
	if v == nil {
		return false
	}

	atDoc, ok := v.(*AtDoc)
	if !ok {
		return false
	}

	if !a.AtDocToken.Equal(atDoc.AtDocToken) {
		return false
	}

	if a.Lp.IsNotNil() {
		if !a.Lp.Equal(atDoc.Lp) {
			return false
		}
	}

	if a.Rp.IsNotNil() {
		if !a.Rp.Equal(atDoc.Rp) {
			return false
		}
	}

	if a.LineDoc != nil {
		if !a.LineDoc.Equal(atDoc.LineDoc) {
			return false
		}
	}

	var expecting, actual []*KvExpr
	expecting = append(expecting, a.Kv...)
	actual = append(actual, atDoc.Kv...)

	if len(expecting) != len(actual) {
		return false
	}

	for index, each := range expecting {
		ac := actual[index]
		if !each.Equal(ac) {
			return false
		}
	}

	return true
}

// Format provides a formatter for api command, now nothing to do
func (a *AtServer) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two AtServer are equal
func (a *AtServer) Equal(v any) bool {
	if v == nil {
		return false
	}

	atServer, ok := v.(*AtServer)
	if !ok {
		return false
	}

	if !a.AtServerToken.Equal(atServer.AtServerToken) {
		return false
	}

	if !a.Lp.Equal(atServer.Lp) {
		return false
	}

	if !a.Rp.Equal(atServer.Rp) {
		return false
	}

	var expecting, actual []*KvExpr
	expecting = append(expecting, a.Kv...)
	actual = append(actual, atServer.Kv...)
	if len(expecting) != len(actual) {
		return false
	}

	sort.Slice(expecting, func(i, j int) bool {
		return expecting[i].Key.Text() < expecting[j].Key.Text()
	})

	sort.Slice(actual, func(i, j int) bool {
		return actual[i].Key.Text() < actual[j].Key.Text()
	})

	for index, each := range expecting {
		ac := actual[index]
		if !each.Equal(ac) {
			return false
		}
	}

	return true
}

// Equal compares whether the element literals in two ServiceRoute are equal
func (s *ServiceRoute) Equal(v any) bool {
	if v == nil {
		return false
	}

	sr, ok := v.(*ServiceRoute)
	if !ok {
		return false
	}

	if !s.AtDoc.Equal(sr.AtDoc) {
		return false
	}

	if s.AtServer != nil {
		if !s.AtServer.Equal(sr.AtServer) {
			return false
		}
	}

	if s.AtHandler != nil {
		if !s.AtHandler.Equal(sr.AtHandler) {
			return false
		}
	}

	return s.Route.Equal(sr.Route)
}

// Format provides a formatter for api command, now nothing to do
func (s *ServiceRoute) Format() error {
	// todo
	return nil
}

// GetHandler returns handler name of api route
func (s *ServiceRoute) GetHandler() Expr {
	if s.AtHandler != nil {
		return s.AtHandler.Name
	}

	return s.AtServer.Kv.Get("handler")
}

// Format provides a formatter for api command, now nothing to do
func (a *ServiceApi) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two ServiceApi are equal
func (a *ServiceApi) Equal(v any) bool {
	if v == nil {
		return false
	}

	api, ok := v.(*ServiceApi)
	if !ok {
		return false
	}

	if !a.ServiceToken.Equal(api.ServiceToken) {
		return false
	}

	if !a.Name.Equal(api.Name) {
		return false
	}

	if !a.Lbrace.Equal(api.Lbrace) {
		return false
	}

	if !a.Rbrace.Equal(api.Rbrace) {
		return false
	}

	var expecting, acutal []*ServiceRoute
	expecting = append(expecting, a.ServiceRoute...)
	acutal = append(acutal, api.ServiceRoute...)
	if len(expecting) != len(acutal) {
		return false
	}

	sort.Slice(expecting, func(i, j int) bool {
		return expecting[i].Route.Path.Text() < expecting[j].Route.Path.Text()
	})

	sort.Slice(acutal, func(i, j int) bool {
		return acutal[i].Route.Path.Text() < acutal[j].Route.Path.Text()
	})

	for index, each := range expecting {
		ac := acutal[index]
		if !each.Equal(ac) {
			return false
		}
	}

	return true
}

// Format provides a formatter for api command, now nothing to do
func (s *Service) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two Service are equal
func (s *Service) Equal(v any) bool {
	if v == nil {
		return false
	}

	service, ok := v.(*Service)
	if !ok {
		return false
	}

	if s.AtServer != nil {
		if !s.AtServer.Equal(service.AtServer) {
			return false
		}
	}

	return s.ServiceApi.Equal(service.ServiceApi)
}

// Get returns the target KV by specified key
func (kv KV) Get(key string) Expr {
	for _, each := range kv {
		if each.Key.Text() == key {
			return each.Value
		}
	}
	return nil
}
