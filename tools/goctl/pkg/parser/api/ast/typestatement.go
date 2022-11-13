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
	Type token.Token
	Expr *TypeExpr
}

func (t *TypeLiteralStmt) Pos() token.Position {
	return t.Type.Position
}

func (t *TypeLiteralStmt) stmtNode() {}
func (t *TypeLiteralStmt) typeNode() {}

type TypeGroupStmt struct {
	Type     token.Token
	LParen   token.Token
	ExprList []*TypeExpr
	RParen   token.Token
}

func (t *TypeGroupStmt) Pos() token.Position {
	return t.Type.Position
}

func (t *TypeGroupStmt) stmtNode() {}
func (t *TypeGroupStmt) typeNode() {}

/*******************TypeStmt End********************/

/*******************TypeExpr Begin********************/

type TypeExpr struct {
	Name     token.Token
	Assign   token.Token
	DataType DataType
}

func (e *TypeExpr) Pos() token.Position {
	return e.Name.Position
}

func (e *TypeExpr) exprNode() {}

func (e *TypeExpr) isStruct() bool {
	return e.DataType.ContainsStruct()
}

/*******************TypeExpr Begin********************/

/*******************Elem Begin********************/

type ElemExpr struct {
	Name     []token.Token
	DataType DataType
	Tag      token.Token
}

func (e *ElemExpr) Pos() token.Position {
	if len(e.Name) > 0 {
		return e.Name[0].Position
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
	LBrace   token.Token
	Elements ElemExprList
	RBrace   token.Token
}

func (t *StructDataType) RawText() string {
	b := bytes.NewBuffer(nil)
	b.WriteRune('{')
	for _, v := range t.Elements {
		b.WriteRune('\n')
		var nameList []string
		for _, n := range v.Name {
			nameList = append(nameList, n.Text)
		}
		b.WriteString(fmt.Sprintf("%s %s %s", strings.Join(nameList, ", "), v.DataType.RawText(), v.Tag.Text))
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
	return t.LBrace.Position
}

func (t *StructDataType) exprNode()     {}
func (t *StructDataType) dataTypeNode() {}

type SliceDataType struct {
	LBrack   token.Token
	RBrack   token.Token
	DataType DataType
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
	return t.LBrack.Position
}

func (t *SliceDataType) exprNode()     {}
func (t *SliceDataType) dataTypeNode() {}

type MapDataType struct {
	Map    token.Token
	LBrack token.Token
	Key    DataType
	RBrack token.Token
	Value  DataType
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
	return t.Map.Position
}

func (t *MapDataType) exprNode()     {}
func (t *MapDataType) dataTypeNode() {}

type ArrayDataType struct {
	LBrack   token.Token
	Length   token.Token
	RBrack   token.Token
	DataType DataType
}

func (t *ArrayDataType) RawText() string {
	return fmt.Sprintf("[%s]%s", t.Length.Text, t.DataType.RawText())
}

func (t *ArrayDataType) ContainsStruct() bool {
	return t.DataType.ContainsStruct()
}

func (t *ArrayDataType) CanEqual() bool {
	return t.DataType.CanEqual()
}

func (t *ArrayDataType) Pos() token.Position {
	return t.LBrack.Position
}

func (t *ArrayDataType) exprNode()     {}
func (t *ArrayDataType) dataTypeNode() {}

type InterfaceDataType struct {
	Interface token.Token
}

func (t *InterfaceDataType) RawText() string {
	return t.Interface.Text
}

func (t *InterfaceDataType) ContainsStruct() bool {
	return false
}

func (t *InterfaceDataType) CanEqual() bool {
	return true
}

func (t *InterfaceDataType) Pos() token.Position {
	return t.Interface.Position
}

func (t *InterfaceDataType) exprNode() {}

func (t *InterfaceDataType) dataTypeNode() {}

type PointerDataType struct {
	Star     token.Token
	DataType DataType
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
	return t.Star.Position
}

func (t *PointerDataType) exprNode()     {}
func (t *PointerDataType) dataTypeNode() {}

type AnyDataType struct {
	Any token.Token
}

func (t *AnyDataType) RawText() string {
	return t.Any.Text
}

func (t *AnyDataType) ContainsStruct() bool {
	return false
}

func (t *AnyDataType) Pos() token.Position {
	return t.Any.Position
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
	Base token.Token
}

func (t *BaseDataType) RawText() string {
	return t.Base.Text
}

func (t *BaseDataType) ContainsStruct() bool {
	return false
}

func (t *BaseDataType) CanEqual() bool {
	return true
}

func (t *BaseDataType) Pos() token.Position {
	return t.Base.Position
}

func (t *BaseDataType) exprNode()     {}
func (t *BaseDataType) dataTypeNode() {}

/*******************DataType End********************/
