# 文本序列化和反序列化

go-zero针对文本的序列化和反序列化主要在三个地方使用

* http api请求体的反序列化
* http api返回体的序列化
* 配置文件的反序列化

本文假定读者已经定义过api文件以及修改过配置文件，如不熟悉，可参照

* [快速构建高并发微服务](shorturl.md)
* [快速构建高并发微服务](bookstore.md)

## 1. http api请求体的反序列化

在反序列化的过程中的针对请求数据的`数据格式`以及`数据校验`需求，go-zero实现了自己的一套反序列化机制

### 1.1 `数据格式`以订单order.api文件为例

```go
type (
	createOrderReq struct {
		token     string `path:"token"`     // 用户token
		productId string `json:"productId"` // 商品ID
		num       int    `json:"num"`       // 商品数量
	}
	createOrderRes struct {
		success bool `json:"success"` // 是否成功
	}
	findOrderReq struct {
		token    string `path:"token"`    // 用户token
		page     int    `form:"page"`     // 页数
		pageSize int8   `form:"pageSize"` // 页大小
	}
	findOrderRes struct {
		orderInfo []orderInfo `json:"orderInfo"` // 商品ID
	}
	orderInfo struct {
		productId   string `json:"productId"`   // 商品ID
		productName string `json:"productName"` // 商品名称
		num         int    `json:"num"`         // 商品数量
	}
	deleteOrderReq struct {
		id string `path:"id"`
	}
	deleteOrderRes struct {
		success bool `json:"success"` // 是否成功
	}
)

service order {
    @doc(
        summary: 创建订单
    )
    @server(
        handler: CreateOrderHandler
    )
    post /order/add/:token(createOrderReq) returns(createOrderRes)

    @doc(
        summary: 获取订单
    )
    @server(
        handler: FindOrderHandler
    )
    get /order/find/:token(findOrderReq) returns(findOrderRes)

    @doc(
        summary: 删除订单
    )
    @server(
        handler: DeleteOrderHandler
    )
    delete /order/:id(deleteOrderReq) returns(deleteOrderRes)
}
```

http api请求体的反序列化的tag有三种：

* `path`：http url 路径中参数反序列化
  * `/order/add/1234567`会解析出来token为1234567
* `form`：http  form表单反序列化，需要 header头添加  Content-Type: multipart/form-data
  * `/order/find/1234567?page=1&pageSize=20`会解析出来token为1234567，page为1，pageSize为20

* `json`：http request json body反序列化，需要 header头添加  Content-Type: application/json
  * `{"productId":"321","num":1}`会解析出来productId为321，num为1

### 1.2 `数据校验`以用户user.api文件为例

```go
type (
	createUserReq struct {
		age    int8   `json:"age,default=20,range=(12:100]"` // 年龄
		name   string `json:"name"`                          // 名字
		alias  string `json:"alias,optional"`                // 别名
		sex    string `json:"sex,options=male|female"`       // 性别
		avatar string `json:"avatar,default=default.png"`    // 头像
	}
	createUserRes struct {
		success bool `json:"success"` // 是否成功
	}
)

service user {
    @doc(
        summary: 创建订单
    )
    @server(
        handler: CreateUserHandler
    )
    post /user/add(createUserReq) returns(createUserRes)
}
```

数据校验有很多种方式，包括以下但不限：

* `age`：默认不输入为20，输入则取值范围为(12:100]，前开后闭
* `name`：必填，不可为空
* `alias`：选填，可为空
* `sex`：必填，取值为`male`或`female`
* `avatar`：选填，默认为`default.png`

更多详情参见[unmarshaler_test.go](../core/mapping/unmarshaler_test.go)

## 2. http api返回体的序列化

* 使用官方默认的`encoding/json`包序列化，在此不再累赘

## 3. 配置文件的反序列化

* `配置文件的反序列化`和`http api请求体的反序列化`使用同一套解析规则，可参照`http api请求体的反序列化`
