syntax = "v1"

info(
    title: "base api"
    desc: "base api"
    author: "Ryan SU"
    email: "yuansu.china.work@gmail.com"
    version: "v1.0"
)

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

// Basic ID request | 基础id参数请求
type IDReq {
    // ID
    // Required: true
    Id  uint64 `json:"id" validate:"number"`
}


// Basic ID request in path | 基础ID地址参数请求
type IDPathReq {
    // ID
    // Required: true
    Id  uint64 `path:"id"`
}

// Basic UUID request | 基础UUID参数请求
type UUIDReq {
    // UUID
    // Required: true
    // Max length: 36
    UUID string `json:"UUID" validate:"len=36"`
}

// The base response data | 基础信息
// swagger:model BaseInfo
type BaseInfo {
    // ID
    Id        uint64    `json:"id"`

    // Create date | 创建日期
    CreatedAt int64     `json:"createdAt,optional"`

    // Update date | 更新日期
    UpdatedAt int64     `json:"updatedAt,optional"`
}

// The request params of setting boolean status | 设置状态参数
type StatusCodeReq {
    // ID
    // Required: true
    Id     uint64  `json:"id" validate:"number"`

    // Status code | 状态码
    // Required: true
    Status uint64  `json:"status" validate:"number"`
}