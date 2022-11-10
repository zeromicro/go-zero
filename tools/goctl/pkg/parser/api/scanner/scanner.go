package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const (
	initMode mode = iota

	// document mode bg
	documentHalfOpen
	documentOpen
	documentHalfClose
	documentClose
	// document mode end

	// string mode bg
	stringOpen
	stringClose
	// string mode end

)

var missingInput = errors.New("missing input")

type mode int
type Scanner struct {
	filename string
	size     int

	data         []rune
	position     int // 当前字符位置
	readPosition int // 当前读取位置（当前字符之后的位置）
	ch           rune

	lines []int
}

func (s *Scanner) NextToken() (token.Token, error) {
	s.skipWhiteSpace()
	switch s.ch {
	case '/':
		peekOne := s.peekRune()
		switch peekOne {
		case '/':
			return s.scanLineComment(), nil
		case '*':
			return s.scanDocument()
		default:
			return s.newToken(token.QUO), nil
		}
	case '-':
		return s.newToken(token.SUB), nil
	case '*':
		return s.newToken(token.MUL), nil
	case '(':
		return s.newToken(token.LPAREN), nil
	case '[':
		return s.newToken(token.LBRACK), nil
	case '{':
		return s.newToken(token.LBRACE), nil
	case ',':
		return s.newToken(token.COMMA), nil
	case '.':
		position := s.position
		peekOne := s.peekRune()
		if peekOne != '.' {
			return s.newToken(token.DOT), nil
		}
		s.readRune()
		peekOne = s.peekRune()
		if peekOne != '.' {
			return s.newToken(token.DOT), nil
		}
		s.readRune()
		s.readRune()
		return token.Token{
			Type:     token.ELLIPSIS,
			Text:     "...",
			Position: s.newPosition(position),
		}, nil
	case ')':
		return s.newToken(token.RPAREN), nil
	case ']':
		return s.newToken(token.RBRACK), nil
	case '}':
		return s.newToken(token.RBRACE), nil
	case ';':
		return s.newToken(token.SEMICOLON), nil
	case ':':
		return s.newToken(token.COLON), nil
	case '=':
		return s.newToken(token.ASSIGN), nil
	case '@':
		return s.scanAt()
	case '"':
		return s.scanString('"', token.STRING)
	case '`':
		return s.scanString('`', token.RAW_STRING)
	case 0:
		return token.EofToken, nil
	default:
		if s.isIdentifierLetter(s.ch) {
			return s.scanIdent(), nil
		}
		if s.isDigit(s.ch) {
			return s.scanInt(), nil
		}
		tok := token.NewIllegalToken(s.ch, s.newPosition(s.position))
		s.readRune()
		return tok, nil
	}
}

func (s *Scanner) newToken(tp token.Type) token.Token {
	tok := token.Token{
		Type:     tp,
		Text:     string(s.ch),
		Position: s.positionAt(),
	}
	s.readRune()
	return tok
}

func (s *Scanner) readRune() {
	if s.readPosition >= s.size {
		s.ch = 0
	} else {
		s.ch = s.data[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition += 1
}

func (s *Scanner) peekRune() rune {
	if s.readPosition >= s.size {
		return 0
	}
	return s.data[s.readPosition]
}

func (s *Scanner) scanString(delim rune, tp token.Type) (token.Token, error) {
	position := s.position
	var stringMode = initMode
	for {
		switch s.ch {
		case delim:
			switch stringMode {
			case initMode:
				stringMode = stringOpen
			case stringOpen:
				stringMode = stringClose
				s.readRune()
				return token.Token{
					Type:     tp,
					Text:     string(s.data[position:s.position]),
					Position: s.newPosition(position),
				}, nil
			}
		case 0:
			switch stringMode {
			case initMode: // assert: dead code
				return token.ErrorToken, s.assertExpected(token.EOF, tp)
			case stringOpen:
				return token.ErrorToken, s.assertExpectedString(token.EOF.String(), string(delim))
			case stringClose: // assert: dead code
				return token.Token{
					Type:     tp,
					Text:     string(s.data[position:s.position]),
					Position: s.newPosition(position),
				}, nil
			}
		}
		s.readRune()
	}
}

func (s *Scanner) scanAt() (token.Token, error) {
	position := s.position
	peek := s.peekRune()
	if !s.isLetter(peek) {
		if peek == 0 {
			return token.NewIllegalToken(s.ch, s.positionAt()), nil
		}
		return token.ErrorToken, s.assertExpectedString(string(peek), token.IDENT.String())
	}

	s.readRune()
	letters := s.scanLetterSet()
	switch letters {
	case "handler":
		return token.Token{
			Type:     token.AT_HANDLER,
			Text:     "@handler",
			Position: s.newPosition(position),
		}, nil
	case "server":
		return token.Token{
			Type:     token.AT_SERVER,
			Text:     "@server",
			Position: s.newPosition(position),
		}, nil
	case "doc":
		return token.Token{
			Type:     token.AT_DOC,
			Text:     "@doc",
			Position: s.newPosition(position),
		}, nil
	default:

		return token.ErrorToken, s.assertExpectedString(
			"@"+letters,
			token.AT_DOC.String(),
			token.AT_HANDLER.String(),
			token.AT_SERVER.String())
	}
}

func (s *Scanner) scanInt() token.Token {
	position := s.position
	for s.isDigit(s.ch) {
		s.readRune()
	}

	integer := string(s.data[position:s.position])
	return token.Token{
		Type:     token.INT,
		Text:     integer,
		Position: s.newPosition(position),
	}
}

func (s *Scanner) scanIdent() token.Token {
	position := s.position
	for s.isIdentifierLetter(s.ch) || s.isDigit(s.ch) {
		s.readRune()
	}

	ident := string(s.data[position:s.position])
	tp, ok := token.LookupKeyword(ident)
	if ok {
		return token.Token{
			Type:     tp,
			Text:     ident,
			Position: s.newPosition(position),
		}
	}

	return token.Token{
		Type:     token.IDENT,
		Text:     ident,
		Position: s.newPosition(position),
	}
}

func (s *Scanner) scanLetterSet() string {
	position := s.position
	for s.isLetter(s.ch) {
		s.readRune()
	}
	return string(s.data[position:s.position])
}

func (s *Scanner) scanLineComment() token.Token {
	position := s.position
	for s.ch != '\n' && s.ch != 0 {
		s.readRune()
	}
	return token.Token{
		Type:     token.COMMENT,
		Text:     string(s.data[position:s.position]),
		Position: s.newPosition(position),
	}
}

func (s *Scanner) scanDocument() (token.Token, error) {
	position := s.position
	var documentMode = initMode
	for {
		switch s.ch {
		case '*':
			switch documentMode {
			case documentHalfOpen:
				documentMode = documentOpen // /*
			case documentOpen, documentHalfClose:
				documentMode = documentHalfClose // (?m)\/\*\*+
			}

		case 0:
			switch documentMode {
			case initMode, documentHalfOpen: // assert: dead code
				return token.ErrorToken, s.assertExpected(token.EOF, token.MUL)
			case documentOpen:
				return token.ErrorToken, s.assertExpected(token.EOF, token.MUL)
			case documentHalfClose:
				return token.ErrorToken, s.assertExpected(token.EOF, token.QUO)
			}
		case '/':
			switch documentMode {
			case initMode: // /
				documentMode = documentHalfOpen
			case documentHalfOpen: // assert: dead code
				return token.ErrorToken, s.assertExpected(token.QUO, token.MUL)
			case documentHalfClose:
				documentMode = documentClose // /*\*+*/
				s.readRune()
				tok := token.Token{
					Type:     token.DOCUMENT,
					Text:     string(s.data[position:s.position]),
					Position: s.newPosition(position),
				}
				return tok, nil
			}
		}
		s.readRune()
	}
}

func (s *Scanner) assertExpected(actual token.Type, expected ...token.Type) error {
	var expects []string
	for _, v := range expected {
		expects = append(expects, fmt.Sprintf("'%s'", v.String()))
	}

	text := fmt.Sprint(s.positionAt().String(), " ", fmt.Sprintf(
		"expected %s, got '%s'",
		strings.Join(expects, " | "),
		actual.String(),
	))
	return errors.New(text)
}

func (s *Scanner) assertExpectedString(actual string, expected ...string) error {
	var expects []string
	for _, v := range expected {
		expects = append(expects, fmt.Sprintf("'%s'", v))
	}

	text := fmt.Sprint(s.positionAt().String(), " ", fmt.Sprintf(
		"expected %s, got '%s'",
		strings.Join(expects, " | "),
		actual,
	))
	return errors.New(text)
}

func (s *Scanner) positionAt() token.Position {
	return s.newPosition(s.position)
}

func (s *Scanner) newPosition(position int) token.Position {
	line := s.lineCount()
	return token.Position{
		Filename: s.filename,
		Line:     line,
		Column:   position - s.lines[line-1],
	}
}

func (s *Scanner) lineCount() int {
	return len(s.lines)
}

func (s *Scanner) skipWhiteSpace() {
	for s.isWhiteSpace(s.ch) {
		s.readRune()
	}
}

func (s *Scanner) isDigit(b rune) bool {
	return b >= '0' && b <= '9'
}

func (s *Scanner) isLetter(b rune) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func (s *Scanner) isIdentifierLetter(b rune) bool {
	if s.isLetter(b) {
		return true
	}
	return b == '_'
}

func (s *Scanner) isWhiteSpace(b rune) bool {
	if b == '\n' {
		s.lines = append(s.lines, s.position)
	}
	return b == ' ' || b == '\t' || b == '\r' || b == '\f' || b == '\v' || b == '\n'
}

func MustNewScanner(filename string, src interface{}) *Scanner {
	sc, err := NewScanner(filename, src)
	if err != nil {
		log.Fatalln(err)
	}
	return sc
}

func NewScanner(filename string, src interface{}) (*Scanner, error) {
	data, err := readData(filename, src)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, missingInput
	}

	var runeList []rune
	for _, r := range string(data) {
		runeList = append(runeList, r)
	}

	filename = filepath.Base(filename)
	s := &Scanner{
		filename:     filename,
		size:         len(runeList),
		data:         runeList,
		lines:        []int{-1},
		readPosition: 0,
	}

	s.readRune()
	return s, nil
}

func readData(filename string, src interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		return data, nil
	}

	switch v := src.(type) {
	case []byte:
		data = append(data, v...)
	case *bytes.Buffer:
		data = v.Bytes()
	case string:
		data = []byte(v)
	default:
		return nil, fmt.Errorf("unsupported type: %T", src)
	}

	return data, nil
}
