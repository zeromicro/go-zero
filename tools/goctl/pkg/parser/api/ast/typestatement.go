package ast

import (
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

func (t *TypeLiteralStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, t.Type, t.Expr.Format(noLead))
	return w.String()
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

func (t *TypeGroupStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, t.Type, t.LParen)
	if len(t.ExprList) > 0 {
		w.NewLine()
	}

	for _, v := range t.ExprList {
		w.Writeln(v.Format(prefix + indent))
	}

	if len(t.ExprList) > 0 {
		w.Write(prefix, t.RParen)
	} else {
		w.Write(t.RParen)
	}
	return w.String()
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

func (e *TypeExpr) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, e.Name, e.Assign, e.DataType.Format(prefix))
	return w.String()
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

func (e *ElemExpr) Format(prefix string) string {
	w := NewWriter()
	w.Write(prefix)
	for idx, v := range e.Name {
		w.Write(v)
		if idx < len(e.Name)-1 {
			w.Write(", ")
		}
	}

	w.Write(whiteSpace, e.DataType.Format(prefix), whiteSpace, e.Tag)
	return w.String()
}

func (e *ElemExpr) exprNode() {}

/*******************Elem End********************/

/*******************ElemExprList Begin********************/

type ElemExprList []*ElemExpr

func (e ElemExprList) Format(prefix string) string {
	if len(e) == 0 {
		return ""
	}
	var list []ElemExprList
	var elems ElemExprList
	for _, v := range e {
		if v.DataType.ContainsStruct() {
			list = append(list, elems)
			list = append(list, ElemExprList{v})
			elems = ElemExprList{}
			continue
		}
		elems = append(elems, v)
	}
	if len(elems) > 0 {
		list = append(list, elems)
	}

	w := NewWriter()
	for _, v := range list {
		tw := w.UseTabWriter()
		for _, val := range v {
			var nameList []string
			for _, n := range val.Name {
				nameList = append(nameList, n.Text)
			}
			if _, ok := val.DataType.(*StructDataType); ok {
				w.WriteWithWhiteSpaceInfixln(prefix, strings.Join(nameList, ", "), val.DataType.Format(prefix), val.Tag)
			} else {
				tw.WriteWithInfixIndentln(prefix, strings.Join(nameList, ", "), val.DataType.Format(prefix), val.Tag)
			}
		}
		tw.Flush()
	}

	return w.String()
}

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
	return t.LBrace.Position
}

func (t *StructDataType) Format(prefix string) string {
	w := NewWriter()
	w.Write(t.LBrace)
	if len(t.Elements) > 0 {
		w.NewLine()
	}
	w.Write(t.Elements.Format(prefix + indent))
	if len(t.Elements) > 0 {
		w.Write(prefix, t.RBrace)
	} else {
		w.Write(t.RBrace)
	}
	return w.String()
}

func (t *StructDataType) exprNode()     {}
func (t *StructDataType) dataTypeNode() {}

type SliceDataType struct {
	LBrack   token.Token
	RBrack   token.Token
	DataType DataType
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
	return t.LBrack.Position
}

func (t *SliceDataType) Format(prefix string) string {
	w := NewWriter()
	w.Write(t.LBrack, t.RBrack, t.DataType.Format(prefix+indent))
	return w.String()
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
	return t.Format("")
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

func (t *MapDataType) Format(prefix string) string {
	w := NewWriter()
	w.Write(t.Map, t.LBrack, t.Key.Format(prefix), t.RBrack, t.Value.Format(prefix))
	return w.String()
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
	return t.Format("")
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

func (t *ArrayDataType) Format(prefix string) string {
	w := NewWriter()
	w.Write(t.LBrack, t.Length, t.RBrack, t.DataType.Format(prefix+indent))
	return w.String()
}

func (t *ArrayDataType) exprNode()     {}
func (t *ArrayDataType) dataTypeNode() {}

type InterfaceDataType struct {
	Interface token.Token
	LBrace    token.Token
	RBrace    token.Token
}

func (t *InterfaceDataType) RawText() string {
	return t.Format("")
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

func (t *InterfaceDataType) Format(_ string) string {
	w := NewWriter()
	w.Write(t.Interface, t.LBrace, t.RBrace)
	return w.String()
}

func (t *InterfaceDataType) exprNode() {}

func (t *InterfaceDataType) dataTypeNode() {}

type PointerDataType struct {
	Star     token.Token
	DataType DataType
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
	return t.Star.Position
}

func (t *PointerDataType) Format(prefix string) string {
	w := NewWriter()
	w.Write(t.Star, t.DataType.Format(prefix+indent))
	return w.String()
}

func (t *PointerDataType) exprNode()     {}
func (t *PointerDataType) dataTypeNode() {}

type AnyDataType struct {
	Any token.Token
}

func (t *AnyDataType) RawText() string {
	return t.Format("")
}

func (t *AnyDataType) ContainsStruct() bool {
	return false
}

func (t *AnyDataType) Format(_ string) string {
	w := NewWriter()
	w.Write(t.Any)
	return w.String()
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
	return t.Format("")
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

func (t *BaseDataType) Format(_ string) string {
	w := NewWriter()
	w.Write(t.Base)
	return w.String()
}

func (t *BaseDataType) exprNode()     {}
func (t *BaseDataType) dataTypeNode() {}

/*******************DataType End********************/
