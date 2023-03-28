package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

/*******************TypeStmt Begin********************/

// TypeStmt is the interface for type statement.
type TypeStmt interface {
	Stmt
	typeNode()
}

// TypeLiteralStmt is the type statement for type literal.
type TypeLiteralStmt struct {
	// Type is the type keyword.
	Type *TokenNode
	// Expr is the type expression.
	Expr *TypeExpr
}

func (t *TypeLiteralStmt) HasHeadCommentGroup() bool {
	return t.Type.HasHeadCommentGroup()
}

func (t *TypeLiteralStmt) HasLeadingCommentGroup() bool {
	return t.Expr.HasLeadingCommentGroup()
}

func (t *TypeLiteralStmt) CommentGroup() (head, leading CommentGroup) {
	_, leading = t.Expr.CommentGroup()
	return t.Type.HeadCommentGroup, leading
}

func (t *TypeLiteralStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(withNode(t.Type, t.Expr), withPrefix(prefix...), expectSameLine())
	return w.String()
}

func (t *TypeLiteralStmt) End() token.Position {
	return t.Expr.End()
}

func (t *TypeLiteralStmt) Pos() token.Position {
	return t.Type.Pos()
}

func (t *TypeLiteralStmt) stmtNode() {}
func (t *TypeLiteralStmt) typeNode() {}

// TypeGroupStmt is the type statement for type group.
type TypeGroupStmt struct {
	// Type is the type keyword.
	Type *TokenNode
	// LParen is the left parenthesis.
	LParen *TokenNode
	// ExprList is the type expression list.
	ExprList []*TypeExpr
	// RParen is the right parenthesis.
	RParen *TokenNode
}

func (t *TypeGroupStmt) HasHeadCommentGroup() bool {
	return t.Type.HasHeadCommentGroup()
}

func (t *TypeGroupStmt) HasLeadingCommentGroup() bool {
	return t.RParen.HasLeadingCommentGroup()
}

func (t *TypeGroupStmt) CommentGroup() (head, leading CommentGroup) {
	return t.Type.HeadCommentGroup, t.RParen.LeadingCommentGroup
}

func (t *TypeGroupStmt) Format(prefix ...string) string {
	if len(t.ExprList) == 0 {
		return ""
	}
	w := NewBufferWriter()
	typeNode := transferTokenNode(t.Type, withTokenNodePrefix(prefix...))
	w.Write(withNode(typeNode, t.LParen), expectSameLine())
	w.NewLine()
	for _, e := range t.ExprList {
		w.Write(withNode(e), withPrefix(peekOne(prefix)+Indent))
		w.NewLine()
	}
	w.WriteText(t.RParen.Format(prefix...))
	return w.String()
}

func (t *TypeGroupStmt) End() token.Position {
	return t.RParen.End()
}

func (t *TypeGroupStmt) Pos() token.Position {
	return t.Type.Pos()
}

func (t *TypeGroupStmt) stmtNode() {}
func (t *TypeGroupStmt) typeNode() {}

/*******************TypeStmt End********************/

/*******************TypeExpr Begin********************/

// TypeExpr is the type expression.
type TypeExpr struct {
	// Name is the type name.
	Name *TokenNode
	// Assign is the assign operator.
	Assign *TokenNode
	// DataType is the data type.
	DataType DataType
}

func (e *TypeExpr) HasHeadCommentGroup() bool {
	return e.Name.HasHeadCommentGroup()
}

func (e *TypeExpr) HasLeadingCommentGroup() bool {
	return e.DataType.HasLeadingCommentGroup()
}

func (e *TypeExpr) CommentGroup() (head, leading CommentGroup) {
	_, leading = e.DataType.CommentGroup()
	return e.Name.HeadCommentGroup, leading
}

func (e *TypeExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	nameNode := transferTokenNode(e.Name, withTokenNodePrefix(prefix...))
	dataTypeNode := transfer2TokenNode(e.DataType, false, withTokenNodePrefix(prefix...))
	if e.Assign != nil {
		w.Write(withNode(nameNode, e.Assign, dataTypeNode), expectSameLine())
	} else {
		w.Write(withNode(nameNode, dataTypeNode), expectSameLine())
	}
	return w.String()
}

func (e *TypeExpr) End() token.Position {
	return e.DataType.End()
}

func (e *TypeExpr) Pos() token.Position {
	return e.Name.Pos()
}

func (e *TypeExpr) exprNode() {}

func (e *TypeExpr) isStruct() bool {
	return e.DataType.ContainsStruct()
}

/*******************TypeExpr Begin********************/

/*******************Elem Begin********************/

// ElemExpr is the element expression.
type ElemExpr struct {
	// Name is the field element name.
	Name []*TokenNode
	// DataType is the field data type.
	DataType DataType
	// Tag is the field tag.
	Tag *TokenNode
}

// IsAnonymous returns true if the element is anonymous.
func (e *ElemExpr) IsAnonymous() bool {
	return len(e.Name) == 0
}

func (e *ElemExpr) HasHeadCommentGroup() bool {
	if e.IsAnonymous() {
		return e.DataType.HasHeadCommentGroup()
	}
	return e.Name[0].HasHeadCommentGroup()
}

func (e *ElemExpr) HasLeadingCommentGroup() bool {
	if e.Tag != nil {
		return e.Tag.HasLeadingCommentGroup()
	}
	return e.DataType.HasLeadingCommentGroup()
}

func (e *ElemExpr) CommentGroup() (head, leading CommentGroup) {
	if e.Tag != nil {
		leading = e.Tag.LeadingCommentGroup
	} else {
		_, leading = e.DataType.CommentGroup()
	}
	if e.IsAnonymous() {
		head, _ := e.DataType.CommentGroup()
		return head, leading
	}
	return e.Name[0].HeadCommentGroup, leading
}

func (e *ElemExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	var nameNodeList []*TokenNode
	for idx, n := range e.Name {
		if idx == 0 {
			nameNodeList = append(nameNodeList,
				transferTokenNode(n, ignoreLeadingComment()))
		} else if idx < len(e.Name)-1 {
			nameNodeList = append(nameNodeList,
				transferTokenNode(n, ignoreLeadingComment(), ignoreHeadComment()))
		} else {
			nameNodeList = append(nameNodeList, transferTokenNode(n, ignoreHeadComment()))
		}
	}

	var dataTypeOption []tokenNodeOption
	if e.DataType.ContainsStruct() {
		dataTypeOption = append(dataTypeOption, withTokenNodePrefix(peekOne(prefix)+Indent))
	} else {
		dataTypeOption = append(dataTypeOption, withTokenNodePrefix(prefix...))
	}
	dataTypeNode := transfer2TokenNode(e.DataType, false, dataTypeOption...)
	if len(nameNodeList) > 0 {
		nameNode := transferNilInfixNode(nameNodeList,
			withTokenNodePrefix(prefix...), withTokenNodeInfix(", "))
		if e.Tag != nil {
			w.Write(withNode(nameNode, dataTypeNode, e.Tag), expectIndentInfix(), expectSameLine())
		} else {
			w.Write(withNode(nameNode, dataTypeNode), expectIndentInfix(), expectSameLine())
		}
	} else {
		if e.Tag != nil {
			w.Write(withNode(dataTypeNode, e.Tag), expectIndentInfix(), expectSameLine())
		} else {
			w.Write(withNode(dataTypeNode), expectIndentInfix(), expectSameLine())
		}
	}
	return w.String()
}

func (e *ElemExpr) End() token.Position {
	if e.Tag != nil {
		return e.Tag.End()
	}
	return e.DataType.End()
}

func (e *ElemExpr) Pos() token.Position {
	if len(e.Name) > 0 {
		return e.Name[0].Pos()
	}
	return token.IllegalPosition
}

func (e *ElemExpr) exprNode() {}

/*******************Elem End********************/

/*******************ElemExprList Begin********************/

// ElemExprList is the element expression list.
type ElemExprList []*ElemExpr

/*******************ElemExprList Begin********************/

/*******************DataType Begin********************/

// DataType represents the data type.
type DataType interface {
	Expr
	dataTypeNode()
	// CanEqual returns true if the data type can be equal.
	CanEqual() bool
	// ContainsStruct returns true if the data type contains struct.
	ContainsStruct() bool
	// RawText returns the raw text of the data type.
	RawText() string
}

// AnyDataType is the any data type.
type AnyDataType struct {
	// Any is the any token node.
	Any     *TokenNode
	isChild bool
}

func (t *AnyDataType) HasHeadCommentGroup() bool {
	return t.Any.HasHeadCommentGroup()
}

func (t *AnyDataType) HasLeadingCommentGroup() bool {
	return t.Any.HasLeadingCommentGroup()
}

func (t *AnyDataType) CommentGroup() (head, leading CommentGroup) {
	return t.Any.HeadCommentGroup, t.Any.LeadingCommentGroup
}

func (t *AnyDataType) Format(prefix ...string) string {
	return t.Any.Format(prefix...)
}

func (t *AnyDataType) End() token.Position {
	return t.Any.End()
}

func (t *AnyDataType) RawText() string {
	return t.Any.Token.Text
}

func (t *AnyDataType) ContainsStruct() bool {
	return false
}

func (t *AnyDataType) Pos() token.Position {
	return t.Any.Pos()
}

func (t *AnyDataType) exprNode() {}

func (t *AnyDataType) dataTypeNode() {}

func (t *AnyDataType) CanEqual() bool {
	return true
}

// ArrayDataType is the array data type.
type ArrayDataType struct {
	// LB is the left bracket token node.
	LBrack *TokenNode
	// Len is the array length.
	Length *TokenNode
	// RB is the right bracket token node.
	RBrack *TokenNode
	// DataType is the array data type.
	DataType DataType
	isChild  bool
}

func (t *ArrayDataType) HasHeadCommentGroup() bool {
	return t.LBrack.HasHeadCommentGroup()
}

func (t *ArrayDataType) HasLeadingCommentGroup() bool {
	return t.DataType.HasLeadingCommentGroup()
}

func (t *ArrayDataType) CommentGroup() (head, leading CommentGroup) {
	_, leading = t.DataType.CommentGroup()
	return t.LBrack.HeadCommentGroup, leading
}

func (t *ArrayDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	lbrack := transferTokenNode(t.LBrack, ignoreLeadingComment())
	lengthNode := transferTokenNode(t.Length, ignoreLeadingComment())
	rbrack := transferTokenNode(t.RBrack, ignoreHeadComment())
	var dataType *TokenNode
	var options []tokenNodeOption
	options = append(options, withTokenNodePrefix(prefix...))
	if t.isChild {
		options = append(options, ignoreComment())
	} else {
		options = append(options, ignoreHeadComment())
	}

	dataType = transfer2TokenNode(t.DataType, false, options...)
	node := transferNilInfixNode([]*TokenNode{lbrack, lengthNode, rbrack, dataType})
	w.Write(withNode(node))
	return w.String()
}

func (t *ArrayDataType) End() token.Position {
	return t.DataType.End()
}

func (t *ArrayDataType) RawText() string {
	return t.Format("")
}

func (t *ArrayDataType) ContainsStruct() bool {
	return t.DataType.ContainsStruct()
}

func (t *ArrayDataType) CanEqual() bool {
	return t.DataType.CanEqual()
}

func (t *ArrayDataType) Pos() token.Position {
	return t.LBrack.Pos()
}

func (t *ArrayDataType) exprNode()     {}
func (t *ArrayDataType) dataTypeNode() {}

// BaseDataType is a common id type which contains bool, uint8, uint16, uint32,
// uint64, int8, int16, int32, int64, float32, float64, complex64, complex128,
// string, int, uint, uintptr, byte, rune, any.
type BaseDataType struct {
	// Base is the base token node.
	Base    *TokenNode
	isChild bool
}

func (t *BaseDataType) HasHeadCommentGroup() bool {
	return t.Base.HasHeadCommentGroup()
}

func (t *BaseDataType) HasLeadingCommentGroup() bool {
	return t.Base.HasLeadingCommentGroup()
}

func (t *BaseDataType) CommentGroup() (head, leading CommentGroup) {
	return t.Base.HeadCommentGroup, t.Base.LeadingCommentGroup
}

func (t *BaseDataType) Format(prefix ...string) string {
	return t.Base.Format(prefix...)
}

func (t *BaseDataType) End() token.Position {
	return t.Base.End()
}

func (t *BaseDataType) RawText() string {
	return t.Base.Token.Text
}

func (t *BaseDataType) ContainsStruct() bool {
	return false
}

func (t *BaseDataType) CanEqual() bool {
	return true
}

func (t *BaseDataType) Pos() token.Position {
	return t.Base.Pos()
}

func (t *BaseDataType) exprNode()     {}
func (t *BaseDataType) dataTypeNode() {}

// InterfaceDataType is the interface data type.
type InterfaceDataType struct {
	// Interface is the interface token node.
	Interface *TokenNode
	isChild   bool
}

func (t *InterfaceDataType) HasHeadCommentGroup() bool {
	return t.Interface.HasHeadCommentGroup()
}

func (t *InterfaceDataType) HasLeadingCommentGroup() bool {
	return t.Interface.HasLeadingCommentGroup()
}

func (t *InterfaceDataType) CommentGroup() (head, leading CommentGroup) {
	return t.Interface.HeadCommentGroup, t.Interface.LeadingCommentGroup
}

func (t *InterfaceDataType) Format(prefix ...string) string {
	return t.Interface.Format(prefix...)
}

func (t *InterfaceDataType) End() token.Position {
	return t.Interface.End()
}

func (t *InterfaceDataType) RawText() string {
	return t.Interface.Token.Text
}

func (t *InterfaceDataType) ContainsStruct() bool {
	return false
}

func (t *InterfaceDataType) CanEqual() bool {
	return true
}

func (t *InterfaceDataType) Pos() token.Position {
	return t.Interface.Pos()
}

func (t *InterfaceDataType) exprNode() {}

func (t *InterfaceDataType) dataTypeNode() {}

// MapDataType is the map data type.
type MapDataType struct {
	// Map is the map token node.
	Map *TokenNode
	// Lbrack is the left bracket token node.
	LBrack *TokenNode
	// Key is the map key data type.
	Key DataType
	// Rbrack is the right bracket token node.
	RBrack *TokenNode
	// Value is the map value data type.
	Value   DataType
	isChild bool
}

func (t *MapDataType) HasHeadCommentGroup() bool {
	return t.Map.HasHeadCommentGroup()
}

func (t *MapDataType) HasLeadingCommentGroup() bool {
	return t.Value.HasLeadingCommentGroup()
}

func (t *MapDataType) CommentGroup() (head, leading CommentGroup) {
	_, leading = t.Value.CommentGroup()
	return t.Map.HeadCommentGroup, leading
}

func (t *MapDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	mapNode := transferTokenNode(t.Map, ignoreLeadingComment())
	lbrack := transferTokenNode(t.LBrack, ignoreLeadingComment())
	rbrack := transferTokenNode(t.RBrack, ignoreComment())
	var keyOption, valueOption []tokenNodeOption
	keyOption = append(keyOption, ignoreComment())
	valueOption = append(valueOption, withTokenNodePrefix(prefix...))

	if t.isChild {
		valueOption = append(valueOption, ignoreComment())
	} else {
		valueOption = append(valueOption, ignoreHeadComment())
	}

	keyDataType := transfer2TokenNode(t.Key, true, keyOption...)
	valueDataType := transfer2TokenNode(t.Value, false, valueOption...)
	node := transferNilInfixNode([]*TokenNode{mapNode, lbrack, keyDataType, rbrack, valueDataType})
	w.Write(withNode(node))
	return w.String()
}

func (t *MapDataType) End() token.Position {
	return t.Value.End()
}

func (t *MapDataType) RawText() string {
	return t.Format("")
}

func (t *MapDataType) ContainsStruct() bool {
	return t.Key.ContainsStruct() || t.Value.ContainsStruct()
}

func (t *MapDataType) CanEqual() bool {
	return false
}

func (t *MapDataType) Pos() token.Position {
	return t.Map.Pos()
}

func (t *MapDataType) exprNode()     {}
func (t *MapDataType) dataTypeNode() {}

// PointerDataType is the pointer data type.
type PointerDataType struct {
	// Star is the star token node.
	Star *TokenNode
	// DataType is the pointer data type.
	DataType DataType
	isChild  bool
}

func (t *PointerDataType) HasHeadCommentGroup() bool {
	return t.Star.HasHeadCommentGroup()
}

func (t *PointerDataType) HasLeadingCommentGroup() bool {
	return t.DataType.HasLeadingCommentGroup()
}

func (t *PointerDataType) CommentGroup() (head, leading CommentGroup) {
	_, leading = t.DataType.CommentGroup()
	return t.Star.HeadCommentGroup, leading
}

func (t *PointerDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	star := transferTokenNode(t.Star, ignoreLeadingComment(), withTokenNodePrefix(prefix...))
	var dataTypeOption []tokenNodeOption
	dataTypeOption = append(dataTypeOption, ignoreHeadComment())
	dataType := transfer2TokenNode(t.DataType, false, dataTypeOption...)
	node := transferNilInfixNode([]*TokenNode{star, dataType})
	w.Write(withNode(node))
	return w.String()
}

func (t *PointerDataType) End() token.Position {
	return t.DataType.End()
}

func (t *PointerDataType) RawText() string {
	return t.Format("")
}

func (t *PointerDataType) ContainsStruct() bool {
	return t.DataType.ContainsStruct()
}

func (t *PointerDataType) CanEqual() bool {
	return t.DataType.CanEqual()
}

func (t *PointerDataType) Pos() token.Position {
	return t.Star.Pos()
}

func (t *PointerDataType) exprNode()     {}
func (t *PointerDataType) dataTypeNode() {}

// SliceDataType is the slice data type.
type SliceDataType struct {
	// Lbrack is the left bracket token node.
	LBrack *TokenNode
	// Rbrack is the right bracket token node.
	RBrack *TokenNode
	// DataType is the slice data type.
	DataType DataType
	isChild  bool
}

func (t *SliceDataType) HasHeadCommentGroup() bool {
	return t.LBrack.HasHeadCommentGroup()
}

func (t *SliceDataType) HasLeadingCommentGroup() bool {
	return t.DataType.HasLeadingCommentGroup()
}

func (t *SliceDataType) CommentGroup() (head, leading CommentGroup) {
	_, leading = t.DataType.CommentGroup()
	return t.LBrack.HeadCommentGroup, leading
}

func (t *SliceDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	lbrack := transferTokenNode(t.LBrack, ignoreLeadingComment())
	rbrack := transferTokenNode(t.RBrack, ignoreHeadComment())
	dataType := transfer2TokenNode(t.DataType, false, withTokenNodePrefix(prefix...), ignoreHeadComment())
	node := transferNilInfixNode([]*TokenNode{lbrack, rbrack, dataType})
	w.Write(withNode(node))
	return w.String()
}

func (t *SliceDataType) End() token.Position {
	return t.DataType.End()
}

func (t *SliceDataType) RawText() string {
	return t.Format("")
}

func (t *SliceDataType) ContainsStruct() bool {
	return t.DataType.ContainsStruct()
}

func (t *SliceDataType) CanEqual() bool {
	return false
}

func (t *SliceDataType) Pos() token.Position {
	return t.LBrack.Pos()
}

func (t *SliceDataType) exprNode()     {}
func (t *SliceDataType) dataTypeNode() {}

// StructDataType is the structure data type.
type StructDataType struct {
	// Lbrace is the left brace token node.
	LBrace *TokenNode
	// Elements is the structure elements.
	Elements ElemExprList
	// Rbrace is the right brace token node.
	RBrace  *TokenNode
	isChild bool
}

func (t *StructDataType) HasHeadCommentGroup() bool {
	return t.LBrace.HasHeadCommentGroup()
}

func (t *StructDataType) HasLeadingCommentGroup() bool {
	return t.RBrace.HasLeadingCommentGroup()
}

func (t *StructDataType) CommentGroup() (head, leading CommentGroup) {
	return t.LBrace.HeadCommentGroup, t.RBrace.LeadingCommentGroup
}

func (t *StructDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	if len(t.Elements) == 0 {
		lbrace := transferTokenNode(t.LBrace, withTokenNodePrefix(prefix...), ignoreLeadingComment())
		rbrace := transferTokenNode(t.RBrace, ignoreHeadComment())
		brace := transferNilInfixNode([]*TokenNode{lbrace, rbrace})
		w.Write(withNode(brace), expectSameLine())
		return w.String()
	}
	w.WriteText(t.LBrace.Format(NilIndent))
	w.NewLine()
	for _, e := range t.Elements {
		var nameNodeList []*TokenNode
		if len(e.Name) > 0 {
			for idx, n := range e.Name {
				if idx == 0 {
					nameNodeList = append(nameNodeList,
						transferTokenNode(n, withTokenNodePrefix(peekOne(prefix)+Indent), ignoreLeadingComment()))
				} else if idx < len(e.Name)-1 {
					nameNodeList = append(nameNodeList,
						transferTokenNode(n, ignoreLeadingComment(), ignoreHeadComment()))
				} else {
					nameNodeList = append(nameNodeList, transferTokenNode(n, ignoreHeadComment()))
				}
			}
		}
		var dataTypeOption []tokenNodeOption
		if e.DataType.ContainsStruct() || e.IsAnonymous() {
			dataTypeOption = append(dataTypeOption, withTokenNodePrefix(peekOne(prefix)+Indent))
		} else {
			dataTypeOption = append(dataTypeOption, withTokenNodePrefix(prefix...))
		}
		dataTypeNode := transfer2TokenNode(e.DataType, false, dataTypeOption...)
		if len(nameNodeList) > 0 {
			nameNode := transferNilInfixNode(nameNodeList, withTokenNodeInfix(", "))
			if e.Tag != nil {
				if e.DataType.ContainsStruct() {
					w.Write(withNode(nameNode, dataTypeNode, e.Tag), expectSameLine())
				} else {
					w.Write(withNode(nameNode, e.DataType, e.Tag), expectIndentInfix(), expectSameLine())
				}
			} else {
				if e.DataType.ContainsStruct() {
					w.Write(withNode(nameNode, dataTypeNode), expectSameLine())
				} else {
					w.Write(withNode(nameNode, e.DataType), expectIndentInfix(), expectSameLine())
				}
			}
		} else {
			if e.Tag != nil {
				if e.DataType.ContainsStruct() {
					w.Write(withNode(dataTypeNode, e.Tag), expectSameLine())
				} else {
					w.Write(withNode(e.DataType, e.Tag), expectIndentInfix(), expectSameLine())
				}
			} else {
				if e.DataType.ContainsStruct() {
					w.Write(withNode(dataTypeNode), expectSameLine())
				} else {
					w.Write(withNode(dataTypeNode), expectIndentInfix(), expectSameLine())
				}
			}
		}
		w.NewLine()
	}
	w.WriteText(t.RBrace.Format(prefix...))
	return w.String()
}

func (t *StructDataType) End() token.Position {
	return t.RBrace.End()
}

func (t *StructDataType) RawText() string {
	return t.Format("")
}

func (t *StructDataType) ContainsStruct() bool {
	return true
}

func (t *StructDataType) CanEqual() bool {
	for _, v := range t.Elements {
		if !v.DataType.CanEqual() {
			return false
		}
	}
	return true
}

func (t *StructDataType) Pos() token.Position {
	return t.LBrace.Pos()
}

func (t *StructDataType) exprNode()     {}
func (t *StructDataType) dataTypeNode() {}

/*******************DataType End********************/
