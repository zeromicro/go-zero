# cztctl

基于 go-zero 的 goctl工具魔改的工具包

cztctl会同步goctl最新功能

# 安装工具

如果本地有goctl,安装cztctl报错，请把go-zero和goctl升级到v1.9.0

```shell

# 设置代理
$ go env -W GOPROXY=https://goproxy.cn/,direct 
# 安装cztctl
$ go install github.com/lerity-yao/go-zero/tools/cztctl@latest

```


# 环境配置

```
#查看gopath
$ go env GOPATH
```

将$GOPATH/bin中的 cztctl 添加到环境变量

# 检查版本

```
$ cztctl --version
cztctl version 1.9.0-alpha linux/amd64
```

# 命令

```SHELL
$ cztctl
A cli tool to generate api, zrpc, model code

GitHub: https://github.com/lerity-yao/go-zero
Site:   https://go-zero.dev

Usage:
  cztctl [command]

Available Commands:
  api         Generate api related files
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help      help for cztctl
  -v, --version   version for cztctl                                                                                                                                                                                                
                                                                                                                                                                                                                                    

Use "cztctl [command] --help" for more information about a command.
```


| 命令        | 简写 | 描述              |
|-----------|----|-----------------|
| --version | -v | 查看当前版本          |
| --help    | -h | 查看帮助提示信息        |
| api       | -h | 根据api文件生成代码，文档等 |

## api 命令

根据api文件生成代码文档等

```shell
$ cztctl api -h
Generate api related files

Usage:
  cztctl api [flags]
  cztctl api [command]

Available Commands:
  swagger     Generate swagger file from api

Flags:
      --branch string   The branch of the remote repo, it does work with --remote
  -h, --help            help for api                                                                                                                                                                                                
      --home string     The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority                                                                          
      --o string        Output a sample api file                                                                                                                                                                                    
      --remote string   The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority                                                                          
                        The git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure                                                                                        
                                                                                                                                                                                                                                    

Use "cztctl api [command] --help" for more information about a command.
```

| 命令      | 简写 | 描述                   |
|---------|----|----------------------|
| --help  | -h | 查看帮助提示信息             |
| swagger | -h | 根据api文件生成 swagger 文档 |

**暂时只支持swagger命令**

### api swagger

根据api文件生成swagger文档。

生成 swagger 命令为：

```
cztctl api  swagger -api ./tools/cztctl/test/test.api -dir .
```


`cztctl api swagger` 是 `goctl api swagger` 为基础扩展的。

`goctl api swagger` 功能请查看 https://go-zero.dev/docs/tutorials/cli/swagger

`cztctl api swagger` 在 `goctl api swagger` 的基础上多了如下功能：

**1、识别 validate tag,自动根据 validate tag 生成注释**

请注意，这里的 validate tag 是指 https://github.com/lerity-yao/param-validator 

```api
// GBoxCommonBoxReq 下拉框入参
GBoxCommonBoxReq {
    Ad string `form:"ad,optional" validate:"omitempty,xStr=1-10"` // ad广告
}
```

生成为文档注释为
```
ad广告
字段校验规则 omitempty，允许为空
字段校验规则，为空不校验，不为空则校验 xStr=1-10，{0}长度{1}，首尾不能有空格
```

行内注释在第一行，后面，接 validate的教研规则，一行一个，包括校验规则

**2、支持字段属性头部注释**

```api
GBoxCommonBoxItem {
    // 注意value 可能存在几种情况
    // 注意value 123
    // 注意value 456
    Value string `json:"value"` // 标签值
    // 注意 label 可能存在几种情况
    // 注意 label 123
    // 注意 label 456
    Label string `json:"label"` //标签
   // ad == 1
    Extra map[string]interface{} `json:"extra"` // 额外信息,对象类型，请把其他数据都放这里
}
```

生成为文档注释为

```
extra：
    description:	
    额外信息,对象类型，请把其他数据都放这里
    ad == 1

label：
    标签
    注意 label 可能存在几种情况
    注意 label 123
    注意 label 456

value：
    标签值
    注意value 可能存在几种情况
    注意value 123
    注意value 456
````