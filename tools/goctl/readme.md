# goctl

English | [简体中文](readme-cn.md)

## goctl introduction

* Define api requests
* Automatically generate golang (backend), java (iOS & Android), typescript (web & desktop app), dart (flutter) based on the defined api
* Generate MySQL CRUD, check [goctl model](model/sql) for details

## goctl usage instructions

### goctl parameter description

  `goctl api [go/java/ts] [-api user/user.api] [-dir ./src]`

  > api followed by the target language, now supports go/java/typescript
  >
  > -api the path to the api file
  >
  > -dir the target dir to generate in

#### API syntax description

```golang
type int userType

type user {
	name string `json:"user"` // user name
}

type student {
	name string `json:"name"` // student's name
}

type teacher {
}

type (
	address {
		city string `json:"city"` // city
	}

	innerType {
		image string `json:"image"`
	}

	createRequest {
		innerType
		name string `form:"name"`
		age int `form:"age,optional"`
		address []address `json:"address,optional"`
	}

	getRequest {
		name string `path:"name"`
		age int `form:"age,optional"`
	}

	getResponse {
		code int `json:"code"`
		desc string `json:"desc,omitempty"`
		address address `json:"address"`
		service int `json:"service"`
	}
)

service user-api {
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
    @handler GetProfileHandler
    get /api/profile/:name(getRequest) returns(getResponse)

    @handler CreateProfileHandler
    post /api/profile/create(createRequest)
}

service user-api {
    @handler PingHandler
    head /api/ping()
}
```

1. type part: type declaration.
3. service part: service represents a set of services, a service can be composed of multiple groups of service with the same name, you can configure the group attribute for each group of service to specify the subdirectory where the service is generated.
   service contains api routes, such as the first route of the first group of service above, GetProfileHandler indicates the handler that handles this route.
   `get /api/profile/:name(getRequest) returns(getResponse)` where get represents the request method of the api (get/post/put/delete), `/api/profile/:name` describes the route path, `:name` is assigned by the
   The request getRequest assigns a value to the property inside, and getResponse is the returned structure.

#### api vscode plugin

Developers can search for the api plugin for goctl in vscode and goland, which provides api syntax highlighting, syntax detection and formatting related functions.

  1. support syntax highlighting and type navigation.
  2. syntax detection, formatting api will automatically detect where the api is written wrong, using vscode default formatting shortcut (option+command+F) or custom ones can be used.
  3. formatting (option+command+F), similar to code formatting, unified style support.

#### Generate golang code based on the defined api file

  The command is as follows.  
  ```goctl api go -api user/user.api -dir user```

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


The generated code can be run directly, there are a few things that need to be changed.

* Add some resources that need to be passed to logic in `servicecontext.go`, such as mysql, redis, rpc, etc.
* Add the code to handle the business logic in the handlers and logic of the defined get/post/put/delete requests

#### Generate java code based on the defined api file

```Plain Text
goctl api java -api user/user.api -dir . /src
```

#### Generate typescript code from the defined api file

```Plain Text
goctl api ts -api user/user.api -dir . /src -webapi ***
```

ts needs to specify the directory where the webapi is located

#### Generate Dart code based on the defined api file

```Plain Text
goctl api dart -api user/user.api -dir . /src
```