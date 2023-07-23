# 配置项管理

| 名称              | 是否可选 | 说明                                          |
|-------------------|----------|-----------------------------------------------|
| namingFormat      | YES      | 文件名称格式化符                      |

# naming-format
`namingFormat`可以用于对生成代码的文件名称进行格式化，和日期格式化符（yyyy-MM-dd）类似，在代码生成时可以根据这些配置项的格式化符进行格式化。

## 格式化符(gozero)
格式化符由`go`,`zero`组成，如常见的三种格式化风格你可以这样编写：
* lower: `gozero`
* camel: `goZero`
* snake: `go_zero`

常见格式化符生成示例
源字符：welcome_to_go_zero

| 格式化符   | 格式化结果            | 说明                      |
|------------|-----------------------|---------------------------|
| gozero     | welcometogozero       | 小写                      |
| goZero     | welcomeToGoZero       | 驼峰                      |
| go_zero    | welcome_to_go_zero    | snake                     |
| Go#zero    | Welcome#to#go#zero    | #号分割Title类型          |
| GOZERO     | WELCOMETOGOZERO       | 大写                      |
| \_go#zero_ | \_welcome#to#go#zero_ | 下划线做前后缀，并且#分割 |

错误格式化符示例
* go
* gOZero
* zero
* goZEro
* goZERo
* goZeRo
* tal

# 使用方法
目前可通过在生成api、rpc、model时通过`--style`参数指定format格式，如：
```shell script
goctl api go test.api -dir . -style gozero
```
```shell script
 goctl rpc proto -src test.proto -dir . -style go_zero
```
```shell script
goctl model mysql datasource -url="" -table="*" -dir ./snake -style GoZero
```

# 默认值
当不指定-style时默认值为`gozero`
