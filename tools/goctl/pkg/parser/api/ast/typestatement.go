package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

/*******************TypeStmt Begin********************/

type TypeStmt interface {
	Stmt
	typeNode()
}

type TypeLiteralStmt struct {
	Type *TokenNode
	Expr *TypeExpr
}

func (t *TypeLiteralStmt) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(WithNode(t.Type, t.Expr), WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
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

type TypeGroupStmt struct {
	Type     *TokenNode
	LParen   *TokenNode
	ExprList []*TypeExpr
	RParen   *TokenNode
}

func (t *TypeGroupStmt) Format(prefix ...string) string {
	if len(t.ExprList) == 0 {
		return ""
	}
	w := NewBufferWriter()
	w.Write(WithNode(t.Type, t.LParen), WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
	w.NewLine()
	for _, e := range t.ExprList {
		w.Write(WithNode(e), WithPrefix(peekOne(prefix)+Indent), WithMode(ModeExpectInSameLine))
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

type TypeExpr struct {
	Name     *TokenNode
	Assign   *TokenNode
	DataType DataType
}

func (e *TypeExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	if e.Assign != nil {
		w.Write(WithNode(e.Name, e.Assign, e.DataType),
			WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
	} else {
		w.Write(WithNode(e.Name, e.DataType),
			WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
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

type ElemExpr struct {
	Name     []*TokenNode
	DataType DataType
	Tag      *TokenNode
}

func (e *ElemExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	var nameList []string
	for _, n := range e.Name {
		nameList = append(nameList, n.Token.Text)
	}

	var nameNode *TokenNode
	if len(e.Name) == 1 {
		nameNode = e.Name[0]
	} else {
		nameNode = NewTokenNode(token.Token{
			Type:     token.IDENT,
			Text:     strings.Join(nameList, ", "),
			Position: e.Name[0].Pos(),
		})
	}

	if e.Tag != nil {
		w.Write(WithNode(nameNode, e.DataType, e.Tag),
			WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
	} else {
		w.Write(WithNode(nameNode, e.DataType),
			WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
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

type ElemExprList []*ElemExpr

func (e ElemExprList) Pos() token.Position {
	if len(e) > 0 {
		return e[0].Pos()
	}
	return token.IllegalPosition
}

func (e ElemExprList) exprNode() {}

/*******************ElemExprList Begin********************/

/*******************DataType Begin********************/

type DataType interface {
	Expr
	dataTypeNode()
	CanEqual() bool
	ContainsStruct() bool
	RawText() string
}

type StructDataType struct {
	LBrace   *TokenNode
	Elements ElemExprList
	RBrace   *TokenNode
}

func (t *StructDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	if len(t.Elements) == 0 {
		w.WriteText("{}")
		return w.String()
	}
	w.WriteText(t.LBrace.Format(NilIndent))
	w.NewLine()
	for _, e := range t.Elements {
		var nameList []string
		for _, n := range e.Name {
			nameList = append(nameList, n.Token.Text)
		}

		var nameNode *TokenNode
		if len(e.Name) == 1 {
			nameNode = e.Name[0]
		} else {
			nameNode = NewTokenNode(token.Token{
				Type:     token.IDENT,
				Text:     strings.Join(nameList, ", "),
				Position: e.Name[0].Pos(),
			})
			nameNode.HeadCommentGroup = e.Name[0].HeadCommentGroup
			nameNode.LeadingCommentGroup = e.Name[len(e.Name)-1].LeadingCommentGroup
		}

		if e.Tag != nil {
			w.Write(WithNode(nameNode, e.DataType, e.Tag),
				WithPrefix(peekOne(prefix)+Indent), WithMode(ModeExpectInSameLine))
		} else {
			w.Write(WithNode(nameNode, e.DataType),
				WithPrefix(peekOne(prefix)+Indent), WithMode(ModeExpectInSameLine))
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
	b := bytes.NewBuffer(nil)
	b.WriteRune('{')
	for _, v := range t.Elements {
		b.WriteRune('\n')
		var nameList []string
		for _, n := range v.Name {
			nameList = append(nameList, n.Token.Text)
		}
		b.WriteString(fmt.Sprintf("%s %s %s", strings.Join(nameList, ", "), v.DataType.RawText(), v.Tag.Token.Text))
	}
	b.WriteRune('\n')
	b.WriteRune('}')
	return b.String()
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

type SliceDataType struct {
	LBrack   *TokenNode
	RBrack   *TokenNode
	DataType DataType
}

func (t *SliceDataType) Format(prefix ...string) string {
	brackWriter := NewBufferWriter()
	brackWriter.Write(WithNode(t.LBrack, t.RBrack),
		WithInfix(NilIndent), WithMode(ModeExpectInSameLine))

	w := NewBufferWriter()
	w.WriteText(brackWriter.String() + t.DataType.Format(prefix...))
	return w.String()
}

func (t *SliceDataType) End() token.Position {
	return t.DataType.End()
}

func (t *SliceDataType) RawText() string {
	return fmt.Sprintf("[]%s", t.DataType.RawText())
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

type MapDataType struct {
	Map    *TokenNode
	LBrack *TokenNode
	Key    DataType
	RBrack *TokenNode
	Value  DataType
}

func (t *MapDataType) Format(prefix ...string) string {
	w1 := NewBufferWriter()
	w1.Write(WithNode(t.Map, t.LBrack),
		WithInfix(NilIndent), WithMode(ModeExpectInSameLine))

	w := NewBufferWriter()
	w.WriteText(w1.String() + strings.TrimPrefix(t.Key.Format(prefix...),peekOne(prefix)) +
		t.RBrack.Format() + t.Value.Format(prefix...))
	return w.String()
}

func (t *MapDataType) End() token.Position {
	return t.Value.End()
}

func (t *MapDataType) RawText() string {
	return fmt.Sprintf("map[%s]%s", t.Key.RawText(), t.Value.RawText())
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

type ArrayDataType struct {
	LBrack   *TokenNode
	Length   *TokenNode
	RBrack   *TokenNode
	DataType DataType
}

func (t *ArrayDataType) Format(prefix ...string) string {
	w1 := NewBufferWriter()
	w1.Write(WithNode(t.LBrack, t.Length, t.RBrack),
		WithInfix(NilIndent), WithMode(ModeExpectInSameLine))
	w := NewBufferWriter()
	w.WriteText(w1.String() + t.DataType.Format(prefix...))
	return w.String()
}

func (t *ArrayDataType) End() token.Position {
	return t.DataType.End()
}

func (t *ArrayDataType) RawText() string {
	return ""
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

type InterfaceDataType struct {
	Interface *TokenNode
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

type PointerDataType struct {
	Star     *TokenNode
	DataType DataType
}

func (t *PointerDataType) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(WithNode(t.Star, t.DataType), WithInfix(NilIndent),
		WithPrefix(prefix...), WithMode(ModeExpectInSameLine))
	return w.String()
}

func (t *PointerDataType) End() token.Position {
	return t.DataType.End()
}

func (t *PointerDataType) RawText() string {
	return "*" + t.DataType.RawText()
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

type AnyDataType struct {
	Any *TokenNode
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

// BaseDataType is a common id type which contains bool, uint8, uint16, uint32,
// uint64, int8, int16, int32, int64, float32, float64, complex64, complex128,
// string, int, uint, uintptr, byte, rune, any.
type BaseDataType struct {
	Base *TokenNode
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

/*******************DataType End********************/
