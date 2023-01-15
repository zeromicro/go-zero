
### 1. 公告列表

1. 路由定义

- Url: /signature/list
- Method: POST
- Request: `SignatureListRequest`
- Response: `SignatureListResponse`


2. 类型定义 



```golang
type SignatureListRequest struct {
	StartTime int64 `json:"startTime"`
	EndTime int64 `json:"endTime"`
	SearchData string `json:"searchData"`
	PageSize int64 `json:"pageSize"`
	NowPage int64 `json:"nowPage"`
}

type SignatureListResponse struct {
	Count int64 `json:"count"`
	NowPage int64 `json:"nowPage"`
	List []*SignatureListItem `json:"list"`
}

type SignatureListItem struct {
	SignatureId int64 `json:"signatureId"`
	Title string `json:"title"`
	PeopleNum int64 `json:"peopleNum"`
	SignNum int64 `json:"signNum"`
	TeacherName string `json:"teacherName"`
	Time int64 `json:"time"`
}
```
  


### 2. 批量下载

1. 路由定义

- Url: /signature/export
- Method: POST
- Request: `SignatureExportRequest`
- Response: `SignatureExportResponse`


2. 类型定义 



```golang
type SignatureExportRequest struct {
	SignatureId int64 `json:"signatureId"`
}

type SignatureExportResponse struct {
	List []*SignatureExportListItem `json:"list"`
}

type SignatureExportListItem struct {
	FileName string `json:"fileName"`
	FileUrl string `json:"fileUrl"`
}
```
  


### 3. 查看公告详情

1. 路由定义

- Url: /signature/info
- Method: POST
- Request: `SignatureInfoRequest`
- Response: `SignatureInfoResponse`


2. 类型定义 



```golang
type SignatureInfoRequest struct {
	SignatureId int64 `json:"signatureId"`
}

type SignatureInfoResponse struct {
	Title string `json:"title"`
	Desc string `json:"desc"`
	DocPath string `json:"docPath"`
	FileName string `json:"fileName"`
	DepartmentList []*SignatureDepartmentListItem `json:"departmentList"`
	ClassList []*SignatureClassListItem `json:"classList"`
}

type SignatureDepartmentListItem struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	PartType int64 `json:"partType"`
	PeopleNum int64 `json:"peopleNum"`
	SignNum int64 `json:"signNum"`
}

type SignatureClassListItem struct {
	Token string `json:"token"`
	Name string `json:"name"`
	PartType int64 `json:"partType"`
	PeopleNum int64 `json:"peopleNum"`
	SignNum int64 `json:"signNum"`
}
```
  


### 4. 查看签名情况

1. 路由定义

- Url: /signature/detail
- Method: POST
- Request: `SignatureDetailRequest`
- Response: `SignatureDetailResponse`


2. 类型定义 



```golang
type SignatureDetailRequest struct {
	Sid int64 `json:"sId"`
	PartId int64 `json:"partId"`
	ClassToken string `json:"classToken"`
	PartType int64 `json:"partType"`
}

type SignatureDetailResponse struct {
	SignList []*SignListItem `json:"signList"`
	NotSignList []*SignListItem `json:"notSignList"`
}

type SignListItem struct {
	Name string `json:"name"`
	Tokens []string `json:"tokens"`
}
```
  


### 5. 提醒所有人

1. 路由定义

- Url: /signature/remind-all
- Method: POST
- Request: `SignRemindAllRequest`
- Response: `SignRemindAllResponse`


2. 类型定义 



```golang
type SignRemindAllRequest struct {
	SId int64 `json:"sId"`
}

type SignRemindAllResponse struct {
}
```
  


### 6. 发短信提醒人

1. 路由定义

- Url: /signature/remind-people
- Method: POST
- Request: `SignRemindPeopleRequest`
- Response: `SignRemindPeopleResponse`


2. 类型定义 



```golang
type SignRemindPeopleRequest struct {
	SId int64 `json:"sId"`
	UserIds []string `json:"userIds"`
	PartType int64 `json:"partType"`
}

type SignRemindPeopleResponse struct {
}
```
  


### 7. 增加公告

1. 路由定义

- Url: /signature/add
- Method: POST
- Request: `SignatureAddRequest`
- Response: `SignatureAddResponse`


2. 类型定义 



```golang
type SignatureAddRequest struct {
	Title string `json:"title"`
	Desc string `json:"desc"`
	WordPath string `json:"wordPath"`
	FileName string `json:"fileName"`
	ClassTokens []string `json:"classTokens"`
	DepartmentList []*DepartmentListItem `json:"departmentList"`
}

type DepartmentListItem struct {
	DepartmentId int64 `json:"departmentId"`
	IdType int64 `json:"idType"`
	TeacherIds []int64 `json:"teacherIds"`
}

type SignatureAddResponse struct {
}
```
  


### 8. 查看公告

1. 路由定义

- Url: /signature/h5/do
- Method: POST
- Request: `SignatureDoRequest`
- Response: `SignatureDoResponse`


2. 类型定义 



```golang
type SignatureDoRequest struct {
	Sid int64 `json:"sId"`
}

type SignatureDoResponse struct {
	WordPath string `json:"wordPath"`
	IsSign int64 `json:"isSign"`
	SignPath string `json:"signPath"`
}
```
  


### 9. 公告签名

1. 路由定义

- Url: /signature/h5/sign
- Method: POST
- Request: `SignatureSignRequest`
- Response: `SignatureSignResponse`


2. 类型定义 



```golang
type SignatureSignRequest struct {
	SId int64 `json:"sId"`
	SingPath string `json:"singPath"`
}

type SignatureSignResponse struct {
	SignPath string `json:"signPath"`
}
```
  

