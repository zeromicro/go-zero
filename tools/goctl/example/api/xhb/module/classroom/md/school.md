
### 1. 查询1公里内的学校

1. 路由定义

- Url: /classroom/search/school/:long/:lat
- Method: GET
- Request: `-`
- Response: `SearchSchoolResponse`

2. 请求定义


3. 返回定义


```golang
type SearchSchoolResponse struct {
	SchoolInfos []SchoolInfo `json:"schoolInfos"`
}
```
  


### 2. 根据学校名称查询学校

1. 路由定义

- Url: /classroom/search/school/by/name
- Method: POST
- Request: `SearchSchoolByNameRequest`
- Response: `SearchSchoolByNameResponse`

2. 请求定义


```golang
type SearchSchoolByNameRequest struct {
	Name string `json:"name"`
	Page int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}
```


3. 返回定义


```golang
type SearchSchoolByNameResponse struct {
	SchoolInfos []SchoolInfo `json:"schoolInfos"`
}
```
  


### 3. 用户主动创建学校

1. 路由定义

- Url: /classroom/create/school/:token
- Method: POST
- Request: `CreateSchoolRequest`
- Response: `CreateSchoolResponse`

2. 请求定义


```golang
type CreateSchoolRequest struct {
	Token string `path:"token"`
	Name string `json:"name"`
	Long float32 `json:"long"`
	Lat float32 `json:"lat"`
	Province string `json:"province"`
	City string `json:"city"`
	District string `json:"district"`
	Address string `json:"address"`
}
```


3. 返回定义


```golang
type CreateSchoolResponse struct {
}
```
  

