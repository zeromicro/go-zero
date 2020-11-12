package parser

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type rootState struct {
	*baseState
}

func newRootState(r *bufio.Reader, lineNumber *int) state {
	var state = newBaseState(r, lineNumber)
	return rootState{
		baseState: state,
	}
}

func (s rootState) process(api *spec.ApiSpec) (state, error) {
	var annos []spec.Annotation
	var builder strings.Builder
	for {
		ch, err := s.readSkipComment()
		if err != nil {
			return nil, err
		}

		switch {
		case isSpace(ch):
			if builder.Len() == 0 {
				continue
			}

			token := builder.String()
			builder.Reset()
			return s.processToken(token, annos)
		case ch == at:
			if builder.Len() > 0 {
				return nil, fmt.Errorf("%q before %q", builder.String(), at)
			}

			var annoName string
		annoLoop:
			for {
				next, err := s.readSkipComment()
				if err != nil {
					return nil, err
				}

				switch {
				case isSpace(next):
					if builder.Len() > 0 {
						annoName = builder.String()
						builder.Reset()
					}
				case next == leftParenthesis:
					if err := s.unread(); err != nil {
						return nil, err
					}

					if builder.Len() > 0 {
						annoName = builder.String()
						builder.Reset()
					}
					attrs, err := s.parseProperties()
					if err != nil {
						return nil, err
					}

					annos = append(annos, spec.Annotation{
						Name:       annoName,
						Properties: attrs,
					})
					break annoLoop
				default:
					builder.WriteRune(next)
				}
			}
		case ch == leftParenthesis:
			if builder.Len() == 0 {
				return nil, fmt.Errorf("incorrect %q at the beginning of the line", leftParenthesis)
			}

			if err := s.unread(); err != nil {
				return nil, err
			}

			token := builder.String()
			builder.Reset()
			return s.processToken(token, annos)
		case isLetterDigit(ch):
			builder.WriteRune(ch)
		case isNewline(ch):
			if builder.Len() > 0 {
				return nil, fmt.Errorf("incorrect newline after %q", builder.String())
			}
		}
	}
}

func (s rootState) processToken(token string, annos []spec.Annotation) (state, error) {
	switch token {
	case infoDirective:
		return newInfoState(s.baseState), nil
	case serviceDirective:
		return newServiceState(s.baseState, annos), nil
	default:
		return nil, fmt.Errorf("wrong directive %q", token)
	}
}
