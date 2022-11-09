package errorx

const (
	// DatabaseError
	// normal database error
	DatabaseError string = "database error occur"

	// RedisError
	// normal redis error
	RedisError string = "redis error occur"

	// request error

	// ApiRequestFailed
	// EN: The interface request failed, please try again later!
	// ZH_CN: 请求出错，请稍候重试
	ApiRequestFailed string = "sys.api.apiRequestFailed"

	// CreateSuccess
	// EN: Create successfully
	// ZH_CN: 新建成功
	CreateSuccess string = "common.createSuccess"

	// CreateFailed
	// EN: Create failed
	// ZH_CN: 新建失败
	CreateFailed string = "common.createFailed"

	// UpdateSuccess
	// EN: Update successfully
	// ZH_CN: 更新成功
	UpdateSuccess string = "common.updateSuccess"

	// UpdateFailed
	// EN: Update failed
	// ZH_CN: 更新失败
	UpdateFailed string = "common.updateFailed"

	// DeleteSuccess
	// EN: Delete successfully
	// ZH_CN: 删除成功
	DeleteSuccess string = "common.deleteSuccess"

	// DeleteFailed
	// EN: Delete failed
	// ZH_CN: 删除失败
	DeleteFailed string = "common.deleteFailed"

	// GetInfoSuccess
	// EN: Get information Successfully
	// ZH_CN: 获取信息成功
	GetInfoSuccess string = "common.getInfoSuccess"

	// GetInfoFailed
	// EN: Get information Failed
	// ZH_CN: 获取信息失败
	GetInfoFailed string = "common.getInfoFailed"

	// TargetNotFound
	// EN: Target does not exist
	// ZH_CN: 目标不存在
	TargetNotFound string = "common.targetNotExist"

	// Success
	// EN: Successful
	// ZH_CN: 成功
	Success string = "common.successful"

	// Failed
	// EN: Failed
	// ZH_CN: 失败
	Failed string = "common.failed"

	// InitRunning
	// EN: The initialization is running...
	// ZH_CN: 正在初始化...
	InitRunning string = "sys.init.initializeIsRunning"

	// AlreadyInit
	// EN: The database had been initialized
	// ZH_CN: 数据库已被初始化
	AlreadyInit string = "sys.init.alreadyInit"
)
