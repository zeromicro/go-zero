package parser

import (
	"bufio"
	"errors"
	"regexp"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

// struct match
const typeRegex = `(?m)(?m)(^ *type\s+[a-zA-Z][a-zA-Z0-9_-]+\s+(((struct)\s*?\{[\w\W]*?[^\{]\})|([a-zA-Z][a-zA-Z0-9_-]+)))|(^ *type\s*?\([\w\W]+\}\s*\))`

var (
	emptyStrcut = errors.New("struct body not found")
	emptyType   spec.Type
)

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

func MatchStruct(api string) (info, structBody, service string, err error) {
	r := regexp.MustCompile(typeRegex)
	indexes := r.FindAllStringIndex(api, -1)
	if len(indexes) == 0 {
		return "", "", "", emptyStrcut
	}
	startIndexes := indexes[0]
	endIndexes := indexes[len(indexes)-1]

	info = api[:startIndexes[0]]
	structBody = api[startIndexes[0]:endIndexes[len(endIndexes)-1]]
	service = api[endIndexes[len(endIndexes)-1]:]

	firstIIndex := strings.Index(info, "i")
	if firstIIndex > 0 {
		info = info[firstIIndex:]
	}

	lastServiceRightBraceIndex := strings.LastIndex(service, "}") + 1
	var firstServiceIndex int
	for index, char := range service {
		if !isSpace(char) && !isNewline(char) {
			firstServiceIndex = index
			break
		}
	}
	service = service[firstServiceIndex:lastServiceRightBraceIndex]
	return
}
