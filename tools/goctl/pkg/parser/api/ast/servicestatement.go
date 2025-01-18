package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

// AtServerStmt represents @server statement.
type AtServerStmt struct {
	// AtServer is the @server token.
	AtServer *TokenNode
	// LParen is the left parenthesis token.
	LParen *TokenNode
	// Values is the key-value pairs.
	Values []*KVExpr
	// RParen is the right parenthesis token.
	RParen *TokenNode
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
	atServerNode := transferTokenNode(a.AtServer, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	w.Write(withNode(atServerNode, a.LParen), expectSameLine())
	w.NewLine()
	for _, v := range a.Values {
		node := transferNilInfixNode([]*TokenNode{v.Key, v.Colon})
		node = transferTokenNode(node, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment())
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
	atDocNode := transferTokenNode(a.AtDoc, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	w.Write(withNode(atDocNode, a.LParen), expectSameLine())
	w.NewLine()
	for _, v := range a.Values {
		node := transferNilInfixNode([]*TokenNode{v.Key, v.Colon})
		node = transferTokenNode(node, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment())
		w.Write(withNode(node, v.Value), expectIndentInfix(), expectSameLine())
		w.NewLine()
	}
	w.Write(withNode(transferTokenNode(a.RParen, withTokenNodePrefix(prefix...))))
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
		text := s.AtServerStmt.Format()
		if len(text) > 0 {
			w.WriteText(text)
			w.NewLine()
		}
	}
	serviceNode := transferTokenNode(s.Service, withTokenNodePrefix(prefix...))
	w.Write(withNode(serviceNode, s.Name, s.LBrace), expectSameLine())
	if len(s.Routes) == 0 {
		w.Write(withNode(transferTokenNode(s.RBrace, withTokenNodePrefix(prefix...))))
		return w.String()
	}
	w.NewLine()
	for idx, route := range s.Routes {
		routeNode := transfer2TokenNode(route, false, withTokenNodePrefix(peekOne(prefix)+Indent))
		w.Write(withNode(routeNode))
		if idx < len(s.Routes)-1 {
			w.NewLine()
		}
	}
	w.Write(withNode(transferTokenNode(s.RBrace, withTokenNodePrefix(prefix...))))
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

func (s *ServiceNameExpr) Format(...string) string {
	w := NewBufferWriter()
	w.WriteText(s.Name.Format())
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
	atDocNode := transferTokenNode(a.AtHandler, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	nameNode := transferTokenNode(a.Name, ignoreHeadComment())
	w.Write(withNode(atDocNode, nameNode), expectSameLine())
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
	routeNode := transfer2TokenNode(s.Route, false, withTokenNodePrefix(prefix...))
	w.Write(withNode(routeNode))
	w.NewLine()
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
	methodNode := transferTokenNode(r.Method, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	if r.Response != nil {
		if r.Response.Body == nil {
			r.Response.RParen = transferTokenNode(r.Response.RParen, ignoreHeadComment())
			if r.Request != nil {
				w.Write(withNode(methodNode, r.Path, r.Request), expectSameLine())
			} else {
				w.Write(withNode(methodNode, r.Path), expectSameLine())
			}
		} else {
			r.Response.RParen = transferTokenNode(r.Response.RParen, ignoreHeadComment())
			if r.Request != nil {
				w.Write(withNode(methodNode, r.Path, r.Request, r.Returns, r.Response), expectSameLine())
			} else {
				w.Write(withNode(methodNode, r.Path, r.Returns, r.Response), expectSameLine())
			}
		}
	} else if r.Request != nil {
		r.Request.RParen = transferTokenNode(r.Request.RParen, ignoreHeadComment())
		w.Write(withNode(methodNode, r.Path, r.Request), expectSameLine())
	} else {
		pathNode := transferTokenNode(r.Path.Value, ignoreHeadComment())
		w.Write(withNode(methodNode, pathNode), expectSameLine())
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
	pathNode := transferTokenNode(p.Value, ignoreComment())
	return pathNode.Format(prefix...)
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

func (b *BodyStmt) Format(...string) string {
	w := NewBufferWriter()
	if b.Body == nil {
		return ""
	}
	w.Write(withNode(b.LParen, b.Body, b.RParen), withInfix(NilIndent), expectSameLine())
	return w.String()
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

func (e *BodyExpr) Format(...string) string {
	w := NewBufferWriter()
	if e.LBrack != nil {
		lbrackNode := transferTokenNode(e.LBrack, ignoreComment())
		rbrackNode := transferTokenNode(e.RBrack, ignoreComment())
		if e.Star != nil {
			starNode := transferTokenNode(e.Star, ignoreComment())
			w.Write(withNode(lbrackNode, rbrackNode, starNode, e.Value), withInfix(NilIndent), expectSameLine())
		} else {
			w.Write(withNode(lbrackNode, rbrackNode, e.Value), withInfix(NilIndent), expectSameLine())
		}
	} else if e.Star != nil {
		starNode := transferTokenNode(e.Star, ignoreComment())
		w.Write(withNode(starNode, e.Value), withInfix(NilIndent), expectSameLine())
	} else {
		w.Write(withNode(e.Value))
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

func (e *BodyExpr) IsArrayType() bool {
	return e.LBrack != nil
}
