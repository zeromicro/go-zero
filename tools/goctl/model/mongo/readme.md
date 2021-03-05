# mongo生成model

## 背景

在业务务开发中，model(dao)数据访问层是一个服务必不可缺的一层，因此数据库访问的CURD也是必须要对外提供的访问方法， 而CURD在go-zero中就仅存在两种情况

* 带缓存model
* 不带缓存model

从代码结构上来看，C-U-R-D四个方法就是固定的结构，因此我们可以将其交给goctl工具去完成，帮助我们提升开发效率。

## 方案设计

mongo的生成不同于mysql，mysql可以从scheme_information库中读取到一张表的信息（字段名称，数据类型，索引等），
而mongo是文档型数据库，我们暂时无法从db中读取某一条记录来实现字段信息获取，就算有也不一定是完整信息（某些字段可能是omitempty修饰，可有可无）， 这里采用type自己编写+代码生成方式实现

## 使用示例

假设我们需要生成一个usermodel.go的代码文件，其包含用户信息字段有

|字段名称|字段类型|
|---|---|
|_id|bson.ObejctId|
|name|string|

### 编写types.go

```shell
$ vim types.go
```

```golang
package model

//go:generate goctl model mongo -t User
import "github.com/globalsign/mgo/bson"

type User struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}
```

### 生成代码

生成代码的方式有两种

* 命令行生成 在types.go所在文件夹执行命令
    ```shell
    $ goctl model mongo -t User -style gozero
    ```
* 在types.go中添加`//go:generate`，然后点击执行按钮即可生成，内容示例如下：
  ```golang
  //go:generate goctl model mongo -t User
  ```

### 生成示例代码

* usermodel.go

  ```golang
  package model
  
  import (
      "context"
  
      "github.com/globalsign/mgo/bson"
      cachec "github.com/tal-tech/go-zero/core/stores/cache"
      "github.com/tal-tech/go-zero/core/stores/mongoc"
  )
  
  type UserModel interface {
      Insert(data *User, ctx context.Context) error
      FindOne(id string, ctx context.Context) (*User, error)
      Update(data *User, ctx context.Context) error
      Delete(id string, ctx context.Context) error
  }
  
  type defaultUserModel struct {
      *mongoc.Model
  }
  
  func NewUserModel(url, collection string, c cachec.CacheConf) UserModel {
      return &defaultUserModel{
          Model: mongoc.MustNewModel(url, collection, c),
      }
  }
  
  func (m *defaultUserModel) Insert(data *User, ctx context.Context) error {
      if !data.ID.Valid() {
          data.ID = bson.NewObjectId()
      }
  
      session, err := m.TakeSession()
      if err != nil {
          return err
      }
  
      defer m.PutSession(session)
      return m.GetCollection(session).Insert(data)
  }
  
  func (m *defaultUserModel) FindOne(id string, ctx context.Context) (*User, error) {
      if !bson.IsObjectIdHex(id) {
          return nil, ErrInvalidObjectId
      }
  
      session, err := m.TakeSession()
      if err != nil {
          return nil, err
      }
  
      defer m.PutSession(session)
      var data User
  
      err = m.GetCollection(session).FindOneIdNoCache(&data, bson.ObjectIdHex(id))
      switch err {
      case nil:
          return &data, nil
      case mongoc.ErrNotFound:
          return nil, ErrNotFound
      default:
          return nil, err
      }
  }
  
  func (m *defaultUserModel) Update(data *User, ctx context.Context) error {
      session, err := m.TakeSession()
      if err != nil {
          return err
      }
  
      defer m.PutSession(session)
  
      return m.GetCollection(session).UpdateIdNoCache(data.ID, data)
  }
  
  func (m *defaultUserModel) Delete(id string, ctx context.Context) error {
      session, err := m.TakeSession()
      if err != nil {
          return err
      }
  
      defer m.PutSession(session)
  
      return m.GetCollection(session).RemoveIdNoCache(bson.ObjectIdHex(id))
  }
  ```

* error.go

  ```golang
  package model

  import "errors"
  
  var ErrNotFound = errors.New("not found")
  var ErrInvalidObjectId = errors.New("invalid objectId")
  ```

### 文件目录预览

```text
.
├── error.go
├── types.go
└── usermodel.go

```

## 命令预览

```text
NAME:
   goctl model - generate model code

USAGE:
   goctl model command [command options] [arguments...]

COMMANDS:
   mysql  generate mysql model
   mongo  generate mongo model

OPTIONS:
   --help, -h  show help
```

```text
NAME:
   goctl model mongo - generate mongo model

USAGE:
   goctl model mongo [command options] [arguments...]

OPTIONS:
   --type value, -t value  specified model type name
   --cache, -c             generate code with cache [optional]
   --dir value, -d value   the target dir
   --style value           the file naming format, see [https://github.com/tal-tech/go-zero/tree/master/tools/goctl/config/readme.md]

```

> 温馨提示
> 
> `--type` 支持slice传值，示例 `goctl model mongo -t=User -t=Class`
## 注意事项

types.go本质上与xxxmodel.go无关，只是将type定义部分交给开发人员自己编写了，在xxxmodel.go中，mongo文档的存储结构必须包含
`_id`字段，对应到types中的field为`ID`，model中的findOne,update均以data.ID来进行操作的，当然，如果不符合你的命名风格，你也 可以修改模板，只要保证`id`
在types中的field名称和模板中一致就行。