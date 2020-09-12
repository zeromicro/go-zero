package format

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/urfave/cli"
)

var (
	reg = regexp.MustCompile("type (?P<name>.*)[\\s]+{")
)

func GoFormatApi(c *cli.Context) error {
	dir := c.String("dir")
	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	printToConsole := c.Bool("p")

	var be errorx.BatchError
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, errBack error) (err error) {
		if strings.HasSuffix(path, ".api") {
			err := ApiFormat(path, printToConsole)
			if err != nil {
				be.Add(util.WrapErr(err, fi.Name()))
			}
		}
		return nil
	})
	be.Add(err)
	if be.NotNil() {
		errs := be.Err().Error()
		fmt.Println(errs)
		os.Exit(1)
	}
	return be.Err()
}

func ApiFormat(path string, printToConsole bool) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	r := reg.ReplaceAllStringFunc(string(data), func(m string) string {
		parts := reg.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}
		if !strings.Contains(m, "struct") {
			return "type " + parts[1] + " struct {"
		}
		return m
	})

	info, st, service, err := parser.MatchStruct(r)
	if err != nil {
		return err
	}
	info = strings.TrimSpace(info)
	if len(service) == 0 || len(st) == 0 {
		return nil
	}

	fs, err := format.Source([]byte(strings.TrimSpace(st)))
	if err != nil {
		str := err.Error()
		lineNumber := strings.Index(str, ":")
		if lineNumber > 0 {
			ln, err := strconv.ParseInt(str[:lineNumber], 10, 64)
			if err != nil {
				return err
			}
			pn := 0
			if len(info) > 0 {
				pn = countRune(info, '\n') + 1
			}
			number := int(ln) + pn + 1
			return errors.New(fmt.Sprintf("line: %d, %s", number, str[lineNumber+1:]))
		}
		return err
	}

	result := strings.Join([]string{info, string(fs), service}, "\n\n")
	if printToConsole {
		_, err := fmt.Print(result)
		return err
	}
	result = strings.TrimSpace(result)
	return ioutil.WriteFile(path, []byte(result), os.ModePerm)
}

func countRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}
