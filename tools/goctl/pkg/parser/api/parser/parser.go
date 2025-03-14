package parser

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/scanner"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const (
	idAPI          = "api"
	groupKeyText   = "group"
	infoTitleKey   = "Title"
	infoDescKey    = "Desc"
	infoVersionKey = "Version"
	infoAuthorKey  = "Author"
	infoEmailKey   = "Email"
)

// Parser is the parser for api file.
type Parser struct {
	s      *scanner.Scanner
	errors []error

	curTok  token.Token
	peekTok token.Token

	headCommentGroup ast.CommentGroup
	api              *ast.AST
	node             map[token.Token]*ast.TokenNode
}

// New creates a new parser.
func New(filename string, src interface{}) *Parser {
	abs, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalln(err)
	}

	p := &Parser{
		s:    scanner.MustNewScanner(abs, src),
		api:  &ast.AST{Filename: abs},
		node: make(map[token.Token]*ast.TokenNode),
	}

	return p
}

// Parse parses the api file.
func (p *Parser) Parse() *ast.AST {
	if !p.init() {
		return nil
	}

	for p.curTokenIsNotEof() {
		stmt := p.parseStmt()
		if isNil(stmt) {
			return nil
		}

		p.appendStmt(stmt)
		if !p.nextToken() {
			return nil
		}
	}

	return p.api
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.curTok.Type {
	case token.IDENT:
		switch {
		case p.curTok.Is(token.Syntax):
			return p.parseSyntaxStmt()
		case p.curTok.Is(token.Info):
			return p.parseInfoStmt()
		case p.curTok.Is(token.Service):
			return p.parseService()
		case p.curTok.Is(token.TypeKeyword):
			return p.parseTypeStmt()
		case p.curTok.Is(token.ImportKeyword):
			return p.parseImportStmt()
		default:
			p.expectIdentError(p.curTok, token.Syntax, token.Info, token.Service, token.TYPE)
			return nil
		}
	case token.AT_SERVER:
		return p.parseService()
	default:
		p.errors = append(p.errors, fmt.Errorf("%s unexpected token '%s'", p.curTok.Position.String(), p.peekTok.Text))
		return nil
	}
}

func (p *Parser) parseService() *ast.ServiceStmt {
	var stmt = &ast.ServiceStmt{}
	if p.curTokenIs(token.AT_SERVER) {
		atServerStmt := p.parseAtServerStmt()
		if atServerStmt == nil {
			return nil
		}

		stmt.AtServerStmt = atServerStmt
		if !p.advanceIfPeekTokenIs(token.Service) {
			return nil
		}
	}
	stmt.Service = p.curTokenNode()

	if !p.advanceIfPeekTokenIs(token.IDENT) {
		return nil
	}

	// service name expr
	nameExpr := p.parseServiceNameExpr()
	if nameExpr == nil {
		return nil
	}

	stmt.Name = nameExpr

	// token '{'
	if !p.advanceIfPeekTokenIs(token.LBRACE) {
		return nil
	}

	stmt.LBrace = p.curTokenNode()

	// service item statements
	routes := p.parseServiceItemsStmt()
	if routes == nil {
		return nil
	}

	stmt.Routes = routes

	// token '}'
	if !p.advanceIfPeekTokenIs(token.RBRACE) {
		return nil
	}

	stmt.RBrace = p.curTokenNode()

	return stmt
}

func (p *Parser) parseServiceItemsStmt() []*ast.ServiceItemStmt {
	var stmt = make([]*ast.ServiceItemStmt, 0)
	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RBRACE) {
		item := p.parseServiceItemStmt()
		if item == nil {
			return nil
		}

		stmt = append(stmt, item)
		if p.peekTokenIs(token.RBRACE) {
			break
		}

		if p.notExpectPeekToken(token.AT_DOC, token.AT_HANDLER, token.RBRACE) {
			return nil
		}
	}

	return stmt
}

func (p *Parser) parseServiceItemStmt() *ast.ServiceItemStmt {
	var stmt = &ast.ServiceItemStmt{}
	// statement @doc
	if p.peekTokenIs(token.AT_DOC) {
		if !p.nextToken() {
			return nil
		}

		atDocStmt := p.parseAtDocStmt()
		if atDocStmt == nil {
			return nil
		}

		stmt.AtDoc = atDocStmt
	}

	// statement @handler
	if !p.advanceIfPeekTokenIs(token.AT_HANDLER, token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.RBRACE) {
		return stmt
	}

	atHandlerStmt := p.parseAtHandlerStmt()
	if atHandlerStmt == nil {
		return nil
	}

	stmt.AtHandler = atHandlerStmt

	// statement route
	route := p.parseRouteStmt()
	if route == nil {
		return nil
	}

	stmt.Route = route

	return stmt
}

func (p *Parser) parseRouteStmt() *ast.RouteStmt {
	var stmt = &ast.RouteStmt{}
	// token http method
	if !p.advanceIfPeekTokenIs(token.HttpMethods...) {
		return nil
	}

	stmt.Method = p.curTokenNode()

	// path expression
	pathExpr := p.parsePathExpr()
	if pathExpr == nil {
		return nil
	}

	stmt.Path = pathExpr

	if p.peekTokenIs(token.AT_DOC, token.AT_HANDLER, token.RBRACE) {
		return stmt
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	if p.notExpectPeekToken(token.Returns, token.LPAREN) {
		return nil
	}

	if p.peekTokenIs(token.LPAREN) {
		// request expression
		requestBodyStmt := p.parseBodyStmt()
		if requestBodyStmt == nil {
			return nil
		}

		stmt.Request = requestBodyStmt
	}

	if p.notExpectPeekToken(token.Returns, token.AT_DOC, token.AT_HANDLER, token.RBRACE, token.SEMICOLON) {
		return nil
	}

	// token 'returns'
	if p.peekTokenIs(token.Returns) {
		if !p.nextToken() {
			return nil
		}

		stmt.Returns = p.curTokenNode()
		responseBodyStmt := p.parseBodyStmt()
		if responseBodyStmt == nil {
			return nil
		}

		stmt.Response = responseBodyStmt
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBodyStmt() *ast.BodyStmt {
	var stmt = &ast.BodyStmt{}
	// token '('
	if !p.advanceIfPeekTokenIs(token.LPAREN) {
		return nil
	}

	stmt.LParen = p.curTokenNode()
	// token ')'
	if p.peekTokenIs(token.RPAREN) {
		if !p.nextToken() {
			return nil
		}

		stmt.RParen = p.curTokenNode()
		return stmt
	}

	expr := p.parseBodyExpr()
	if expr == nil {
		return nil
	}

	stmt.Body = expr

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseBodyExpr() *ast.BodyExpr {
	var expr = &ast.BodyExpr{}
	switch {
	case p.peekTokenIs(token.LBRACK): // token '['
		if !p.nextToken() {
			return nil
		}

		expr.LBrack = p.curTokenNode()

		// token ']'
		if !p.advanceIfPeekTokenIs(token.RBRACK) {
			return nil
		}

		expr.RBrack = p.curTokenNode()

		switch {
		case p.peekTokenIs(token.MUL):
			if !p.nextToken() {
				return nil
			}

			expr.Star = p.curTokenNode()
			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}

			expr.Value = p.curTokenNode()
			return expr
		case p.peekTokenIs(token.IDENT):
			if !p.nextToken() {
				return nil
			}

			expr.Value = p.curTokenNode()
			return expr
		default:
			p.expectPeekToken(token.MUL, token.IDENT)
			return nil
		}
	case p.peekTokenIs(token.MUL):
		if !p.nextToken() {
			return nil
		}

		expr.Star = p.curTokenNode()
		if !p.advanceIfPeekTokenIs(token.IDENT) {
			return nil
		}

		expr.Value = p.curTokenNode()
		return expr
	case p.peekTokenIs(token.IDENT):
		if !p.nextToken() {
			return nil
		}

		expr.Value = p.curTokenNode()
		return expr
	default:
		p.expectPeekToken(token.LBRACK, token.MUL, token.IDENT)
		return nil
	}
}

func (p *Parser) parsePathExpr() *ast.PathExpr {
	var expr = &ast.PathExpr{}

	var values []token.Token
	for p.curTokenIsNotEof() &&
		p.peekTokenIsNot(token.LPAREN, token.Returns, token.AT_DOC, token.AT_HANDLER, token.SEMICOLON, token.RBRACE) {
		// token '/'
		if !p.advanceIfPeekTokenIs(token.QUO) {
			return nil
		}

		values = append(values, p.curTok)
		if p.peekTokenIs(token.LPAREN, token.Returns, token.AT_DOC, token.AT_HANDLER, token.SEMICOLON, token.RBRACE) {
			break
		}
		if p.notExpectPeekTokenGotComment(p.curTokenNode().PeekFirstLeadingComment(), token.COLON, token.IDENT, token.INT) {
			return nil
		}

		// token ':' or IDENT
		if p.notExpectPeekToken(token.COLON, token.IDENT, token.INT) {
			return nil
		}

		if p.notExpectPeekTokenGotComment(p.curTokenNode().PeekFirstLeadingComment(), token.COLON) {
			return nil
		}

		// token ':'
		if p.peekTokenIs(token.COLON) {
			if !p.nextToken() {
				return nil
			}

			values = append(values, p.curTok)
		}

		// path id tokens
		pathTokens := p.parsePathItem()
		if pathTokens == nil {
			return nil
		}

		values = append(values, pathTokens...)
		if p.notExpectPeekToken(token.QUO, token.LPAREN, token.Returns, token.AT_DOC, token.AT_HANDLER, token.SEMICOLON, token.RBRACE) {
			return nil
		}
	}

	var textList []string
	for _, v := range values {
		textList = append(textList, v.Text)
	}

	node := ast.NewTokenNode(token.Token{
		Type:     token.PATH,
		Text:     strings.Join(textList, ""),
		Position: values[0].Position,
	})
	node.SetLeadingCommentGroup(p.curTokenNode().LeadingCommentGroup)
	expr.Value = node

	return expr
}

func (p *Parser) parsePathItem() []token.Token {
	var list []token.Token
	if !p.advanceIfPeekTokenIs(token.IDENT, token.INT) {
		return nil
	}
	list = append(list, p.curTok)

	for p.curTokenIsNotEof() &&
		p.peekTokenIsNot(token.QUO, token.LPAREN, token.Returns, token.AT_DOC, token.AT_HANDLER, token.RBRACE, token.SEMICOLON, token.EOF) {
		if p.peekTokenIs(token.SUB) {
			if !p.nextToken() {
				return nil
			}

			list = append(list, p.curTok)

			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}
			list = append(list, p.curTok)
		} else {
			if p.peekTokenIs(token.LPAREN, token.Returns, token.AT_DOC, token.AT_HANDLER, token.SEMICOLON, token.RBRACE) {
				return list
			}

			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}

			list = append(list, p.curTok)
		}
	}

	return list
}

func (p *Parser) parseServiceNameExpr() *ast.ServiceNameExpr {
	var expr = &ast.ServiceNameExpr{}
	var text = p.curTok.Text

	pos := p.curTok.Position
	if p.peekTokenIs(token.SUB) {
		if !p.nextToken() {
			return nil
		}

		text += p.curTok.Text
		if !p.expectPeekToken(idAPI) {
			return nil
		}

		if !p.nextToken() {
			return nil
		}

		text += p.curTok.Text
	}

	node := ast.NewTokenNode(token.Token{
		Type:     token.IDENT,
		Text:     text,
		Position: pos,
	})
	node.SetLeadingCommentGroup(p.curTokenNode().LeadingCommentGroup)
	expr.Name = node

	return expr
}

func (p *Parser) parseAtDocStmt() ast.AtDocStmt {
	if p.notExpectPeekToken(token.LPAREN, token.STRING) {
		return nil
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseAtDocGroupStmt()
	}

	return p.parseAtDocLiteralStmt()
}

func (p *Parser) parseAtDocGroupStmt() ast.AtDocStmt {
	var stmt = &ast.AtDocGroupStmt{}
	stmt.AtDoc = p.curTokenNode()

	// token '('
	if !p.advanceIfPeekTokenIs(token.LPAREN) {
		return nil
	}

	stmt.LParen = p.curTokenNode()

	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RPAREN) {
		expr := p.parseKVExpression()
		if expr == nil {
			return nil
		}

		stmt.Values = append(stmt.Values, expr)
		if p.notExpectPeekToken(token.RPAREN, token.IDENT) {
			return nil
		}
	}

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseAtDocLiteralStmt() ast.AtDocStmt {
	var stmt = &ast.AtDocLiteralStmt{}
	stmt.AtDoc = p.curTokenNode()

	if !p.advanceIfPeekTokenIs(token.STRING) {
		return nil
	}

	stmt.Value = p.curTokenNode()

	return stmt
}

func (p *Parser) parseAtHandlerStmt() *ast.AtHandlerStmt {
	var stmt = &ast.AtHandlerStmt{}
	stmt.AtHandler = p.curTokenNode()

	// token IDENT
	if !p.advanceIfPeekTokenIs(token.IDENT) {
		return nil
	}

	stmt.Name = p.curTokenNode()

	return stmt
}

func (p *Parser) parseAtServerStmt() *ast.AtServerStmt {
	var stmt = &ast.AtServerStmt{}
	stmt.AtServer = p.curTokenNode()

	// token '('
	if !p.advanceIfPeekTokenIs(token.LPAREN) {
		return nil
	}

	stmt.LParen = p.curTokenNode()

	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RPAREN) {
		expr := p.parseAtServerKVExpression()
		if expr == nil {
			return nil
		}

		stmt.Values = append(stmt.Values, expr)
		if p.notExpectPeekToken(token.RPAREN, token.IDENT) {
			return nil
		}
	}

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseTypeStmt() ast.TypeStmt {
	switch {
	case p.peekTokenIs(token.LPAREN):
		return p.parseTypeGroupStmt()
	case p.peekTokenIs(token.IDENT):
		return p.parseTypeLiteralStmt()
	default:
		p.expectPeekToken(token.LPAREN, token.IDENT)
		return nil
	}
}

func (p *Parser) parseTypeLiteralStmt() ast.TypeStmt {
	var stmt = &ast.TypeLiteralStmt{}
	stmt.Type = p.curTokenNode()

	expr := p.parseTypeExpr()
	if expr == nil {
		return nil
	}

	stmt.Expr = expr

	return stmt
}

func (p *Parser) parseTypeGroupStmt() ast.TypeStmt {
	var stmt = &ast.TypeGroupStmt{}
	stmt.Type = p.curTokenNode()

	// token '('
	if !p.nextToken() {
		return nil
	}

	stmt.LParen = p.curTokenNode()

	exprList := p.parseTypeExprList()
	if exprList == nil {
		return nil
	}

	stmt.ExprList = exprList

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseTypeExprList() []*ast.TypeExpr {
	if !p.expectPeekToken(token.IDENT, token.RPAREN) {
		return nil
	}

	var exprList = make([]*ast.TypeExpr, 0)
	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RPAREN, token.EOF) {
		expr := p.parseTypeExpr()
		if expr == nil {
			return nil
		}

		exprList = append(exprList, expr)
		if !p.expectPeekToken(token.IDENT, token.RPAREN) {
			return nil
		}
	}

	return exprList
}

func (p *Parser) parseTypeExpr() *ast.TypeExpr {
	var expr = &ast.TypeExpr{}
	// token IDENT
	if !p.advanceIfPeekTokenIs(token.IDENT) {
		return nil
	}

	if p.curTokenIsKeyword() {
		return nil
	}

	expr.Name = p.curTokenNode()

	// token '='
	if p.peekTokenIs(token.ASSIGN) {
		if !p.nextToken() {
			return nil
		}

		expr.Assign = p.curTokenNode()
	}

	dt := p.parseDataType()
	if isNil(dt) {
		return nil
	}

	expr.DataType = dt

	return expr
}

func (p *Parser) parseDataType() ast.DataType {
	switch {
	case p.peekTokenIs(token.Any):
		return p.parseAnyDataType()
	case p.peekTokenIs(token.LBRACE):
		return p.parseStructDataType()
	case p.peekTokenIs(token.IDENT):
		if p.peekTokenIs(token.MapKeyword) {
			return p.parseMapDataType()
		}
		if !p.nextToken() {
			return nil
		}

		if p.curTokenIsKeyword() {
			return nil
		}

		node := p.curTokenNode()
		baseDT := &ast.BaseDataType{Base: node}

		return baseDT
	case p.peekTokenIs(token.LBRACK):
		if !p.nextToken() {
			return nil
		}

		switch {
		case p.peekTokenIs(token.RBRACK):
			return p.parseSliceDataType()
		case p.peekTokenIs(token.INT, token.ELLIPSIS):
			return p.parseArrayDataType()
		default:
			p.expectPeekToken(token.RBRACK, token.INT, token.ELLIPSIS)
			return nil
		}
	case p.peekTokenIs(token.ANY):
		return p.parseInterfaceDataType()
	case p.peekTokenIs(token.MUL):
		return p.parsePointerDataType()
	default:
		p.expectPeekToken(token.IDENT, token.LBRACK, token.ANY, token.MUL, token.LBRACE)
		return nil
	}
}
func (p *Parser) parseStructDataType() *ast.StructDataType {
	var tp = &ast.StructDataType{}
	if !p.nextToken() {
		return nil
	}

	tp.LBrace = p.curTokenNode()

	if p.notExpectPeekToken(token.IDENT, token.MUL, token.RBRACE) {
		return nil
	}

	// ElemExprList
	elems := p.parseElemExprList()
	if elems == nil {
		return nil
	}

	tp.Elements = elems

	if !p.advanceIfPeekTokenIs(token.RBRACE) {
		return nil
	}

	tp.RBrace = p.curTokenNode()

	return tp
}

func (p *Parser) parseElemExprList() ast.ElemExprList {
	var list = make(ast.ElemExprList, 0)
	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RBRACE, token.EOF) {
		if p.notExpectPeekToken(token.IDENT, token.MUL, token.RBRACE) {
			return nil
		}

		expr := p.parseElemExpr()
		if expr == nil {
			return nil
		}

		list = append(list, expr)
		if p.notExpectPeekToken(token.IDENT, token.MUL, token.RBRACE) {
			return nil
		}
	}

	return list
}

func (p *Parser) parseElemExpr() *ast.ElemExpr {
	var expr = &ast.ElemExpr{}
	if !p.advanceIfPeekTokenIs(token.IDENT, token.MUL) {
		return nil
	}

	if p.curTokenIsKeyword() {
		return nil
	}

	identNode := p.curTokenNode()
	if p.curTokenIs(token.MUL) {
		star := p.curTokenNode()
		if !p.advanceIfPeekTokenIs(token.IDENT) {
			return nil
		}

		var dt ast.DataType
		if p.curTokenIs(token.Any) {
			dt = &ast.AnyDataType{Any: p.curTokenNode()}
		} else {
			dt = &ast.BaseDataType{Base: p.curTokenNode()}
		}
		expr.DataType = &ast.PointerDataType{
			Star:     star,
			DataType: dt,
		}
	} else if p.peekTok.Line() > identNode.Token.Line() || p.peekTokenIs(token.RAW_STRING) {
		if p.curTokenIs(token.Any) {
			expr.DataType = &ast.AnyDataType{Any: identNode}
		} else {
			expr.DataType = &ast.BaseDataType{Base: identNode}
		}
	} else {
		expr.Name = append(expr.Name, identNode)
		if p.notExpectPeekToken(token.COMMA, token.IDENT, token.LBRACK, token.ANY, token.MUL, token.LBRACE) {
			return nil
		}

		for p.peekTokenIs(token.COMMA) {
			if !p.nextToken() {
				return nil
			}

			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}

			if p.curTokenIsKeyword() {
				return nil
			}

			expr.Name = append(expr.Name, p.curTokenNode())
		}

		dt := p.parseDataType()
		if isNil(dt) {
			return nil
		}

		expr.DataType = dt
	}

	if p.notExpectPeekToken(token.RAW_STRING, token.MUL, token.IDENT, token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.RAW_STRING) {
		if !p.nextToken() {
			return nil
		}

		expr.Tag = p.curTokenNode()
	}

	return expr
}

func (p *Parser) parseAnyDataType() *ast.AnyDataType {
	var tp = &ast.AnyDataType{}
	if !p.nextToken() {
		return nil
	}

	tp.Any = p.curTokenNode()

	return tp
}

func (p *Parser) parsePointerDataType() *ast.PointerDataType {
	var tp = &ast.PointerDataType{}
	if !p.nextToken() {
		return nil
	}

	tp.Star = p.curTokenNode()

	if p.notExpectPeekToken(token.IDENT, token.LBRACK, token.ANY, token.MUL) {
		return nil
	}

	// DataType
	dt := p.parseDataType()
	if isNil(dt) {
		return nil
	}

	tp.DataType = dt

	return tp
}

func (p *Parser) parseInterfaceDataType() *ast.InterfaceDataType {
	var tp = &ast.InterfaceDataType{}
	if !p.nextToken() {
		return nil
	}

	tp.Interface = p.curTokenNode()

	return tp
}

func (p *Parser) parseMapDataType() *ast.MapDataType {
	var tp = &ast.MapDataType{}
	if !p.nextToken() {
		return nil
	}

	tp.Map = p.curTokenNode()

	// token '['
	if !p.advanceIfPeekTokenIs(token.LBRACK) {
		return nil
	}

	tp.LBrack = p.curTokenNode()

	// DataType
	dt := p.parseDataType()
	if isNil(dt) {
		return nil
	}

	tp.Key = dt

	// token  ']'
	if !p.advanceIfPeekTokenIs(token.RBRACK) {
		return nil
	}

	tp.RBrack = p.curTokenNode()

	// DataType
	dt = p.parseDataType()
	if isNil(dt) {
		return nil
	}

	tp.Value = dt

	return tp
}

func (p *Parser) parseArrayDataType() *ast.ArrayDataType {
	var tp = &ast.ArrayDataType{}
	tp.LBrack = p.curTokenNode()

	// token INT | ELLIPSIS
	if !p.nextToken() {
		return nil
	}

	tp.Length = p.curTokenNode()

	// token ']'
	if !p.advanceIfPeekTokenIs(token.RBRACK) {
		return nil
	}

	tp.RBrack = p.curTokenNode()

	// DataType
	dt := p.parseDataType()
	if isNil(dt) {
		return nil
	}

	tp.DataType = dt

	return tp
}

func (p *Parser) parseSliceDataType() *ast.SliceDataType {
	var tp = &ast.SliceDataType{}
	tp.LBrack = p.curTokenNode()

	// token ']'
	if !p.advanceIfPeekTokenIs(token.RBRACK) {
		return nil
	}

	tp.RBrack = p.curTokenNode()

	// DataType
	dt := p.parseDataType()
	if isNil(dt) {
		return nil
	}

	tp.DataType = dt

	return tp
}

func (p *Parser) parseImportStmt() ast.ImportStmt {
	if p.notExpectPeekToken(token.LPAREN, token.STRING) {
		return nil
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseImportGroupStmt()
	}

	return p.parseImportLiteralStmt()
}

func (p *Parser) parseImportLiteralStmt() ast.ImportStmt {
	var stmt = &ast.ImportLiteralStmt{}
	stmt.Import = p.curTokenNode()

	// token STRING
	if !p.advanceIfPeekTokenIs(token.STRING) {
		return nil
	}

	stmt.Value = p.curTokenNode()

	return stmt
}

func (p *Parser) parseImportGroupStmt() ast.ImportStmt {
	var stmt = &ast.ImportGroupStmt{}
	stmt.Import = p.curTokenNode()

	// token '('
	if !p.advanceIfPeekTokenIs(token.LPAREN) { // assert: dead code
		return nil
	}

	stmt.LParen = p.curTokenNode()

	// token STRING
	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RPAREN) {
		if !p.advanceIfPeekTokenIs(token.STRING) {
			return nil
		}

		stmt.Values = append(stmt.Values, p.curTokenNode())

		if p.notExpectPeekToken(token.RPAREN, token.STRING) {
			return nil
		}
	}

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseInfoStmt() *ast.InfoStmt {
	var stmt = &ast.InfoStmt{}
	stmt.Info = p.curTokenNode()

	// token '('
	if !p.advanceIfPeekTokenIs(token.LPAREN) {
		return nil
	}

	stmt.LParen = p.curTokenNode()

	for p.curTokenIsNotEof() && p.peekTokenIsNot(token.RPAREN) {
		expr := p.parseKVExpression()
		if expr == nil {
			return nil
		}

		stmt.Values = append(stmt.Values, expr)
		if p.notExpectPeekToken(token.RPAREN, token.IDENT) {
			return nil
		}
	}

	// token ')'
	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		return nil
	}

	stmt.RParen = p.curTokenNode()

	return stmt
}

func (p *Parser) parseAtServerKVExpression() *ast.KVExpr {
	var expr = &ast.KVExpr{}

	// token IDENT
	if !p.advanceIfPeekTokenIs(token.IDENT, token.RPAREN) {
		return nil
	}

	expr.Key = p.curTokenNode()

	if !p.advanceIfPeekTokenIs(token.COLON) {
		return nil
	}
	expr.Colon = p.curTokenNode()

	var valueTok token.Token
	var leadingCommentGroup ast.CommentGroup
	if p.notExpectPeekToken(token.QUO, token.DURATION, token.IDENT, token.INT, token.STRING) {
		return nil
	}

	if p.peekTokenIs(token.QUO) {
		if !p.nextToken() {
			return nil
		}

		slashTok := p.curTok
		var pathText = slashTok.Text
		if !p.advanceIfPeekTokenIs(token.IDENT) {
			return nil
		}

		pathText += p.curTok.Text
		if p.peekTokenIs(token.SUB) { //  parse abc-efg format
			if !p.nextToken() {
				return nil
			}

			pathText += p.curTok.Text
			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}

			pathText += p.curTok.Text
		}

		valueTok = token.Token{
			Text:     pathText,
			Position: slashTok.Position,
		}
		leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
	} else if p.peekTokenIs(token.DURATION) {
		if !p.nextToken() {
			return nil
		}

		valueTok = p.curTok
		leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
		node := ast.NewTokenNode(valueTok)
		node.SetLeadingCommentGroup(leadingCommentGroup)
		expr.Value = node
		return expr
	} else if p.peekTokenIs(token.INT) {
		if !p.nextToken() {
			return nil
		}

		valueTok = p.curTok
		leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
		node := ast.NewTokenNode(valueTok)
		node.SetLeadingCommentGroup(leadingCommentGroup)
		expr.Value = node
		return expr
	} else if p.peekTokenIs(token.STRING) {
		if !p.nextToken() {
			return nil
		}

		valueTok = p.curTok
		leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
		node := ast.NewTokenNode(valueTok)
		node.SetLeadingCommentGroup(leadingCommentGroup)
		expr.Value = node
		return expr
	} else {
		if !p.advanceIfPeekTokenIs(token.IDENT) {
			return nil
		}

		valueTok = p.curTok
		leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
		if p.peekTokenIs(token.COMMA) {
			for {
				if p.peekTokenIs(token.COMMA) {
					if !p.nextToken() {
						return nil
					}

					slashTok := p.curTok
					if !p.advanceIfPeekTokenIs(token.IDENT) {
						return nil
					}

					idTok := p.curTok
					valueTok = token.Token{
						Text:     valueTok.Text + slashTok.Text + idTok.Text,
						Position: valueTok.Position,
					}
					leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
				} else {
					break
				}
			}

			valueTok.Type = token.PATH
			node := ast.NewTokenNode(valueTok)
			node.SetLeadingCommentGroup(leadingCommentGroup)
			expr.Value = node
			return expr
		} else if p.peekTokenIs(token.SUB) {
			for {
				if p.peekTokenIs(token.SUB) {
					if !p.nextToken() {
						return nil
					}

					subTok := p.curTok
					if !p.advanceIfPeekTokenIs(token.IDENT) {
						return nil
					}

					idTok := p.curTok
					valueTok = token.Token{
						Text:     valueTok.Text + subTok.Text + idTok.Text,
						Position: valueTok.Position,
					}
					leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
				} else {
					break
				}
			}

			valueTok.Type = token.PATH
			node := ast.NewTokenNode(valueTok)
			node.SetLeadingCommentGroup(leadingCommentGroup)
			expr.Value = node
			return expr
		}
	}

	for {
		if p.peekTokenIs(token.QUO) {
			if !p.nextToken() {
				return nil
			}

			slashTok := p.curTok
			var pathText = valueTok.Text
			pathText += slashTok.Text
			if !p.advanceIfPeekTokenIs(token.IDENT) {
				return nil
			}

			pathText += p.curTok.Text
			if p.peekTokenIs(token.SUB) { //  parse abc-efg format
				if !p.nextToken() {
					return nil
				}

				pathText += p.curTok.Text
				if !p.advanceIfPeekTokenIs(token.IDENT) {
					return nil
				}

				pathText += p.curTok.Text
			}

			valueTok = token.Token{
				Text:     pathText,
				Position: valueTok.Position,
			}
			leadingCommentGroup = p.curTokenNode().LeadingCommentGroup
		} else {
			break
		}
	}

	valueTok.Type = token.PATH
	node := ast.NewTokenNode(valueTok)
	node.SetLeadingCommentGroup(leadingCommentGroup)
	expr.Value = node

	return expr
}

func (p *Parser) parseKVExpression() *ast.KVExpr {
	var expr = &ast.KVExpr{}

	// token IDENT
	if !p.advanceIfPeekTokenIs(token.IDENT) {
		return nil
	}

	expr.Key = p.curTokenNode()

	// token COLON
	if !p.advanceIfPeekTokenIs(token.COLON) {
		return nil
	}
	expr.Colon = p.curTokenNode()

	// token STRING
	if !p.advanceIfPeekTokenIs(token.STRING) {
		return nil
	}

	expr.Value = p.curTokenNode()

	return expr
}

// syntax = "v1"
func (p *Parser) parseSyntaxStmt() *ast.SyntaxStmt {
	var stmt = &ast.SyntaxStmt{}
	stmt.Syntax = p.curTokenNode()

	// token '='
	if !p.advanceIfPeekTokenIs(token.ASSIGN) {
		return nil
	}

	stmt.Assign = p.curTokenNode()

	// token STRING
	if !p.advanceIfPeekTokenIs(token.STRING) {
		return nil
	}

	stmt.Value = p.curTokenNode()

	return stmt
}

func (p *Parser) curTokenIsNotEof() bool {
	return p.curTokenIsNot(token.EOF)
}

func (p *Parser) curTokenIsNot(expected token.Type) bool {
	return p.curTok.Type != expected
}

func (p *Parser) curTokenIsKeyword() bool {
	tp, ok := token.LookupKeyword(p.curTok.Text)
	if ok {
		p.curTokenIs()
		p.expectIdentError(p.curTok.Fork(tp), token.IDENT)
		return true
	}

	return false
}

func (p *Parser) curTokenIs(expected ...interface{}) bool {
	for _, v := range expected {
		switch val := v.(type) {
		case token.Type:
			if p.curTok.Type == val {
				return true
			}
		case string:
			if p.curTok.Text == val {
				return true
			}
		}
	}

	return false
}

func (p *Parser) advanceIfPeekTokenIs(expected ...interface{}) bool {
	if p.expectPeekToken(expected...) {
		if !p.nextToken() {
			return false
		}

		return true
	}

	return false
}

func (p *Parser) peekTokenIs(expected ...interface{}) bool {
	for _, v := range expected {
		switch val := v.(type) {
		case token.Type:
			if p.peekTok.Type == val {
				return true
			}
		case string:
			if p.peekTok.Text == val {
				return true
			}
		}
	}

	return false
}

func (p *Parser) peekTokenIsNot(expected ...interface{}) bool {
	for _, v := range expected {
		switch val := v.(type) {
		case token.Type:
			if p.peekTok.Type == val {
				return false
			}
		case string:
			if p.peekTok.Text == val {
				return false
			}
		}
	}

	return true
}

func (p *Parser) notExpectPeekToken(expected ...interface{}) bool {
	if !p.peekTokenIsNot(expected...) {
		return false
	}

	var expectedString []string
	for _, v := range expected {
		expectedString = append(expectedString, fmt.Sprintf("'%s'", v))
	}

	var got string
	if p.peekTok.Type == token.ILLEGAL {
		got = p.peekTok.Text
	} else {
		got = p.peekTok.Type.String()
	}

	var err error
	if p.peekTok.Type == token.EOF {
		position := p.curTok.Position
		position.Column = position.Column + len(p.curTok.Text)
		err = fmt.Errorf(
			"%s syntax error: expected %s, got '%s'",
			position,
			strings.Join(expectedString, " | "),
			got)
	} else {
		err = fmt.Errorf(
			"%s syntax error: expected %s, got '%s'",
			p.peekTok.Position,
			strings.Join(expectedString, " | "),
			got)
	}
	p.errors = append(p.errors, err)

	return true
}

func (p *Parser) notExpectPeekTokenGotComment(actual *ast.CommentStmt, expected ...interface{}) bool {
	if actual == nil {
		return false
	}

	var expectedString []string
	for _, v := range expected {
		switch val := v.(type) {
		case token.Token:
			expectedString = append(expectedString, fmt.Sprintf("'%s'", val.Text))
		default:
			expectedString = append(expectedString, fmt.Sprintf("'%s'", v))
		}
	}

	got := actual.Comment.Type.String()
	p.errors = append(p.errors, fmt.Errorf(
		"%s syntax error: expected %s, got '%s'",
		p.peekTok.Position,
		strings.Join(expectedString, " | "),
		got))

	return true
}

func (p *Parser) expectPeekToken(expected ...interface{}) bool {
	if p.peekTokenIs(expected...) {
		return true
	}

	var expectedString []string
	for _, v := range expected {
		expectedString = append(expectedString, fmt.Sprintf("'%s'", v))
	}

	var got string
	if p.peekTok.Type == token.ILLEGAL {
		got = p.peekTok.Text
	} else {
		got = p.peekTok.Type.String()
	}

	var err error
	if p.peekTok.Type == token.EOF {
		position := p.curTok.Position
		position.Column = position.Column + len(p.curTok.Text)
		err = fmt.Errorf(
			"%s syntax error: expected %s, got '%s'",
			position,
			strings.Join(expectedString, " | "),
			got)
	} else {
		err = fmt.Errorf(
			"%s syntax error: expected %s, got '%s'",
			p.peekTok.Position,
			strings.Join(expectedString, " | "),
			got)
	}
	p.errors = append(p.errors, err)

	return false
}

func (p *Parser) expectIdentError(tok token.Token, expected ...interface{}) {
	var expectedString []string
	for _, v := range expected {
		expectedString = append(expectedString, fmt.Sprintf("'%s'", v))
	}

	p.errors = append(p.errors, fmt.Errorf(
		"%s syntax error: expected %s, got '%s'",
		tok.Position,
		strings.Join(expectedString, " | "),
		tok.Type.String()))
}

func (p *Parser) init() bool {
	if !p.nextToken() {
		return false
	}

	return p.nextToken()
}

func (p *Parser) nextToken() bool {
	var err error
	p.curTok = p.peekTok
	var line = -1
	if p.curTok.Valid() {
		if p.curTokenIs(token.EOF) {
			for _, v := range p.headCommentGroup {
				p.appendStmt(v)
			}

			p.headCommentGroup = ast.CommentGroup{}
			return true
		}

		node := ast.NewTokenNode(p.curTok)
		if p.headCommentGroup.Valid() {
			node.HeadCommentGroup = append(node.HeadCommentGroup, p.headCommentGroup...)
			p.headCommentGroup = ast.CommentGroup{}
		}

		p.node[p.curTok] = node
		line = p.curTok.Line()
	}

	p.peekTok, err = p.s.NextToken()
	if err != nil {
		p.errors = append(p.errors, err)
		return false
	}

	var leadingCommentGroup ast.CommentGroup
	for p.peekTok.Type == token.COMMENT || p.peekTok.Type == token.DOCUMENT {
		commentStmt := &ast.CommentStmt{Comment: p.peekTok}
		if p.peekTok.Line() == line && line > -1 {
			leadingCommentGroup = append(leadingCommentGroup, commentStmt)
		} else {
			p.headCommentGroup = append(p.headCommentGroup, commentStmt)
		}

		p.peekTok, err = p.s.NextToken()
		if err != nil {
			p.errors = append(p.errors, err)
			return false
		}
	}

	if len(leadingCommentGroup) > 0 {
		p.curTokenNode().SetLeadingCommentGroup(leadingCommentGroup)
	}

	return true
}

func (p *Parser) curTokenNode() *ast.TokenNode {
	return p.getNode(p.curTok)
}

func (p *Parser) getNode(tok token.Token) *ast.TokenNode {
	return p.node[tok]
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	vo := reflect.ValueOf(v)
	if vo.Kind() == reflect.Ptr {
		return vo.IsNil()
	}

	return false
}

// CheckErrors check parser errors.
func (p *Parser) CheckErrors() error {
	if len(p.errors) == 0 {
		return nil
	}

	var errors []string
	for _, e := range p.errors {
		errors = append(errors, e.Error())
	}

	return fmt.Errorf(strings.Join(errors, "\n"))
}

func (p *Parser) appendStmt(stmt ...ast.Stmt) {
	p.api.Stmts = append(p.api.Stmts, stmt...)
}

func (p *Parser) hasNoErrors() bool {
	return len(p.errors) == 0
}

/************************EXPERIMENTAL CODE BG************************/
// The following code block are experimental, do not call it out of unit test.

// ParseForUintTest parse the source code for unit test.
func (p *Parser) ParseForUintTest() *ast.AST {
	api := &ast.AST{}
	if !p.init() {
		return nil
	}

	for p.curTokenIsNotEof() {
		stmt := p.parseStmtForUniTest()
		if isNil(stmt) {
			return nil
		}

		api.Stmts = append(api.Stmts, stmt)
		if !p.nextToken() {
			return nil
		}
	}

	return api
}

func (p *Parser) parseStmtForUniTest() ast.Stmt {
	switch p.curTok.Type {
	case token.AT_SERVER:
		return p.parseAtServerStmt()
	case token.AT_HANDLER:
		return p.parseAtHandlerStmt()
	case token.AT_DOC:
		return p.parseAtDocStmt()
	}
	return nil
}

/************************EXPERIMENTAL CODE END************************/
