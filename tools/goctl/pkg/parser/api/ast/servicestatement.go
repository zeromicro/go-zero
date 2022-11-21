package ast

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type AtServerStmt struct {
	AtServer *TokenNode
	LParen   *TokenNode
	Values   []*KVExpr
	RParen   *TokenNode
}

func (a *AtServerStmt) HasHeadCommentGroup() bool {
	return a.AtServer.HasHeadCommentGroup()
}

func (a *AtServerStmt) HasLeadingCommentGroup() bool {
	return a.RParen.HasLeadingCommentGroup()
}

func (a *AtServerStmt) CommentGroup() (head, leading CommentGroup) {
	return a.AtServer.HeadCommentGroup, a.RParen.LeadingCommentGroup
}

func (a *AtServerStmt) Format(prefix ...string) string {
	if len(a.Values) == 0 {
		return ""
	}
	var textList []string
	for _, v := range a.Values {
		if v.Value.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format())
	}
	if len(textList) == 0 {
		return ""
	}

	w := NewBufferWriter()
	w.Write(withNode(a.AtServer, a.LParen), withPrefix(prefix...), expectSameLine())
	w.NewLine()
	for _, v := range a.Values {
		node := transferTokenNode(v.Key, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment())
		w.Write(withNode(node, v.Value), expectIndentInfix(), expectSameLine())
		w.NewLine()
	}
	w.Write(withNode(transferTokenNode(a.RParen, withTokenNodePrefix(prefix...))))
	return w.String()
}

func (a *AtServerStmt) End() token.Position {
	return a.RParen.End()
}

func (a *AtServerStmt) Pos() token.Position {
	return a.AtServer.Pos()
}

func (a *AtServerStmt) stmtNode() {}

type AtDocStmt interface {
	Stmt
	atDocNode()
}

type AtDocLiteralStmt struct {
	AtDoc *TokenNode
	Value *TokenNode
}

func (a *AtDocLiteralStmt) HasHeadCommentGroup() bool {
	return a.AtDoc.HasHeadCommentGroup()
}

func (a *AtDocLiteralStmt) HasLeadingCommentGroup() bool {
	return a.Value.HasLeadingCommentGroup()
}

func (a *AtDocLiteralStmt) CommentGroup() (head, leading CommentGroup) {
	return a.AtDoc.HeadCommentGroup, a.Value.LeadingCommentGroup
}

func (a *AtDocLiteralStmt) Format(prefix ...string) string {
	if a.Value.IsZeroString() {
		return ""
	}
	w := NewBufferWriter()
	atDocNode := transferTokenNode(a.AtDoc, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	valueNode := transferTokenNode(a.Value, ignoreHeadComment())
	w.Write(withNode(atDocNode, valueNode), expectSameLine())
	return w.String()
}

func (a *AtDocLiteralStmt) End() token.Position {
	return a.Value.End()
}

func (a *AtDocLiteralStmt) atDocNode() {}

func (a *AtDocLiteralStmt) Pos() token.Position {
	return a.AtDoc.Pos()
}

func (a *AtDocLiteralStmt) stmtNode() {}

type AtDocGroupStmt struct {
	AtDoc  *TokenNode
	LParen *TokenNode
	Values []*KVExpr
	RParen *TokenNode
}

func (a *AtDocGroupStmt) HasHeadCommentGroup() bool {
	return a.AtDoc.HasHeadCommentGroup()
}

func (a *AtDocGroupStmt) HasLeadingCommentGroup() bool {
	return a.RParen.HasLeadingCommentGroup()
}

func (a *AtDocGroupStmt) CommentGroup() (head, leading CommentGroup) {
	return a.AtDoc.HeadCommentGroup, a.RParen.LeadingCommentGroup
}

func (a *AtDocGroupStmt) Format(prefix ...string) string {
	if len(a.Values) == 0 {
		return ""
	}
	var textList []string
	for _, v := range a.Values {
		if v.Value.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format(peekOne(prefix)+Indent))
	}
	if len(textList) == 0 {
		return ""
	}

	w := NewBufferWriter()
	w.WriteText(a.AtDoc.Format(prefix...) + WhiteSpace + a.LParen.Format())
	w.NewLine()
	w.WriteText(strings.Join(textList, NewLine))
	w.NewLine()
	w.WriteText(a.RParen.Format(prefix...))
	return w.String()
}

func (a *AtDocGroupStmt) End() token.Position {
	return a.RParen.End()
}

func (a *AtDocGroupStmt) atDocNode() {}

func (a *AtDocGroupStmt) Pos() token.Position {
	return a.AtDoc.Pos()
}

func (a *AtDocGroupStmt) stmtNode() {}

type ServiceStmt struct {
	AtServerStmt *AtServerStmt
	Service      *TokenNode
	Name         *ServiceNameExpr
	LBrace       *TokenNode
	Routes       []*ServiceItemStmt
	RBrace       *TokenNode
}

func (s *ServiceStmt) HasHeadCommentGroup() bool {
	if s.AtServerStmt != nil {
		return s.AtServerStmt.HasHeadCommentGroup()
	}
	return s.Service.HasHeadCommentGroup()
}

func (s *ServiceStmt) HasLeadingCommentGroup() bool {
	return s.RBrace.HasLeadingCommentGroup()
}

func (s *ServiceStmt) CommentGroup() (head, leading CommentGroup) {
	if s.AtServerStmt != nil {
		head, _ = s.AtServerStmt.CommentGroup()
		return head, s.RBrace.LeadingCommentGroup
	}
	return s.Service.HeadCommentGroup, s.RBrace.LeadingCommentGroup
}

func (s *ServiceStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	if s.AtServerStmt != nil {
		w.WriteText(s.AtServerStmt.Format(prefix...))
		w.NewLine()
	}
	w.Write(withNode(s.Service, s.Name, s.LBrace),
		withPrefix(prefix...), withMode(ModeExpectInSameLine))
	w.NewLine()
	var textList []string
	for _, v := range s.Routes {
		textList = append(textList, v.Format(peekOne(prefix)+Indent))
	}
	w.WriteText(strings.Join(textList, NewLine))
	w.NewLine()
	w.WriteText(s.RBrace.Format(prefix...))
	return w.String()
}

func (s *ServiceStmt) End() token.Position {
	return s.RBrace.End()
}

func (s *ServiceStmt) Pos() token.Position {
	if s.AtServerStmt != nil {
		return s.AtServerStmt.Pos()
	}
	return s.Service.Pos()
}

func (s *ServiceStmt) stmtNode() {}

type ServiceNameExpr struct {
	Name *TokenNode
}

func (s *ServiceNameExpr) HasHeadCommentGroup() bool {
	return s.Name.HasHeadCommentGroup()
}

func (s *ServiceNameExpr) HasLeadingCommentGroup() bool {
	return s.Name.HasLeadingCommentGroup()
}

func (s *ServiceNameExpr) CommentGroup() (head, leading CommentGroup) {
	return s.Name.HeadCommentGroup, s.Name.LeadingCommentGroup
}

func (s *ServiceNameExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.WriteText(s.Name.Format(prefix...))
	return w.String()
}

func (s *ServiceNameExpr) End() token.Position {
	return s.Name.End()
}

func (s *ServiceNameExpr) Pos() token.Position {
	return s.Name.Pos()
}

func (s *ServiceNameExpr) exprNode() {}

type AtHandlerStmt struct {
	AtHandler *TokenNode
	Name      *TokenNode
}

func (a *AtHandlerStmt) HasHeadCommentGroup() bool {
	return a.AtHandler.HasHeadCommentGroup()
}

func (a *AtHandlerStmt) HasLeadingCommentGroup() bool {
	return a.Name.HasLeadingCommentGroup()
}

func (a *AtHandlerStmt) CommentGroup() (head, leading CommentGroup) {
	return a.AtHandler.HeadCommentGroup, a.Name.LeadingCommentGroup
}

func (a *AtHandlerStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.WriteText(a.AtHandler.Format(prefix...) + WhiteSpace + a.Name.Format())
	return w.String()
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
}

func (s *ServiceItemStmt) HasHeadCommentGroup() bool {
	if s.AtDoc != nil {
		return s.AtDoc.HasHeadCommentGroup()
	}
	return s.AtHandler.HasHeadCommentGroup()
}

func (s *ServiceItemStmt) HasLeadingCommentGroup() bool {
	return s.Route.HasLeadingCommentGroup()
}

func (s *ServiceItemStmt) CommentGroup() (head, leading CommentGroup) {
	_, leading = s.Route.CommentGroup()
	if s.AtDoc != nil {
		head, _ = s.AtDoc.CommentGroup()
		return head, leading
	}
	head, _ = s.AtHandler.CommentGroup()
	return head, leading
}

func (s *ServiceItemStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	if s.AtDoc != nil {
		w.WriteText(s.AtDoc.Format(prefix...))
		w.NewLine()
	}
	w.WriteText(s.AtHandler.Format(prefix...))
	w.NewLine()
	w.WriteText(s.Route.Format(prefix...))
	return w.String()
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
	Method   *TokenNode
	Path     *PathExpr
	Request  *BodyStmt
	Returns  *TokenNode
	Response *BodyStmt
}

func (r *RouteStmt) HasHeadCommentGroup() bool {
	return r.Method.HasHeadCommentGroup()
}

func (r *RouteStmt) HasLeadingCommentGroup() bool {
	if r.Response != nil {
		return r.Response.HasLeadingCommentGroup()
	} else if r.Returns != nil {
		return r.Returns.HasLeadingCommentGroup()
	} else if r.Request != nil {
		return r.Request.HasLeadingCommentGroup()
	}
	return r.Path.HasLeadingCommentGroup()
}

func (r *RouteStmt) CommentGroup() (head, leading CommentGroup) {
	head, _ = r.Method.CommentGroup()
	if r.Response != nil {
		_, leading = r.Response.CommentGroup()
	} else if r.Returns != nil {
		_, leading = r.Returns.CommentGroup()
	} else if r.Request != nil {
		_, leading = r.Request.CommentGroup()
	}
	return head, leading
}

func (r *RouteStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	if r.Request == nil {
		w.Write(withNode(r.Method, r.Path), withPrefix(prefix...),
			withMode(ModeExpectInSameLine))
	} else if r.Returns == nil {
		w.Write(withNode(r.Method, r.Path, r.Request),
			withPrefix(prefix...), withMode(ModeExpectInSameLine))
	} else {
		w.Write(withNode(r.Method, r.Path, r.Request,
			r.Returns, r.Response), withPrefix(prefix...),
			withMode(ModeExpectInSameLine))
	}
	return w.String()
}

func (r *RouteStmt) End() token.Position {
	if r.Response != nil {
		return r.Response.End()
	}
	if r.Returns != nil {
		return r.Returns.Pos()
	}
	if r.Request != nil {
		return r.Request.End()
	}
	return r.Path.End()
}

func (r *RouteStmt) Pos() token.Position {
	return r.Method.Pos()
}

func (r *RouteStmt) stmtNode() {}

type PathExpr struct {
	Value *TokenNode
}

func (p *PathExpr) HasHeadCommentGroup() bool {
	return p.Value.HasHeadCommentGroup()
}

func (p *PathExpr) HasLeadingCommentGroup() bool {
	return p.Value.HasLeadingCommentGroup()
}

func (p *PathExpr) CommentGroup() (head, leading CommentGroup) {
	return p.Value.CommentGroup()
}

func (p *PathExpr) Format(prefix ...string) string {
	return p.Value.Format(prefix...)
}

func (p *PathExpr) End() token.Position {
	return p.Value.End()
}

func (p *PathExpr) Pos() token.Position {
	return p.Value.Pos()
}

func (p *PathExpr) exprNode() {}

type BodyStmt struct {
	LParen *TokenNode
	Body   *BodyExpr
	RParen *TokenNode
}

func (b *BodyStmt) HasHeadCommentGroup() bool {
	return b.LParen.HasHeadCommentGroup()
}

func (b *BodyStmt) HasLeadingCommentGroup() bool {
	return b.RParen.HasLeadingCommentGroup()
}

func (b *BodyStmt) CommentGroup() (head, leading CommentGroup) {
	return b.LParen.HeadCommentGroup, b.RParen.LeadingCommentGroup
}

func (b *BodyStmt) Format(prefix ...string) string {
	return "(" + b.Body.Format() + ")"
}

func (b *BodyStmt) End() token.Position {
	return b.RParen.End()
}

func (b *BodyStmt) Pos() token.Position {
	return b.LParen.Pos()
}

func (b *BodyStmt) stmtNode() {}

type BodyExpr struct {
	LBrack *TokenNode
	RBrack *TokenNode
	Star   *TokenNode
	Value  *TokenNode
}

func (e *BodyExpr) HasHeadCommentGroup() bool {
	if e.LBrack != nil {
		return e.LBrack.HasHeadCommentGroup()
	} else if e.Star != nil {
		return e.Star.HasHeadCommentGroup()
	} else {
		return e.Value.HasHeadCommentGroup()
	}
}

func (e *BodyExpr) HasLeadingCommentGroup() bool {
	return e.Value.HasLeadingCommentGroup()
}

func (e *BodyExpr) CommentGroup() (head, leading CommentGroup) {
	if e.LBrack != nil {
		head = e.LBrack.HeadCommentGroup
	} else if e.Star != nil {
		head = e.Star.HeadCommentGroup
	} else {
		head = e.Value.HeadCommentGroup
	}
	return head, e.Value.LeadingCommentGroup
}

func (e *BodyExpr) End() token.Position {
	return e.Value.End()
}

func (e *BodyExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	if e.LBrack != nil {
		if e.Star != nil {
			w.WriteText(peekOne(prefix) + "[]*" + e.Value.Format(prefix...))
		} else {
			w.WriteText(peekOne(prefix) + "[]" + e.Value.Format(prefix...))
		}
	} else if e.Star != nil {
		w.WriteText(peekOne(prefix) + "*" + e.Value.Format(prefix...))
	} else {
		w.WriteText(peekOne(prefix) + e.Value.Format(prefix...))
	}
	return w.String()
}

func (e *BodyExpr) Pos() token.Position {
	if e.LBrack != nil {
		return e.LBrack.Pos()
	}
	if e.Star != nil {
		return e.Star.Pos()
	}
	return e.Value.Pos()
}

func (e *BodyExpr) exprNode() {}
