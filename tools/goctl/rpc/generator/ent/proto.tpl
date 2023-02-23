{{.groupName}}  rpc create{{.modelName}} ({{.modelName}}Info) returns (Base{{if .useUUID}}UU{{end}}IDResp);
{{.groupName}}  rpc update{{.modelName}} ({{.modelName}}Info) returns (BaseResp);
{{.groupName}}  rpc get{{.modelName}}List ({{.modelName}}ListReq) returns ({{.modelName}}ListResp);
{{.groupName}}  rpc get{{.modelName}}ById ({{if .useUUID}}UU{{end}}IDReq) returns ({{.modelName}}Info);
{{.groupName}}  rpc delete{{.modelName}} ({{if .useUUID}}UU{{end}}IDsReq) returns (BaseResp);
