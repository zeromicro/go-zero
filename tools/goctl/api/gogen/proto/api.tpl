import "base.api"

type (
    // The response data of {{.modelNameSpace}} information | {{.modelName}}信息
    {{.modelName}}Info {
        Base{{if .useUUID}}UU{{end}}IDInfo{{.infoData}}
    }

    // The response data of {{.modelNameSpace}} list | {{.modelName}}列表数据
    {{.modelName}}ListResp {
        BaseDataInfo

        // {{.modelName}} list data | {{.modelName}}列表数据
        Data {{.modelName}}ListInfo `json:"data"`
    }

    // {{.modelName}} list data | {{.modelName}}列表数据
    {{.modelName}}ListInfo {
        BaseListInfo

        // The API list data | {{.modelName}}列表数据
        Data  []{{.modelName}}Info  `json:"data"`
    }

    // Get {{.modelNameSpace}} list request params | {{.modelName}}列表请求参数
    {{.modelName}}ListReq {
        PageInfo{{.listData}}
    }

    // {{.modelName}} information response | {{.modelName}}信息返回体
    {{.modelName}}InfoResp {
        BaseDataInfo

        // {{.modelName}} information | {{.modelName}}数据
        Data {{.modelName}}Info `json:"data"`
    }
)

@server(
    jwt: Auth
    group: {{.modelNameLowerCase}}
    middleware: Authority
)

service {{.apiServiceName}} {
    // Create {{.modelNameSpace}} information | 创建{{.modelName}}
    @handler create{{.modelName}}
    post /{{.modelNameSnake}}/create ({{.modelName}}Info) returns (BaseMsgResp)

    // Update {{.modelNameSpace}} information | 更新{{.modelName}}
    @handler update{{.modelName}}
    post /{{.modelNameSnake}}/update ({{.modelName}}Info) returns (BaseMsgResp)

    // Delete {{.modelNameSpace}} information | 删除{{.modelName}}信息
    @handler delete{{.modelName}}
    post /{{.modelNameSnake}}/delete ({{if .useUUID}}UU{{end}}IDsReq) returns (BaseMsgResp)

    // Get {{.modelNameSpace}} list | 获取{{.modelName}}列表
    @handler get{{.modelName}}List
    post /{{.modelNameSnake}}/list ({{.modelName}}ListReq) returns ({{.modelName}}ListResp)

    // Get {{.modelNameSpace}} by ID | 通过ID获取{{.modelName}}
    @handler get{{.modelName}}ById
    post /{{.modelNameSnake}} ({{if .useUUID}}UU{{end}}IDReq) returns ({{.modelName}}InfoResp)
}
