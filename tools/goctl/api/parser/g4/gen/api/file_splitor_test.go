// DO NOT EDIT.
// Tool: split apiparser_parser.go
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.
package api

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSplitor(t *testing.T) {
	t.Skip("skip this test because it is used to split the apiparser_parser.go file by developer.")
	dir := "."
	data, err := os.ReadFile(filepath.Join(dir, "apiparser_parser.go"))
	if err != nil {
		log.Fatalln(err)
	}

	r := bytes.NewReader(data)
	reader := bufio.NewReader(r)
	var lines, files int
	buffer := bytes.NewBuffer(nil)

	for {
		fn, part := "apiparser_parser0.go", "main"
		if files > 0 {
			fn = fmt.Sprintf("apiparser_parser%d.go", files)
			part = fmt.Sprintf("%d", files)
		}
		fp := filepath.Join(dir, fn)
		buffer.Reset()
		if files > 0 {
			buffer.WriteString(fmt.Sprintf(`package api
import "github.com/zeromicro/antlr"

// Part %s
// The apiparser_parser.go file was split into multiple files because it
// was too large and caused a possible memory overflow during goctl installation.
`, part))
		}

		var exit bool
		for {
			line, _, err := reader.ReadLine()
			buffer.Write(line)
			buffer.WriteRune('\n')
			if err != nil {
				fmt.Printf("%+v\n", err)
				exit = true
				break
			}
			lines += 1
			if string(line) == "}" && lines >= 650 {
				break
			}
		}

		src, err := format.Source(buffer.Bytes())
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		err = os.WriteFile(fp, src, os.ModePerm)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if exit {
			break
		}
		lines = 0
		files += 1
	}

	err = os.Rename(filepath.Join(dir, "apiparser_parser0.go"), filepath.Join(dir, "apiparser_parser.go"))
	if err != nil {
		log.Fatalln(err)
	}
}
