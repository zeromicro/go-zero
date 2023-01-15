
### 1. user title

1. 路由定义

- Url: /api/users/user/:name
- Method: GET
- Request: `getRequest`
- Response: `getResponse`


2. 类型定义 



```golang
type GetRequest struct {
	Name string `path:"name"`
	Age int `form:"age,optional"`
}

type GetResponse struct {
	Code int `json:"code"`
	Desc string `json:"desc,omitempty"`
	Address address `json:"address"`
}
```
  


### 2. N/A

1. 路由定义

- Url: /api/users/create
- Method: POST
- Request: `createRequest`
- Response: `-`


2. 类型定义 



```golang
type CreateRequest struct {
	Name string `form:"name"`
	Age int `form:"age,optional"`
	Address []address `json:"address,optional"`
}
```
  

