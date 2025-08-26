package parser

import (
	"errors"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/emicklei/proto"
)

type (
	// DefaultProtoParser types an empty struct
	DefaultProtoParser struct{}
)

var ErrGoPackage = errors.New(`option go_package = "" field is not filled in`)

// NewDefaultProtoParser creates a new instance
func NewDefaultProtoParser() *DefaultProtoParser {
	return &DefaultProtoParser{}
}

// Parse provides to parse the proto file into a golang structure,
// which is convenient for subsequent rpc generation and use
func (p *DefaultProtoParser) Parse(src string, multiple ...bool) (Proto, error) {
	var ret Proto

	abs, err := filepath.Abs(src)
	if err != nil {
		return Proto{}, err
	}

	r, err := os.Open(abs)
	if err != nil {
		return ret, err
	}
	defer r.Close()

	parser := proto.NewParser(r)
	set, err := parser.Parse()
	if err != nil {
		return ret, err
	}

	var serviceList Services
	proto.Walk(
		set,
		proto.WithImport(func(i *proto.Import) {
			ret.Import = append(ret.Import, Import{Import: i})
		}),
		proto.WithMessage(func(message *proto.Message) {
			ret.Message = append(ret.Message, Message{Message: message})
		}),
		proto.WithPackage(func(p *proto.Package) {
			ret.Package = Package{Package: p}
		}),
		proto.WithService(func(service *proto.Service) {
			serv := Service{Service: service}
			elements := service.Elements
			for _, el := range elements {
				v, _ := el.(*proto.RPC)
				if v == nil {
					continue
				}
				serv.RPC = append(serv.RPC, &RPC{RPC: v})
			}

			serviceList = append(serviceList, serv)
		}),
		proto.WithOption(func(option *proto.Option) {
			if option.Name == "go_package" {
				ret.GoPackage = option.Constant.Source
			}
		}),
	)
	if err = serviceList.validate(abs, multiple...); err != nil {
		return ret, err
	}

	if len(ret.GoPackage) == 0 {
		if ret.Package.Package == nil {
			return ret, ErrGoPackage
		}
		ret.GoPackage = ret.Package.Name
	}

	ret.PbPackage = GoSanitized(filepath.Base(ret.GoPackage))
	ret.Src = abs
	ret.Name = filepath.Base(abs)
	ret.Service = serviceList

	return ret, nil
}

// GoSanitized copy from protobuf, for more information, please see google.golang.org/protobuf@v1.25.0/internal/strs/strings.go:71
func GoSanitized(s string) string {
	// Sanitize the input to the set of valid characters,
	// which must be '_' or be in the Unicode L or N categories.
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	// Prepend '_' in the event of a Go keyword conflict or if
	// the identifier is invalid (does not start in the Unicode L category).
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}

// CamelCase copy from protobuf, for more information, please see github.com/golang/protobuf@v1.4.2/protoc-gen-go/generator/generator.go:2648
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
