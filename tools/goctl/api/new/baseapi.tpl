syntax = "v1"

// The basic response with data | 基础带数据信息
type BaseDataInfo {
    // Error code | 错误代码
    Code int    `json:"code"`

    // Message | 提示信息
    Msg  string `json:"msg"`

    // Data | 数据
    Data string `json:"data,omitempty"`
}

// The basic response with data | 基础带数据信息
type BaseListInfo {
    // The total number of data | 数据总数
    Total uint64 `json:"total"`

    // Data | 数据
    Data string `json:"data,omitempty"`
}

// The basic response without data | 基础不带数据信息
type BaseMsgResp {
    // Error code | 错误代码
    Code int    `json:"code"`

    // Message | 提示信息
    Msg  string `json:"msg"`
}

// The simplest message | 最简单的信息
// swagger:response SimpleMsg
type SimpleMsg {
    // Message | 信息
    Msg string `json:"msg"`
}

// The page request parameters | 列表请求参数
type PageInfo {
    // Page number | 第几页
    // Required: true
    Page   uint64    `json:"page" validate:"number"`

    // Page size | 单页数据行数
    // Required: true
    // Maximum: 100000
    PageSize  uint64    `json:"pageSize" validate:"number,max=100000"`
}

// Basic ID request | 基础ID参数请求
type IDReq {
    // ID
    // Required: true
    Id  uint64 `json:"id" validate:"number"`
}

// Basic IDs request | 基础ID数组参数请求
type IDsReq {
    // IDs
    // Required: true
    Ids  []uint64 `json:"ids"`
}


// Basic ID request | 基础ID地址参数请求
type IDPathReq {
    // ID
    // Required: true
    Id  uint64 `path:"id"`
}

// Basic UUID request | 基础UUID参数请求
type UUIDReq {
    // ID
    // Required: true
    // Max length: 36
    Id string `json:"id" validate:"len=36"`
}

// Basic UUID array request | 基础UUID数组参数请求
type UUIDsReq {
    // Ids
    // Required: true
    Ids []string `json:"ids"`
}

// The base ID response data | 基础ID信息
type BaseIDInfo {
    // ID
    Id        uint64    `json:"id,optional"`

    // Create date | 创建日期
    CreatedAt int64     `json:"createdAt,optional"`

    // Update date | 更新日期
    UpdatedAt int64     `json:"updatedAt,optional"`
}

// The base UUID response data | 基础UUID信息
type BaseUUIDInfo {
    // ID
    Id        string    `json:"id,optional"`

    // Create date | 创建日期
    CreatedAt int64     `json:"createdAt,optional"`

    // Update date | 更新日期
    UpdatedAt int64     `json:"updatedAt,optional"`
}


@server(
	group: base
)

service {{.name}} {
	// Initialize database | 初始化数据库
	@handler initDatabase
	get /init/database returns (BaseMsgResp)
}
