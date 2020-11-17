package format

import (
	"bufio"
	"errors"
	"fmt"
	"go/scanner"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/urfave/cli"
)

const (
	leftParenthesis  = "("
	rightParenthesis = ")"
	leftBrace        = "{"
	rightBrace       = "}"
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
	_, err := parser.ParseApi(data)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	s := bufio.NewScanner(strings.NewReader(data))
	var tapCount = 0
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		noCommentLine := util.RemoveComment(line)
		if noCommentLine == rightParenthesis || noCommentLine == rightBrace {
			tapCount -= 1
		}
		if tapCount < 0 {
			line = strings.TrimSuffix(line, rightBrace)
			line = strings.TrimSpace(line)
			if strings.HasSuffix(line, leftBrace) {
				tapCount += 1
			}
		}
		util.WriteIndent(&builder, tapCount)
		builder.WriteString(line + "\n")
		if strings.HasSuffix(noCommentLine, leftParenthesis) || strings.HasSuffix(noCommentLine, leftBrace) {
			tapCount += 1
		}
	}
	return strings.TrimSpace(builder.String()), nil
}
