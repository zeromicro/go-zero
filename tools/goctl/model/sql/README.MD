<div style="text-align: center;"><h1>Sql生成工具说明文档</h1></div>

<h2>前言</h2>
在当前Sql代码生成工具是基于sqlc生成的逻辑。

<h2>关键字</h2>

+ 查询类型(前暂不支持同一字段多种类型混合生成，如按照campus_id查询单结果又查询All或者Limit)
  - 单结果查询
    - FindOne(主键特有)
    - FindOneByXxx
  - 多结果查询
    - FindAllByXxx
    - FindLimitByXxx
- withCache
- withoutCache

<h2>准备工作</h2>

- table

    ```
    CREATE TABLE `user_info` (
    `id` bigint(20) NOT NULL COMMENT '主键',
    `campus_id` bigint(20) DEFAULT NULL COMMENT '整校id',
    `name` varchar(255) DEFAULT NULL COMMENT '用户姓名',
    `id_number` varchar(255) DEFAULT NULL COMMENT '身份证',
    `age` int(10) DEFAULT NULL COMMENT '年龄',
    `gender` tinyint(1) DEFAULT NULL COMMENT '性别，0-男，1-女，2-不限',
    `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
    ```

<h2>imports生成</h2>
imports代码生成对应model中包的引入管理，仅使用于晓黑板项目中（非相对路径动态生成），目前受`withCache`参数的影响，除此之外其实为固定代码。

- withCache
  
    ```
    import (
        "database/sql""fmt"
        "strings"
        "time"

        "github.com/tal-tech/go-zero/core/stores/sqlc"
        "github.com/tal-tech/go-zero/core/stores/sqlx"
        "github.com/tal-tech/go-zero/core/stringx"
        "xiao/service/shared/builderx"
    )
    ```

- withoutCache
  
    ```
    import (
        "database/sql""fmt"
        "strings"
        "time"

        "github.com/tal-tech/go-zero/core/stores/sqlx"
        "github.com/tal-tech/go-zero/core/stringx"
        "xiao/service/shared/builderx"
    )
    ```

<h2>vars生成</h2>

vars部分对应model中var声明的包含的代码块,由`table`名和`withCache`来决定其中的代码生成内容，`withCache`决定是否要生成缓存key变量的声明。

- withCache

    ```
    var (
        UserInfoFieldNames          = builderx.FieldNames(&UserInfo{})
        UserInfoRows                = strings.Join(UserInfoFieldNames, ",")
        UserInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(UserInfoFieldNames, "id", "create_time", "update_time"), ",")
        UserInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(UserInfoFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
    
        cacheUserInfoIdPrefix = "cache#userInfo#id#"
        cacheUserInfoCampusIdPrefix = "cache#userInfo#campusId#"
        cacheUserInfoNamePrefix = "cache#userInfo#name#"
        cacheUserInfoMobilePrefix = "cache#userInfo#mobile#"
    )
    ```

- withoutCache

    ```
    var (
        UserInfoFieldNames          = builderx.FieldNames(&UserInfo{})
        UserInfoRows                = strings.Join(UserInfoFieldNames, ",")
        UserInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(UserInfoFieldNames, "id", "create_time", "update_time"), ",")
        UserInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(UserInfoFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
    )
    ```

<h2>types生成</h2>

ypes部分对应model中type声明的包含的代码块,由`table`名和`withCache`来决定其中的代码生成内容，`withCache`决定引入sqlc还是sqlx。

- withCache
    ```
    type (
        UserInfoModel struct {
            conn  sqlc.CachedConn
            table string
        }

        UserInfo struct {
            Id int64 `db:"id"` // 主键id
            CampusId int64 `db:"campus_id"` // 整校id
            Name string `db:"name"` // 用户姓名
            IdNumber string `db:"id_number"` // 身份证
            Age int64 `db:"age"` // 年龄
            Gender int64 `db:"gender"` // 性别，0-男，1-女，2-不限
            Mobile string `db:"mobile"` // 手机号
            CreateTime time.Time `db:"create_time"` // 创建时间
            UpdateTime time.Time `db:"update_time"` // 更新时间
        }
    )    
    ```

- withoutCache
    ```
    type (
        UserInfoModel struct {
            conn  sqlx.SqlConn
            table string
        }

        UserInfo struct {
            Id int64 `db:"id"` // 主键id
            CampusId int64 `db:"campus_id"` // 整校id
            Name string `db:"name"` // 用户姓名
            IdNumber string `db:"id_number"` // 身份证
            Age int64 `db:"age"` // 年龄
            Gender int64 `db:"gender"` // 性别，0-男，1-女，2-不限
            Mobile string `db:"mobile"` // 手机号
            CreateTime time.Time `db:"create_time"` // 创建时间
            UpdateTime time.Time `db:"update_time"` // 更新时间
        }
    )
    ```
<h2>New生成</h2>
new生成对应model中struct的New函数，受`withCache`影响决定是否要引入cacheRedis

- withCache
    ```
    func NewUserInfoModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *UserInfoModel {
        return &UserInfoModel{
            CachedConn: sqlc.NewConn(conn, c),
            table:      table,
        }
    }
    ```
- withoutCache
    ```
    func NewUserInfoModel(conn sqlx.SqlConn, table string) *UserInfoModel {
	    return &UserInfoModel{conn: conn, table: table}
    }
    ```


<h2>FindOne查询生成</h2>
FindOne查询代码生成仅对主键有效。如`user_info`中生成的FindOne如下：

- withCache

    ```
    func (m *UserInfoModel) FindOne(id int64) (*UserInfo, error) {
        idKey := fmt.Sprintf("%s%v", cacheUserInfoIdPrefix, id)
        var resp UserInfo
        err := m.QueryRow(&resp, idKey, func(conn sqlx.SqlConn, v interface{}) error {
            query := `select ` + userInfoRows + ` from ` + m.table  + `where id = ? limit 1`
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
    ```

- withoutCache

    ```
    func (m *UserInfoModel) FindOne(id int64) (*UserInfo, error) {
        
        query := `select ` + userInfoRows + ` from ` + m.table  + `where id = ? limit 1`
        var resp UserInfo
        err := m.conn.QueryRow(&resp, query, id)
        switch err {
        case nil:
            return &resp, nil
        case sqlx.ErrNotFound:
            return nil, ErrNotFound
        default:
            return nil, err
        
    }
    ```

<h2>FindOneByXxx查询生成</h2>

FindOneByXxx查询生成可以按照单个字段查询、多个字段以AND关系且表达式符号为`=`的查询（下称：组合查询）,对除主键之外的字段有效，对于单个字段可以用`withCache`来控制是否需要缓存，这里的缓存只缓存主键，并不缓存整个struct，注意：这里有一个隐藏的规则，如果单个字段查询需要cache，那么主键一定有cache；多个字段组成的`组合查询`一律没有缓存处理,<strong><i>且组合查询不能相互嵌套</i></strong>，否则会报`circle query with other fields`错误，下面我们按场景来依次查看对应代码生成后的示例。

>注：目前暂不支持除equals之外的条件查询。

+ 单字段查询  
    以name查询为例
    - withCache
        ```
        func (m *UserInfoModel) FindOneByName(name string) (*UserInfo, error) {
            nameKey := fmt.Sprintf("%s%v", cacheUserInfoNamePrefix, name)
            var id string
            err := m.GetCache(key, &id)
            if err != nil {
                return nil, err
            }
            if id != "" {
                return m.FindOne(id)
            }
            var resp UserInfo
            query := `select ` + userInfoRows + ` from ` + m.table  + `where name = ? limit 1`
            err = m.QueryRowNoCache(&resp, query, name)
            switch err {
            case nil:
                err = m.SetCache(nameKey, resp.Id)
                if err != nil {
                    logx.Error(err)
                }
                return &resp, nil
            case sqlc.ErrNotFound:
                return nil, ErrNotFound
            default:
                return nil, err
            }
        }    
        ```
    - withoutCache

        ```
        func (m *UserInfoModel) FindOneByName(name string) (*UserInfo, error) {
            var resp UserInfo
            query := `select ` + userInfoRows + ` from ` + m.table  + `where name = ? limit 1`
            err = m.conn.QueryRow(&resp, query, name)
            switch err {
            case nil:
                return &resp, nil
            case sqlx.ErrNotFound:
                return nil, ErrNotFound
            default:
                return nil, err
            }
        }
        ```

- 组合查询  
    以`campus_id`和`id_number`查询为例。
    
        ```
        func (m *UserInfoModel) FindOneByCampusIdAndIdNumber(campusId int64,idNumber string) (*UserInfo, error) {
            var resp UserInfo
            query := `select ` + userInfoRows + ` from ` + m.table  + `where campus_id = ? AND id_number = ? limit 1`
            err = m.QueryRowNoCache(&resp, query, campusId, idNumber)
            // err = m.conn.QueryRows(&resp, query, campusId, idNumber)
            switch err {
            case nil:
                return &resp, nil
            case sqlx.ErrNotFound:
                return nil, ErrNotFound
            default:
                return nil, err
            }
        }
        ```
<h2>FindAllByXxx生成</h2>
FindAllByXxx查询和FindOneByXxx功能相似，只是FindOneByXxx限制了limit等于1，而FindAllByXxx是查询所有,以两个例子来说明

- 查询单个字段`name`等于某值的所有数据
    ```
    func (m *UserInfoModel) FindAllByName(name string) ([]*UserInfo, error) {
        var resp []*UserInfo
        query := `select ` + userInfoRows + ` from ` + m.table  + `where name = ?`
        err := m.QueryRowsNoCache(&resp, query, name)
        // err := m.conn.QueryRows(&resp, query, name)
        if err != nil {
            return nil, err
        }
        return resp, nil
    }    
    ```
- 查询多个组合字段`campus_id`等于某值且`gender`等于某值的所有数据
    ```
    func (m *UserInfoModel) FindAllByCampusIdAndGender(campusId int64,gender int64) ([]*UserInfo, error) {
        var resp []*UserInfo
        query := `select ` + userInfoRows + ` from ` + m.table  + `where campus_id = ? AND gender = ?`
        err := m.QueryRowsNoCache(&resp, query, campusId, gender)
        // err := m.conn.QueryRows(&resp, query, campusId, gender)
        if err != nil {
            return nil, err
        }
        return resp, nil
    }    
    ```

<h2>FindLimitByXxx生成</h2>
FindLimitByXxx查询和FindAllByXxx功能相似，只是FindAllByXxx限制了limit，除此之外还会生成查询对应Count总数的代码，而FindAllByXxx是查询所有数据,以几个例子来说明

- 查询`gender`等于某值的分页数据，按照`create_time`降序
    ```
    func (m *UserInfoModel) FindLimitByGender(gender int64, page, limit int) ([]*UserInfo, error) {
        var resp []*UserInfo
        query := `select ` + userInfoRows + `from ` + m.table  + `where gender = ? order by create_time DESC limit ?,?`
        err := m.QueryRowsNoCache(&resp, query, gender, (page-1)*limit, limit)
        // err := m.conn.QueryRows(&resp, query, gender, (page-1)*limit, limit)
        if err != nil {
            return nil, err
        }
        return resp, nil
    }

    func (m *UserInfoModel) FindAllCountByGender(gender int64) (int64, error) {
        var count int64
        query := `select count(1)  from ` + m.table  + `where gender = ? `
        err := m.QueryRowsNoCache(&count, query, gender)
        // err := m.conn.QueryRow(&count, query, gender)
        if err != nil {
            return 0, err
        }
        return count, nil
    }       
    ```
- 查询`gender`等于某值的分页数据，按照`create_time`降序、`update_time`生序排序
    ```
    func (m *UserInfoModel) FindLimitByGender(gender int64, page, limit int) ([]*UserInfo, error) {
        var resp []*UserInfo
        query := `select ` + userInfoRows + `from ` + m.table  + `where gender = ? order by create_time DESC,update_time ASC limit ?,?`
        err := m.QueryRowsNoCache(&resp, query, gender, (page-1)*limit, limit)
       // err := m.conn.QueryRows(&resp, query, gender, (page-1)*limit, limit)
        if err != nil {
            return nil, err
        }
        return resp, nil
    }

    func (m *UserInfoModel) FindAllCountByGender(gender int64) (int64, error) {
        var count int64
        query := `select count(1)  from ` + m.table  + `where gender = ? `
        err := m.QueryRowNoCache(&count, query, gender)
        // err := m.conn.QueryRow(&count, query, gender)
        if err != nil {
            return 0, err
        }
        return count, nil
    }          
    ```
- 查询`gender`等于某值且`campus_id`为某值按照`create_time`降序的分页数据
    ```
    func (m *UserInfoModel) FindLimitByGenderAndCampusId(gender int64,campusId int64, page, limit int) ([]*UserInfo, error) {
        var resp []*UserInfo
        query := `select ` + userInfoRows + `from ` + m.table  + `where gender = ? AND campus_id = ? order by create_time DESC limit ?,?`
        err := m.QueryRowsNoCache(&resp, query, gender, campusId, (page-1)*limit, limit)
        // err := m.conn.QueryRows(&resp, query, gender, campusId, (page-1)*limit, limit)
        if err != nil {
            return nil, err
        }
        return resp, nil
    }

    func (m *UserInfoModel) FindAllCountByGenderAndCampusId(gender int64,campusId int64) (int64, error) {
        var count int64
        query := `select count(1)  from ` + m.table  + `where gender = ? AND campus_id = ? `
        err := m.QueryRowsNoCache(&count, query, gender, campusId)
        // err := m.conn.QueryRow(&count, query, gender, campusId)
        if err != nil {
            return 0, err
        }
        return count, nil
    }      
    ```

<h2>Delete生成</h2>
Delete代码根据`withCache`的不同可以生成带缓存逻辑代码和不带缓存逻辑代码，<strong><i>Delete代码生成仅按照主键删除</i></strong>。从FindOneByXxx方法描述得知，非主键`withCache`了那么主键会强制被cache，因此在delete时也会删除主键cache。

- withCache  
  根据`mobile`查询用户信息

    ```
    func (m *UserInfoModel) Delete(userId int64) error {
        userIdKey := fmt.Sprintf("%s%v", cacheUserInfoUserIdPrefix, userId)
        mobileKey := fmt.Sprintf("%s%v", cacheUserInfoMobilePrefix, mobile)
        _, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
            query := `delete from ` + m.table +  + `where user_id = ?`
            return conn.Exec(query, userId)
        }, userIdKey, mobileKey)
        return err
        }
    ```
- withoutCache
    ```
    func (m *UserInfoModel) Delete(userId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table +  + `where user_id = ?`
		return conn.Exec(query, userId)
	}, )
	return err
}
    ```
<h2>Insert生成</h2>

<h2>Update生成</h2>

<h2>待完善（TODO）</h2>

- 同一字段多种查询方式代码生成(优先级较高)
- 条件查询
- 范围查询
- ...

<h2>反馈与建议</h2>

- 无