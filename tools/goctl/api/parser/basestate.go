package parser

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	startState = iota
	attrNameState
	attrValueState
	attrColonState
	multilineState
)

type baseState struct {
	r          *bufio.Reader
	lineNumber *int
}

func newBaseState(r *bufio.Reader, lineNumber *int) *baseState {
	return &baseState{
		r:          r,
		lineNumber: lineNumber,
	}
}

func (s *baseState) parseProperties() (map[string]string, error) {
	var r = s.r
	var attributes = make(map[string]string)
	var builder strings.Builder
	var key string
	var st = startState

	for {
		ch, err := s.readSkipComment()
		if err != nil {
			return nil, err
		}

		switch st {
		case startState:
			switch {
			case isNewline(ch):
				return nil, fmt.Errorf("%q should be on the same line with %q", leftParenthesis, infoDirective)
			case isSpace(ch):
				continue
			case ch == leftParenthesis:
				st = attrNameState
			default:
				return nil, fmt.Errorf("unexpected char %q after %q", ch, infoDirective)
			}
		case attrNameState:
			switch {
			case isNewline(ch):
				if builder.Len() > 0 {
					return nil, fmt.Errorf("unexpected newline after %q", builder.String())
				}
			case isLetterDigit(ch):
				builder.WriteRune(ch)
			case isSpace(ch):
				if builder.Len() > 0 {
					key = builder.String()
					builder.Reset()
					st = attrColonState
				}
			case ch == colon:
				if builder.Len() == 0 {
					return nil, fmt.Errorf("unexpected leading %q", ch)
				}
				key = builder.String()
				builder.Reset()
				st = attrValueState
			case ch == rightParenthesis:
				return attributes, nil
			}
		case attrColonState:
			switch {
			case isSpace(ch):
				continue
			case ch == colon:
				st = attrValueState
			default:
				return nil, fmt.Errorf("bad char %q after %q in %q", ch, key, infoDirective)
			}
		case attrValueState:
			switch {
			case ch == multilineBeginTag:
				if builder.Len() > 0 {
					return nil, fmt.Errorf("%q before %q", builder.String(), multilineBeginTag)
				} else {
					st = multilineState
				}
			case isSpace(ch):
				if builder.Len() > 0 {
					builder.WriteRune(ch)
				}
			case isNewline(ch):
				attributes[key] = builder.String()
				builder.Reset()
				st = attrNameState
			case ch == rightParenthesis:
				attributes[key] = builder.String()
				builder.Reset()
				return attributes, nil
			default:
				builder.WriteRune(ch)
			}
		case multilineState:
			switch {
			case ch == multilineEndTag:
				attributes[key] = builder.String()
				builder.Reset()
				st = attrNameState
			case isNewline(ch):
				var multipleNewlines bool
			loopAfterNewline:
				for {
					next, err := read(r)
					if err != nil {
						return nil, err
					}

					switch {
					case isSpace(next):
						continue
					case isNewline(next):
						multipleNewlines = true
					default:
						if err := unread(r); err != nil {
							return nil, err
						}
						break loopAfterNewline
					}
				}

				if multipleNewlines {
					fmt.Fprintln(&builder)
				} else {
					builder.WriteByte(' ')
				}
			case ch == rightParenthesis:
				if builder.Len() > 0 {
					attributes[key] = builder.String()
					builder.Reset()
				}
				return attributes, nil
			default:
				builder.WriteRune(ch)
			}
		}
	}
}

func (s *baseState) read() (rune, error) {
	value, err := read(s.r)
	if err != nil {
		return 0, err
	}
	if isNewline(value) {
		*s.lineNumber++
	}
	return value, nil
}

func (s *baseState) readSkipComment() (rune, error) {
	ch, err := s.read()
	if err != nil {
		return 0, err
	}

	if isSlash(ch) {
		value, err := s.mayReadToEndOfLine()
		if err != nil {
			return 0, err
		}

		if value > 0 {
			ch = value
		}
	}
	return ch, nil
}

func (s *baseState) mayReadToEndOfLine() (rune, error) {
	ch, err := s.read()
	if err != nil {
		return 0, err
	}

	if isSlash(ch) {
		for {
			value, err := s.read()
			if err != nil {
				return 0, err
			}

			if isNewline(value) {
				return value, nil
			}
		}
	}
	err = s.unread()
	return 0, err
}

func (s *baseState) readLineSkipComment() (string, error) {
	line, err := s.readLine()
	if err != nil {
		return "", err
	}

	var commentIdx = strings.Index(line, "//")
	if commentIdx >= 0 {
		return line[:commentIdx], nil
	}
	return line, nil
}

func (s *baseState) readLine() (string, error) {
	line, _, err := s.r.ReadLine()
	if err != nil {
		return "", err
	}
	*s.lineNumber++
	return string(line), nil
}

func (s *baseState) skipSpaces() error {
	return skipSpaces(s.r)
}

func (s *baseState) unread() error {
	return unread(s.r)
}
