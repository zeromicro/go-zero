# 1.8.4-beta

## swagger
  - [features] Supported operation id for swagger
## Other
  - Updated version to 1.8.4-beta


# 1.8.4-alpha

## swagger
1. [bug fix] remove example generation when request body are `query`, `path` and `header`
- it not supported in api spec 2.0
- it's will generate example when request body is json format.
2. [features] swagger generation supported definitions
- supported response definitions
- supported json request body definitions
- do not support query and form definitions, use parameters instead.

**How to use?**
Use the `useDefinitions` keyword in the info code block of the API file to declare the enable. This value is a boolean value. When set to `true`, it will enable the generation of definitions. Otherwise, it will be generated according to properties, and the default is `false`, for example:

```go
syntax = "v1"

info (
  ...
  wrapCodeMsg: true
  useDefinitions: true
)
...
```

the demo result of swagger.json

```json
{
  ...
  "responses": {
          "200": {
            "description": "",
            "schema": {
              "type": "object",
              "properties": {
                "code": {
                  "description": "1001-User not login\u003cbr\u003e1002-User permission denied",
                  "type": "integer",
                  "example": 0
                },
                "data": {
                  "$ref": "#/definitions/FormResp"
                },
                "msg": {
                  "description": "business message",
                  "type": "string",
                  "example": "ok"
                }
              }
            }
          }
        }
  ...
}
```

For a complete API example, please refer to the `api/swagger/example/example.api` file in pr. For a complete swagger result example, please refer to the `api/swagger/example/example.swagger.json` file in pr.

## 2. `goctl api go` code generation
- [bug-fix] Add flag `--type-group` to control the output of types(deprecated: experimental switch control type grouping is no longer used), if true, the types in only one group will separate by file.
- example `goctl api go --api demo.api --dir demo --type-group`
-  use `group` keyword in @server block to define  the group name in api file, for example
```go
@server(
  group: user
)
service demo{
  ...
}
```
the example of separated types by file
```
.
└── types
    ├── common.go
    ├── gotoolexport.go
    ├── importfile.go
    ├── process.go
    └── types.go
```

## 3 API Parser
- supported identifier value for info key-value in api parser
  for example

```
syntax = "v1"

info(
  enable: true
  disable: false
)
...
```