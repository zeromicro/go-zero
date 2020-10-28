package format

import (
	"errors"
	"fmt"
	"go/format"
	"go/scanner"
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
	useStdin := c.Bool("stdin")

	var be errorx.BatchError
	if useStdin {
		if err := ApiFormatByStdin(); err != nil {
			be.Add(err)
		}
	} else {
		dir := c.String("dir")
		if len(dir) == 0 {
			return errors.New("missing -dir")
		}

		_, err := os.Lstat(dir)
		if err != nil {
			return errors.New(dir + ": No such file or directory")
		}

		err = filepath.Walk(dir, func(path string, fi os.FileInfo, errBack error) (err error) {
			if strings.HasSuffix(path, ".api") {
				if err := ApiFormatByPath(path); err != nil {
					be.Add(util.WrapErr(err, fi.Name()))
				}
			}
			return nil
		})
		be.Add(err)
	}
	if be.NotNil() {
		scanner.PrintError(os.Stderr, be.Err())
		os.Exit(1)
	}
	return be.Err()
}

func ApiFormatByStdin() error {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	result, err := apiFormat(string(data))
	if err != nil {
		return err
	}

	_, err = fmt.Print(result)
	if err != nil {
		return err
	}
	return nil
}

func ApiFormatByPath(apiFilePath string) error {
	data, err := ioutil.ReadFile(apiFilePath)
	if err != nil {
		return err
	}

	result, err := apiFormat(string(data))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(apiFilePath, []byte(result), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func apiFormat(data string) (string, error) {

	r := reg.ReplaceAllStringFunc(data, func(m string) string {
		parts := reg.FindStringSubmatch(m)
		if len(parts) < 2 {
			return m
		}
		if !strings.Contains(m, "struct") {
			return "type " + parts[1] + " struct {"
		}
		return m
	})

	apiStruct, err := parser.ParseApi(r)
	if err != nil {
		return "", err
	}
	info := strings.TrimSpace(apiStruct.Info)
	if len(apiStruct.Service) == 0 {
		return data, nil
	}

	fs, err := format.Source([]byte(strings.TrimSpace(apiStruct.StructBody)))
	if err != nil {
		str := err.Error()
		lineNumber := strings.Index(str, ":")
		if lineNumber > 0 {
			ln, err := strconv.ParseInt(str[:lineNumber], 10, 64)
			if err != nil {
				return "", err
			}
			pn := 0
			if len(info) > 0 {
				pn = countRune(info, '\n') + 1
			}
			number := int(ln) + pn + 1
			return "", errors.New(fmt.Sprintf("line: %d, %s", number, str[lineNumber+1:]))
		}
		return "", err
	}

	var result string
	if len(strings.TrimSpace(info)) > 0 {
		result += strings.TrimSpace(info) + "\n\n"
	}
	if len(strings.TrimSpace(apiStruct.Imports)) > 0 {
		result += strings.TrimSpace(apiStruct.Imports) + "\n\n"
	}
	if len(strings.TrimSpace(string(fs))) > 0 {
		result += strings.TrimSpace(string(fs)) + "\n\n"
	}
	if len(strings.TrimSpace(apiStruct.Service)) > 0 {
		result += strings.TrimSpace(apiStruct.Service) + "\n\n"
	}

	return result, nil
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
