
### 1. 新建模板

1. 路由定义

- Url: /studenttask/template/new
- Method: POST
- Request: `TemplateRequest`
- Response: `-`


2. 类型定义 



```golang
type TemplateRequest struct {
	Subject string `json:"subject"` // 适用学科
	ApplicableGrade string `json:"applicableGrade"` // 适用年级
	Title string `json:"title"` // 模板名称
	Description string `json:"description"` // 模板描述
	Audio *Audio `json:"audio,optional"` // 语音
	Link string `json:"link,optional"` // 链接
	Images []*Image `json:"images,optional"` // 图片列表
	Recommend bool `json:"recommend"` // 是否推荐给其他人
}
```
  


### 2. 模板列表

1. 路由定义

- Url: /studenttask/template/list
- Method: POST
- Request: `TemplateListRequest`
- Response: `TemplateListResponse`


2. 类型定义 



```golang
type TemplateListRequest struct {
	Type int `json:"type,optional"` // 模板类型，0-全部，1-推荐模板，2-我的模板
	Subject string `json:"subject,optional"` // 适用学科
	ApplicableGrade int64 `json:"applicableGrade,optional"` // 适用年级
}

type TemplateListResponse struct {
	Total int64 `json:"total"` // 模板总数
	List []*Template `json:"list"` // 模板列表
}
```
  


### 3. 模板详情

1. 路由定义

- Url: /studenttask/template/:id
- Method: GET
- Request: `TemplateDetailRequest`
- Response: `Template`


2. 类型定义 



```golang
type TemplateDetailRequest struct {
	Id int64 `path:"id"` // 模板id
}
```
  


### 4. 学生接收人员列表

1. 路由定义

- Url: /studenttask/receiver/list
- Method: GET
- Request: `-`
- Response: `ReceiverResponse`


2. 类型定义 



```golang
type ReceiverResponse struct {
	List []*ReceiverResponseItem `json:"list"` // 接收人员列表
}

type ReceiverResponseItem struct {
	ClassId int64 `json:"classId"`
	ClassName string `json:"className"`
	ClassBadge string `json:"classBadge"`
	Students []*Student `json:"students"`
}

type Student struct {
	Id int64 `json:"id"` // 学生id
	Name string `json:"name"` // 学生姓名
	Avatar string `json:"avatar"` // 学生头像
}
```
  


### 5. 使用学科列表

1. 路由定义

- Url: /studenttask/subject/list
- Method: GET
- Request: `-`
- Response: `SubjectResponse`


2. 类型定义 



```golang
type SubjectResponse struct {
	List []string `json:"list"` // 适用学科列表
}
```
  


### 6. 使用年级列表

1. 路由定义

- Url: /studenttask/applicable/grade/list
- Method: GET
- Request: `-`
- Response: `ApplicableGradeResponse`


2. 类型定义 



```golang
type ApplicableGradeResponse struct {
	List []string `json:"list"` // 适用年级列表
}
```
  


### 7. 布置作业

1. 路由定义

- Url: /studenttask/task/publish
- Method: POST
- Request: `TaskPublishRequest`
- Response: `-`


2. 类型定义 



```golang
type TaskPublishRequest struct {
	Id int64 `json:"id,optional"` // 作业ID，当id不为空时为编辑作业，否则为新布置作业
	Receivers []*ReceiverItem `json:"receivers"` // 接收人列表
	Subject string `json:"subject"` // 适用学科
	ApplicableGrade string `json:"applicableGrade"` // 适用年级
	Title string `json:"title"` // 作业标题
	Description string `json:"description"` // 作业描述
	Audio *Audio `json:"audio,optional"` // 语音
	Link string `json:"link,optional"` // 链接
	Images []*Image `json:"images,optional"` // 图片列表
	Recommend bool `json:"recommend"` // 是否推荐给其他老师
	AsTemplate bool `json:"asTemplate"` // 是否保存为我的模板
	EndTime int64 `json:"endTime"` // 结束时间
	TemplateId int64 `json:"templateId,optional"` // 模板id，非必填
}

type ReceiverItem struct {
	ClassId int64 `json:"classId"` // 班级id
	StudentIds []int64 `json:"studentIds"` // 学生id数组
}
```
  


### 8. 已布置的作业列表

1. 路由定义

- Url: /studenttask/task/list
- Method: POST
- Request: `TaskListRequest`
- Response: `TaskListResponse`


2. 类型定义 



```golang
type TaskListRequest struct {
	ClassId int64 `json:"classId"` // 班级id
	Subject string `json:"subject,optional"` // 学科
}

type TaskListResponse struct {
	Total int64 `json:"total"` // 作业总数
	List []*Task `json:"list"` // 作业列表
}

type Task struct {
	Id int64 `json:"id"` // 作业id
	Subject string `json:"subject"` // 学科
	Title string `json:"title"` // 作业标题
	Description string `json:"description"` // 作业描述
	CreateTime int64 `json:"createTime"` // 创建时间时间戳，单位：秒
	Creator string `json:"creator"` // 创建人
	CompleteCount int64 `json:"completeCount"` // 提交人数
	Total int64 `json:"total"` // 学生接收总人数
}
```
  


### 9. 撤回已布置的作业

1. 路由定义

- Url: /studenttask/task/recall/:id
- Method: POST
- Request: `TaskRecallRequest`
- Response: `-`


2. 类型定义 



```golang
type TaskRecallRequest struct {
	Id int64 `path:"id"` // 作业id
}
```
  


### 10. 教师布置作业详情

1. 路由定义

- Url: /studenttask/task/detail/:id
- Method: GET
- Request: `TaskDetailRequest`
- Response: `TaskDetailResponse`


2. 类型定义 



```golang
type TaskDetailRequest struct {
	Id int64 `path:"id"` // 作业id
}

type TaskDetailResponse struct {
	Subject string `json:"subject"` // 适用学科
	ApplicableGrade string `json:"applicableGrade"` // 适用年级
	Title string `json:"title"` // 标题
	Description string `json:"description"` // 描述
	Audio *Audio `json:"audio,optional"` // 语音文件
	Link string `json:"link,optional"` // 链接
	Images []*Image `json:"images,optional"` // 图片列表
	Creator *User `json:"creator"` // 创建人
	EndTime int64 `json:"endTime"` // 结束时间
	CompleteList []*CompleteUser `json:"completeList"` // 已提交人员列表
	UnCompleteList []*User `json:"unCompleteList"` // 未提交人员列表
	CreateTime int64 `json:"createTime"` // 学生上传作业内容
}

type CompleteUser struct {
	Name string `json:"name"` // 姓名
	Avatar string `json:"avatar"` // 头像
	ClassName string `json:"className"` // 年级班级名称
	StudentTaskId int `json:"studentTaskId"` // 作业学生绑定信息ID
}
```
  


### 11. 教师对学生作业进行备注

1. 路由定义

- Url: /studenttask/task/remark
- Method: POST
- Request: `TaskRemarkRequest`
- Response: `-`


2. 类型定义 



```golang
type TaskRemarkRequest struct {
	Id int64 `json:"id"` // 作业id
	StudentId int64 `json:"studentId"` // 学生ID
	Remark string `json:"remark"` // 备注
}
```
  


### 12. 对未提交人员进行晓叮当提醒

1. 路由定义

- Url: /studenttask/task/ding
- Method: POST
- Request: `TaskDingRequest`
- Response: `-`


2. 类型定义 



```golang
type TaskDingRequest struct {
	Id int64 `json:"id"` // 作业id
	Sms bool `json:"sms"` // 是否短信提醒
	StudentIds []int64 `json:"studentIds"`
}
```
  


### 13. 学生作业列表

1. 路由定义

- Url: /studenttask/student/task/list/:studentId
- Method: POST
- Request: `StudentTaskListRequest`
- Response: `StudentTaskListResponse`


2. 类型定义 



```golang
type StudentTaskListRequest struct {
	StudentId int64 `path:"studentId"` // 学生id
	Page int `json:"page,optional"` // 页码，非必填，默认1
}

type StudentTaskListResponse struct {
	Total int64 `json:"total"` // 总数
	List []*StudentTaskItem `json:"list"` // 学生作业列表
}

type StudentTaskItem struct {
	Id int64 `json:"id"` // 作业id
	Subject string `json:"subject"` // 适用学科
	Title string `json:"title"` // 标题
	Description string `json:"description"` // 描述
	CreateTime int64 `json:"createTime"` // 创建时间时间戳，单位：秒
	Creator string `json:"creator"` // 创建人
}
```
  


### 14. 学生作业详情

1. 路由定义

- Url: /studenttask/student/task/detail
- Method: POST
- Request: `StudentTaskDetailRequest`
- Response: `StudentTaskDetailResponse`


2. 类型定义 



```golang
type StudentTaskDetailRequest struct {
	Id int64 `json:"id"` // 作业id
	StudentId int64 `json:"studentId"` // 学生id
}

type StudentTaskDetailResponse struct {
	Subject string `json:"subject"` // 适用学科
	ApplicableGrade string `json:"applicableGrade"` // 适用年级
	Title string `json:"title"` // 作业标题
	Description string `json:"description"` // 作业描述
	Audio *Audio `json:"audio,optional"` // 语音文件
	Link string `json:"link,optional"` // 链接
	Images []*Image `json:"images,optional"` // 图片列表
	Creator *User `json:"creator"` // 创建人
	EndTime int64 `json:"endTime"` // 结束时间时间戳，单位：秒
	Upload *StudentTaskUpload `json:"studentTaskUpload,omitempty"` // 上传的作业
	CreateTime int64 `json:"createTime"` // 创建时间
}
```
  


### 15. 上传作业

1. 路由定义

- Url: /studenttask/student/task/upload
- Method: POST
- Request: `StudentTaskUploadRequest`
- Response: `-`


2. 类型定义 



```golang
type StudentTaskUploadRequest struct {
	Id int64 `json:"id"` // 作业id
	Time int64 `json:"time"` // 完成时间耗时，单位：分钟
	Description string `json:"description"` // 描述
	Audio *Audio `json:"audio,optional"` // 语音文件
	Link string `json:"link,optional"` // 链接
	Images []*Image `json:"images,optional"` // 图片列表
}
```
  


### 16. 教师端学生作业详情

1. 路由定义

- Url: /studenttask/task/student-task-complete-detail/:studentTaskId
- Method: GET
- Request: `StudentTaskCompleteDetailRequest`
- Response: `StudentTaskCompleteDetailResponse`


2. 类型定义 



```golang
type StudentTaskCompleteDetailRequest struct {
	StudentTaskId int `path:"studentTaskId"` // 学生作业任务 ID
}

type StudentTaskCompleteDetailResponse struct {
	StudentTaskId int `json:"studentTaskId"` // 学生作业任务 ID
	StudentTaskUpload *StudentTaskUpload `json:"studentTaskUpload"` // 学生作业上传的内容
}
```
  

