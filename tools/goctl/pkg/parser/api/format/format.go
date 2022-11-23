package format

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func File(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(nil)
	if err := Source(data, buffer); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buffer.Bytes(), 0666)
}

func Source(source []byte, w io.Writer) error {
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
