# Api语法描述

简体中文 | [English](readme_en.md)

## api示例

``` golang
/**
 * api语法示例及语法说明
 */

// api语法版本
syntax = "v1"

// import literal
import "foo.api"

// import group
import (
    
)
info(
    author: "songmeizi"
    date:   "2020-01-08"
    desc:   "api语法示例及语法说明"
)

// type literal

type Foo{
    Foo int `json:"foo"`
}

// type group

type(
    Bar{
        Bar int `json:"bar"`
    }
)

// service block
@server(
    jwt:   Auth
    group: foo
)
service foo-api{
    @doc "foo"
    @handler foo
    post /foo (Foo) returns (Bar)
}
```

## api语法结构
* syntax语法声明
* import语法块
* info语法块
* type语法块
* service语法块
* 隐藏通道

> ### 温馨提示️
> 在以上语法结构中，各个语法块从语法上来说，按照语法块为单位，可以在.api文件中任意位置声明，
> 但是为了提高阅读效率，我们建议按照以上顺序进行声明，因为在将来可能会通过严格模式来控制语法块的顺序。

### syntax语法声明
syntax是新加入的语法结构，该语法的引入可以解决：
* 快速针对api版本定位存在问题的语法结构
* 针对版本做语法解析
* 防止api语法大版本升级导致前后不能向前兼容

> ### 警告 ⚠️
> 在被import的api必须要和main api的syntax版本一致。

**语法定义**
``` antlrv4
'syntax'={checkVersion(p)}STRING
```

**语法说明**
> syntax：固定token，标志一个syntax语法结构的开始
> 
> checkVersion：自定义go方法，检测`STRING`是否为一个合法的版本号，目前检测逻辑为，STRING必须是满足`(?m)"v[1-9][0-9]*"`正则。
> 
> STRING：一串英文双引号包裹的字符串，如"v1"
> 
> 一个api语法文件只能有0或者1个syntax语法声明，如果没有syntax，则默认为v1版本
> 


**正确语法示例** ✅

eg1：不规范写法
``` api
syntax="v1"
```

eg2：规范写法(推荐)
``` api
syntax = "v2"
```

**错误语法示例** ❌

eg1：
``` api
syntax = "v0"
```

eg2：
``` api
syntax = v1
```

eg3：
``` api
syntax = "V1"
```

### import语法块
随着业务规模增大，api中定义的结构体和服务越来越多，所有的语法描述均为一个api文件，这是多么糟糕的一个问题，
其会大大增加了阅读难度和维护难度，import语法块可以帮助我们解决这个问题，通过拆分api文件，
不同的api文件按照一定规则声明，可以降低阅读难度和维护难度。

> ### 警告 ⚠️
> 这里import不像golang那样包含package声明，仅仅是一个文件路径的引入，最终解析后会把所有的声明都汇聚到一个spec.Spec中。
> 不能import多个相同路径，否则会解析错误。

**语法定义**
``` antlrv4
'import' {checkImportValue(p)}STRING  
|'import' '(' ({checkImportValue(p)}STRING)+ ')'
```

**语法说明**
> import：固定token，标志一个import语法的开始
>
> checkImportValue：自定义go方法，检测`STRING`是否为一个合法的文件路径，目前检测逻辑为，STRING必须是满足`(?m)"(/?[a-zA-Z0-9_#-])+\.api"`正则。
>
> STRING：一串英文双引号包裹的字符串，如"foo.api"
> 


**正确语法示例** ✅

eg：
``` api
import "foo.api"
import "foo/bar.api"

import(
    "bar.api"
    "foo/bar/foo.api"
)
```


**错误语法示例** ❌

eg：
``` api
import foo.api
import "foo.txt"
import (
    bar.api
    bar.api
)
```

### info语法块
info语法块是一个包含了多个键值对的语法体，其作用相当于一个api服务的描述，解析器会将其映射到spec.Spec中，
以备用于翻译成其他语言(golang、java等)时需要携带的meta元素。如果仅仅是对当前api的一个说明，而不考虑其翻译
时传递到其他语言，则使用简单的多行注释或者java风格的文档注释即可，关于注释说明请参考下文的 **隐藏通道**。

> ### 警告 ⚠️
> 不能使用重复的key，每个api文件只能有0或者1个info语法块

**语法定义**
``` antlrv4
'info' '(' (ID {checkKeyValue(p)}VALUE)+ ')'
```

**语法说明**
> info：固定token，标志一个info语法块的开始
>
> checkKeyValue：自定义go方法，检测`VALUE`是否为一个合法值。
>
> VALUE：key对应的值，可以为单行的除'\r','\n','/'后的任意字符，多行请以""包裹，不过强烈建议所有都以""包裹
>

**正确语法示例** ✅

eg1：不规范写法
``` api
info(
foo: foo value
bar:"bar value"
    desc:"long long long long
long long text"
)
```

eg2：规范写法(推荐)
``` api
info(
    foo: "foo value"
    bar: "bar value"
    desc: "long long long long long long text"
)
```

**错误语法示例** ❌

eg1：没有key-value内容
``` api
info()
```

eg2：不包含冒号
``` api
info(
    foo value
)
```

eg3：key-value没有换行
``` api
info(foo:"value")
```

eg4：没有key
``` api
info(
    : "value"
)
```

eg5：非法的key
``` api
info(
    12: "value"
)
```

eg6：移除旧版本多行语法
``` api
info(
    foo: >
    some text
    <
)
```

### type语法块
在api服务中，我们需要用到一个结构体(类)来作为请求体，响应体的载体，因此我们需要声明一些结构体来完成这件事情，
type语法块由golang的type演变而来，当然也保留着一些golang type的特性，沿用golang特性有：
* 保留了golang内置数据类型`bool`,`int`,`int8`,`int16`,`int32`,`int64`,`uint`,`uint8`,`uint16`,`uint32`,`uint64`,`uintptr`,`float32`,`float64`,`complex64`,`complex128`,`string`,`byte`,`rune`,
* 兼容golang struct风格声明
* 保留golang关键字

> ### 警告 ⚠️
> * 不支持alias
> * 不支持time.Time数据类型
> * 结构体名称、字段名称、不能为golang关键字

**语法定义**
> 限于篇幅，请查看tools/goctl/api/parser/g4/ApiParser.g4中查看typeSpec定义。
### service语法块
### 隐藏通道