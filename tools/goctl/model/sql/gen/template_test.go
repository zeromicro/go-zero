package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/template"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func TestGenTemplates(t *testing.T) {
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)
	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)
	file := filepath.Join(dir, "model-new.tpl")
	data, err := os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, string(data), template.New)
}

func TestRevertTemplate(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, RevertTemplate(name))

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}

func TestClean(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	assert.Nil(t, Clean())

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	_, err = os.ReadFile(file)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	name := "model-new.tpl"
	err := pathx.InitTemplates(category, templates)
	assert.Nil(t, err)

	dir, err := pathx.GetTemplateDir(category)
	assert.Nil(t, err)

	file := filepath.Join(dir, name)
	data, err := os.ReadFile(file)
	assert.Nil(t, err)

	modifyData := string(data) + "modify"
	err = pathx.CreateTemplate(category, name, modifyData)
	assert.Nil(t, err)

	data, err = os.ReadFile(file)
	assert.Nil(t, err)

	assert.Equal(t, string(data), modifyData)

	assert.Nil(t, Update())

	data, err = os.ReadFile(file)
	assert.Nil(t, err)
	assert.Equal(t, template.New, string(data))
}
