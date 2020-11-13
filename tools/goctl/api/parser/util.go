package parser

import (
	"bufio"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

var emptyType spec.Type

func GetType(api *spec.ApiSpec, t string) spec.Type {
	for _, tp := range api.Types {
		if tp.Name == t {
			return tp
		}
	}

	return emptyType
}

func isLetterDigit(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || ('0' <= r && r <= '9')
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isSlash(r rune) bool {
	return r == '/'
}

func isNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

func read(r *bufio.Reader) (rune, error) {
	ch, _, err := r.ReadRune()
	return ch, err
}

func readLine(r *bufio.Reader) (string, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return "", err
	} else {
		return string(line), nil
	}
}

func skipSpaces(r *bufio.Reader) error {
	for {
		next, err := read(r)
		if err != nil {
			return err
		}
		if !isSpace(next) {
			return unread(r)
		}
	}
}

func unread(r *bufio.Reader) error {
	return r.UnreadRune()
}
