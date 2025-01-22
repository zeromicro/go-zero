package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

// InfoStmt is the info statement.
type InfoStmt struct {
	// Info is the info keyword.
	Info *TokenNode
	// LParen is the left parenthesis.
	LParen *TokenNode
	// Values is the info values.
	Values []*KVExpr
	// RParen is the right parenthesis.
	RParen *TokenNode
}

func (i *InfoStmt) HasHeadCommentGroup() bool {
	return i.Info.HasHeadCommentGroup()
}

func (i *InfoStmt) HasLeadingCommentGroup() bool {
	return i.RParen.HasLeadingCommentGroup()
}

func (i *InfoStmt) CommentGroup() (head, leading CommentGroup) {
	return i.Info.HeadCommentGroup, i.RParen.LeadingCommentGroup
}

func (i *InfoStmt) Format(prefix ...string) string {
	if len(i.Values) == 0 {
		return ""
	}
	var textList []string
	for _, v := range i.Values {
		if v.Value.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format(Indent))
	}
	if len(textList) == 0 {
		return ""
	}

	w := NewBufferWriter()
	infoNode := transferTokenNode(i.Info, withTokenNodePrefix(prefix...), ignoreLeadingComment())
	w.Write(withNode(infoNode, i.LParen))
	w.NewLine()
	for _, v := range i.Values {
		node := transferNilInfixNode([]*TokenNode{v.Key, v.Colon})
		node = transferTokenNode(node, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment())
		w.Write(withNode(node, v.Value), expectIndentInfix(), expectSameLine())
		w.NewLine()
	}
	w.Write(withNode(transferTokenNode(i.RParen, withTokenNodePrefix(prefix...))))
	return w.String()
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.End()
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Pos()
}

func (i *InfoStmt) stmtNode() {}
