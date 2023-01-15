
### 1. 课程表是否有新提醒

1. 路由定义

- Url: /classroom/have/remind/table/:token/:classroomId
- Method: GET
- Request: `QueryHaveRemindTableRequest`
- Response: `QueryHaveRemindTableResponse`

2. 请求定义


```golang
type QueryHaveRemindTableRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
}
```


3. 返回定义


```golang
type QueryHaveRemindTableResponse struct {
	HaveRemind bool `json:"haveRemind"`
}
```
  


### 2. 课程表点击小红点消失行为

1. 路由定义

- Url: /classroom/remind/table/:token
- Method: POST
- Request: `SetClassTableRemindRequest`
- Response: `SetClassTableRemindResponse`

2. 请求定义


```golang
type SetClassTableRemindRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
}
```


3. 返回定义


```golang
type SetClassTableRemindResponse struct {
}
```
  


### 3. 返回班级的课程表的state

1. 路由定义

- Url: /classroom/table/state/:token/:classroomId
- Method: GET
- Request: `QueryClassTableStateRequest`
- Response: `QueryClassTableStateResponse`

2. 请求定义


```golang
type QueryClassTableStateRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
}
```


3. 返回定义


```golang
type QueryClassTableStateResponse struct {
	Id int64 `json:"id"`
	ClassroomId string `json:"classroomId"`
	State string `json:"state"` // state有三种状态，NONE为课程表不存在|EXIST为课程表存在|WAITING为等待审核
	TableImages []*ImageRequest `json:"tableImages"`
	Name string `json:"name"`
	CreateTime int64 `json:"createTime"`
}
```
  


### 4. 管理员审核课程表图片

1. 路由定义

- Url: /classroom/review/table/:token
- Method: POST
- Request: `ReviewClassTableRequest`
- Response: `ReviewClassTableRequest`

2. 请求定义


```golang
type ReviewClassTableRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	ReviewState bool `json:"reviewState"`
}
```


3. 返回定义


```golang
type ReviewClassTableRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	ReviewState bool `json:"reviewState"`
}
```
  


### 5. 上传课程表图片

1. 路由定义

- Url: /classroom/upload/table/:token
- Method: POST
- Request: `UploadClassTableRequest`
- Response: `UploadClassTableResponse`

2. 请求定义


```golang
type UploadClassTableRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	TableImages []*ImageRequest `json:"tableImages"`
}
```


3. 返回定义


```golang
type UploadClassTableResponse struct {
	Id int64 `json:"id,omitempty"`
	ClassroomId string `json:"classroomId,omitempty"`
	State string `json:"state,omitempty"`
	TableImages []*ImageRequest `json:"tableImages,omitempty"`
	UserId string `json:"userId,omitempty"`
	CreateTime int64 `json:"createTime,omitempty"`
}
```
  

