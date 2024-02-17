# Goctl Model

goctl model 为go-zero下的工具模块中的组件之一，目前支持识别mysql ddl进行model层代码生成，通过命令行或者idea插件（即将支持）可以有选择地生成带redis cache或者不带redis cache的代码逻辑。

## 快速开始

* 通过ddl生成

    ```shell script
    goctl model mysql ddl -src="./*.sql" -dir="./sql/model" -c
    ```

    执行上述命令后即可快速生成CURD代码。

    ```Plain Text
    model
    │   ├── error.go
    │   └── usermodel.go
    ```

* 通过datasource生成

    ```shell script
    goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="*"  -dir="./model"
    ```

* 生成代码示例
  
	```go

	package model
    
    import (
    	"database/sql"
    	"fmt"
    	"strings"
    	"time"
        
        "github.com/zeromicro/go-zero/core/stores/builder"
    	"github.com/zeromicro/go-zero/core/stores/cache"
    	"github.com/zeromicro/go-zero/core/stores/sqlc"
    	"github.com/zeromicro/go-zero/core/stores/sqlx"
    	"github.com/zeromicro/go-zero/core/stringx"
    )
    
    var (
    	userFieldNames          = builder.FieldNames(&User{})
    	userRows                = strings.Join(userFieldNames, ",")
    	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), ",")
    	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
    
    	cacheUserPrefix       = "cache:User:user:"
    	cacheUserNamePrefix   = "cache:User:name:"
    	cacheUserMobilePrefix = "cache:User:mobile:"
    	cacheUserIdPrefix     = "cache:User:id:"
    )
    
    type (
    	UserModel interface {
    		Insert(data User) (sql.Result, error)
    		FindOne(id int64) (*User, error)
    		FindOneByUser(user string) (*User, error)
    		FindOneByName(name string) (*User, error)
    		FindOneByMobile(mobile string) (*User, error)
    		Update(data User) error
    		Delete(id int64) error
    	}
    
    	defaultUserModel struct {
    		sqlc.CachedConn
    		table string
    	}
    
    	User struct {
    		Id         int64     `db:"id"`
    		User       string    `db:"user"`     // 用户
    		Name       string    `db:"name"`     // 用户名称
    		Password   string    `db:"password"` // 用户密码
    		Mobile     string    `db:"mobile"`   // 手机号
    		Gender     string    `db:"gender"`   // 男｜女｜未公开
    		Nickname   sql.NullString    `db:"nickname"` // 用户昵称
    		CreateTime sql.NullTime `db:"create_time"`
    		UpdateTime time.Time `db:"update_time"`
    	}
    )
    
    func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
    	return &defaultUserModel{
    		CachedConn: sqlc.NewConn(conn, c),
    		table:      "user",
    	}
    }
    
    func (m *defaultUserModel) Insert(data User) (sql.Result, error) {
    	userKey := fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
    	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
    	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
    	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
    		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
    		return conn.Exec(query, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
    	}, userMobileKey, userKey, userNameKey)
    	return ret, err
    }
    
    func (m *defaultUserModel) FindOne(id int64) (*User, error) {
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
    	var resp User
    	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v any) error {
    		query := fmt.Sprintf("select %s from %s where id = ? limit 1", userRows, m.table)
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
    
    func (m *defaultUserModel) FindOneByUser(user string) (*User, error) {
    	userKey := fmt.Sprintf("%s%v", cacheUserPrefix, user)
    	var resp User
    	err := m.QueryRowIndex(&resp, userKey, m.formatPrimary, func(conn sqlx.SqlConn, v any) (i any, e error) {
    		query := fmt.Sprintf("select %s from %s where user = ? limit 1", userRows, m.table)
    		if err := conn.QueryRow(&resp, query, user); err != nil {
    			return nil, err
    		}
    		return resp.Id, nil
    	}, m.queryPrimary)
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }
    
    func (m *defaultUserModel) FindOneByName(name string) (*User, error) {
    	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
    	var resp User
    	err := m.QueryRowIndex(&resp, userNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v any) (i any, e error) {
    		query := fmt.Sprintf("select %s from %s where name = ? limit 1", userRows, m.table)
    		if err := conn.QueryRow(&resp, query, name); err != nil {
    			return nil, err
    		}
    		return resp.Id, nil
    	}, m.queryPrimary)
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }
    
    func (m *defaultUserModel) FindOneByMobile(mobile string) (*User, error) {
    	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
    	var resp User
    	err := m.QueryRowIndex(&resp, userMobileKey, m.formatPrimary, func(conn sqlx.SqlConn, v any) (i any, e error) {
    		query := fmt.Sprintf("select %s from %s where mobile = ? limit 1", userRows, m.table)
    		if err := conn.QueryRow(&resp, query, mobile); err != nil {
    			return nil, err
    		}
    		return resp.Id, nil
    	}, m.queryPrimary)
    	switch err {
    	case nil:
    		return &resp, nil
    	case sqlc.ErrNotFound:
    		return nil, ErrNotFound
    	default:
    		return nil, err
    	}
    }
    
    func (m *defaultUserModel) Update(data User) error {
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
    	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
    		query := fmt.Sprintf("update %s set %s where id = ?", m.table, userRowsWithPlaceHolder)
    		return conn.Exec(query, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
    	}, userIdKey)
    	return err
    }
    
    func (m *defaultUserModel) Delete(id int64) error {
    	data, err := m.FindOne(id)
    	if err != nil {
    		return err
    	}
    
    	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
    	userKey := fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
    	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
    	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
    	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
    		query := fmt.Sprintf("delete from %s where id = ?", m.table)
    		return conn.Exec(query, id)
    	}, userIdKey, userKey, userNameKey, userMobileKey)
    	return err
    }
    
    func (m *defaultUserModel) formatPrimary(primary any) string {
    	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
    }
    
    func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary any) error {
    	query := fmt.Sprintf("select %s from %s where id = ? limit 1", userRows, m.table)
    	return conn.QueryRow(v, query, primary)
    }

	```

## 用法

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
	goctl model mysql -src={patterns} -dir={dir} -cache
	```

	help

	```
	NAME:
       goctl model mysql ddl - generate mysql model from ddl
    
    USAGE:
       goctl model mysql ddl [command options] [arguments...]
    
    OPTIONS:
       --src value, -s value  the path or path globbing patterns of the ddl
       --dir value, -d value  the target dir
       --style value          the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]
       --cache, -c            generate code with cache [optional]
       --idea                 for idea plugin [optional]
       --database, -db        the name of database [optional]
	```

  * datasource

	```shell script
	goctl model mysql datasource -url={datasource} -table={patterns}  -dir={dir} -cache=true
	```

	help

	```
	NAME:
       goctl model mysql datasource - generate model from datasource
    
    USAGE:
       goctl model mysql datasource [command options] [arguments...]
    
    OPTIONS:
       --url value              the data source of database,like "root:password@tcp(127.0.0.1:3306)/database
       --table value, -t value  the table or table globbing patterns in the database
       --cache, -c              generate code with cache [optional]
       --dir value, -d value    the target dir
       --style value            the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]
       --idea                   for idea plugin [optional]


	```

	示例用法请参考[用法](./example/generator.sh)
  
	> NOTE: goctl model mysql ddl/datasource 均新增了一个`--style`参数，用于标记文件命名风格。

  目前仅支持redis缓存，如果选择带缓存模式，即生成的`FindOne(ByXxx)`&`Delete`代码会生成带缓存逻辑的代码，目前仅支持单索引字段（除全文索引外），对于联合索引我们默认认为不需要带缓存，且不属于通用型代码，因此没有放在代码生成行列，如example中user表中的`id`、`name`、`mobile`字段均属于单字段索引。

* 不带缓存模式

  * ddl
  
      ```shell script
        goctl model -src={patterns} -dir={dir}
      ```

  * datasource
  
      ```shell script
        goctl model mysql datasource -url={datasource} -table={patterns}  -dir={dir}
      ```

  or
  * ddl

      ```shell script
        goctl model -src={patterns} -dir={dir}
      ```

  * datasource
  
      ```shell script
        goctl model mysql datasource -url={datasource} -table={patterns}  -dir={dir}
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

* 为什么不支持`findPageLimit`、`findAll`这种模型代码生成？
  
  目前，我认为除了基本的CURD外，其他的代码均属于<i>业务型</i>代码，这个我觉得开发人员根据业务需要进行编写更好。

# 类型转换规则

| mysql dataType | golang dataType | golang Null dataType |
|----------------|-----------------|----------------------------------------|
| bool           | int64           | sql.NullInt64                          |
| boolean        | int64           | sql.NullInt64                          |
| tinyint        | int64           | sql.NullInt64                          |
| smallint       | int64           | sql.NullInt64                          |
| mediumint      | int64           | sql.NullInt64                          |
| int            | int64           | sql.NullInt64                          |
| integer        | int64           | sql.NullInt64                          |
| bigint         | int64           | sql.NullInt64                          |
| float          | float64         | sql.NullFloat64                        |
| double         | float64         | sql.NullFloat64                        |
| decimal        | float64         | sql.NullFloat64                        |
| date           | time.Time       | sql.NullTime                           |
| datetime       | time.Time       | sql.NullTime                           |
| timestamp      | time.Time       | sql.NullTime                           |
| time           | string          | sql.NullString                         |
| year           | time.Time       | sql.NullInt64                          |
| char           | string          | sql.NullString                         |
| varchar        | string          | sql.NullString                         |
| binary         | string          | sql.NullString                         |
| varbinary      | string          | sql.NullString                         |
| tinytext       | string          | sql.NullString                         |
| text           | string          | sql.NullString                         |
| mediumtext     | string          | sql.NullString                         |
| longtext       | string          | sql.NullString                         |
| enum           | string          | sql.NullString                         |
| set            | string          | sql.NullString                         |
| json           | string          | sql.NullString                         |

