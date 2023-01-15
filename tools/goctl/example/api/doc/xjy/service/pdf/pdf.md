
### 1. N/A

1. 路由定义

- Url: /pdf/handle/progress/:taskId
- Method: GET
- Request: `ProgressRequest`
- Response: `-`


2. 类型定义 



```golang
type ProgressRequest struct {
	TaskId string `path:"taskId"`
}
```
  

