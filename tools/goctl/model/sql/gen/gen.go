package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	pwd             = "."
	createTableFlag = `(?m)^(?i)CREATE\s+TABLE` // ignore case
)

type (
	defaultGenerator struct {
		source string
		dir    string
		console.Console
	}
	Option func(generator *defaultGenerator)
)

func NewDefaultGenerator(source, dir string, opt ...Option) *defaultGenerator {
	if dir == "" {
		dir = pwd
	}
	generator := &defaultGenerator{source: source, dir: dir}
	var optionList []Option
	optionList = append(optionList, newDefaultOption())
	optionList = append(optionList, opt...)
	for _, fn := range optionList {
		fn(generator)
	}
	return generator
}

func WithConsoleOption(c console.Console) Option {
	return func(generator *defaultGenerator) {
		generator.Console = c
	}
}

func newDefaultOption() Option {
	return func(generator *defaultGenerator) {
		generator.Console = console.NewColorConsole()
	}
}

func (g *defaultGenerator) Start(withCache bool) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}
	err = util.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}
	modelList, err := g.genFromDDL(withCache)
	if err != nil {
		return err
	}

	for tableName, code := range modelList {
		name := fmt.Sprintf("%smodel.go", strings.ToLower(stringx.From(tableName).ToCamel()))
		filename := filepath.Join(dirAbs, name)
		if util.FileExists(filename) {
			g.Warning("%s already exists, ignored.", name)
			continue
		}
		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}
	// generate error file
	filename := filepath.Join(dirAbs, "vars.go")
	if !util.FileExists(filename) {
		err = ioutil.WriteFile(filename, []byte(template.Error), os.ModePerm)
		if err != nil {
			return err
		}
	}
	g.Success("Done.")
	return nil
}

// ret1: key-table name,value-code
func (g *defaultGenerator) genFromDDL(withCache bool) (map[string]string, error) {
	ddlList := g.split()
	m := make(map[string]string)
	for _, ddl := range ddlList {
		table, err := parser.Parse(ddl)
		if err != nil {
			return nil, err
		}
		code, err := g.genModel(*table, withCache)
		if err != nil {
			return nil, err
		}
		m[table.Name.Source()] = code
	}
	return m, nil
}

type (
	Table struct {
		parser.Table
		CacheKey map[string]Key
	}
)

func (g *defaultGenerator) genModel(in parser.Table, withCache bool) (string, error) {
	t := util.With("model").
		Parse(template.Model).
		GoFmt(true)

	m, err := genCacheKeys(in)
	if err != nil {
		return "", err
	}

	importsCode, err := genImports(withCache, in.ContainsTime())
	if err != nil {
		return "", err
	}

	var table Table
	table.Table = in
	table.CacheKey = m

	varsCode, err := genVars(table, withCache)
	if err != nil {
		return "", err
	}

	typesCode, err := genTypes(table, withCache)
	if err != nil {
		return "", err
	}

	newCode, err := genNew(table, withCache)
	if err != nil {
		return "", err
	}

	insertCode, err := genInsert(table, withCache)
	if err != nil {
		return "", err
	}

	var findCode = make([]string, 0)
	findOneCode, err := genFindOne(table, withCache)
	if err != nil {
		return "", err
	}

	findOneByFieldCode, err := genFineOneByField(table, withCache)
	if err != nil {
		return "", err
	}

	findCode = append(findCode, findOneCode, findOneByFieldCode)
	updateCode, err := genUpdate(table, withCache)
	if err != nil {
		return "", err
	}

	deleteCode, err := genDelete(table, withCache)
	if err != nil {
		return "", err
	}

	output, err := t.Execute(map[string]interface{}{
		"imports": importsCode,
		"vars":    varsCode,
		"types":   typesCode,
		"new":     newCode,
		"insert":  insertCode,
		"find":    strings.Join(findCode, "\n"),
		"update":  updateCode,
		"delete":  deleteCode,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
