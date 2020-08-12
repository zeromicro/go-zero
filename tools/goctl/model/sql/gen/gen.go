package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

const (
	pwd             = "."
	createTableFlag = `(?m)CREATE\s+TABLE`
)

type (
	defaultGenerator struct {
		source string
		src    string
		dir    string
	}
)

func NewDefaultGenerator(src, dir string) *defaultGenerator {
	if src == "" {
		src = pwd
	}
	if dir == "" {
		dir = pwd
	}
	return &defaultGenerator{src: src, dir: dir}
}

func (g *defaultGenerator) Start(withCache bool) error {
	fileSrc, err := filepath.Abs(g.src)
	if err != nil {
		return err
	}
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}
	err = util.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(fileSrc)
	if err != nil {
		return err
	}
	g.source = string(data)
	modelList, err := g.genFromDDL(withCache)
	if err != nil {
		return err
	}

	for tableName, code := range modelList {
		filename := filepath.Join(dirAbs, fmt.Sprintf("%smodel.go", stringx.From(tableName).Lower()))
		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}
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
	t := templatex.With("model").
		Parse(sqltemplate.Model).
		GoFmt(true)

	m, err := genCacheKeys(in)
	if err != nil {
		return "", err
	}
	importsCode, err := genImports(withCache)
	if err != nil {
		return "", err
	}
	var table Table
	table.Table = in
	table.CacheKey = m

	varsCode, err := genVars(table)
	if err != nil {
		return "", err
	}
	typesCode, err := genTypes(table)
	if err != nil {
		return "", err
	}
	newCode, err := genNew(table)
	if err != nil {
		return "", err
	}
	insertCode, err := genInsert(table)
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
		"find":    strings.Join(findCode, "\r\n"),
		"update":  updateCode,
		"delete":  deleteCode,
	})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
