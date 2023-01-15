
### 1. N/A

1. 路由定义

- Url: /api/goctl/gen
- Method: POST
- Request: `GenRequest`
- Response: `GenResponse`


2. 类型定义 



```golang
type GenRequest struct {
	Table *Table `json:"table"`
}

type GenResponse struct {
	Src string `json:"src"`
}
```
  


### 2. N/A

1. 路由定义

- Url: /api/goctl/schema/list
- Method: GET
- Request: `-`
- Response: `-`


2. 类型定义 



```golang

```
  


### 3. N/A

1. 路由定义

- Url: /api/goctl/table/search
- Method: POST
- Request: `TableSearchRequest`
- Response: `TableSearchResponse`


2. 类型定义 



```golang
type TableSearchRequest struct {
	Schema string `json:"schema,default=campus_test"`
	Keyword string `json:"keyword,optional"`
}

type TableSearchResponse struct {
	Tables []*Table `json:"tables"`
}
```
  

