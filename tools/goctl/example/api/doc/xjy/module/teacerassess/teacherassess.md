
### 1. 新增评价组

1. 路由定义

- Url: /teacherassess/object/add
- Method: POST
- Request: `AddObjectRequest`
- Response: `-`


2. 类型定义 



```golang
type AddObjectRequest struct {
	ObjectName string `json:"objectName"` // 评价组名称
	MondayTeachers []int64 `json:"mondayTeachers,optional"` // 周一被评价教师id
	TuesdayTeachers []int64 `json:"tuesdayTeachers,optional"`
	WednesdayTeachers []int64 `json:"wednesdayTeachers,optional"`
	ThursdayTeachers []int64 `json:"thursdayTeachers,optional"`
	FridayTeachers []int64 `json:"fridayTeachers,optional"`
	SaturdayTeachers []int64 `json:"saturdayTeachers,optional"`
	SundayTeachers []int64 `json:"sundayTeachers,optional"`
	WeekTeachers []int64 `json:"weekTeachers,optional"` // 固定评价组被评价教师id
	Type int64 `json:"type"` // 评价组类型。0-自定义评价组，1-固定评价组
}
```
  


### 2. 编辑评价组

1. 路由定义

- Url: /teacherassess/object/update
- Method: POST
- Request: `UpdateObjectRequest`
- Response: `-`


2. 类型定义 



```golang
type UpdateObjectRequest struct {
	ObjectName string `json:"objectName"` // 评价组名称
	ObjectId int64 `json:"objectId"` // 评价组id
	MondayTeachers []int64 `json:"mondayTeachers,optional"` // 周一被评价教师id
	TuesdayTeachers []int64 `json:"tuesdayTeachers,optional"`
	WednesdayTeachers []int64 `json:"wednesdayTeachers,optional"`
	ThursdayTeachers []int64 `json:"thursdayTeachers,optional"`
	FridayTeachers []int64 `json:"fridayTeachers,optional"`
	SaturdayTeachers []int64 `json:"saturdayTeachers,optional"`
	SundayTeachers []int64 `json:"sundayTeachers,optional"`
	WeekTeachers []int64 `json:"weekTeachers,optional"` // 固定评价组被评价教师id
	Type int64 `json:"type"` // 评价组类型。0-自定义评价组，1-固定评价组
}
```
  


### 3. 评价任务信息

1. 路由定义

- Url: /teacherassess/task/info
- Method: POST
- Request: `AssessTaskInfoRequest`
- Response: `TaskInfoResponse`


2. 类型定义 



```golang
type AssessTaskInfoRequest struct {
	TaskId int64 `json:"taskId"` // 评价任务id
}

type TaskInfoResponse struct {
	TaskId int64 `json:"taskId"` // 评价任务id
	TaskName string `json:"taskName"` // 评价任务名称
	TaskTypeNameList []string `json:"taskTypeNameList"` // 指标名称--原评价分类名称
	TaskTypeIds string `json:"taskTypeIds"` // 指标名称id--原评价分类名称id。以;为分割
	SyncStu int64 `json:"syncStu"` // 是否同步学生评价。0-不同步，1-同步
	ClassList []*TaskLevel `json:"classList"` // 班级评价等级
	TeacherList []*TaskLevel `json:"teacherList"` // 教师评价等级
}
```
  


### 4. 新增评价任务

1. 路由定义

- Url: /teacherassess/task/add
- Method: POST
- Request: `AddTaskRequest`
- Response: `-`


2. 类型定义 



```golang
type AddTaskRequest struct {
	TaskName string `json:"taskName"` // 评价任务名称
	TaskTypeIdList []int64 `json:"taskTypeIdList"` // 指标名称id--原评价分类id
	SyncStu int64 `json:"syncStu"` // 是否同步学生评价。0-不同步，1-同步
	ClassList []*TaskLevel `json:"classList"` // 班级评价等级
	TeacherList []*TaskLevel `json:"teacherList"` // 教师评价等级
}
```
  


### 5. 编辑评价任务

1. 路由定义

- Url: /teacherassess/task/update
- Method: POST
- Request: `UpdateTaskRequest`
- Response: `-`


2. 类型定义 



```golang
type UpdateTaskRequest struct {
	TaskId int64 `json:"taskId"` // 评价任务id
	TaskName string `json:"taskName"` // 评价任务名称
	TaskTypeIdList []int64 `json:"taskTypeIdList"` // 指标名称id--原评价分类id
	SyncStu int64 `json:"syncStu"` // 是否同步学生评价。0-不同步，1-同步
	ClassList []*TaskLevel `json:"classList"` // 班级评价等级
	TeacherList []*TaskLevel `json:"teacherList"` // 教师评价等级
}
```
  

