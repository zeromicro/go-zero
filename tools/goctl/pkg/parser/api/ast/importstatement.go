package ast

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type ImportStmt interface {
	Stmt
	importNode()
}

type ImportLiteralStmt struct {
	Import *TokenNode
	Value  *TokenNode
}

func (i *ImportLiteralStmt) HasHeadCommentGroup() bool {
	return i.Import.HasHeadCommentGroup()
}

func (i *ImportLiteralStmt) HasLeadingCommentGroup() bool {
	return i.Value.HasLeadingCommentGroup()
}

func (i *ImportLiteralStmt) CommentGroup() (head, leading CommentGroup) {
	return i.Import.HeadCommentGroup, i.Value.LeadingCommentGroup
}

func (i *ImportLiteralStmt) Format(prefix ...string) (result string) {
	if i.Value.IsZeroString() {
		return ""
	}
	w := NewBufferWriter()
	importNode := transferTokenNode(i.Import, ignoreLeadingComment(), withTokenNodePrefix(prefix...))
	w.Write(withNode(importNode, i.Value), withMode(ModeExpectInSameLine))
	return w.String()
}

func (i *ImportLiteralStmt) End() token.Position {
	return i.Value.End()
}

func (i *ImportLiteralStmt) importNode() {}

func (i *ImportLiteralStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportLiteralStmt) stmtNode() {}

type ImportGroupStmt struct {
	Import *TokenNode
	LParen *TokenNode
	Values []*TokenNode
	RParen *TokenNode
}

func (i *ImportGroupStmt) HasHeadCommentGroup() bool {
	return i.Import.HasHeadCommentGroup()
}

func (i *ImportGroupStmt) HasLeadingCommentGroup() bool {
	return i.RParen.HasLeadingCommentGroup()
}

func (i *ImportGroupStmt) CommentGroup() (head, leading CommentGroup) {
	return i.Import.HeadCommentGroup, i.RParen.LeadingCommentGroup
}

func (i *ImportGroupStmt) Format(prefix ...string) string {
	var textList []string
	for _, v := range i.Values {
		if v.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format(Indent))
	}
	if len(textList) == 0 {
		return ""
	}

	importNode := transferTokenNode(i.Import, ignoreLeadingComment(), withTokenNodePrefix(prefix...))
	w := NewBufferWriter()
	w.Write(withNode(importNode, i.LParen), expectSameLine())
	w.NewLine()
	for _, v := range i.Values {
		node := transferTokenNode(v, withTokenNodePrefix(peekOne(prefix)+Indent))
		w.Write(withNode(node), expectSameLine())
		w.NewLine()
	}
	w.Write(withNode(transferTokenNode(i.RParen, withTokenNodePrefix(prefix...))))
	return w.String()
}

func (i *ImportGroupStmt) End() token.Position {
	return i.RParen.End()
}

func (i *ImportGroupStmt) importNode() {}

func (i *ImportGroupStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportGroupStmt) stmtNode() {}
