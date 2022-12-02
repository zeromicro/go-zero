import "base.api"

type (
    // The response data of {{.modelName}} information | {{.modelName}}信息
    {{.modelName}}Info {
        BaseInfo{{.infoData}}
    }

    // Create or update {{.modelName}} information request | 创建或更新{{.modelName}}信息
    CreateOrUpdate{{.modelName}}Req {
        // ID
        // Required: true
        Id            uint64 `json:"id"`{{.infoData}}
    }

    // The response data of {{.modelName}} list | {{.modelName}}列表数据
    {{.modelName}}ListResp {
        BaseDataInfo

        // {{.modelName}} list data | API 列表数据
        Data {{.modelName}}ListInfo `json:"data"`
    }

    // {{.modelName}} list data | {{.modelName}} 列表数据
    {{.modelName}}ListInfo {
        BaseListInfo

        // The API list data | API列表数据
        Data  []{{.modelName}}Info  `json:"data"`
    }

    // Get {{.modelName}} list request params | {{.modelName}}列表请求参数
    {{.modelName}}ListReq {
        PageInfo{{.listData}}
    }
)

@server(
    jwt: Auth
    group: {{.modelNameLowerCase}}
    middleware: Authority
)

service {{.serviceName}} {
    // Create or update {{.modelName}} information | 创建或更新{{.modelName}}
    @handler createOrUpdate{{.modelName}}
    post /{{.modelNameLowerCase}} (CreateOrUpdate{{.modelName}}Req) returns (BaseMsgResp)

    // Delete {{.modelName}} information | 删除{{.modelName}}信息
    @handler delete{{.modelName}}
    delete /{{.modelNameLowerCase}} (IDReq) returns (BaseMsgResp)

    // Get {{.modelName}} list | 获取{{.modelName}}列表
    @handler get{{.modelName}}List
    post /{{.modelNameLowerCase}}/list ({{.modelName}}ListReq) returns ({{.modelName}}ListResp)
}
