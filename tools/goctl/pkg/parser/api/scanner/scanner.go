package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
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

type mode int

// Scanner is a lexical scanner.
type Scanner struct {
	filename string
	size     int

	data         []rune
	position     int // current position in input (points to current char)
	readPosition int // current reading position in input (after current char)
	ch           rune

	lines []int
}

// NextToken returns the next token.
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
			return s.scanIntOrDuration(), nil
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

func (s *Scanner) scanIntOrDuration() token.Token {
	position := s.position
	for s.isDigit(s.ch) {
		s.readRune()
	}
	switch s.ch {
	case 'n', 'µ', 'm', 's', 'h':
		return s.scanDuration(position)
	default:
		return token.Token{
			Type:     token.INT,
			Text:     string(s.data[position:s.position]),
			Position: s.newPosition(position),
		}
	}
}

// scanDuration scans a duration literal, for example "1ns", "1µs", "1ms", "1s", "1m", "1h".
func (s *Scanner) scanDuration(bgPos int) token.Token {
	switch s.ch {
	case 'n':
		return s.scanNanosecond(bgPos)
	case 'µ':
		return s.scanMicrosecond(bgPos)
	case 'm':
		return s.scanMillisecondOrMinute(bgPos)
	case 's':
		return s.scanSecond(bgPos)
	case 'h':
		return s.scanHour(bgPos)
	default:
		return s.illegalToken()
	}
}

func (s *Scanner) scanNanosecond(bgPos int) token.Token {
	s.readRune()
	if s.ch != 's' {
		return s.illegalToken()
	}
	s.readRune()

	return token.Token{
		Type:     token.DURATION,
		Text:     string(s.data[bgPos:s.position]),
		Position: s.newPosition(bgPos),
	}
}

func (s *Scanner) scanMicrosecond(bgPos int) token.Token {
	s.readRune()
	if s.ch != 's' {
		return s.illegalToken()
	}

	s.readRune()
	if !s.isDigit(s.ch) {
		return token.Token{
			Type:     token.DURATION,
			Text:     string(s.data[bgPos:s.position]),
			Position: s.newPosition(bgPos),
		}
	}

	for s.isDigit(s.ch) {
		s.readRune()
	}

	if s.ch != 'n' {
		return s.illegalToken()
	}

	return s.scanNanosecond(bgPos)

}

func (s *Scanner) scanMillisecondOrMinute(bgPos int) token.Token {
	s.readRune()
	if s.ch != 's' { // minute
		if s.ch == 0 || !s.isDigit(s.ch) {
			return token.Token{
				Type:     token.DURATION,
				Text:     string(s.data[bgPos:s.position]),
				Position: s.newPosition(bgPos),
			}
		}

		return s.scanMinute(bgPos)
	}

	return s.scanMillisecond(bgPos)
}

func (s *Scanner) scanMillisecond(bgPos int) token.Token {
	s.readRune()
	if !s.isDigit(s.ch) {
		return token.Token{
			Type:     token.DURATION,
			Text:     string(s.data[bgPos:s.position]),
			Position: s.newPosition(bgPos),
		}
	}

	for s.isDigit(s.ch) {
		s.readRune()
	}

	switch s.ch {
	case 'n':
		return s.scanNanosecond(bgPos)
	case 'µ':
		return s.scanMicrosecond(bgPos)
	default:
		return s.illegalToken()
	}
}

func (s *Scanner) scanSecond(bgPos int) token.Token {
	s.readRune()
	if !s.isDigit(s.ch) {
		return token.Token{
			Type:     token.DURATION,
			Text:     string(s.data[bgPos:s.position]),
			Position: s.newPosition(bgPos),
		}
	}

	for s.isDigit(s.ch) {
		s.readRune()
	}

	switch s.ch {
	case 'n':
		return s.scanNanosecond(bgPos)
	case 'µ':
		return s.scanMicrosecond(bgPos)
	case 'm':
		s.readRune()
		if s.ch != 's' {
			return s.illegalToken()
		}
		return s.scanMillisecond(bgPos)
	default:
		return s.illegalToken()
	}
}

func (s *Scanner) scanMinute(bgPos int) token.Token {
	if !s.isDigit(s.ch) {
		return token.Token{
			Type:     token.DURATION,
			Text:     string(s.data[bgPos:s.position]),
			Position: s.newPosition(bgPos),
		}
	}

	for s.isDigit(s.ch) {
		s.readRune()
	}

	switch s.ch {
	case 'n':
		return s.scanNanosecond(bgPos)
	case 'µ':
		return s.scanMicrosecond(bgPos)
	case 'm':
		s.readRune()
		if s.ch != 's' {
			return s.illegalToken()
		}
		return s.scanMillisecond(bgPos)
	case 's':
		return s.scanSecond(bgPos)
	default:
		return s.illegalToken()
	}
}

func (s *Scanner) scanHour(bgPos int) token.Token {
	s.readRune()
	if !s.isDigit(s.ch) {
		return token.Token{
			Type:     token.DURATION,
			Text:     string(s.data[bgPos:s.position]),
			Position: s.newPosition(bgPos),
		}
	}

	for s.isDigit(s.ch) {
		s.readRune()
	}

	switch s.ch {
	case 'n':
		return s.scanNanosecond(bgPos)
	case 'µ':
		return s.scanMicrosecond(bgPos)
	case 'm':
		return s.scanMillisecondOrMinute(bgPos)
	case 's':
		return s.scanSecond(bgPos)
	default:
		return s.illegalToken()
	}
}

func (s *Scanner) illegalToken() token.Token {
	tok := token.NewIllegalToken(s.ch, s.newPosition(s.position))
	s.readRune()
	return tok
}

func (s *Scanner) scanIdent() token.Token {
	position := s.position
	for s.isIdentifierLetter(s.ch) || s.isDigit(s.ch) {
		s.readRune()
	}

	ident := string(s.data[position:s.position])
	if ident == "interface" && s.ch == '{' && s.peekRune() == '}' {
		s.readRune()
		s.readRune()
		return token.Token{
			Type:     token.ANY,
			Text:     string(s.data[position:s.position]),
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

// MustNewScanner returns a new scanner for the given filename and data.
func MustNewScanner(filename string, src interface{}) *Scanner {
	sc, err := NewScanner(filename, src)
	if err != nil {
		log.Fatalln(err)
	}
	return sc
}

// NewScanner returns a new scanner for the given filename and data.
func NewScanner(filename string, src interface{}) (*Scanner, error) {
	data, err := readData(filename, src)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("filename: %s, missing input", filename)
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
	if strings.HasSuffix(filename, ".api") && pathx.FileExists(filename) {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	switch v := src.(type) {
	case []byte:
		return v, nil
	case *bytes.Buffer:
		return v.Bytes(), nil
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", src)
	}
}
