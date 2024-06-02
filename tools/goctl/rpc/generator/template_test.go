package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func TestGenTemplates(t *testing.T) {
	_ = Clean()
	err := GenTemplates()
	assert.Nil(t, err)
}

func TestRevertTemplate(t *testing.T) {
	_ = Clean()
	err := GenTemplates()
	assert.Nil(t, err)
	fp, err := pathx.GetTemplateDir(category)
	if err != nil {
		return
	}
	mainTpl := filepath.Join(fp, mainTemplateFile)
	data, err := os.ReadFile(mainTpl)
	if err != nil {
		return
	}
	assert.Equal(t, templates[mainTemplateFile], string(data))

	err = RevertTemplate("test")
	if err != nil {
		assert.Equal(t, "test: no such file name", err.Error())
	}

	err = os.WriteFile(mainTpl, []byte("modify"), os.ModePerm)
	if err != nil {
		return
	}

	data, err = os.ReadFile(mainTpl)
	if err != nil {
		return
	}
	assert.Equal(t, "modify", string(data))

	err = RevertTemplate(mainTemplateFile)
	assert.Nil(t, err)

	data, err = os.ReadFile(mainTpl)
	if err != nil {
		return
	}
	assert.Equal(t, templates[mainTemplateFile], string(data))
}

func TestClean(t *testing.T) {
	_ = Clean()
	err := GenTemplates()
	assert.Nil(t, err)
	fp, err := pathx.GetTemplateDir(category)
	if err != nil {
		return
	}
	mainTpl := filepath.Join(fp, mainTemplateFile)
	_, err = os.Stat(mainTpl)
	assert.Nil(t, err)

	err = Clean()
	assert.Nil(t, err)

	_, err = os.Stat(mainTpl)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	_ = Clean()
	err := GenTemplates()
	assert.Nil(t, err)
	fp, err := pathx.GetTemplateDir(category)
	if err != nil {
		return
	}
	mainTpl := filepath.Join(fp, mainTemplateFile)

	err = os.WriteFile(mainTpl, []byte("modify"), os.ModePerm)
	if err != nil {
		return
	}

	data, err := os.ReadFile(mainTpl)
	if err != nil {
		return
	}
	assert.Equal(t, "modify", string(data))

	assert.Nil(t, Update())

	data, err = os.ReadFile(mainTpl)
	if err != nil {
		return
	}
	assert.Equal(t, templates[mainTemplateFile], string(data))
}

func TestGetCategory(t *testing.T) {
	_ = Clean()
	result := Category()
	assert.Equal(t, category, result)
}
