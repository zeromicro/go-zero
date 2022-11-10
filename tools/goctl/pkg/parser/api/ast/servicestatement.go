package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type AtServerStmt struct {
	AtServer token.Token
	LParen   token.Token
	Values   []*KVExpr
	RParen   token.Token
}

func (a *AtServerStmt) Pos() token.Position {
	return a.AtServer.Position
}

func (a *AtServerStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, a.AtServer, a.LParen)
	if len(a.Values) > 0 {
		w.NewLine()
	}
	tw := w.UseTabWriter()
	for _, v := range a.Values {
		tw.WriteWithInfixIndentln("", v.Format(prefix+indent))
	}
	tw.Flush()
	if len(a.Values) > 0 {
		w.Write(prefix, a.RParen)
	} else {
		w.Write(a.RParen)
	}
	return w.String()
}

func (a *AtServerStmt) stmtNode() {}

type AtDocStmt interface {
	Formatter
	Stmt
	atDocNode()
}

type AtDocLiteralStmt struct {
	AtDoc token.Token
	Value token.Token
}

func (a *AtDocLiteralStmt) atDocNode() {}

func (a *AtDocLiteralStmt) Pos() token.Position {
	return a.AtDoc.Position
}

func (a *AtDocLiteralStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, a.AtDoc, a.Value)
	return w.String()
}

func (a *AtDocLiteralStmt) stmtNode() {}

type AtDocGroupStmt struct {
	AtDoc  token.Token
	LParen token.Token
	Values []*KVExpr
	RParen token.Token
}

func (a *AtDocGroupStmt) atDocNode() {}

func (a *AtDocGroupStmt) Pos() token.Position {
	return a.AtDoc.Position
}

func (a *AtDocGroupStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, a.AtDoc, a.LParen)
	if len(a.Values) > 0 {
		w.NewLine()
	}
	tw := w.UseTabWriter()
	for _, val := range a.Values {
		tw.WriteWithInfixIndentln(prefix, val.Format(prefix))
	}
	tw.Flush()
	if len(a.Values) > 0 {
		w.Write(prefix, a.RParen)
	} else {
		w.Write(a.RParen)
	}
	return w.String()
}

func (a *AtDocGroupStmt) stmtNode() {}

type ServiceStmt struct {
	AtServerStmt *AtServerStmt
	Service      token.Token
	Name         *ServiceNameExpr
	LBrace       token.Token
	Routes       []*ServiceItemStmt
	RBrace       token.Token
}

func (s *ServiceStmt) Pos() token.Position {
	if s.AtServerStmt != nil {
		return s.AtServerStmt.Pos()
	}
	return s.Service.Position
}

func (s *ServiceStmt) Format(prefix string) string {
	w := NewWriter()
	if s.AtServerStmt != nil {
		w.Writeln(s.AtServerStmt.Format(noLead))
	}
	w.WriteWithWhiteSpaceInfix(prefix, s.Service, s.Name.Format(noLead), s.LBrace)
	if len(s.Routes) > 0 {
		w.NewLine()
	}
	for _, v := range s.Routes {
		w.Writeln(v.Format(prefix + indent))
	}
	if len(s.Routes) > 0 {
		w.Write(prefix, s.RBrace)
	} else {
		w.Write(s.RBrace)
	}

	return w.String()
}

func (s *ServiceStmt) stmtNode() {}

type ServiceNameExpr struct {
	ID     token.Token
	Joiner token.Token // optional
	API    token.Token // optional
}

func (s *ServiceNameExpr) Pos() token.Position {
	return s.ID.Position
}

func (s *ServiceNameExpr) Format(prefix string) string {
	w := NewWriter()
	w.Write(prefix, s.ID, s.Joiner, s.API)
	return w.String()
}

func (s *ServiceNameExpr) exprNode() {}

type AtHandlerStmt struct {
	AtHandler token.Token
	Name      token.Token
}

func (a *AtHandlerStmt) Pos() token.Position {
	return a.AtHandler.Position
}

func (a *AtHandlerStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, a.AtHandler, a.Name)
	return w.String()
}

func (a *AtHandlerStmt) stmtNode() {}

type ServiceItemStmt struct {
	AtDoc     AtDocStmt
	AtHandler *AtHandlerStmt
	Route     *RouteStmt
}

func (s *ServiceItemStmt) Pos() token.Position {
	if s.AtDoc != nil {
		return s.AtDoc.Pos()
	}
	return s.AtHandler.Pos()
}

func (s *ServiceItemStmt) Format(prefix string) string {
	w := NewWriter()
	if !isNil(s.AtDoc) {
		w.Writeln(s.AtDoc.Format(prefix))
	}
	if s.AtHandler != nil {
		w.Writeln(s.AtHandler.Format(prefix))
	}
	if s.Route != nil {
		w.Write(s.Route.Format(prefix))
	}
	return w.String()
}

func (s *ServiceItemStmt) stmtNode() {}

type RouteStmt struct {
	Method   token.Token
	Path     *PathExpr
	Request  *BodyStmt
	Returns  token.Token
	Response *BodyStmt
}

func (r *RouteStmt) Pos() token.Position {
	return r.Method.Position
}

func (r *RouteStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, r.Method, r.Path.Format(noLead))
	if r.Request != nil {
		w.Write(whiteSpace, r.Request.Format(noLead))
	}
	if r.Returns.Valid() {
		w.Write(whiteSpace, r.Returns)
	}
	if r.Response != nil {
		w.Write(whiteSpace, r.Response.Format(noLead))
	}
	return w.String()
}

func (r *RouteStmt) stmtNode() {}

type PathExpr struct {
	Values []token.Token
}

func (p *PathExpr) Pos() token.Position {
	if len(p.Values) == 0 {
		return token.Position{}
	}
	return p.Values[0].Position
}

func (p *PathExpr) Format(prefix string) string {
	w := NewWriter()
	w.Write(prefix)
	for _, v := range p.Values {
		w.Write(v)
	}
	return w.String()
}

func (p *PathExpr) exprNode() {}

type BodyStmt struct {
	LParen token.Token
	Body   *BodyExpr
	RParen token.Token
}

func (b *BodyStmt) Pos() token.Position {
	return b.LParen.Position
}

func (b *BodyStmt) Format(_ string) string {
	w := NewWriter()
	w.Write(b.LParen)
	if b.Body != nil {
		w.Write(b.Body.Format(noLead))
	}
	w.Write(b.RParen)
	return w.String()
}

func (b *BodyStmt) stmtNode() {}

type BodyExpr struct {
	LBrack token.Token
	RBrack token.Token
	Star   token.Token
	Value  token.Token
}

func (e *BodyExpr) Format(prefix string) string {
	w := NewWriter()
	w.Write(prefix, e.LBrack, e.RBrack, e.Star, e.Value)
	return w.String()
}

func (e *BodyExpr) Pos() token.Position {
	if e.LBrack.Valid() {
		return e.LBrack.Position
	}
	if e.Star.Valid() {
		return e.Star.Position
	}
	return e.Value.Position
}

func (e *BodyExpr) exprNode() {}
