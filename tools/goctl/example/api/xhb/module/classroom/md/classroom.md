
### 1. 用户在班级信息

1. 路由定义

- Url: /classroom/users/info
- Method: POST
- Request: `ClassroomUsersInfoRequest`
- Response: `ClassroomUsersInfoResponse`

2. 请求定义


```golang
type ClassroomUsersInfoRequest struct {
	Token string `json:"token"`
	ClassroomId string `json:"classroomId,optional"`
	UserIds []string `json:"userIds"`
}
```


3. 返回定义


```golang
type ClassroomUsersInfoResponse struct {
	UserView []*UserView `json:"users,omitempty"`
}
```
  


### 2. 查找班级和班级里的成员列表

1. 路由定义

- Url: /classroom/and/members/:token/:speakable
- Method: GET
- Request: `QueryClassroomsAndMembersRequest`
- Response: `QueryClassroomsAndMembersResponse`

2. 请求定义


```golang
type QueryClassroomsAndMembersRequest struct {
	Token string `path:"token"`
	Speakable string `path:"speakable"`
	Cursor int64 `form:"cursor,optional"`
	Limit int8 `form:"limit,optional"`
}
```


3. 返回定义


```golang
type QueryClassroomsAndMembersResponse struct {
	Classrooms []*ClassroomView `json:"classroomViews,omitempty"`
	Cursor int64 `json:"cursor,omitempty"`
	HaveNextPage bool `json:"haveNextPage"`
}
```
  


### 3. 创建一个班级并且加入人员

1. 路由定义

- Url: /classroom/create/add/members/:token
- Method: POST
- Request: `CreateClassroomAddMembersRequest`
- Response: `ClassroomView`

2. 请求定义


```golang
type CreateClassroomAddMembersRequest struct {
	Token string `path:"token"`
	Name string `json:"name,optional"`
	BadgeId string `json:"badgeId,optional"` // 班徽id
	BadgeType string `json:"badgeType,optional"` // 班徽类型
	ClassType string `json:"classType,optional"` // 班级类型
	ManagerName string `json:"managerName,optional"` // 创建者显示名
	ClassroomIds []string `json:"classroomIds,optional"` // 班级列表
	Users []string `json:"users,optional"` // 人员列表
}
```


3. 返回定义


```golang
type ClassroomView struct {
	Id string `json:"id"`
	Name string `json:"name"` // 班级名称
	Code string `json:"code"` // 班级号,唯一
	BadgeId string `json:"badgeId,omitempty"` // 班徽id
	BadgeType string `json:"badgeType"` // 班徽类型
	Members []*UserView `json:"members,omitempty"` // 班级成员列表
	WantJoinMembers []string `json:"wantJoinMembers,omitempty"` // 待加入的成员姓名
	ApplyMembers []*ApplyMember `json:"applyMembers,omitempty"` //新的申请的成员信息(申请加入和申请退出)
	MembersCount int `json:"membersCount"` // 成员数量
	TeacherName string `json:"teacherName"` // 老师姓名
	JoinState string `json:"joinState,omitempty"` // 是否状态为申请中
	SpeakableInClassroom bool `json:"speakableInClassroom"` // 是否可在班级发言
	TouchTime string `json:"touchTime"` // 根据触碰时间智能排序班级
	ClassType string `json:"classType"` // 班级类型
	ManagerId string `json:"managerId"` // 班级创建者id
	Stage string `json:"stage,omitempty"` // 学段(高中，初中)
	Grade string `json:"grade,omitempty"` // 年级
	ClassNo string `json:"classNo,omitempty"` // 班级号（1-20班）
	Subject string `json:"subject,omitempty"` // 学科
	School string `json:"school,omitempty"` // 学校
	SchoolId string `json:"schoolId,omitempty"` // 学校
	MutualDisclosure bool `json:"mutualDisclosure"` // 任课教师相互可见
	Hide bool `json:"hide"` // 搜索隐藏
	Role string `json:"role,omitempty"` // 班级角色
	Vip bool `json:"vip"` // 是否是整校班级
	BgImage string `json:"bgImage,omitempty"` // 班级背景图片
	Locking bool `json:"locking"` // 班级是否锁定
	BindCloud bool `json:"bindCloud"` //是否绑定晓教云
	IsFull bool `json:"isFull"` //是否达到班级人数限制
	UploadClassTableMembers []*UploadClassTableMember `json:"uploadClassTableMembers"` //申请课程表的人员
	ApplyLeaveMembers []*ApplyLeaveMember `json:"applyLeaveMembers"` // 提出请假申请的人员
	LastUpdateTime int64 `json:"lastUpdateTime,omitempty"` // 上次更新时间
	CreateTime int64 `json:"createTime,omitempty"`
	JoinTime int64 `json:"joinTime,omitempty"`
}
```
  


### 4. 创建班级

1. 路由定义

- Url: /classroom/create/:token
- Method: POST
- Request: `CreateClassroomRequest`
- Response: `ClassroomView`

2. 请求定义


```golang
type CreateClassroomRequest struct {
	Token string `path:"token"`
	Name string `json:"name,optional"` // 班级名称
	BadgeId string `json:"badgeId,optional"` // 班徽id
	BadgeType string `json:"badgeType,optional"` // 班徽类型
	ClassType string `json:"classType,optional"` // 班级类型
	ManagerName string `json:"managerName,optional"` // 创建者显示名
	School string `json:"school,optional"` // 学校
	Grade string `json:"grade,optional"` // 年级
	Stage string `json:"stage,optional"` // 学段(高中，初中)
	Subject string `json:"subject,optional"` // 学科
	ClassNo string `json:"classNo,optional"` // 班级号（1-20班)
	BindCloud bool `json:"bindCloud,optional"` // 是否绑定晓教云
}
```


3. 返回定义


```golang
type ClassroomView struct {
	Id string `json:"id"`
	Name string `json:"name"` // 班级名称
	Code string `json:"code"` // 班级号,唯一
	BadgeId string `json:"badgeId,omitempty"` // 班徽id
	BadgeType string `json:"badgeType"` // 班徽类型
	Members []*UserView `json:"members,omitempty"` // 班级成员列表
	WantJoinMembers []string `json:"wantJoinMembers,omitempty"` // 待加入的成员姓名
	ApplyMembers []*ApplyMember `json:"applyMembers,omitempty"` //新的申请的成员信息(申请加入和申请退出)
	MembersCount int `json:"membersCount"` // 成员数量
	TeacherName string `json:"teacherName"` // 老师姓名
	JoinState string `json:"joinState,omitempty"` // 是否状态为申请中
	SpeakableInClassroom bool `json:"speakableInClassroom"` // 是否可在班级发言
	TouchTime string `json:"touchTime"` // 根据触碰时间智能排序班级
	ClassType string `json:"classType"` // 班级类型
	ManagerId string `json:"managerId"` // 班级创建者id
	Stage string `json:"stage,omitempty"` // 学段(高中，初中)
	Grade string `json:"grade,omitempty"` // 年级
	ClassNo string `json:"classNo,omitempty"` // 班级号（1-20班）
	Subject string `json:"subject,omitempty"` // 学科
	School string `json:"school,omitempty"` // 学校
	SchoolId string `json:"schoolId,omitempty"` // 学校
	MutualDisclosure bool `json:"mutualDisclosure"` // 任课教师相互可见
	Hide bool `json:"hide"` // 搜索隐藏
	Role string `json:"role,omitempty"` // 班级角色
	Vip bool `json:"vip"` // 是否是整校班级
	BgImage string `json:"bgImage,omitempty"` // 班级背景图片
	Locking bool `json:"locking"` // 班级是否锁定
	BindCloud bool `json:"bindCloud"` //是否绑定晓教云
	IsFull bool `json:"isFull"` //是否达到班级人数限制
	UploadClassTableMembers []*UploadClassTableMember `json:"uploadClassTableMembers"` //申请课程表的人员
	ApplyLeaveMembers []*ApplyLeaveMember `json:"applyLeaveMembers"` // 提出请假申请的人员
	LastUpdateTime int64 `json:"lastUpdateTime,omitempty"` // 上次更新时间
	CreateTime int64 `json:"createTime,omitempty"`
	JoinTime int64 `json:"joinTime,omitempty"`
}
```
  


### 5. 根据班级码查找班级

1. 路由定义

- Url: /classroom/lookup/:code
- Method: GET
- Request: `QueryClassroomByCodeRequest`
- Response: `ClassroomView`

2. 请求定义


```golang
type QueryClassroomByCodeRequest struct {
	Code string `path:"code"`
}
```


3. 返回定义


```golang
type ClassroomView struct {
	Id string `json:"id"`
	Name string `json:"name"` // 班级名称
	Code string `json:"code"` // 班级号,唯一
	BadgeId string `json:"badgeId,omitempty"` // 班徽id
	BadgeType string `json:"badgeType"` // 班徽类型
	Members []*UserView `json:"members,omitempty"` // 班级成员列表
	WantJoinMembers []string `json:"wantJoinMembers,omitempty"` // 待加入的成员姓名
	ApplyMembers []*ApplyMember `json:"applyMembers,omitempty"` //新的申请的成员信息(申请加入和申请退出)
	MembersCount int `json:"membersCount"` // 成员数量
	TeacherName string `json:"teacherName"` // 老师姓名
	JoinState string `json:"joinState,omitempty"` // 是否状态为申请中
	SpeakableInClassroom bool `json:"speakableInClassroom"` // 是否可在班级发言
	TouchTime string `json:"touchTime"` // 根据触碰时间智能排序班级
	ClassType string `json:"classType"` // 班级类型
	ManagerId string `json:"managerId"` // 班级创建者id
	Stage string `json:"stage,omitempty"` // 学段(高中，初中)
	Grade string `json:"grade,omitempty"` // 年级
	ClassNo string `json:"classNo,omitempty"` // 班级号（1-20班）
	Subject string `json:"subject,omitempty"` // 学科
	School string `json:"school,omitempty"` // 学校
	SchoolId string `json:"schoolId,omitempty"` // 学校
	MutualDisclosure bool `json:"mutualDisclosure"` // 任课教师相互可见
	Hide bool `json:"hide"` // 搜索隐藏
	Role string `json:"role,omitempty"` // 班级角色
	Vip bool `json:"vip"` // 是否是整校班级
	BgImage string `json:"bgImage,omitempty"` // 班级背景图片
	Locking bool `json:"locking"` // 班级是否锁定
	BindCloud bool `json:"bindCloud"` //是否绑定晓教云
	IsFull bool `json:"isFull"` //是否达到班级人数限制
	UploadClassTableMembers []*UploadClassTableMember `json:"uploadClassTableMembers"` //申请课程表的人员
	ApplyLeaveMembers []*ApplyLeaveMember `json:"applyLeaveMembers"` // 提出请假申请的人员
	LastUpdateTime int64 `json:"lastUpdateTime,omitempty"` // 上次更新时间
	CreateTime int64 `json:"createTime,omitempty"`
	JoinTime int64 `json:"joinTime,omitempty"`
}
```
  


### 6. 查找班级，通过班级号或者老师手机号码

1. 路由定义

- Url: /classroom/search/by/v2/:token/:keyword
- Method: GET
- Request: `SearchClassroomByCodeOrMobileRequest`
- Response: `SearchClassroomByCodeOrMobileResponse`

2. 请求定义


```golang
type SearchClassroomByCodeOrMobileRequest struct {
	Token string `path:"token"`
	Keyword string `path:"keyword"`
}
```


3. 返回定义


```golang
type SearchClassroomByCodeOrMobileResponse struct {
	Classrooms []*ClassroomView `json:"classrooms,omitempty"`
}
```
  


### 7. 修改班级信息

1. 路由定义

- Url: /classroom/info/update/:token
- Method: POST
- Request: `UpdateClassroomInfoRequest`
- Response: `UpdateClassroomInfoResponse`

2. 请求定义


```golang
type UpdateClassroomInfoRequest struct {
	Token string `path:"token"`
	Name string `json:"name,optional"` // 班级名称
	BadgeId string `json:"badgeId,optional"` // 班徽id
	BadgeType string `json:"badgeType,optional"` // 班徽类型
	ClassType string `json:"classType,optional"` // 班级类型
	ManagerName string `json:"managerName,optional"` // 创建者显示名
	Stage string `json:"stage,optional"` // 学段(高中，初中)
	Grade string `json:"grade,optional"` // 年级
	ClassNo string `json:"classNo,optional"` // 班级号（1-20班)
	Subject string `json:"subject,optional"` // 学科
	School string `json:"school,optional"` // 学校
	ClassroomId string `json:"classroomId,optional"` // 班级id
	UserId string `json:"userId,optional"` // 用户id
	BgImage string `json:"bgImage,optional"` // 班级背景图
	MutualDisclosure bool `json:"mutualDisclosure,optional"` // 任课教师相互可见
	Hide bool `json:"hide,optional"` // 是否可以根据手机号查询
	Locking bool `json:"locking,optional"` //是否锁定此班级
	BindCloud bool `json:"bindCloud,optional"` //是否绑定晓教云
}
```


3. 返回定义


```golang
type UpdateClassroomInfoResponse struct {
	Successful bool `json:"successful"`
}
```
  


### 8. 以小孩子的长辈身份申请加入班级

1. 路由定义

- Url: /classroom/join/with/child/:token
- Method: POST
- Request: `ApplyJoinClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type ApplyJoinClassroomRequest struct {
	Token string `path:"token"`
	Message string `json:"message,optional"` // 验证信息
	DisplayName string `json:"displayName,optional"` // 用户在此班级的显示名称(以老师身份加入)
	ClassroomId string `json:"classroomId"` // 需要加入的班级id
	ChildName string `json:"childName,optional"` // 孩子姓名（以家长身份加入）
	Relationship string `json:"relationship,optional"` // 关系
	JToken string `json:"token,optional"`
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 9. 解散班级

1. 路由定义

- Url: /classroom/remove/:token/:classroomId
- Method: POST
- Request: `RemoveClassroomRequest`
- Response: `BoolResponse`

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
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 10. 获取待审核的班级成员

1. 路由定义

- Url: /classroom/audit/members/:token/:classroomId
- Method: GET
- Request: `GetClassroomAuditMembersRequest`
- Response: `GetClassroomAuditMembersResponse`

2. 请求定义


```golang
type GetClassroomAuditMembersRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type GetClassroomAuditMembersResponse struct {
	ApplyJoin []ApplyUser `json:"applyJoin"` // 申请加入班级的成员
	ApplyQuit []ApplyUser `json:"applyQuit"` // 申请退出班级的成员
}
```
  


### 11. 同意/拒绝批量用户加入班级

1. 路由定义

- Url: /classroom/batch/audit/:token
- Method: POST
- Request: `BatchAuditJoinClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type BatchAuditJoinClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	Members []string `json:"members"`
	Agree bool `json:"agree,optional"`
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 12. 同意/拒绝批量用户退出班级

1. 路由定义

- Url: /classroom/batch/quit/:token
- Method: POST
- Request: `BatchQuitClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type BatchQuitClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	Members []string `json:"members"`
	Agree bool `json:"agree,optional"`
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 13. 获取班级成员和教师

1. 路由定义

- Url: /classroom/members/and/teacher/:token/:classroomId
- Method: GET
- Request: `QueryClassroomMembersAndTeachersRequest`
- Response: `QueryClassroomMembersAndTeachersResponse`

2. 请求定义


```golang
type QueryClassroomMembersAndTeachersRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
}
```


3. 返回定义


```golang
type QueryClassroomMembersAndTeachersResponse struct {
	Members []*UserView `json:"members"` // 班级成员列表
	Teachers []*UserView `json:"teachers"` // 班级教师列表
}
```
  


### 14. 移除班级成员(班级创建者移除教师和普通成员)

1. 路由定义

- Url: /classroom/member/remove/:token
- Method: POST
- Request: `RemoveClassroomMemberRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type RemoveClassroomMemberRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	MemberId string `json:"memberId"` // 成员id
	Message string `json:"message,optional"` // 移除原因
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 15. 班级内权限设置，普通成员，可发言成员，任课老师，转让班级相互设置

1. 路由定义

- Url: /classroom/teacher/setting/:token
- Method: POST
- Request: `SetClassroomPermissionRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type SetClassroomPermissionRequest struct {
	Token string `path:"token"`
	UserId string `json:"userId"` // 对方用户id
	ClassroomId string `json:"classroomId"` // 所在班级id
	Role string `json:"role"` // 将要改变的角色  COMMON SPEAKABLE ADMIN MANAGER
	Subject string `json:"subject,optional"` // 学科
	Code string `json:"code,optional"` //转让班级验证码
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 16. 申请退出班级

1. 路由定义

- Url: /classroom/quit/:token
- Method: POST
- Request: `ApplyQuitClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type ApplyQuitClassroomRequest struct {
	Token string `path:"token"`
	Message string `json:"message"` // 验证信息
	ClassroomId string `json:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 17. 任课教师退出班级/班级创建者转让班级

1. 路由定义

- Url: /teacher/quit/classroom/:token/:classroomId
- Method: PUT
- Request: `TeacherQuitClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type TeacherQuitClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"` // 班级id
	TargetTeacher string `form:"targetTeacher,optional"` // 转让的目的老师
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 18. 设置用户在班级显示名称

1. 路由定义

- Url: /classroom/member/displayname/:token
- Method: POST
- Request: `SetClassroomUserInfoRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type SetClassroomUserInfoRequest struct {
	Token string `path:"token"`
	UserId string `json:"userId"`
	ClassroomId string `json:"classroomId"`
	Owner bool `json:"owner,optional"`
	ChildName string `json:"childName,optional"`
	Relationship string `json:"relationship,optional"` // 关系
	DisplayName string `json:"displayName,optional"` // 别名
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 19. 查找用户自己孩子信息

1. 路由定义

- Url: /find/child/:token
- Method: GET
- Request: `FindChildByUserIdRequest`
- Response: `FindChildByUserIdResponse`

2. 请求定义


```golang
type FindChildByUserIdRequest struct {
	Token string `path:"token"`
}
```


3. 返回定义


```golang
type FindChildByUserIdResponse struct {
	Classroom *ClassroomView `json:"classroom,omitempty"` // 班级
	Children []string `json:"children,omitempty"` // 孩子列表
	Childrens []ChildrenStruct `json:"childrens,omitempty"` // 孩子列表
}
```
  


### 20. 验证班级中是否有重名的人

1. 路由定义

- Url: /exist/displayname/in/classroom/:token/:classroomId
- Method: POST
- Request: `DisplayNameExistInClassroomRequest`
- Response: `DisplayNameExistInClassroomResponse`

2. 请求定义


```golang
type DisplayNameExistInClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"` // 班级id
	Realname string `json:"realname,optional"` // 孩子姓名
	UserId string `json:"userId,optional"` // 用户id
}
```


3. 返回定义


```golang
type DisplayNameExistInClassroomResponse struct {
	Exist bool `json:"exist"` // 是否重名
	Desc string `json:"desc"` //4.7.3
}
```
  


### 21. 验证班级中是否有重名的孩子

1. 路由定义

- Url: /exist/childname/in/classroom/:token
- Method: POST
- Request: `ChildNameExistInClassroomRequest`
- Response: `ChildNameExistInClassroomResponse`

2. 请求定义


```golang
type ChildNameExistInClassroomRequest struct {
	Token string `path:"token"`
	ChildName string `json:"childName"` // 孩子姓名
	ClassroomId string `json:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type ChildNameExistInClassroomResponse struct {
	Exist bool `json:"exist"` // 是否重名
	Relationship string `json:"relationship"` // 关系
	Mobile string `json:"mobile"` // 手机号码
	UserId string `json:"userId"` // 用户id
}
```
  


### 22. 班级成员转移

1. 路由定义

- Url: /classroom/transfer/:token
- Method: POST
- Request: `TransferRoomMemberRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type TransferRoomMemberRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroom"` // 转让的目的老师
	MembersIds []string `json:"members"` // 转移的用户
	FromClassroomId string `json:"fromClassroom"` // 原班级
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 23. 生成班级相关验证码

1. 路由定义

- Url: /classroom/activation/:token
- Method: POST
- Request: `GenerateClassroomVerificationCodeRequest`
- Response: `GenerateClassroomVerificationCodeResponse`

2. 请求定义


```golang
type GenerateClassroomVerificationCodeRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
}
```


3. 返回定义


```golang
type GenerateClassroomVerificationCodeResponse struct {
}
```
  


### 24. 用户班级（全量接口）

1. 路由定义

- Url: /user/classrooms/:token
- Method: GET
- Request: `UserClassroomsRequest`
- Response: `UserClassroomsResponse`

2. 请求定义


```golang
type UserClassroomsRequest struct {
	Token string `path:"token"`
}
```


3. 返回定义


```golang
type UserClassroomsResponse struct {
	LastTimestamp int64 `json:"lastTimestamp"` // 最近拉取时间
	ClassroomViews []*ClassroomView `json:"classrooms"` // 班级列表
}
```
  


### 25. 用户在班级信息

1. 路由定义

- Url: /v2/classroom/user/info/:token/:classroomId/:userId
- Method: GET
- Request: `ClassroomUserInfoRequest`
- Response: `UserView`

2. 请求定义


```golang
type ClassroomUserInfoRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
	UserId string `path:"userId"`
}
```


3. 返回定义


```golang
type UserView struct {
	Id string `json:"id"`
	Username string `json:"username,omitempty"`
	Realname string `json:"realname,omitempty"`
	Email string `json:"email,omitempty"`
	Mobile string `json:"mobile"`
	Address string `json:"address,omitempty"`
	Role string `json:"role"`
	Avatar string `json:"avatar,omitempty"`
	Pinyin string `json:"pinyin,omitempty"`
	School string `json:"school,omitempty"`
	DisplayName string `json:"displayName"`
	UserType string `json:"userType,omitempty"`
	SpeakableInClassroom bool `json:"speakableInClassroom"`
	ClassroomRole string `json:"classroomRole"`
	Subject string `json:"subject,omitempty"`
	CreateTime int64 `json:"createTime,omitempty"`
}
```
  


### 26. 获取未注册用户

1. 路由定义

- Url: /classroom/unregister/user/:token/:classroomId
- Method: GET
- Request: `ClassroomUnregisterUserRequest`
- Response: `ClassroomUnregisterUserResponse`

2. 请求定义


```golang
type ClassroomUnregisterUserRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"` // 班级名称
}
```


3. 返回定义


```golang
type ClassroomUnregisterUserResponse struct {
	Users []*RoomUnregisterUser `json:"result,omitempty"`
}
```
  


### 27. 删除未注册用户

1. 路由定义

- Url: /classroom/remove/unregister/user/:token
- Method: DELETE
- Request: `RemoveClassroomUnregisterUserRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type RemoveClassroomUnregisterUserRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	Mobile string `json:"mobile"`
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 28. 老师ding一下该班级未注册人员

1. 路由定义

- Url: /classroom/remove/unregister/user/:token
- Method: POST
- Request: `UnregisterDingRequest`
- Response: `UnregisterDingResponse`

2. 请求定义


```golang
type UnregisterDingRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"`
	Mobiles []string `json:"mobiles"`
	Type string `json:"type"`
}
```


3. 返回定义


```golang
type UnregisterDingResponse struct {
}
```
  


### 29. 邀请未注册用户

1. 路由定义

- Url: /invite/user/join/classroom/:mobile/:code
- Method: PUT
- Request: `InvitedJoinClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type InvitedJoinClassroomRequest struct {
	Mobile string `path:"mobile"` // 验证信息
	Code string `path:"code"` // 用户在此班级的显示名称(以老师身份加入)
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 30. 班级添加成员

1. 路由定义

- Url: /classroom/add/members/:token/:classroomId
- Method: POST
- Request: `AddClassroomMemberRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type AddClassroomMemberRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"`
	ClassroomIds []string `json:"classroomIds,optional"` // 班级列表
	Users []string `json:"users,optional"` // 人员列表
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 31. 班级催一下老师同意申请

1. 路由定义

- Url: /remind/teacher/for/apply/:token/:classroomId
- Method: GET
- Request: `RemindTeacherApplyRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type RemindTeacherApplyRequest struct {
	Token string `path:"token"`
	ClassroomId string `path:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 32. 用户取消申请

1. 路由定义

- Url: /classroom/applied/:token/:classroomId
- Method: DELETE
- Request: `CancelApplyRequest`
- Response: `CancelApplyResponse`

2. 请求定义


```golang
type CancelApplyRequest struct {
	ClassroomId string `path:"classroomId"`
	Token string `path:"token"`
}
```


3. 返回定义


```golang
type CancelApplyResponse struct {
	Successful bool `json:"successful"`
}
```
  


### 33. 删除退出班级申请

1. 路由定义

- Url: /classroom/quit-apply/:token
- Method: POST
- Request: `CancelQuitClassroomRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type CancelQuitClassroomRequest struct {
	Token string `path:"token"`
	ClassroomId string `json:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 34. 用户班级(增量接口)

1. 路由定义

- Url: /user/classrooms/change/:token
- Method: GET
- Request: `UserClassroomsChangeRequest`
- Response: `UserClassroomsChangeResponse`

2. 请求定义


```golang
type UserClassroomsChangeRequest struct {
	Token string `path:"token"`
	Timestamp int64 `form:"timestamp,optional"`
}
```


3. 返回定义


```golang
type UserClassroomsChangeResponse struct {
UserClassroomsChangeTimestamp
	CheckMessageId string `json:"checkMessageId,omitempty"`
	TotalSize int `json:"totalSize"` // 总记录数
	ClassroomMemberWrong []string `json:"classroomMemberWrong,omitempty"`
	Events []*Event `json:"events,omitempty"` // 事件
}

type UserClassroomsChangeTimestamp struct {
	HasMore bool `json:"hasMore"` // 是否还有更多消息
	LastTimestamp int64 `json:"lastTimestamp"` // 本次获取的最后一条timestamp
}
```
  


### 35. 站外用户注册

1. 路由定义

- Url: /classroom/outside-register
- Method: POST
- Request: `OutsideRegisterRequest`
- Response: `OutsideRegisterResponse`

2. 请求定义


```golang
type OutsideRegisterRequest struct {
	Name string `json:"name"`
	Classroom string `json:"classroom"`
	Mobile string `json:"mobile"`
	Role string `json:"role"` // TEACHER PARENT STUDENT
	Relation string `json:"relation,optional"` // FATHER:爸爸 MOTHER:妈妈 GRANDFATHER:爷爷 GRANDMOTHER:奶奶 MATERNALGRANDFATHER:外公 MATERNALGRANDMOTHER:外婆 TEACHER:老师 STUDENT:学生 UNCLE:叔叔 MOTHER_SISTER:阿姨 UNCLE_BROTHER:大伯 AUNT:姑姑 MOTHER_BROTHER:舅舅 SISTER:姐姐 BROTHER:哥哥 OTHER:其它
	JtToken string `json:"jtToken,optional"` //48小时免审核JtToken
	StudentNumber string `json:"student_number"` // 学号
}
```


3. 返回定义


```golang
type OutsideRegisterResponse struct {
}
```
  


### 36. 班级管理员查看免审核入班成员

1. 路由定义

- Url: /classroom/free/audit/members/:classroomId
- Method: GET
- Request: `GetClassroomFreeAuditMembersRequest`
- Response: `GetClassroomFreeAuditMembersResponse`

2. 请求定义


```golang
type GetClassroomFreeAuditMembersRequest struct {
	ClassroomId string `path:"classroomId"` // 班级id
}
```


3. 返回定义


```golang
type GetClassroomFreeAuditMembersResponse struct {
	FreeAuditMembers []string `json:"freeAuditMembers"` // 班级免审核成员名字
}
```
  


### 37. 修改任课老师任教学科

1. 路由定义

- Url: /user/subject/update/:token
- Method: POST
- Request: `UpdateClassroomSubjectRequest`
- Response: `BoolResponse`

2. 请求定义


```golang
type UpdateClassroomSubjectRequest struct {
	Token string `path:"token"`
	UserId string `json:"userId,optional"` // 对方用户id
	ClassroomId string `json:"classroomId"` // 所在班级id
	Subject string `json:"subject"` // 科目
}
```


3. 返回定义


```golang
type BoolResponse struct {
	Successful bool `json:"successful"` // 是否成功应答
}
```
  


### 38. 用户关联用户(增量接口)

1. 路由定义

- Url: /related/users/change/:token
- Method: GET
- Request: `UserClassroomRelationChangeRequest`
- Response: `UserClassroomsChangeResponse`

2. 请求定义


```golang
type UserClassroomRelationChangeRequest struct {
	Token string `path:"token"`
	Timestamp int64 `form:"timestamp,optional"`
}
```


3. 返回定义


```golang
type UserClassroomsChangeResponse struct {
UserClassroomsChangeTimestamp
	CheckMessageId string `json:"checkMessageId,omitempty"`
	TotalSize int `json:"totalSize"` // 总记录数
	ClassroomMemberWrong []string `json:"classroomMemberWrong,omitempty"`
	Events []*Event `json:"events,omitempty"` // 事件
}

type UserClassroomsChangeTimestamp struct {
	HasMore bool `json:"hasMore"` // 是否还有更多消息
	LastTimestamp int64 `json:"lastTimestamp"` // 本次获取的最后一条timestamp
}
```
  


### 39. 班级通讯录

1. 路由定义

- Url: /classroom/communication/books
- Method: GET
- Request: `CommBooksRequest`
- Response: `CommBooksResponse`

2. 请求定义


```golang
type CommBooksRequest struct {
	Uid string `form:"uid"`
	Cid string `form:"cid"`
}
```


3. 返回定义


```golang
type CommBooksResponse struct {
	Students *ClassMembers `json:"students"`
	Teachers *ClassMembers `json:"teachers"`
}
```
  


### 40. 学生名称搜索

1. 路由定义

- Url: /classroom/student/search
- Method: GET
- Request: `SearchStudentRequest`
- Response: `SearchStudentResponse`

2. 请求定义


```golang
type SearchStudentRequest struct {
	Cid string `form:"cid"`
	Uid string `form:"uid"`
	Name string `form:"name"`
}
```


3. 返回定义


```golang
type SearchStudentResponse struct {
	Students []*SearchStudent `json:"students"`
}
```
  


### 41. 入班

1. 路由定义

- Url: /classroom/apply/join
- Method: POST
- Request: `JoinClassRequest`
- Response: `JoinClassResponse`

2. 请求定义


```golang
type JoinClassRequest struct {
	Uid string `json:"uid"`
	Cid string `json:"cid"`
	Name string `json:"name"`
	ClassNumber string `json:"class_number,optional"`
	Relation string `json:"relation"`
	Mobile string `json:"mobile,optional"`
	CheckMobile string `json:"check_mobile,optional"`
	JToken string `json:"token,optional"`
	Message string `json:"message,optional"`
}
```


3. 返回定义


```golang
type JoinClassResponse struct {
	IsRepeat bool `json:"isRepeat"`
	Name string `json:"name"`
	Relation string `json:"relation"`
	Mobile string `json:"mobile"`
}
```
  


### 42. 用户身份学生信息

1. 路由定义

- Url: /classroom/user/student/info
- Method: GET
- Request: `UserStudentInfoRequest`
- Response: `SearchStudent`

2. 请求定义


```golang
type UserStudentInfoRequest struct {
	Uid string `form:"uid"`
}
```


3. 返回定义


```golang
type SearchStudent struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Alias string `json:"alias"`
	Avatar string `json:"avatar"`
	ClassCount int `json:"classCount"`
	ClassNumber string `json:"classNumber"`
	Mobile string `json:"mobile"`
	Relation string `json:"relation"`
}
```
  


### 43. 编辑学生信息

1. 路由定义

- Url: /classroom/student/update
- Method: POST
- Request: `UpdateStudentRequest`
- Response: `UpdateStudentResponse`

2. 请求定义


```golang
type UpdateStudentRequest struct {
	Uid string `json:"uid"`
	Cid string `json:"cid"`
	Sid string `json:"sid"`
	ClassNumber string `json:"classNumber,optional"`
	Alias string `json:"alias,optional"`
}
```


3. 返回定义


```golang
type UpdateStudentResponse struct {
}
```
  


### 44. 请出班级

1. 路由定义

- Url: /classroom/leave
- Method: POST
- Request: `LeaveClassRequest`
- Response: `LeaveClassResponse`

2. 请求定义


```golang
type LeaveClassRequest struct {
	UserId string `json:"uid"`
	ClassId string `json:"cid"`
	StudentId string `json:"sid,optional"`
	MemberId string `json:"mid,optional"`
	Role string `json:"role"`
	Message string `json:"msg,optional"`
}
```


3. 返回定义


```golang
type LeaveClassResponse struct {
}
```
  


### 45. 用户关联学生(全量接口)

1. 路由定义

- Url: /classroom/related/students/:token
- Method: GET
- Request: `RelatedStudentsRequest`
- Response: `RelatedStudentsResponse`

2. 请求定义


```golang
type RelatedStudentsRequest struct {
	Token string `path:"token"`
}
```


3. 返回定义


```golang
type RelatedStudentsResponse struct {
	LastTimestamp int64 `json:"lastTimestamp"`
	ClassStudents []*ClassStudent `json:"classStudents"`
}
```
  


### 46. 用户关联学生(增量接口)

1. 路由定义

- Url: /classroom/related/students/change/:token
- Method: GET
- Request: `RelatedStudentsChangeRequest`
- Response: `RelatedStudentsResponse`

2. 请求定义


```golang
type RelatedStudentsChangeRequest struct {
	Token string `path:"token"`
	Sid string `form:"sid"`
	Timestamp int64 `form:"timestamp,optional"`
}
```


3. 返回定义


```golang
type RelatedStudentsResponse struct {
	LastTimestamp int64 `json:"lastTimestamp"`
	ClassStudents []*ClassStudent `json:"classStudents"`
}
```
  


### 47. 批量更新学生信息

1. 路由定义

- Url: /classroom/student/batch/update
- Method: POST
- Request: `BatchUpdateStudentRequest`
- Response: `UpdateStudentResponse`

2. 请求定义


```golang
type BatchUpdateStudentRequest struct {
	Uid string `json:"uid"`
	Cid string `json:"cid"`
	StudentBases []*StudentBase `json:"studentBases,optional"`
}
```


3. 返回定义


```golang
type UpdateStudentResponse struct {
}
```
  

