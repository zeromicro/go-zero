# Goctl Model

goctl model 为go-zero下的工具模块中的组件之一，目前支持识别mysql ddl进行model层代码生成，通过命令行或者idea插件（即将支持）可以有选择地生成带redis cache或者不带redis cache的代码逻辑。

## 快速开始

* 通过ddl生成

    ```shell script
    goctl model mysql ddl -src="./sql/user.sql" -dir="./sql/model" -c=true
    ```

    执行上述命令后即可快速生成CURD代码。

    ```Plain Text
    model
    │   ├── error.go
    │   └── usermodel.go
    ```

* 通过datasource生成

    ```shell script
    goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="table1,table2"  -dir="./model"
    ```

* 生成代码示例
  
  ``` go
    package model

    import (
    	"database/sql"
    	"fmt"
    	"strings"
    	"time"

    	"github.com/tal-tech/go-zero/core/stores/cache"
    	"github.com/tal-tech/go-zero/core/stores/sqlc"
    	"github.com/tal-tech/go-zero/core/stores/sqlx"
    	"github.com/tal-tech/go-zero/core/stringx"
    	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
    )

    var (
    	userFieldNames          = builderx.FieldNames(&User{})
    	userRows                = strings.Join(userFieldNames, ",")
    	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), ",")
    	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

    	cacheUserMobilePrefix = "cache#User#mobile#"
    	cacheUserIdPrefix     = "cache#User#id#"
    	cacheUserNamePrefix   = "cache#User#name#"
    )

    type (
    	UserModel struct {
    		sqlc.CachedConn
    		table string
    	}

    	User struct {
    		Id         int64     `db:"id"`
    		Name       string    `db:"name"`     // 用户名称
    		Password   string    `db:"password"` // 用户密码
    		Mobile     string    `db:"mobile"`   // 手机号
    		Gender     string    `db:"gender"`   // 男｜女｜未公开
    		Nickname   string    `db:"nickname"` // 用户昵称
    		CreateTime time.Time `db:"create_time"`
    		UpdateTime time.Time `db:"update_time"`
    	}
    )

    func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *UserModel {
    	return &UserModel{
    		CachedConn: sqlc.NewConn(conn, c),
    		table:      table,
    	}
    }

    func (m *UserModel) Insert(data User) (sql.Result, error) {
    	query := `insert into ` + m.table + `(` + userRowsExpectAutoSet + `) value (?, ?, ?, ?, ?)`
    	return m.ExecNoCache(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
    }

    func (m *UserModel) FindOne(id int64) (*User, error) {
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
    	var resp User
    	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
    		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
    		return conn.QueryRow(v, query, id)
    	})
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }

    func (m *UserModel) FindOneByName(name string) (*User, error) {
    	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
    	var resp User
    	err := m.QueryRowIndex(&resp, userNameKey, func(primary interface{}) string {
    		return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
    	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
    		query := `select ` + userRows + ` from ` + m.table + ` where name = ? limit 1`
    		if err := conn.QueryRow(&resp, query, name); err != nil {
    			return nil, err
    		}
    		return resp.Id, nil
    	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
    		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
    		return conn.QueryRow(v, query, primary)
    	})
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }

    func (m *UserModel) FindOneByMobile(mobile string) (*User, error) {
    	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
    	var resp User
    	err := m.QueryRowIndex(&resp, userMobileKey, func(primary interface{}) string {
    		return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
    	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
    		query := `select ` + userRows + ` from ` + m.table + ` where mobile = ? limit 1`
    		if err := conn.QueryRow(&resp, query, mobile); err != nil {
    			return nil, err
    		}
    		return resp.Id, nil
    	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
    		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
    		return conn.QueryRow(v, query, primary)
    	})
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }

    func (m *UserModel) Update(data User) error {
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
    	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
    		query := `update ` + m.table + ` set ` + userRowsWithPlaceHolder + ` where id = ?`
    		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
    	}, userIdKey)
    	return err
    }

    func (m *UserModel) Delete(id int64) error {
    	data, err := m.FindOne(id)
    	if err != nil {
    		return err
    	}
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
    	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
    	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
    	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
    		query := `delete from ` + m.table + ` where id = ?`
    		return conn.Exec(query, id)
    	}, userIdKey, userNameKey, userMobileKey)
    	return err
    }
  ```

### 用法

```Plain Text
goctl model mysql -h
```

```Plain Text
NAME:
   goctl model mysql - generate mysql model"

USAGE:
   goctl model mysql command [command options] [arguments...]

COMMANDS:
   ddl         generate mysql model from ddl"
   datasource  generate model from datasource"

OPTIONS:
   --help, -h  show help
```

## 生成规则

* 默认规则
  
  我们默认用户在建表时会创建createTime、updateTime字段(忽略大小写、下划线命名风格)且默认值均为`CURRENT_TIMESTAMP`，而updateTime支持`ON UPDATE CURRENT_TIMESTAMP`，对于这两个字段生成`insert`、`update`时会被移除，不在赋值范畴内，当然，如果你不需要这两个字段那也无大碍。
* 带缓存模式
  * ddl

	```shell script
	goctl model mysql -src={filename} -dir={dir} -cache=true
	```

  * datasource

	```shell script
	goctl model mysql datasource -url={datasource} -table={tables}  -dir={dir} -cache=true
	```

  目前仅支持redis缓存，如果选择带缓存模式，即生成的`FindOne(ByXxx)`&`Delete`代码会生成带缓存逻辑的代码，目前仅支持单索引字段（除全文索引外），对于联合索引我们默认认为不需要带缓存，且不属于通用型代码，因此没有放在代码生成行列，如example中user表中的`id`、`name`、`mobile`字段均属于单字段索引。

* 不带缓存模式
  
  * ddl
  
      ```shell script
      goctl model -src={filename} -dir={dir}
      ```

  * datasource
  
      ```shell script
          goctl model mysql datasource -url={datasource} -table={tables}  -dir={dir}
      ```

  or
  * ddl

      ```shell script
        goctl model -src={filename} -dir={dir} -cache=false
      ```

  * datasource
  
      ```shell script
        goctl model mysql datasource -url={datasource} -table={tables}  -dir={dir} -cache=false
      ```
  
  生成代码仅基本的CURD结构。

## 缓存

  对于缓存这一块我选择用一问一答的形式进行罗列。我想这样能够更清晰的描述model中缓存的功能。

* 缓存会缓存哪些信息？

  对于主键字段缓存，会缓存整个结构体信息，而对于单索引字段（除全文索引）则缓存主键字段值。

* 数据有更新（`update`）操作会清空缓存吗？
  
  会，但仅清空主键缓存的信息，why？这里就不做详细赘述了。

* 为什么不按照单索引字段生成`updateByXxx`和`deleteByXxx`的代码？
  
  理论上是没任何问题，但是我们认为，对于model层的数据操作均是以整个结构体为单位，包括查询，我不建议只查询某部分字段（不反对），否则我们的缓存就没有意义了。

* 为什么不支持`findPageLimit`、`findAll`这么模式代码生层？
  
  目前，我认为除了基本的CURD外，其他的代码均属于<i>业务型</i>代码，这个我觉得开发人员根据业务需要进行编写更好。

## QA

* goctl model除了命令行模式，支持插件模式吗？

  很快支持idea插件。
