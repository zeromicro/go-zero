package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const defaultSyntax = `syntax = "v1"`

type Parser struct {
	r   *bufio.Reader
	api *ApiStruct
}

func NewParser(filename string) (*Parser, error) {
	apiAbsPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	api, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	apiStruct, err := ParseApi(string(api))
	if err != nil {
		return nil, err
	}

	apiSyntax := strings.TrimSpace(apiStruct.Syntax)
	if len(apiSyntax) == 0 {
		apiStruct.Syntax = defaultSyntax
	}

	for _, item := range strings.Split(apiStruct.Imports, "\n") {
		importLine := strings.TrimSpace(item)
		if len(importLine) > 0 {
			item := strings.TrimPrefix(importLine, "import")
			item = strings.TrimSpace(item)
			item = strings.TrimPrefix(item, `"`)
			item = strings.TrimSuffix(item, `"`)
			var path = item
			if !util.FileExists(item) {
				path = filepath.Join(filepath.Dir(apiAbsPath), item)
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, errors.New("import api file not exist: " + item)
			}

			importStruct, err := ParseApi(string(content))
			if err != nil {
				return nil, err
			}

			if len(importStruct.Imports) > 0 {
				return nil, errors.New("import api should not import another api file recursive")
			}
			syntax := strings.TrimSpace(importStruct.Syntax)
			if len(syntax) == 0 {
				importStruct.Syntax = defaultSyntax
			}

			if apiStruct.Syntax != importStruct.Syntax {
				return nil, fmt.Errorf("inconsistent api syntax: %s,%s", apiStruct.Syntax, importStruct.Syntax)
			}

			apiStruct.Type += "\n" + importStruct.Type
			apiStruct.Service += "\n" + importStruct.Service
		}
	}

	if len(strings.TrimSpace(apiStruct.Service)) == 0 {
		return nil, errors.New("api has no service defined")
	}

	var buffer = new(bytes.Buffer)
	buffer.WriteString(apiStruct.Service)
	return &Parser{
		r:   bufio.NewReader(buffer),
		api: apiStruct,
	}, nil
}

func (p *Parser) Parse() (api *spec.ApiSpec, err error) {
	api = new(spec.ApiSpec)
	var sp = StructParser{Src: p.api.Type}
	types, err := sp.Parse()
	if err != nil {
		return nil, err
	}

	syntaxState := newSyntaxState(p.api.Syntax)
	version, err := syntaxState.process()
	if err != nil {
		return nil, err
	}

	api.Syntax = version

	api.Types = types
	var lineNumber = p.api.serviceBeginLine
	st := newRootState(p.r, &lineNumber)
	for {
		st, err = st.process(api)
		if err == io.EOF {
			return api, p.validate(api)
		}
		if err != nil {
			return nil, fmt.Errorf("near line: %d, %s", lineNumber, err.Error())
		}
		if st == nil {
			return api, p.validate(api)
		}
	}
}
