package ast

import (
	"fmt"
	"sort"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

type (
	// TypeAlias、 TypeStruct
	TypeExpr interface {
		Doc() []Expr
		Format() error
		Equal(v interface{}) bool
		NameExpr() Expr
	}
	TypeAlias struct {
		Name        Expr
		Assign      Expr
		DataType    DataType
		DocExpr     []Expr
		CommentExpr Expr
	}

	TypeStruct struct {
		Name    Expr
		Struct  Expr
		LBrace  Expr
		RBrace  Expr
		DocExpr []Expr
		Fields  []*TypeField
	}

	TypeField struct {
		IsAnonymous bool
		// Name is nil if IsAnonymous
		Name        Expr
		DataType    DataType
		Tag         Expr
		DocExpr     []Expr
		CommentExpr Expr
	}

	// Literal, Interface, Map, Array, Time, Pointer
	DataType interface {
		Expr() Expr
		Equal(dt DataType) bool
		Format() error
		IsNotNil() bool
	}

	// int, bool, Foo,...
	Literal struct {
		Literal Expr
	}

	Interface struct {
		Literal Expr
	}

	Map struct {
		MapExpr Expr
		Map     Expr
		LBrack  Expr
		RBrack  Expr
		Key     Expr
		Value   DataType
	}

	Array struct {
		ArrayExpr Expr
		LBrack    Expr
		RBrack    Expr
		Literal   DataType
	}

	Time struct {
		Literal Expr
	}

	Pointer struct {
		PointerExpr Expr
		Star        Expr
		Name        Expr
	}
)

func (v *ApiVisitor) VisitTypeSpec(ctx *api.TypeSpecContext) interface{} {
	if ctx.TypeLit() != nil {
		return []TypeExpr{ctx.TypeLit().Accept(v).(TypeExpr)}
	}
	return ctx.TypeBlock().Accept(v)
}

func (v *ApiVisitor) VisitTypeLit(ctx *api.TypeLitContext) interface{} {
	typeLit := ctx.TypeLitBody().Accept(v)
	alias, ok := typeLit.(*TypeAlias)
	if ok {
		return alias
	}

	st, ok := typeLit.(*TypeStruct)
	if ok {
		return st
	}

	return typeLit
}

func (v *ApiVisitor) VisitTypeBlock(ctx *api.TypeBlockContext) interface{} {
	list := ctx.AllTypeBlockBody()
	var types []TypeExpr
	for _, each := range list {
		types = append(types, each.Accept(v).(TypeExpr))

	}
	return types
}

func (v *ApiVisitor) VisitTypeLitBody(ctx *api.TypeLitBodyContext) interface{} {
	if ctx.TypeAlias() != nil {
		return ctx.TypeAlias().Accept(v)
	}
	return ctx.TypeStruct().Accept(v)
}

func (v *ApiVisitor) VisitTypeBlockBody(ctx *api.TypeBlockBodyContext) interface{} {
	if ctx.TypeBlockAlias() != nil {
		return ctx.TypeBlockAlias().Accept(v).(*TypeAlias)
	}
	return ctx.TypeBlockStruct().Accept(v).(*TypeStruct)
}

func (v *ApiVisitor) VisitTypeStruct(ctx *api.TypeStructContext) interface{} {
	var st TypeStruct
	st.Name = v.newExprWithToken(ctx.GetStructName())
	v.exportCheck(st.Name)

	if util.UnExport(ctx.GetStructName().GetText()) {

	}
	if ctx.GetStructToken() != nil {
		structExpr := v.newExprWithToken(ctx.GetStructToken())
		structTokenText := ctx.GetStructToken().GetText()
		if structTokenText != "struct" {
			v.panic(structExpr, fmt.Sprintf("expecting 'struct', found input '%s'", structTokenText))
		}

		if api.IsGolangKeyWord(structTokenText, "struct") {
			v.panic(structExpr, fmt.Sprintf("expecting 'struct', but found golang keyword '%s'", structTokenText))
		}

		st.Struct = structExpr
	}

	st.LBrace = v.newExprWithToken(ctx.GetLbrace())
	st.RBrace = v.newExprWithToken(ctx.GetRbrace())
	fields := ctx.AllField()
	for _, each := range fields {
		f := each.Accept(v)
		if f == nil {
			continue
		}
		st.Fields = append(st.Fields, f.(*TypeField))
	}
	return &st
}

func (v *ApiVisitor) VisitTypeBlockStruct(ctx *api.TypeBlockStructContext) interface{} {
	var st TypeStruct
	st.Name = v.newExprWithToken(ctx.GetStructName())
	v.exportCheck(st.Name)

	if ctx.GetStructToken() != nil {
		structExpr := v.newExprWithToken(ctx.GetStructToken())
		structTokenText := ctx.GetStructToken().GetText()
		if structTokenText != "struct" {
			v.panic(structExpr, fmt.Sprintf("expecting 'struct', found imput '%s'", structTokenText))
		}

		if api.IsGolangKeyWord(structTokenText, "struct") {
			v.panic(structExpr, fmt.Sprintf("expecting 'struct', but found golang keyword '%s'", structTokenText))
		}

		st.Struct = structExpr
	}
	st.DocExpr = v.getDoc(ctx)
	st.LBrace = v.newExprWithToken(ctx.GetLbrace())
	st.RBrace = v.newExprWithToken(ctx.GetRbrace())
	fields := ctx.AllField()
	for _, each := range fields {
		f := each.Accept(v)
		if f == nil {
			continue
		}
		st.Fields = append(st.Fields, f.(*TypeField))
	}
	return &st
}

func (v *ApiVisitor) VisitTypeBlockAlias(ctx *api.TypeBlockAliasContext) interface{} {
	var alias TypeAlias
	alias.Name = v.newExprWithToken(ctx.GetAlias())
	alias.Assign = v.newExprWithToken(ctx.GetAssign())
	alias.DataType = ctx.DataType().Accept(v).(DataType)
	alias.DocExpr = v.getDoc(ctx)
	alias.CommentExpr = v.getComment(ctx)
	// todo: reopen if necessary
	v.panic(alias.Name, "unsupport alias")
	return &alias
}

func (v *ApiVisitor) VisitTypeAlias(ctx *api.TypeAliasContext) interface{} {
	var alias TypeAlias
	alias.Name = v.newExprWithToken(ctx.GetAlias())
	alias.Assign = v.newExprWithToken(ctx.GetAssign())
	alias.DataType = ctx.DataType().Accept(v).(DataType)
	alias.DocExpr = v.getDoc(ctx)
	alias.CommentExpr = v.getComment(ctx)
	// todo: reopen if necessary
	v.panic(alias.Name, "unsupport alias")
	return &alias
}

func (v *ApiVisitor) VisitField(ctx *api.FieldContext) interface{} {
	iAnonymousFiled := ctx.AnonymousFiled()
	iNormalFieldContext := ctx.NormalField()
	if iAnonymousFiled != nil {
		return iAnonymousFiled.Accept(v).(*TypeField)
	}
	if iNormalFieldContext != nil {
		return iNormalFieldContext.Accept(v).(*TypeField)
	}
	return nil
}

func (v *ApiVisitor) VisitNormalField(ctx *api.NormalFieldContext) interface{} {
	var field TypeField
	field.Name = v.newExprWithToken(ctx.GetFieldName())
	v.exportCheck(field.Name)

	iDataTypeContext := ctx.DataType()
	if iDataTypeContext != nil {
		field.DataType = iDataTypeContext.Accept(v).(DataType)
		field.CommentExpr = v.getComment(ctx)
	}
	if ctx.GetTag() != nil {
		tagText := ctx.GetTag().GetText()
		tagExpr := v.newExprWithToken(ctx.GetTag())
		if !api.MatchTag(tagText) {
			v.panic(tagExpr, fmt.Sprintf("mismatched tag, found input '%s'", tagText))
		}
		field.Tag = tagExpr
		field.CommentExpr = v.getComment(ctx)
	}
	field.DocExpr = v.getDoc(ctx)
	return &field
}

func (v *ApiVisitor) VisitAnonymousFiled(ctx *api.AnonymousFiledContext) interface{} {
	start := ctx.GetStart()
	stop := ctx.GetStop()
	var field TypeField
	field.IsAnonymous = true
	if ctx.GetStar() != nil {
		nameExpr := v.newExprWithTerminalNode(ctx.ID())
		v.exportCheck(nameExpr)
		field.DataType = &Pointer{
			PointerExpr: v.newExprWithText(ctx.GetStar().GetText()+ctx.ID().GetText(), start.GetLine(), start.GetColumn(), start.GetStart(), stop.GetStop()),
			Star:        v.newExprWithToken(ctx.GetStar()),
			Name:        nameExpr,
		}
	} else {
		nameExpr := v.newExprWithTerminalNode(ctx.ID())
		v.exportCheck(nameExpr)
		field.DataType = &Literal{Literal: nameExpr}
	}
	field.DocExpr = v.getDoc(ctx)
	field.CommentExpr = v.getComment(ctx)
	return &field
}

func (v *ApiVisitor) VisitDataType(ctx *api.DataTypeContext) interface{} {
	if ctx.ID() != nil {
		idExpr := v.newExprWithTerminalNode(ctx.ID())
		v.exportCheck(idExpr)
		return &Literal{Literal: idExpr}
	}
	if ctx.MapType() != nil {
		t := ctx.MapType().Accept(v)
		return t
	}
	if ctx.ArrayType() != nil {
		return ctx.ArrayType().Accept(v)
	}
	if ctx.GetInter() != nil {
		return &Interface{Literal: v.newExprWithToken(ctx.GetInter())}
	}
	if ctx.GetTime() != nil {
		// todo: reopen if it is necessary
		timeExpr := v.newExprWithToken(ctx.GetTime())
		v.panic(timeExpr, "unsupport time.Time")
		return &Time{Literal: timeExpr}
	}
	if ctx.PointerType() != nil {
		return ctx.PointerType().Accept(v)
	}
	return ctx.TypeStruct().Accept(v)
}

func (v *ApiVisitor) VisitPointerType(ctx *api.PointerTypeContext) interface{} {
	nameExpr := v.newExprWithTerminalNode(ctx.ID())
	v.exportCheck(nameExpr)
	return &Pointer{
		PointerExpr: v.newExprWithText(ctx.GetText(), ctx.GetStar().GetLine(), ctx.GetStar().GetColumn(), ctx.GetStar().GetStart(), ctx.ID().GetSymbol().GetStop()),
		Star:        v.newExprWithToken(ctx.GetStar()),
		Name:        nameExpr,
	}
}

func (v *ApiVisitor) VisitMapType(ctx *api.MapTypeContext) interface{} {
	return &Map{
		MapExpr: v.newExprWithText(ctx.GetText(), ctx.GetMapToken().GetLine(), ctx.GetMapToken().GetColumn(),
			ctx.GetMapToken().GetStart(), ctx.GetValue().GetStop().GetStop()),
		Map:    v.newExprWithToken(ctx.GetMapToken()),
		LBrack: v.newExprWithToken(ctx.GetLbrack()),
		RBrack: v.newExprWithToken(ctx.GetRbrack()),
		Key:    v.newExprWithToken(ctx.GetKey()),
		Value:  ctx.GetValue().Accept(v).(DataType),
	}
}

func (v *ApiVisitor) VisitArrayType(ctx *api.ArrayTypeContext) interface{} {
	return &Array{
		ArrayExpr: v.newExprWithText(ctx.GetText(), ctx.GetLbrack().GetLine(), ctx.GetLbrack().GetColumn(), ctx.GetLbrack().GetStart(), ctx.DataType().GetStop().GetStop()),
		LBrack:    v.newExprWithToken(ctx.GetLbrack()),
		RBrack:    v.newExprWithToken(ctx.GetRbrack()),
		Literal:   ctx.DataType().Accept(v).(DataType),
	}
}

func (a *TypeAlias) NameExpr() Expr {
	return a.Name
}

func (a *TypeAlias) Doc() []Expr {
	return a.DocExpr
}

func (a *TypeAlias) Comment() Expr {
	return a.CommentExpr
}

func (a *TypeAlias) Format() error {
	return nil
}

func (a *TypeAlias) Equal(v interface{}) bool {
	if v == nil {
		return false
	}

	alias := v.(*TypeAlias)
	if !a.Name.Equal(alias.Name) {
		return false
	}

	if !a.Assign.Equal(alias.Assign) {
		return false
	}

	if !a.DataType.Equal(alias.DataType) {
		return false
	}

	return EqualDoc(a, alias)
}

func (l *Literal) Expr() Expr {
	return l.Literal
}

func (l *Literal) Format() error {
	// todo
	return nil
}

func (l *Literal) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Literal)
	if !ok {
		return false
	}

	return l.Literal.Equal(v.Literal)
}

func (l *Literal) IsNotNil() bool {
	return l != nil
}

func (i *Interface) Expr() Expr {
	return i.Literal
}

func (i *Interface) Format() error {
	// todo
	return nil
}

func (i *Interface) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Interface)
	if !ok {
		return false
	}

	return i.Literal.Equal(v.Literal)
}

func (i *Interface) IsNotNil() bool {
	return i != nil
}

func (m *Map) Expr() Expr {
	return m.MapExpr
}

func (m *Map) Format() error {
	// todo
	return nil
}

func (m *Map) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Map)
	if !ok {
		return false
	}

	if !m.Key.Equal(v.Key) {
		return false
	}

	if !m.Value.Equal(v.Value) {
		return false
	}

	if !m.MapExpr.Equal(v.MapExpr) {
		return false
	}

	return m.Map.Equal(v.Map)
}

func (m *Map) IsNotNil() bool {
	return m != nil
}

func (a *Array) Expr() Expr {
	return a.ArrayExpr
}

func (a *Array) Format() error {
	// todo
	return nil
}

func (a *Array) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Array)
	if !ok {
		return false
	}

	if !a.ArrayExpr.Equal(v.ArrayExpr) {
		return false
	}

	return a.Literal.Equal(v.Literal)
}

func (a *Array) IsNotNil() bool {
	return a != nil
}

func (t *Time) Expr() Expr {
	return t.Literal
}

func (t *Time) Format() error {
	// todo
	return nil
}

func (t *Time) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Time)
	if !ok {
		return false
	}

	return t.Literal.Equal(v.Literal)
}

func (t *Time) IsNotNil() bool {
	return t != nil
}

func (p *Pointer) Expr() Expr {
	return p.PointerExpr
}

func (p *Pointer) Format() error {
	return nil
}

func (p *Pointer) Equal(dt DataType) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*Pointer)
	if !ok {
		return false
	}

	if !p.PointerExpr.Equal(v.PointerExpr) {
		return false
	}

	if !p.Star.Equal(v.Star) {
		return false
	}

	return p.Name.Equal(v.Name)
}

func (p *Pointer) IsNotNil() bool {
	return p != nil
}

func (s *TypeStruct) NameExpr() Expr {
	return s.Name
}

func (s *TypeStruct) Equal(dt interface{}) bool {
	if dt == nil {
		return false
	}

	v, ok := dt.(*TypeStruct)
	if !ok {
		return false
	}

	if !s.Name.Equal(v.Name) {
		return false
	}

	var expectDoc, actualDoc []Expr
	expectDoc = append(expectDoc, s.DocExpr...)
	actualDoc = append(actualDoc, v.DocExpr...)
	sort.Slice(expectDoc, func(i, j int) bool {
		return expectDoc[i].Line() < expectDoc[j].Line()
	})

	for index, each := range actualDoc {
		if !each.Equal(actualDoc[index]) {
			return false
		}
	}

	if s.Struct != nil {
		if s.Struct != nil {
			if !s.Struct.Equal(v.Struct) {
				return false
			}
		}
	}

	if len(s.Fields) != len(v.Fields) {
		return false
	}

	var expected, acual []*TypeField
	expected = append(expected, s.Fields...)
	acual = append(acual, v.Fields...)

	sort.Slice(expected, func(i, j int) bool {
		return expected[i].DataType.Expr().Line() < expected[j].DataType.Expr().Line()
	})
	sort.Slice(acual, func(i, j int) bool {
		return acual[i].DataType.Expr().Line() < acual[j].DataType.Expr().Line()
	})

	for index, each := range expected {
		ac := acual[index]
		if !each.Equal(ac) {
			return false
		}
	}

	return true
}

func (s *TypeStruct) Doc() []Expr {
	return s.DocExpr
}

func (s *TypeStruct) Format() error {
	// todo
	return nil
}

func (t *TypeField) Equal(v interface{}) bool {
	if v == nil {
		return false
	}

	f, ok := v.(*TypeField)
	if !ok {
		return false
	}

	if t.IsAnonymous != f.IsAnonymous {
		return false
	}

	if !t.DataType.Equal(f.DataType) {
		return false
	}

	if !t.IsAnonymous {
		if !t.Name.Equal(f.Name) {
			return false
		}

		if t.Tag != nil {
			if !t.Tag.Equal(f.Tag) {
				return false
			}
		}
	}

	return EqualDoc(t, f)
}

func (t *TypeField) Doc() []Expr {
	return t.DocExpr
}

func (t *TypeField) Comment() Expr {
	return t.CommentExpr
}

func (t *TypeField) Format() error {
	// todo
	return nil
}
