package modelgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/modelctl/model"
)

var (
	queryRows = `COLUMN_NAME AS name,ORDINAL_POSITION AS position,DATA_TYPE AS type,COLUMN_KEY AS k,COLUMN_COMMENT AS comment`
)

type (
	Config struct {
		// 是否需要生成redis缓存代码逻辑
		WithCache bool
		// 是否强制覆盖已有文件,如果是将导致原已修改文件找不回
		Force bool
		// mysql访问用户
		Username string
		// mysql访问密码
		Password string
		// mysql连接地址
		Address string
		// 库名
		TableSchema string
		// 待生成model所依赖的表
		Tables []string `json:"Tables,omitempty"`
	}
)

// 生成model相关go文件
func genModelWithConfigFile(path string) error {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var c Config
	err = json.Unmarshal(bts, &c)
	if err != nil {
		return err
	}
	dataSourceTemplate := `{{.Username}}:{{.Password}}@tcp({{.Address}})/information_schema`
	tl, err := template.New("").Parse(dataSourceTemplate)
	if err != nil {
		return err
	}
	var dataSourceBuffer = new(bytes.Buffer)
	err = tl.Execute(dataSourceBuffer, c)
	if err != nil {
		return err
	}
	err = genModelWithDataSource(dataSourceBuffer.String(), c.TableSchema, c.Force, c.WithCache, c.Tables)
	if err != nil {
		return err
	}
	return nil
}

func genModelWithDataSource(dataSource, schema string, force, redis bool, tables []string) error {
	fieldModel := NewFieldModel(dataSource, schema)
	if len(tables) == 0 {
		tableList, err := fieldModel.findTables()
		if err != nil {
			return err
		}
		tables = append(tables, tableList...)
	}
	// 暂定package为model
	packageName := "model"
	utilTemplate := &Template{Package: packageName, WithCache: redis}

	err := generateUtilModel(force, utilTemplate)
	if err != nil {
		return err
	}
	for _, table := range tables {
		fieldList, err := fieldModel.findColumns(table)
		if err != nil {
			return err
		}
		modelTemplate, err := generateModelTemplate(packageName, table, fieldList)
		if err != nil {
			return err
		}
		modelTemplate.WithCache = redis
		err = generateSqlModel(force, modelTemplate)
		if err != nil {
			return err
		}
	}
	fmt.Println("model generate done ...")
	return nil
}

// 生成util model
func generateUtilModel(force bool, data *Template) error {
	tl, err := template.New("").Parse(utilTemplateText)
	if err != nil {
		return err
	}
	fileName := "util.go"
	_, err = os.Stat(fileName)
	if err == nil {
		if !force {
			return nil
		}
		os.Remove(fileName)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, model.ModeDirPerm)
	if err != nil {
		return err
	}
	defer file.Close()
	err = tl.Execute(file, data)
	if err != nil {
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 生成sql对应model
func generateSqlModel(force bool, data *Template) error {
	tl, err := template.New("").Parse(modelTemplateText)
	if err != nil {
		return err
	}
	fileName := strings.ToLower(data.ModelCamelWithLowerStart + "model.go")
	_, err = os.Stat(fileName)
	if err == nil {
		if !force {
			fmt.Println(fileName + " already exists")
			return nil
		}
		os.Remove(fileName)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	err = tl.Execute(file, data)
	if err != nil {
		return err
	}
	cmd := exec.Command("goimports", "-w", fileName)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
