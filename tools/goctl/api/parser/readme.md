# Api语法描述

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
    "bar.api"
    "foo/bar.api"
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
> 被import的api必须要和main api的syntax版本一致。

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

## import语法块

随着业务规模增大，api中定义的结构体和服务越来越多，所有的语法描述均为一个api文件，这是多么糟糕的一个问题， 其会大大增加了阅读难度和维护难度，import语法块可以帮助我们解决这个问题，通过拆分api文件，
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

## info语法块

info语法块是一个包含了多个键值对的语法体，其作用相当于一个api服务的描述，解析器会将其映射到spec.Spec中， 以备用于翻译成其他语言(golang、java等)
时需要携带的meta元素。如果仅仅是对当前api的一个说明，而不考虑其翻译 时传递到其他语言，则使用简单的多行注释或者java风格的文档注释即可，关于注释说明请参考下文的 **隐藏通道**。

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

## type语法块

在api服务中，我们需要用到一个结构体(类)来作为请求体，响应体的载体，因此我们需要声明一些结构体来完成这件事情， type语法块由golang的type演变而来，当然也保留着一些golang type的特性，沿用golang特性有：

* 保留了golang内置数据类型`bool`,`int`,`int8`,`int16`,`int32`,`int64`,`uint`,`uint8`,`uint16`,`uint32`,`uint64`,`uintptr`
  ,`float32`,`float64`,`complex64`,`complex128`,`string`,`byte`,`rune`,
* 兼容golang struct风格声明
* 保留golang关键字

> ### 警告 ⚠️
> * 不支持alias
> * 不支持time.Time数据类型
> * 结构体名称、字段名称、不能为golang关键字

**语法定义**
> 由于其和golang相似，因此不做详细说明，具体语法定义请在[ApiParser.g4](g4/ApiParser.g4)中查看typeSpec定义。

**语法说明**

> 参考golang写法

**正确语法示例** ✅

eg1：不规范写法

``` api
type Foo struct{
    Id int `path:"id"` // ①
    Foo int `json:"foo"`
}

type Bar struct{
    // 非导出型字段
    bar int `form:"bar"`
}

type(
    // 非导出型结构体
    fooBar struct{
        FooBar int
    }
)
```

eg2：规范写法（推荐）

``` api
type Foo{
    Id int `path:"id"`
    Foo int `json:"foo"`
}

type Bar{
    Bar int `form:"bar"`
}

type(
    FooBar{
        FooBar int
    }
)
```

**错误语法示例** ❌

eg

``` api
type Gender int // 不支持

// 非struct token
type Foo structure{ 
  CreateTime time.Time // 不支持time.Time
}

// golang关键字 var
type var{} 

type Foo{
  // golang关键字 interface
  Foo interface 
}


type Foo{
  foo int 
  // map key必须要golang内置数据类型
  m map[Bar]string
}
```

**① tag说明**
> tag定义和golang中json tag语法一样，除了json tag外，go-zero还提供了另外一些tag来实现对字段的描述，
> 详情见下表。

* tag表

  |tag key |描述 |提供方 |有效范围 |示例 |
    |:--- |:--- |:--- |:--- |:--- |
  |json|json序列化tag|golang|request、response|`json:"fooo"`|
  |path|路由path，如`/foo/:id`|go-zero|request|`path:"id"`|
  |form|标志请求体是一个form（POST方法时）或者一个query(GET方法时`/search?name=keyword`)|go-zero|request|`form:"name"`|

* tag修饰符
  > 常见参数校验描述

  |tag key |描述 |提供方 |有效范围 |示例 |
    |:--- |:--- |:--- |:--- |:--- |
  |optional|定义当前字段为可选参数|go-zero|request|`json:"name,optional"`|
  |options|定义当前字段的枚举值,多个以竖线②隔开|go-zero|request|`json:"gender,options=male"`|
  |default|定义当前字段默认值|go-zero|request|`json:"gender,default=male"`|
  |range|定义当前字段数值范围|go-zero|request|`json:"age,range=[0:120]"`|

  ② 竖线：|
  > ### 温馨提示
  > tag修饰符需要在tag value后以引文逗号,隔开

## service语法块

service语法块用于定义api服务，包含服务名称，服务metadata，中间件声明，路由，handler等。

> ### 警告 ⚠️
> * main api和被import的api服务名称必须一致，不能出现服务名称歧义。
> * handler名称不能重复
> * 路由（请求方法+请求path）名称不能重复
> * 请求体必须声明为普通（非指针）struct，响应体做了一些向前兼容处理，详请见下文说明
>

**语法定义**

``` antlrv4
serviceSpec:    atServer? serviceApi;
atServer:       '@server' lp='(' kvLit+ rp=')';
serviceApi:     {match(p,"service")}serviceToken=ID serviceName lbrace='{' serviceRoute* rbrace='}';
serviceRoute:   atDoc? (atServer|atHandler) route;
atDoc:          '@doc' lp='('? ((kvLit+)|STRING) rp=')'?;
atHandler:      '@handler' ID;
route:          {checkHttpMethod(p)}httpMethod=ID path request=body? returnToken=ID? response=replybody?;
body:           lp='(' (ID)? rp=')';
replybody:      lp='(' dataType? rp=')';
// kv
kvLit:          key=ID {checkKeyValue(p)}value=LINE_VALUE;

serviceName:    (ID '-'?)+;
path:           (('/' (ID ('-' ID)*))|('/:' (ID ('-' ID)?)))+;
```

**语法说明**

> serviceSpec：包含了一个可选语法块`atServer`和`serviceApi`语法块，其遵循序列模式（编写service必须要按照顺序，否则会解析出错）
>
> atServer： 可选语法块，定义key-value结构的server metadata，'@server'表示这一个server语法块的开始，其可以用于描述serviceApi或者route语法块，其用于描述不同语法块时有一些特殊关键key
> 需要值得注意，见 **atServer关键key描述说明**。
>
> serviceApi：包含了1到多个`serviceRoute`语法块
>
> serviceRoute：按照序列模式包含了`atDoc`,handler和`route`
>
> atDoc：可选语法块，一个路由的key-value描述，其在解析后会传递到spec.Spec结构体，如果不关心传递到spec.Spec,
> 推荐用单行注释替代。
>
> handler：是对路由的handler层描述，可以通过atServer指定`handler` key来指定handler名称，
> 也可以直接用atHandler语法块来定义handler名称
>
> atHandler：'@handler' 固定token，后接一个遵循正则`[_a-zA-Z][a-zA-Z_-]*`)的值，用于声明一个handler名称
>
> route：路由，有`httpMethod`、`path`、可选`request`、可选`response`组成，`httpMethod`是必须是小写。
>
> body：api请求体语法定义，必须要由()包裹的可选的ID值
>
> replyBody：api响应体语法定义，必须由()包裹的struct、~~array(向前兼容处理，后续可能会废弃，强烈推荐以struct包裹，不要直接用array作为响应体)~~
>
> kvLit： 同info key-value
>
> serviceName: 可以有多个'-'join的ID值
>
> path：api请求路径，必须以'/'或者'/:'开头，切不能以'/'结尾，中间可包含ID或者多个以'-'join的ID字符串

**atServer关键key描述说明**

修饰service时

|key|描述|示例|
|:---|:---|:---|
|jwt|声明当前service下所有路由需要jwt鉴权，且会自动生成包含jwt逻辑的代码|`jwt: Auth`|
|group|声明当前service或者路由文件分组|`group: login`|
|middleware|声明当前service需要开启中间件|`middleware: AuthMiddleware`|

修饰route时

|key|描述|示例|
|:---|:---|:---|
|handler|声明一个handler|-|

**正确语法示例** ✅

eg1：不规范写法

``` api
@server(
  jwt: Auth
  group: foo
  middleware: AuthMiddleware
)
service foo-api{
  @doc(
    summary: foo
  )
  @server(
    handler: foo
  )
  // 非导出型body
  post /foo/:id (foo) returns (bar)
  
  @doc "bar"
  @handler bar
  post /bar returns ([]int)// 不推荐数组作为响应体
  
  @handler fooBar
  post /foo/bar (Foo) returns // 可以省略'returns'
}
```

eg2：规范写法（推荐）

``` api
@server(
  jwt: Auth
  group: foo
  middleware: AuthMiddleware
)
service foo-api{
  @doc "foo"
  @handler: foo
  post /foo/:id (Foo) returns (Bar)
}

service foo-api{
  @handler ping
  get /ping
  
  @doc "foo"
  @handler: bar
  post /bar/:id (Foo)
}

```

**错误语法示例** ❌

``` api
// 不支持空的server语法块
@server(
)
// 不支持空的service语法块
service foo-api{
}

service foo-api{
  @doc kkkk // 简版doc必须用英文双引号引起来
  @handler foo
  post /foo
  
  @handler foo // 重复的handler
  post /bar
  
  @handler fooBar
  post /bar // 重复的路由
  
  // @handler和@doc顺序错误
  @handler someHandler
  @doc "some doc"
  post /some/path
  
  // handler缺失
  post /some/path/:id
  
  @handler reqTest
  post /foo/req (*Foo) // 不支持除普通结构体外的其他数据类型作为请求体
  
  @handler replyTest
  post /foo/reply returns (*Foo) // 不支持除普通结构体、数组(向前兼容，后续考虑废弃)外的其他数据类型作为响应体
}
```

## 隐藏通道

隐藏通道目前主要为空百符号，换行符号以及注释，这里我们只说注释，因为空白符号和换行符号我们目前拿来也无用。

### 单行注释

**语法定义**

``` antlrv4
'//' ~[\r\n]*
```

**语法说明**
由语法定义可知道，单行注释必须要以`//`开头，内容为不能包含换行符

**正确语法示例** ✅

``` api
// doc
// comment
```

**错误语法示例** ❌

``` api
// break
line comments
```

### java风格文档注释

**语法定义**

``` antlrv4
'/*' .*? '*/'
```

**语法说明**

由语法定义可知道，单行注释必须要以`/*`开头，`*/`结尾的任意字符。

**正确语法示例** ✅

``` api
/**
 * java-style doc
 */
```

**错误语法示例** ❌

``` api
/*
 * java-style doc */
 */
```

## Doc&Comment

如果想获取某一个元素的doc或者comment开发人员需要怎么定义？

**Doc**
> 我们规定上一个语法块（非隐藏通道内容）的行数line+1到当前语法块第一个元素前的所有注释(当行，或者多行)均为doc， 且保留了`//`、`/*`、`*/`原始标记。

**Comment**
> 我们规定当前语法块最后一个元素所在行开始的一个注释块(当行，或者多行)为comment 且保留了`//`、`/*`、`*/`原始标记。

语法块Doc和Comment的支持情况

|语法块|parent语法块|Doc|Comment|
|:---|:---|:---|:---|
|syntaxLit|api|✅|✅|
|kvLit|infoSpec|✅|✅|
|importLit|importSpec|✅|✅|
|typeLit|api|✅|❌|
|typeLit|typeBlock|✅|❌|
|field|typeLit|✅|✅|
|key-value|atServer|✅|✅|
|atHandler|serviceRoute|✅|✅|
|route|serviceRoute|✅|✅|

以下为对应语法块解析后细带doc和comment的写法
``` api
// syntaxLit doc
syntax = "v1" // syntaxLit commnet

info(
  // kvLit doc
  author: songmeizi // kvLit comment
)

// typeLit doc
type Foo {}

type(
  // typeLit doc
  Bar{}
  
  FooBar{
    // filed doc
    Name int // filed comment
  }
)

@server(
  /**
   * kvLit doc
   * 开启jwt鉴权
   */
  jwt: Auth /**kvLit comment*/
)
service foo-api{
  // atHandler doc
  @handler foo //atHandler comment
  
  /*
   * route doc
   * post请求
   * path为 /foo
   * 请求体：Foo
   * 响应体：Foo
   */
  post /foo (Foo) returns (Foo) // route comment
}
```