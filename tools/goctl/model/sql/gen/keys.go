package gen

import (
	"bytes"
	"strings"
	"text/template"
)

var (
	cacheKeyExpressionTemplate = `cache{{.upperCamelTable}}{{.upperCamelField}}Prefix = "cache#{{.lowerCamelTable}}#{{.lowerCamelField}}#"`
	keyTemplate                = `{{.lowerCamelField}}Key := fmt.Sprintf("%s%v", {{.define}}, {{.lowerCamelField}})`
	keyRespTemplate            = `{{.lowerCamelField}}Key := fmt.Sprintf("%s%v", {{.define}}, resp.{{.upperCamelField}})`
	keyDataTemplate            = `{{.lowerCamelField}}Key := fmt.Sprintf("%s%v", {{.define}}, data.{{.upperCamelField}})`
)

type (
	Key struct {
		Define      string // cacheKey define,如：cacheUserUserIdPrefix
		Value       string // cacheKey value expression,如：cache#user#userId#
		Expression  string // cacheKey expression，如:cacheUserUserIdPrefix="cache#user#userId#"
		KeyVariable string // cacheKey 声明变量，如：userIdKey
		Key         string // 缓存key的代码,如 userIdKey:=fmt.Sprintf("%s%v", cacheUserUserIdPrefix, userId)
		DataKey     string // 缓存key的代码,如 userIdKey:=fmt.Sprintf("%s%v", cacheUserUserIdPrefix, data.userId)
		RespKey     string // 缓存key的代码,如 userIdKey:=fmt.Sprintf("%s%v", cacheUserUserIdPrefix, resp.userId)
	}
)

// key-数据库原始字段名,value-缓存key对象
func genCacheKeys(table *InnerTable) (map[string]Key, error) {
	fields := table.Fields
	var m = make(map[string]Key)
	if !table.ContainsCache {
		return m, nil
	}
	for _, field := range fields {
		if !field.Cache && !field.IsPrimaryKey {
			continue
		}
		t, err := template.New("keyExpression").Parse(cacheKeyExpressionTemplate)
		if err != nil {
			return nil, err
		}
		var expressionBuffer = new(bytes.Buffer)
		err = t.Execute(expressionBuffer, map[string]string{
			"upperCamelTable": table.UpperCamelCase,
			"lowerCamelTable": table.LowerCamelCase,
			"upperCamelField": field.UpperCamelCase,
			"lowerCamelField": field.LowerCamelCase,
		})
		if err != nil {
			return nil, err
		}
		expression := expressionBuffer.String()
		expressionAr := strings.Split(expression, "=")
		define := strings.TrimSpace(expressionAr[0])
		value := strings.TrimSpace(expressionAr[1])
		t, err = template.New("key").Parse(keyTemplate)
		if err != nil {
			return nil, err
		}
		var keyBuffer = new(bytes.Buffer)
		err = t.Execute(keyBuffer, map[string]string{
			"lowerCamelField": field.LowerCamelCase,
			"define":          define,
		})
		if err != nil {
			return nil, err
		}
		t, err = template.New("keyData").Parse(keyDataTemplate)
		if err != nil {
			return nil, err
		}
		var keyDataBuffer = new(bytes.Buffer)
		err = t.Execute(keyDataBuffer, map[string]string{
			"lowerCamelField": field.LowerCamelCase,
			"upperCamelField": field.UpperCamelCase,
			"define":          define,
		})
		if err != nil {
			return nil, err
		}
		t, err = template.New("keyResp").Parse(keyRespTemplate)
		if err != nil {
			return nil, err
		}
		var keyRespBuffer = new(bytes.Buffer)
		err = t.Execute(keyRespBuffer, map[string]string{
			"lowerCamelField": field.LowerCamelCase,
			"upperCamelField": field.UpperCamelCase,
			"define":          define,
		})
		if err != nil {
			return nil, err
		}
		m[field.SnakeCase] = Key{
			Define:      define,
			Value:       value,
			Expression:  expression,
			KeyVariable: field.LowerCamelCase + "Key",
			Key:         keyBuffer.String(),
			DataKey:     keyDataBuffer.String(),
			RespKey:     keyRespBuffer.String(),
		}
	}
	return m, nil
}
