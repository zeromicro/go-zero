
### 1. 解散班级生成验证码

1. 路由定义

- Url: /classroom/remove/activation/:token
- Method: POST
- Request: `RemoveClassroomActivationRequest`
- Response: `RemoveClassroomActivationResponse`

2. 请求定义


```golang
type RemoveClassroomActivationRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
}
```


3. 返回定义


```golang
type RemoveClassroomActivationResponse struct {
	Code string `json:"code"`
}
```
  


### 2. 验证解散班级生成验证码

1. 路由定义

- Url: /classroom/remove/:token/:classroomId
- Method: DELETE
- Request: `RemoveClassroomRequest`
- Response: `RemoveClassroomResponse`

2. 请求定义


```golang
type RemoveClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
	Code string `json:"code,optional"`
}
```


3. 返回定义


```golang
type RemoveClassroomResponse struct {
}
```
  

