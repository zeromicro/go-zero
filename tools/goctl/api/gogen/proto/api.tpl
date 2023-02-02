import "base.api"

type (
    // The response data of {{.modelName}} information | {{.modelName}}信息
    {{.modelName}}Info {
        Base{{if .useUUID}}UUID{{end}}Info{{.infoData}}
    }

    // Create or update {{.modelName}} information request | 创建或更新{{.modelName}}信息
    CreateOrUpdate{{.modelName}}Req {
        // ID
        // Required: true
        Id    {{if .useUUID}}string{{else}}uint64{{end}}    `json:"id"`{{.infoData}}
    }

    // The response data of {{.modelName}} list | {{.modelName}}列表数据
    {{.modelName}}ListResp {
        BaseDataInfo

        // {{.modelName}} list data | {{.modelName}} 列表数据
        Data {{.modelName}}ListInfo `json:"data"`
    }

    // {{.modelName}} list data | {{.modelName}} 列表数据
    {{.modelName}}ListInfo {
        BaseListInfo

        // The API list data | {{.modelName}} 列表数据
        Data  []{{.modelName}}Info  `json:"data"`
    }

    // Get {{.modelNameLowerCase}} list request params | {{.modelName}}列表请求参数
    {{.modelName}}ListReq {
        PageInfo{{.listData}}
    }
)

@server(
    jwt: Auth
    group: {{.modelNameLowerCase}}
    middleware: Authority
)

service {{.apiServiceName}} {
    // Create or update {{.modelName}} information | 创建或更新{{.modelName}}
    @handler createOrUpdate{{.modelName}}
    post /{{.modelNameLowerCase}}/create_or_update (CreateOrUpdate{{.modelName}}Req) returns (BaseMsgResp)

    // Delete {{.modelName}} information | 删除{{.modelName}}信息
    @handler delete{{.modelName}}
    post /{{.modelNameLowerCase}}/delete ({{if .useUUID}}UU{{end}}IDReq) returns (BaseMsgResp)

    // Get {{.modelName}} list | 获取{{.modelName}}列表
    @handler get{{.modelName}}List
    post /{{.modelNameLowerCase}}/list ({{.modelName}}ListReq) returns ({{.modelName}}ListResp)

    // Delete {{.modelName}} information | 删除{{.modelName}}信息
    @handler batchDelete{{.modelName}}
    post /{{.modelNameLowerCase}}/batch_delete ({{if .useUUID}}UU{{end}}IDsReq) returns (BaseMsgResp)
{{if .hasStatus}}
    // Set {{.modelNameLowerCase}}'s status | 更新{{.modelName}}状态
    @handler update{{.modelName}}Status
    post /{{.modelNameLowerCase}}/status (StatusCode{{if .useUUID}}UUID{{end}}Req) returns (BaseMsgResp)
{{end}}
}
