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

func (a *AtServerStmt) stmtNode() {}

type AtDocStmt interface {
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

func (s *ServiceStmt) stmtNode() {}

type ServiceNameExpr struct {
	ID     token.Token
	Joiner token.Token // optional
	API    token.Token // optional
}

func (s *ServiceNameExpr) Pos() token.Position {
	return s.ID.Position
}

func (s *ServiceNameExpr) exprNode() {}

type AtHandlerStmt struct {
	AtHandler token.Token
	Name      token.Token
}

func (a *AtHandlerStmt) Pos() token.Position {
	return a.AtHandler.Position
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

func (p *PathExpr) exprNode() {}

type BodyStmt struct {
	LParen token.Token
	Body   *BodyExpr
	RParen token.Token
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
