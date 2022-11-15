package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type AtServerStmt struct {
	AtServer token.Token
	LParen   token.Token
	Values   []*KVExpr
	RParen   token.Token

	fw *Writer
}

func (a *AtServerStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (a *AtServerStmt) End() token.Position {
	return a.RParen.Position
}

func (a *AtServerStmt) Pos() token.Position {
	return a.AtServer.Position
}

func (a *AtServerStmt) stmtNode() {}

type AtDocStmt interface {
	Stmt
	atDocNode()
}

type AtDocLiteralStmt struct {
	AtDoc token.Token
	Value token.Token

	fw *Writer
}

func (a *AtDocLiteralStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (a *AtDocLiteralStmt) End() token.Position {
	return a.Value.Position
}

func (a *AtDocLiteralStmt) atDocNode() {}

func (a *AtDocLiteralStmt) Pos() token.Position {
	return a.AtDoc.Position
}

func (a *AtDocLiteralStmt) stmtNode() {}

type AtDocGroupStmt struct {
	AtDoc  token.Token
	LParen token.Token
	Values []*KVExpr
	RParen token.Token

	fw *Writer
}

func (a *AtDocGroupStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (a *AtDocGroupStmt) End() token.Position {
	return a.RParen.Position
}

func (a *AtDocGroupStmt) atDocNode() {}

func (a *AtDocGroupStmt) Pos() token.Position {
	return a.AtDoc.Position
}

func (a *AtDocGroupStmt) stmtNode() {}

type ServiceStmt struct {
	AtServerStmt *AtServerStmt
	Service      token.Token
	Name         *ServiceNameExpr
	LBrace       token.Token
	Routes       []*ServiceItemStmt
	RBrace       token.Token

	fw *Writer
}

func (s *ServiceStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceStmt) End() token.Position {
	return s.RBrace.Position
}

func (s *ServiceStmt) Pos() token.Position {
	if s.AtServerStmt != nil {
		return s.AtServerStmt.Pos()
	}
	return s.Service.Position
}

func (s *ServiceStmt) stmtNode() {}

type ServiceNameExpr struct {
	ID     token.Token
	Joiner token.Token // optional
	API    token.Token // optional

	fw *Writer
}

func (s *ServiceNameExpr) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceNameExpr) End() token.Position {
	if s.API.Valid() {
		return s.API.Position
	}
	if s.Joiner.Valid() {
		return s.Joiner.Position
	}
	return s.ID.Position
}

func (s *ServiceNameExpr) Pos() token.Position {
	return s.ID.Position
}

func (s *ServiceNameExpr) exprNode() {}

type AtHandlerStmt struct {
	AtHandler *TokenNode
	Name      *TokenNode

	fw *Writer
}

func (a *AtHandlerStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (a *AtHandlerStmt) End() token.Position {
	return a.Name.End()
}

func (a *AtHandlerStmt) Pos() token.Position {
	return a.AtHandler.Pos()
}

func (a *AtHandlerStmt) stmtNode() {}

type ServiceItemStmt struct {
	AtDoc     AtDocStmt
	AtHandler *AtHandlerStmt
	Route     *RouteStmt

	fw *Writer
}

func (s *ServiceItemStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceItemStmt) End() token.Position {
	return s.Route.End()
}

func (s *ServiceItemStmt) Pos() token.Position {
	if s.AtDoc != nil {
		return s.AtDoc.Pos()
	}
	return s.AtHandler.Pos()
}

func (s *ServiceItemStmt) stmtNode() {}

type RouteStmt struct {
	Method   token.Token
	Path     *PathExpr
	Request  *BodyStmt
	Returns  token.Token
	Response *BodyStmt

	fw *Writer
}

func (r *RouteStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (r *RouteStmt) End() token.Position {
	if r.Response != nil {
		return r.Response.End()
	}
	if r.Returns.Valid() {
		return r.Returns.Position
	}
	if r.Request != nil {
		return r.Request.End()
	}
	return r.Path.End()
}

func (r *RouteStmt) Pos() token.Position {
	return r.Method.Position
}

func (r *RouteStmt) stmtNode() {}

type PathExpr struct {
	Values []token.Token

	fw *Writer
}

func (p *PathExpr) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (p *PathExpr) End() token.Position {
	if len(p.Values) == 0 {
		return token.Position{}
	}
	return p.Values[len(p.Values)-1].Position
}

func (p *PathExpr) Pos() token.Position {
	if len(p.Values) == 0 {
		return token.Position{}
	}
	return p.Values[0].Position
}

func (p *PathExpr) exprNode() {}

type BodyStmt struct {
	LParen token.Token
	Body   *BodyExpr
	RParen token.Token

	fw *Writer
}

func (b *BodyStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (b *BodyStmt) End() token.Position {
	return b.RParen.Position
}

func (b *BodyStmt) Pos() token.Position {
	return b.LParen.Position
}

func (b *BodyStmt) stmtNode() {}

type BodyExpr struct {
	LBrack token.Token
	RBrack token.Token
	Star   token.Token
	Value  token.Token

	fw *Writer
}

func (e *BodyExpr) End() token.Position {
	return e.Value.Position
}

func (e *BodyExpr) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
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
