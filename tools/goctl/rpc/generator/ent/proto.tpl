  rpc createOrUpdate{{.modelName}} ({{.modelName}}Info) returns (BaseResp);
  rpc get{{.modelName}}List ({{.modelName}}ListReq) returns ({{.modelName}}ListResp);
  rpc delete{{.modelName}} (IDReq) returns (BaseResp);
  rpc batchDelete{{.modelName}} (IDsReq) returns (BaseResp);