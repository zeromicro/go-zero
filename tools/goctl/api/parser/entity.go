package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type (
	entity struct {
		state  *baseState
		api    *spec.ApiSpec
		parser entityParser
	}

	entityParser interface {
		parseLine(line string, api *spec.ApiSpec, annos []spec.Annotation) error
		setEntityName(name string)
	}
)

func newEntity(state *baseState, api *spec.ApiSpec, parser entityParser) entity {
	return entity{
		state:  state,
		api:    api,
		parser: parser,
	}
}

func (s *entity) process() error {
	line, err := s.state.readLineSkipComment()
	if err != nil {
		return err
	}

	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("invalid type definition for %q",
			strings.TrimSpace(strings.Trim(string(line), "{")))
	}

	if len(fields) == 2 {
		if fields[1] != leftBrace {
			return fmt.Errorf("bad string %q after type", fields[1])
		}
	} else if len(fields) == 3 {
		if fields[1] != typeStruct {
			return fmt.Errorf("bad string %q after type", fields[1])
		}
		if fields[2] != leftBrace {
			return fmt.Errorf("bad string %q after type", fields[2])
		}
	}

	s.parser.setEntityName(fields[0])

	var annos []spec.Annotation
memberLoop:
	for {
		ch, err := s.state.readSkipComment()
		if err != nil {
			return err
		}

		var annoName string
		var builder strings.Builder
		switch {
		case ch == at:
		annotationLoop:
			for {
				next, err := s.state.readSkipComment()
				if err != nil {
					return err
				}
				switch {
				case isSpace(next):
					if builder.Len() > 0 && annoName == "" {
						annoName = builder.String()
						builder.Reset()
					}
				case isNewline(next):
					if builder.Len() == 0 {
						return errors.New("invalid annotation format")
					}

					if len(annoName) > 0 {
						value := builder.String()
						if value != string(leftParenthesis) {
							builder.Reset()
							annos = append(annos, spec.Annotation{
								Name:  annoName,
								Value: value,
							})
							annoName = ""
							break annotationLoop
						}
					}
				case next == leftParenthesis:
					if builder.Len() == 0 {
						return errors.New("invalid annotation format")
					}
					annoName = builder.String()
					builder.Reset()
					if err := s.state.unread(); err != nil {
						return err
					}
					attrs, err := s.state.parseProperties()
					if err != nil {
						return err
					}
					annos = append(annos, spec.Annotation{
						Name:       annoName,
						Properties: attrs,
					})
					annoName = ""
					break annotationLoop
				default:
					builder.WriteRune(next)
				}
			}
		case ch == rightBrace:
			break memberLoop
		case isLetterDigit(ch):
			if err := s.state.unread(); err != nil {
				return err
			}

			var line string
			line, err = s.state.readLineSkipComment()
			if err != nil {
				return err
			}

			line = strings.TrimSpace(line)
			if err := s.parser.parseLine(line, s.api, annos); err != nil {
				return err
			}

			annos = nil
		}
	}

	return nil
}
