package scanner

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

func Test_readData(t *testing.T) {
	testData := []struct {
		input    interface{}
		expected interface{}
	}{
		{
			input:    []byte("foo"),
			expected: []byte("foo"),
		},
		{
			input:    bytes.NewBufferString("foo"),
			expected: []byte("foo"),
		},
		{
			input:    "foo",
			expected: []byte("foo"),
		},
		{
			input:    "",
			expected: []byte{},
		},
		{
			input:    strings.NewReader("foo"),
			expected: fmt.Errorf("unsupported type: *strings.Reader"),
		},
	}
	for _, v := range testData {
		actual, err := readData("", v.input)
		if err != nil {
			assert.Equal(t, v.expected.(error).Error(), err.Error())
		} else {
			assert.Equal(t, v.expected, actual)
		}
	}
}

func TestNewScanner(t *testing.T) {
	testData := []struct {
		filename string
		src      interface{}
		expected interface{}
	}{
		{
			filename: "foo",
			src:      "foo",
			expected: "foo",
		},
		{
			filename: "foo",
			src:      "",
			expected: missingInput,
		},
	}
	for _, v := range testData {
		s, err := NewScanner(v.filename, v.src)
		if err != nil {
			assert.Equal(t, v.expected.(error).Error(), err.Error())
		} else {
			assert.Equal(t, v.expected, s.filename)
		}
	}
}

func TestScanner_NextToken_lineComment(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.COMMENT,
			Text: "//",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.COMMENT,
			Text: "//foo",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type: token.COMMENT,
			Text: "//bar",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
		{
			Type: token.COMMENT,
			Text: "///",
			Position: token.Position{
				Filename: "foo.api",
				Line:     4,
				Column:   1,
			},
		},
		{
			Type: token.COMMENT,
			Text: "////",
			Position: token.Position{
				Filename: "foo.api",
				Line:     5,
				Column:   1,
			},
		},
		{
			Type: token.QUO,
			Text: "/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     6,
				Column:   1,
			},
		},
		token.EofToken,
	}
	var input = "//\n//foo\n//bar\n///\n////\n/"
	s, err := NewScanner("foo.api", input)
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_document(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.DOCUMENT,
			Text: "/**/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/***/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   6,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*-*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   12,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*/*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   18,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*////*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   24,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*foo*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*---*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
		{
			Type: token.DOCUMENT,
			Text: "/*\n*/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     4,
				Column:   1,
			},
		},
		{
			Type: token.QUO,
			Text: "/",
			Position: token.Position{
				Filename: "foo.api",
				Line:     5,
				Column:   1,
			},
		},
		token.EofToken,
	}
	var input = "/**/ /***/ /*-*/ /*/*/ /*////*/  \n/*foo*/\n/*---*/\n/*\n*/\n/"
	s, err := NewScanner("foo.api", input)
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_invalid_document(t *testing.T) {
	var testData = []string{
		"/*",
		"/**",
		"/***",
		"/*/",
		"/*/*",
		"/*/**",
	}
	for _, v := range testData {
		s, err := NewScanner("foo.api", v)
		assert.NoError(t, err)
		_, err = s.NextToken()
		assertx.Error(t, err)
	}
}

func TestScanner_NextToken_operator(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.SUB,
			Text: "-",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.MUL,
			Text: "*",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   2,
			},
		},
		{
			Type: token.LPAREN,
			Text: "(",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   3,
			},
		},
		{
			Type: token.LBRACE,
			Text: "{",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   4,
			},
		},
		{
			Type: token.COMMA,
			Text: ",",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   5,
			},
		},
		{
			Type: token.DOT,
			Text: ".",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   6,
			},
		},
		{
			Type: token.RPAREN,
			Text: ")",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   7,
			},
		},
		{
			Type: token.RBRACE,
			Text: "}",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   8,
			},
		},
		{
			Type: token.SEMICOLON,
			Text: ";",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   9,
			},
		},
		{
			Type: token.COLON,
			Text: ":",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   10,
			},
		},
		{
			Type: token.ASSIGN,
			Text: "=",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   11,
			},
		},
		{
			Type: token.ELLIPSIS,
			Text: "...",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   12,
			},
		},
	}
	s, err := NewScanner("foo.api", "-*({,.)};:=...")
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_at(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []token.Token{
			{
				Type: token.AT_DOC,
				Text: "@doc",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   1,
				},
			},
			{
				Type: token.AT_HANDLER,
				Text: "@handler",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   5,
				},
			},
			{
				Type: token.AT_SERVER,
				Text: "@server",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   13,
				},
			},
			{
				Type: token.AT_HANDLER,
				Text: "@handler",
				Position: token.Position{
					Filename: "foo.api",
					Line:     2,
					Column:   1,
				},
			},
			{
				Type: token.AT_SERVER,
				Text: "@server",
				Position: token.Position{
					Filename: "foo.api",
					Line:     3,
					Column:   1,
				},
			},
		}
		s, err := NewScanner("foo.api", "@doc@handler@server\n@handler\n@server")
		assert.NoError(t, err)
		for _, expected := range testData {
			actual, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			"@h",
			"@ha",
			"@han",
			"@hand",
			"@handl",
			"@handle",
			"@handlerr",
			"@hhandler",
			"@foo",
			"@sserver",
			"@serverr",
			"@d",
			"@do",
			"@docc",
		}
		for _, v := range testData {
			s, err := NewScanner("foo.api", v)
			assert.NoError(t, err)
			_, err = s.NextToken()
			assertx.Error(t, err)
		}
	})
}

func TestScanner_NextToken_ident(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.IDENT,
			Text: "foo",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "bar",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   5,
			},
		},
		{
			Type: token.GO,
			Text: "go",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type: token.FUNC,
			Text: "func",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   4,
			},
		},
		{
			Type: token.IDENT,
			Text: "_",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "_go",
			Position: token.Position{
				Filename: "foo.api",
				Line:     4,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "info",
			Position: token.Position{
				Filename: "foo.api",
				Line:     5,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "goo",
			Position: token.Position{
				Filename: "foo.api",
				Line:     6,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "vvar",
			Position: token.Position{
				Filename: "foo.api",
				Line:     6,
				Column:   5,
			},
		},
		{
			Type: token.IDENT,
			Text: "imports",
			Position: token.Position{
				Filename: "foo.api",
				Line:     6,
				Column:   10,
			},
		},
		{
			Type: token.IDENT,
			Text: "go1",
			Position: token.Position{
				Filename: "foo.api",
				Line:     7,
				Column:   1,
			},
		},
	}
	s, err := NewScanner("foo.api", "foo bar\ngo func\n_\n_go\ninfo\ngoo vvar imports\ngo1")
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_Key(t *testing.T) {
	var testData =[]token.Token{
		{
			Type:     token.IDENT,
			Text:     "foo",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type:     token.KEY,
			Text:     "foo:",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type:     token.KEY,
			Text:     "bar:",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
		{
			Type:     token.COLON,
			Text:     ":",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   5,
			},
		},
		{
			Type:     token.INTERFACE,
			Text:     "interface",
			Position: token.Position{
				Filename: "foo.api",
				Line:     4,
				Column:   1,
			},
		},
		{
			Type:     token.ANY,
			Text:     "interface{}",
			Position: token.Position{
				Filename: "foo.api",
				Line:     5,
				Column:   1,
			},
		},
		{
			Type:     token.LBRACE,
			Text:     "{",
			Position: token.Position{
				Filename: "foo.api",
				Line:     5,
				Column:   12,
			},
		},
	}
	s, err := NewScanner("foo.api", "foo\nfoo:\nbar::\ninterface\ninterface{}{")
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_int(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.INT,
			Text: `123`,
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.INT,
			Text: `234`,
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   5,
			},
		},
		{
			Type: token.INT,
			Text: `123`,
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type: token.INT,
			Text: `234`,
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
	}
	s, err := NewScanner("foo.api", "123 234\n123\n234a")
	assert.NoError(t, err)
	for _, expected := range testData {
		actual, err := s.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestScanner_NextToken_string(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []token.Token{
			{
				Type: token.STRING,
				Text: `""`,
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   1,
				},
			},
			{
				Type: token.STRING,
				Text: `"foo"`,
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   3,
				},
			},
			{
				Type: token.STRING,
				Text: `"foo\nbar"`,
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   8,
				},
			},
		}
		s, err := NewScanner("foo.api", `"""foo""foo\nbar"`)
		assert.NoError(t, err)
		for _, expected := range testData {
			actual, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`"`,
			`"foo`,
			`"
`,
		}
		for _, v := range testData {
			s, err := NewScanner("foo.api", v)
			assert.NoError(t, err)
			_, err = s.NextToken()
			assertx.Error(t, err)
		}
	})
}

func TestScanner_NextToken_rawString(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []token.Token{
			{
				Type: token.RAW_STRING,
				Text: "``",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   1,
				},
			},
			{
				Type: token.RAW_STRING,
				Text: "`foo`",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   3,
				},
			},
			{
				Type: token.RAW_STRING,
				Text: "`foo bar`",
				Position: token.Position{
					Filename: "foo.api",
					Line:     1,
					Column:   8,
				},
			},
		}
		s, err := NewScanner("foo.api", "```foo``foo bar`")
		assert.NoError(t, err)
		for _, expected := range testData {
			actual, err := s.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			"`",
			"`foo",
			"`    ",
		}
		for _, v := range testData {
			s, err := NewScanner("foo.api", v)
			assert.NoError(t, err)
			_, err = s.NextToken()
			assertx.Error(t, err)
		}
	})
}

func TestScanner_NextToken_anyCase(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		var testData = []string{
			"#",
			"$",
			"^",
			"好",
			"|",
		}
		for _, v := range testData {
			s, err := NewScanner("foo.api", v)
			assert.NoError(t, err)
			tok, err := s.NextToken()
			assert.NoError(t, err)
			fmt.Println(tok.String())
			assert.Equal(t, token.ILLEGAL, tok.Type)
		}
	})

	t.Run("case2", func(t *testing.T) {
		s, err := NewScanner("foo.api", `好の`)
		assert.NoError(t, err)
		for {
			tok, err := s.NextToken()
			if tok.Type == token.EOF {
				break
			}
			assert.NoError(t, err)
			fmt.Println(tok)
		}
	})

	t.Run("case3", func(t *testing.T) {
		s, err := NewScanner("foo.api", `foo`)
		assert.NoError(t, err)
		for {
			tok, err := s.NextToken()
			if tok.Type == token.EOF {
				break
			}
			assert.NoError(t, err)
			fmt.Println(tok)
		}
	})
}

//go:embed test.api
var testInput string

func TestScanner_NextToken(t *testing.T) {
	position := func(line, column int) token.Position {
		return token.Position{
			Filename: "test.api",
			Line:     line,
			Column:   column,
		}
	}
	var testData = []token.Token{
		{
			Type:     token.IDENT,
			Text:     "syntax",
			Position: position(1, 1),
		},
		{
			Type:     token.ASSIGN,
			Text:     "=",
			Position: position(1, 8),
		},
		{
			Type:     token.STRING,
			Text:     `"v1"`,
			Position: position(1, 10),
		},
		{
			Type:     token.IDENT,
			Text:     `info`,
			Position: position(3, 1),
		},
		{
			Type:     token.LPAREN,
			Text:     `(`,
			Position: position(3, 5),
		},
		{
			Type:     token.KEY,
			Text:     `title:`,
			Position: position(4, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"type title here"`,
			Position: position(4, 12),
		},
		{
			Type:     token.KEY,
			Text:     `desc:`,
			Position: position(5, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"type desc here"`,
			Position: position(5, 11),
		},
		{
			Type:     token.KEY,
			Text:     `author:`,
			Position: position(6, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"type author here"`,
			Position: position(6, 13),
		},
		{
			Type:     token.KEY,
			Text:     `email:`,
			Position: position(7, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"type email here"`,
			Position: position(7, 12),
		},
		{
			Type:     token.KEY,
			Text:     `version:`,
			Position: position(8, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"type version here"`,
			Position: position(8, 14),
		},
		{
			Type:     token.RPAREN,
			Text:     `)`,
			Position: position(9, 1),
		},
		{
			Type:     token.TYPE,
			Text:     `type`,
			Position: position(12, 1),
		},
		{
			Type:     token.IDENT,
			Text:     `request`,
			Position: position(12, 6),
		},
		{
			Type:     token.LBRACE,
			Text:     `{`,
			Position: position(12, 14),
		},
		{
			Type:     token.COMMENT,
			Text:     `// TODO: add members here and delete this comment`,
			Position: position(13, 5),
		},
		{
			Type:     token.RBRACE,
			Text:     `}`,
			Position: position(14, 1),
		},
		{
			Type:     token.TYPE,
			Text:     `type`,
			Position: position(16, 1),
		},
		{
			Type:     token.IDENT,
			Text:     `response`,
			Position: position(16, 6),
		},
		{
			Type:     token.LBRACE,
			Text:     `{`,
			Position: position(16, 15),
		},
		{
			Type:     token.COMMENT,
			Text:     `// TODO: add members here and delete this comment`,
			Position: position(17, 5),
		},
		{
			Type:     token.RBRACE,
			Text:     `}`,
			Position: position(18, 1),
		},
		{
			Type:     token.AT_SERVER,
			Text:     `@server`,
			Position: position(20, 1),
		},
		{
			Type:     token.LPAREN,
			Text:     `(`,
			Position: position(20, 8),
		},
		{
			Type:     token.KEY,
			Text:     `jwt:`,
			Position: position(21, 5),
		},
		{
			Type:     token.IDENT,
			Text:     `Auth`,
			Position: position(21, 10),
		},
		{
			Type:     token.KEY,
			Text:     `group:`,
			Position: position(22, 5),
		},
		{
			Type:     token.IDENT,
			Text:     `template`,
			Position: position(22, 12),
		},
		{
			Type:     token.RPAREN,
			Text:     `)`,
			Position: position(23, 1),
		},
		{
			Type:     token.IDENT,
			Text:     `service`,
			Position: position(24, 1),
		},
		{
			Type:     token.IDENT,
			Text:     `template`,
			Position: position(24, 9),
		},
		{
			Type:     token.LBRACE,
			Text:     `{`,
			Position: position(24, 18),
		},
		{
			Type:     token.AT_DOC,
			Text:     `@doc`,
			Position: position(25, 5),
		},
		{
			Type:     token.STRING,
			Text:     `"foo"`,
			Position: position(25, 10),
		},
		{
			Type:     token.DOCUMENT,
			Text:     `/*foo*/`,
			Position: position(25, 16),
		},
		{
			Type:     token.AT_HANDLER,
			Text:     `@handler`,
			Position: position(26, 5),
		},
		{
			Type:     token.IDENT,
			Text:     `handlerName`,
			Position: position(26, 14),
		},
		{
			Type:     token.COMMENT,
			Text:     `// TODO: replace handler name and delete this comment`,
			Position: position(26, 26),
		},
		{
			Type:     token.IDENT,
			Text:     `get`,
			Position: position(27, 5),
		},
		{
			Type:     token.QUO,
			Text:     `/`,
			Position: position(27, 9),
		},
		{
			Type:     token.IDENT,
			Text:     `users`,
			Position: position(27, 10),
		},
		{
			Type:     token.QUO,
			Text:     `/`,
			Position: position(27, 15),
		},
		{
			Type:     token.IDENT,
			Text:     `id`,
			Position: position(27, 16),
		},
		{
			Type:     token.QUO,
			Text:     `/`,
			Position: position(27, 18),
		},
		{
			Type:     token.COLON,
			Text:     `:`,
			Position: position(27, 19),
		},
		{
			Type:     token.IDENT,
			Text:     `userId`,
			Position: position(27, 20),
		},
		{
			Type:     token.LPAREN,
			Text:     `(`,
			Position: position(27, 27),
		},
		{
			Type:     token.IDENT,
			Text:     `request`,
			Position: position(27, 28),
		},
		{
			Type:     token.RPAREN,
			Text:     `)`,
			Position: position(27, 35),
		},
		{
			Type:     token.IDENT,
			Text:     `returns`,
			Position: position(27, 37),
		},
		{
			Type:     token.LPAREN,
			Text:     `(`,
			Position: position(27, 45),
		},
		{
			Type:     token.IDENT,
			Text:     `response`,
			Position: position(27, 46),
		},
		{
			Type:     token.RPAREN,
			Text:     `)`,
			Position: position(27, 54),
		},
	}
	scanner, err := NewScanner("test.api", testInput)
	assert.NoError(t, err)
	for _, v := range testData {
		actual, err := scanner.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestScanner_NextToken_type(t *testing.T) {
	var testData = []token.Token{
		{
			Type: token.IDENT,
			Text: "foo",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   1,
			},
		},
		{
			Type: token.IDENT,
			Text: "string",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   5,
			},
		},
		{
			Type: token.RAW_STRING,
			Text: "`json:\"foo\"`",
			Position: token.Position{
				Filename: "foo.api",
				Line:     1,
				Column:   12,
			},
		},
		{
			Type: token.IDENT,
			Text: "bar",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   1,
			},
		},
		{
			Type: token.LBRACK,
			Text: "[",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   5,
			},
		},
		{
			Type: token.RBRACK,
			Text: "]",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   6,
			},
		},
		{
			Type: token.IDENT,
			Text: "int",
			Position: token.Position{
				Filename: "foo.api",
				Line:     2,
				Column:   7,
			},
		},
		{
			Type: token.IDENT,
			Text: "baz",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   1,
			},
		},
		{
			Type: token.MAP,
			Text: "map",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   5,
			},
		},
		{
			Type: token.LBRACK,
			Text: "[",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   8,
			},
		},
		{
			Type: token.IDENT,
			Text: "string",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   9,
			},
		},
		{
			Type: token.RBRACK,
			Text: "]",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   15,
			},
		},
		{
			Type: token.IDENT,
			Text: "int",
			Position: token.Position{
				Filename: "foo.api",
				Line:     3,
				Column:   16,
			},
		},
	}
	var input = "foo string `json:\"foo\"`\nbar []int\nbaz map[string]int"
	scanner, err := NewScanner("foo.api", input)
	assert.NoError(t, err)
	for _, v := range testData {
		actual, err := scanner.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}
