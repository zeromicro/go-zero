import "base.api"

type (
{{.typeData}}
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
