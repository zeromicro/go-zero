package gen

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"

	"zero/core/logx"
	sqltemplate "zero/tools/goctl/model/sql/template"
)

func GenModel(table *InnerTable) (string, error) {
	t, err := template.New("model").Parse(sqltemplate.Model)
	if err != nil {
		return "", nil
	}
	modelBuffer := new(bytes.Buffer)
	importsCode, err := genImports(table)
	if err != nil {
		return "", err
	}
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
	findOneCode, err := genFindOne(table)
	if err != nil {
		return "", err
	}
	findOneByFieldCode, err := genFineOneByField(table)
	if err != nil {
		return "", err
	}
	findAllCode, err := genFindAllByField(table)
	if err != nil {
		return "", err
	}
	findLimitCode, err := genFindLimitByField(table)
	if err != nil {
		return "", err
	}
	findCode = append(findCode, findOneCode, findOneByFieldCode, findAllCode, findLimitCode)
	updateCode, err := genUpdate(table)
	if err != nil {
		return "", err
	}
	deleteCode, err := genDelete(table)
	if err != nil {
		return "", err
	}

	err = t.Execute(modelBuffer, map[string]interface{}{
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
	result := modelBuffer.String()
	bts, err := format.Source([]byte(result))
	if err != nil {
		logx.Errorf("%+v", err)
		return "", err
	}
	return string(bts), nil
}
