package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type Parser struct {
	r  *bufio.Reader
	st string
}

func NewParser(filename string) (*Parser, error) {
	api, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	info, body, service, err := MatchStruct(string(api))
	if err != nil {
		return nil, err
	}
	var buffer = new(bytes.Buffer)
	buffer.WriteString(info)
	buffer.WriteString(service)
	return &Parser{
		r:  bufio.NewReader(buffer),
		st: body,
	}, nil
}

func (p *Parser) Parse() (api *spec.ApiSpec, err error) {
	api = new(spec.ApiSpec)
	types, err := parseStructAst(p.st)
	if err != nil {
		return nil, err
	}
	api.Types = types
	var lineNumber = 1
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
