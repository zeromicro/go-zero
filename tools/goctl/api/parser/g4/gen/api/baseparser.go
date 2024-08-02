package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/zeromicro/antlr"
)

const (
	versionRegex     = `(?m)"v[1-9][0-9]*"`
	importValueRegex = `(?m)"\/?(?:[^/]+\/)*[^/]+.api"`
	tagRegex         = `(?m)\x60[a-z]+:".+"\x60`
)

var (
	holder = struct{}{}
	kind   = map[string]struct{}{
		"bool":       holder,
		"int":        holder,
		"int8":       holder,
		"int16":      holder,
		"int32":      holder,
		"int64":      holder,
		"uint":       holder,
		"uint8":      holder,
		"uint16":     holder,
		"uint32":     holder,
		"uint64":     holder,
		"uintptr":    holder,
		"float32":    holder,
		"float64":    holder,
		"complex64":  holder,
		"complex128": holder,
		"string":     holder,
		"byte":       holder,
		"rune":       holder,
	}
)

func match(p *ApiParserParser, text string) {
	v := getCurrentTokenText(p)

	if v != text {
		notifyErrorListeners(p, expecting(text, v))
	}
}

func checkVersion(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if !matchRegex(v, versionRegex) {
		notifyErrorListeners(p, mismatched("version", v))
	}
}

func checkImportValue(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if !matchRegex(v, importValueRegex) {
		notifyErrorListeners(p, mismatched("import value", v))
	}
}

func checkKeyValue(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if !strings.HasPrefix(v, ":") {
		notifyErrorListeners(p, mismatched(":", v))
	}

	v = strings.TrimPrefix(v, ":")
	v = strings.TrimFunc(v, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	setCurrentTokenText(p, v)
}

func checkHTTPMethod(p *ApiParserParser) {
	method := getCurrentTokenText(p)
	uppler := strings.ToUpper(method)
	switch uppler {
	case http.MethodPost, http.MethodGet, http.MethodHead,
		http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
		if method != strings.ToLower(method) {
			notifyErrorListeners(p, expecting("http method lower case", method))
		}
	default:
		notifyErrorListeners(p, expecting("http method", method))
	}
}

func checkKeyword(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if IsGolangKeyWord(v) {
		notifyErrorListeners(p, fmt.Sprintf("expecting ID, found golang keyword: '%s'", v))
	}
}

func checkKey(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if IsGolangKeyWord(v) {
		notifyErrorListeners(p, fmt.Sprintf("expecting ID, found golang keyword: '%s'", v))
	}

	if _, ok := kind[v]; !ok {
		notifyErrorListeners(p, fmt.Sprintf("expecting golang basic type, found : '%s'", v))
	}
}

// IsBasicType returns true if the input argument is basic golang type
func IsBasicType(text string) bool {
	_, ok := kind[text]
	return ok
}

// IsGolangKeyWord returns true if input argument is golang keyword, but it will be ignored which in excepts
func IsGolangKeyWord(text string, excepts ...string) bool {
	for _, each := range excepts {
		if text == each {
			return false
		}
	}

	switch text {
	case "var", "const", "package", "import", "func", "return",
		"defer", "go", "select", "interface", "struct", "break", "case",
		"continue", "for", "fallthrough", "else", "if", "switch", "goto",
		"default", "chan", "type", "map", "range":
		return true
	default:
		return false
	}
}

func isNormal(p *ApiParserParser) bool {
	ct := p.GetTokenStream().(*antlr.CommonTokenStream)
	line := p.GetCurrentToken().GetLine()
	tokens := ct.GetAllTokens()
	var list []string
	for _, token := range tokens {
		if token.GetLine() == line {
			text := token.GetText()
			if strings.HasPrefix(text, "//") {
				continue
			}
			if strings.HasPrefix(text, "/*") {
				continue
			}
			if text == "<EOF>" {
				continue
			}
			if strings.TrimSpace(text) == "" {
				continue
			}
			list = append(list, text)
		}
	}
	if len(list) == 1 {
		t := strings.TrimPrefix(list[0], "*")
		if IsGolangKeyWord(t) {
			notifyErrorListeners(p, fmt.Sprintf("expecting ID, found golang keyword: '%s'", t))
		}
	}
	if len(list) > 1 {
		if list[0] == "*" {
			t := strings.TrimPrefix(list[1], "*")
			if IsGolangKeyWord(t) {
				notifyErrorListeners(p, fmt.Sprintf("expecting ID, found golang keyword: '%s'", t))
			}
			return false
		}
	}
	return len(list) > 1
}

// MatchTag returns a Boolean value, which returns true if it does matched, otherwise returns false
func MatchTag(v string) bool {
	return matchRegex(v, tagRegex)
}

func isInterface(p *ApiParserParser) {
	v := getCurrentTokenText(p)
	if IsGolangKeyWord(v) {
		notifyErrorListeners(p, fmt.Sprintf("expecting ID, found golang keyword: '%s'", v))
	}
}

func getCurrentTokenText(p *ApiParserParser) string {
	token := p.GetCurrentToken()
	if token == nil {
		return ""
	}

	return token.GetText()
}

func setCurrentTokenText(p *ApiParserParser, text string) {
	token := p.GetCurrentToken()
	if token == nil {
		return
	}

	token.SetText(text)
}

func notifyErrorListeners(p *ApiParserParser, msg string) {
	p.NotifyErrorListeners(msg, nil, nil)
}

func matchRegex(text, str string) bool {
	re := regexp.MustCompile(str)
	v := re.FindString(text)
	text = strings.TrimFunc(text, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	return v == text
}

func expecting(expecting, found string) string {
	return fmt.Sprintf(`expecting '%s', found input '%s'`, expecting, found)
}

func mismatched(expecting, found string) string {
	return fmt.Sprintf(`mismatched '%s', found input '%s'`, expecting, found)
}
