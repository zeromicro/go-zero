package format

import (
	"bufio"
	"errors"
	"fmt"
	"go/format"
	"go/scanner"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
	apiF "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	leftParenthesis  = "("
	rightParenthesis = ")"
	leftBrace        = "{"
	rightBrace       = "}"
)

var (
	// VarBoolUseStdin describes whether to use stdin or not.
	VarBoolUseStdin bool
	// VarBoolSkipCheckDeclare describes whether to skip.
	VarBoolSkipCheckDeclare bool
	// VarStringDir describes the directory.
	VarStringDir string
	// VarBoolIgnore describes whether to ignore.
	VarBoolIgnore bool
)

// GoFormatApi format api file
func GoFormatApi(_ *cobra.Command, _ []string) error {
	var be errorx.BatchError
	if VarBoolUseStdin {
		if err := apiFormatReader(os.Stdin, VarStringDir, VarBoolSkipCheckDeclare); err != nil {
			be.Add(err)
		}
	} else {
		if len(VarStringDir) == 0 {
			return errors.New("missing -dir")
		}

		_, err := os.Lstat(VarStringDir)
		if err != nil {
			return errors.New(VarStringDir + ": No such file or directory")
		}

		err = filepath.Walk(VarStringDir, func(path string, fi os.FileInfo, errBack error) (err error) {
			if strings.HasSuffix(path, ".api") {
				if err := ApiFormatByPath(path, VarBoolSkipCheckDeclare); err != nil {
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

// apiFormatReader
// filename is needed when there are `import` literals.
func apiFormatReader(reader io.Reader, filename string, skipCheckDeclare bool) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	result, err := apiFormat(string(data), skipCheckDeclare, filename)
	if err != nil {
		return err
	}

	_, err = fmt.Print(result)
	return err
}

// ApiFormatByPath format api from file path
func ApiFormatByPath(apiFilePath string, skipCheckDeclare bool) error {
	if env.UseExperimental() {
		return apiF.File(apiFilePath)
	}

	data, err := os.ReadFile(apiFilePath)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(apiFilePath)
	if err != nil {
		return err
	}

	result, err := apiFormat(string(data), skipCheckDeclare, abs)
	if err != nil {
		return err
	}

	_, err = parser.ParseContentWithParserSkipCheckTypeDeclaration(result, abs)
	if err != nil {
		return err
	}

	return os.WriteFile(apiFilePath, []byte(result), os.ModePerm)
}

func apiFormat(data string, skipCheckDeclare bool, filename ...string) (string, error) {
	var err error
	if skipCheckDeclare {
		_, err = parser.ParseContentWithParserSkipCheckTypeDeclaration(data, filename...)
	} else {
		_, err = parser.ParseContent(data, filename...)
	}
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	s := bufio.NewScanner(strings.NewReader(data))
	tapCount := 0
	newLineCount := 0
	var preLine string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			if newLineCount > 0 {
				continue
			}
			newLineCount++
		} else {
			if preLine == rightBrace {
				builder.WriteString(pathx.NL)
			}
			newLineCount = 0
		}

		if tapCount == 0 {
			ft, err := formatGoTypeDef(line, s, &builder)
			if err != nil {
				return "", err
			}

			if ft {
				continue
			}
		}

		noCommentLine := util.RemoveComment(line)
		if noCommentLine == rightParenthesis || noCommentLine == rightBrace {
			tapCount--
		}
		if tapCount < 0 {
			line := strings.TrimSuffix(noCommentLine, rightBrace)
			line = strings.TrimSpace(line)
			if strings.HasSuffix(line, leftBrace) {
				tapCount++
			}
		}
		if line != "" {
			util.WriteIndent(&builder, tapCount)
		}
		builder.WriteString(line + pathx.NL)
		if strings.HasSuffix(noCommentLine, leftParenthesis) || strings.HasSuffix(noCommentLine, leftBrace) {
			tapCount++
		}
		preLine = line
	}

	return strings.TrimSpace(builder.String()), nil
}

func formatGoTypeDef(line string, scanner *bufio.Scanner, builder *strings.Builder) (bool, error) {
	noCommentLine := util.RemoveComment(line)
	tokenCount := 0
	if strings.HasPrefix(noCommentLine, "type") && (strings.HasSuffix(noCommentLine, leftParenthesis) ||
		strings.HasSuffix(noCommentLine, leftBrace)) {
		var typeBuilder strings.Builder
		typeBuilder.WriteString(mayInsertStructKeyword(line, &tokenCount) + pathx.NL)
		for scanner.Scan() {
			noCommentLine := util.RemoveComment(scanner.Text())
			typeBuilder.WriteString(mayInsertStructKeyword(scanner.Text(), &tokenCount) + pathx.NL)
			if noCommentLine == rightBrace || noCommentLine == rightParenthesis {
				tokenCount--
			}
			if tokenCount == 0 {
				ts, err := format.Source([]byte(typeBuilder.String()))
				if err != nil {
					return false, errors.New("error format \n" + typeBuilder.String())
				}

				result := strings.ReplaceAll(string(ts), " struct ", " ")
				result = strings.ReplaceAll(result, "type ()", "")
				builder.WriteString(result)
				break
			}
		}
		return true, nil
	}

	return false, nil
}

func mayInsertStructKeyword(line string, token *int) string {
	insertStruct := func() string {
		if strings.Contains(line, " struct") {
			return line
		}
		index := strings.Index(line, leftBrace)
		return line[:index] + " struct " + line[index:]
	}

	noCommentLine := util.RemoveComment(line)
	if strings.HasSuffix(noCommentLine, leftBrace) {
		*token++
		return insertStruct()
	}
	if strings.HasSuffix(noCommentLine, rightBrace) {
		noCommentLine = strings.TrimSuffix(noCommentLine, rightBrace)
		noCommentLine = util.RemoveComment(noCommentLine)
		if strings.HasSuffix(noCommentLine, leftBrace) {
			return insertStruct()
		}
	}
	if strings.HasSuffix(noCommentLine, leftParenthesis) {
		*token++
	}

	if strings.Contains(noCommentLine, "`") {
		return util.UpperFirst(strings.TrimSpace(line))
	}

	return line
}
