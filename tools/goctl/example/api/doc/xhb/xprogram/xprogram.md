
### 1. 获取晓程序信息

1. 路由定义

- Url: /xprogram/info/:token
- Method: POST
- Request: `GetProgramInfoRequest`
- Response: `GetXProgramResponse`


2. 类型定义 



```golang
type GetProgramInfoRequest struct {
	AppId string `json:"id"`
	Env string `json:"env,optional"`
}

type GetXProgramResponse struct {
	Id string `json:"id"`
	AppId string `json:"appId"`
	Name string `json:"name"`
	Logo string `json:"logo"`
	Version string `json:"version"`
	Url string `json:"url"`
	AppMiniVersion string `json:"appMiniVersion"`
	VersionUpdateTime int64 `json:"versionUpdateTime"`
}
```
  

