{{.groupName}}  rpc createOrUpdate{{.modelName}} ({{.modelName}}Info) returns (BaseResp);
{{.groupName}}  rpc get{{.modelName}}List ({{.modelName}}PageReq) returns ({{.modelName}}ListResp);
{{.groupName}}  rpc delete{{.modelName}} (IDReq) returns (BaseResp);
{{.groupName}}  rpc batchDelete{{.modelName}} (IDsReq) returns (BaseResp);