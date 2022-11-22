package format

import (
	"io"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func Format(source []byte, w io.Writer) error {
	p := parser.New("", source)
	result := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return err
	}

	result.Format(w)
	return nil
}

func formatForUnitTest(source []byte, w io.Writer) error {
	p := parser.New("", source)
	result := p.Parse()
	if err := p.CheckErrors(); err != nil {
		return err
	}

	result.FormatForUnitTest(w)
	return nil
}
