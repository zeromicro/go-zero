# goctl使用

## goctl用途

* 定义api请求
* 根据定义的api自动生成golang(后端), java(iOS & Android), typescript(web & 晓程序)，dart(flutter)
* 生成MySQL CURD+Cache
* 生成MongoDB CURD+Cache

## goctl使用说明

### 快速生成服务

* api: goctl api new xxxx
* rpc: goctl rpc new xxxx

#### goctl参数说明

  `goctl api [go/java/ts] [-api user/user.api] [-dir ./src]`

  > api 后面接生成的语言，现支持go/java/typescript
  >
  > -api 自定义api所在路径
  >
  > -dir 自定义生成目录

#### 保持goctl总是最新版

  第一次运行会在~/.goctl里增加下面两行：

  ```Plain Text
  url = http://47.97.184.41:7777/
  ```

#### API 语法说明

``` golang
info(
    title: doc title
    desc: >
    doc description first part,
    doc description second part<
    version: 1.0
)

type int userType

type user struct {
	name string `json:"user"` // 用户姓名
}

type student struct {
	name string `json:"name"` // 学生姓名
}

type teacher struct {
}

type (
	address struct {
		city string `json:"city"`
	}

	innerType struct {
		image string `json:"image"`
	}

	createRequest struct {
		innerType
		name    string    `form:"name"`
		age     int       `form:"age,optional"`
		address []address `json:"address,optional"`
	}

	getRequest struct {
		name string `path:"name"`
		age  int    `form:"age,optional"`
	}

	getResponse struct {
		code    int     `json:"code"`
		desc    string  `json:"desc,omitempty"`
		address address `json:"address"`
		service int     `json:"service"`
	}
)

service user-api {
    @doc(
        summary: user title
        desc: >
        user description first part,
        user description second part,
        user description second line
    )
    @server(
        handler: GetUserHandler
        group: user
    )
    get /api/user/:name(getRequest) returns(getResponse)

    @server(
        handler: CreateUserHandler
        group: user
    )
    post /api/users/create(createRequest)
}

@server(
    jwt: Auth
    group: profile
)
service user-api {
    @doc(summary: user title)
    @server(
        handler: GetProfileHandler
    )
    get /api/profile/:name(getRequest) returns(getResponse)

    @server(
        handler: CreateProfileHandler
    )
    post /api/profile/create(createRequest)
}

service user-api {
    @doc(summary: desc in one line)
    @server(
        handler: PingHandler
    )
    head /api/ping()
}

```

1. info部分：描述了api基本信息，比如Auth，api是哪个用途。
2. type部分：type类型声明和golang语法兼容。
3. service部分：service代表一组服务，一个服务可以由多组名称相同的service组成，可以针对每一组service配置jwt和auth认证，另外通过group属性可以指定service生成所在子目录。
   service里面包含api路由，比如上面第一组service的第一个路由，doc用来描述此路由的用途，GetProfileHandler表示处理这个路由的handler，
   `get /api/profile/:name(getRequest) returns(getResponse)` 中get代表api的请求方式（get/post/put/delete）, `/api/profile/:name` 描述了路由path，`:name`通过
   请求getRequest里面的属性赋值，getResponse为返回的结构体，这两个类型都定义在2描述的类型中。

#### api vscode插件

开发者可以在vscode中搜索goctl的api插件，它提供了api语法高亮，语法检测和格式化相关功能。

 1. 支持语法高亮和类型导航。
 2. 语法检测，格式化api会自动检测api编写错误地方，用vscode默认的格式化快捷键(option+command+F)或者自定义的也可以。
 3. 格式化(option+command+F)，类似代码格式化，统一样式支持。

#### 根据定义好的api文件生成golang代码

  命令如下：  
  `goctl api go -api user/user.api -dir user`

  ```Plain Text

	.
    ├── internal
    │   ├── config
    │   │   └── config.go
    │   ├── handler
    │   │   ├── pinghandler.go
    │   │   ├── profile
    │   │   │   ├── createprofilehandler.go
    │   │   │   └── getprofilehandler.go
    │   │   ├── routes.go
    │   │   └── user
    │   │       ├── createuserhandler.go
    │   │       └── getuserhandler.go
    │   ├── logic
    │   │   ├── pinglogic.go
    │   │   ├── profile
    │   │   │   ├── createprofilelogic.go
    │   │   │   └── getprofilelogic.go
    │   │   └── user
    │   │       ├── createuserlogic.go
    │   │       └── getuserlogic.go
    │   ├── svc
    │   │   └── servicecontext.go
    │   └── types
    │       └── types.go
    └── user.go

  ```

  生成的代码可以直接跑，有几个地方需要改：

* 在`servicecontext.go`里面增加需要传递给logic的一些资源，比如mysql, redis，rpc等
* 在定义的get/post/put/delete等请求的handler和logic里增加处理业务逻辑的代码

#### 根据定义好的api文件生成java代码

```shell
goctl api java -api user/user.api -dir ./src
```

#### 根据定义好的api文件生成typescript代码

```shell
goctl api ts -api user/user.api -dir ./src -webapi ***

ts需要指定webapi所在目录
```

#### 根据定义好的api文件生成Dart代码

```shell
goctl api dart -api user/user.api -dir ./src
```

## 根据mysql ddl或者datasource生成model文件

```shell script
goctl model mysql -src={filename} -dir={dir} -cache={true|false}
```

详情参考[model文档](https://github.com/tal-tech/go-zero/blob/master/tools/goctl/model/sql/README.MD)

## 根据定义好的简单go文件生成mongo代码文件(仅限golang使用)  

```shell
goctl model mongo -src {{yourDir}}/xiao/service/xhb/user/model/usermodel.go -cache yes
```

* src需要提供简单的usermodel.go文件，里面只需要提供一个结构体即可
* cache 控制是否需要缓存 yes=需要 no=不需要

src 示例代码如下

  ```go
    package model

    type User struct {
      Name string `o:"find,get,set" c:"姓名"`
      Age int `o:"find,get,set" c:"年纪"`
      School string `c:"学校"`
    }
  ```

结构体中不需要提供Id,CreateTime,UpdateTime三个字段，会自动生成
结构体中每个tag有两个可选标签 c 和 o
c 是该字段的注释
o 是该字段需要生产的操作函数 可以取得get,find,set 分别表示生成返回单个对象的查询方法，返回多个对象的查询方法，设置该字段方法
生成的目标文件会覆盖该简单go文件

## goctl rpc生成（业务剥离中，暂未开放）

  命令 `goctl rpc proto -proto ${proto} -service ${serviceName} -project ${projectName} -dir ${directory} -shared ${shared}`  
  如： `goctl rpc proto -proto test.proto  -service test -project xjy  -dir .`  

  参数说明：

* ${proto}: proto文件
* ${serviceName}: rpc服务名称
* ${projectName}: 所属项目，如xjy,xhb,crm,hera，具体查看help，主要为了根据不同项目服务往redis注册key，可选
* ${directory}: 输出目录
* ${shared}: shared文件生成目录，可选，默认为${pwd}/shared

生成目录结构示例：

  ```Plain Text
	.
    ├── shared [示例目录，可自己指定，强制覆盖更新]
    │   └── contentservicemodel.go
    ├── test
    │   ├── etc
    │   │   └── test.json
    │   ├── internal
    │   │   ├── config
    │   │   │   └── config.go
    │   │   ├── handler [强制覆盖更新]
    │   │   │   ├── changeavatarhandler.go
    │   │   │   ├── changebirthdayhandler.go
    │   │   │   ├── changenamehandler.go
    │   │   │   ├── changepasswordhandler.go
    │   │   │   ├── changeuserinfohandler.go
    │   │   │   ├── getuserinfohandler.go
    │   │   │   ├── loginhandler.go
    │   │   │   ├── logouthandler.go
    │   │   │   └── testhandler.go
    │   │   ├── logic
    │   │   │   ├── changeavatarlogic.go
    │   │   │   ├── changebirthdaylogic.go
    │   │   │   ├── changenamelogic.go
    │   │   │   ├── changepasswordlogic.go
    │   │   │   ├── changeuserinfologic.go
    │   │   │   ├── getuserinfologic.go
    │   │   │   ├── loginlogic.go
    │   │   │   └── logoutlogic.go
    │   │   └── svc
    │   │       └── servicecontext.go
    │   ├── pb
    │   │   └── test.pb.go
    │   └── test.go [强制覆盖更新]
    └── test.proto
  ```

注意 ：目前rpc目录生成的proto文件暂不支持import外部proto文件
