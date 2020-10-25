package gen

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen/testmodel"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen/testnocachemodel"
)

var (
	source = "CREATE TABLE `test` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `nanosecond` bigint NOT NULL DEFAULT '0',\n  `data` varchar(255) DEFAULT '',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `nanosecond_unique` (`nanosecond`)\n) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
)

func TestCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()
	g := NewDefaultGenerator(source, "./testmodel")
	err := g.Start(true)
	assert.Nil(t, err)

	username := os.Getenv("mysql_user")
	password := os.Getenv("mysql_pass")
	if len(username) == 0 || len(password) == 0 {
		fmt.Println("empty user name or password")
		return
	}
	// mysql conn
	url := username + ":" + password + "@tcp(127.0.0.1:3306)/gozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	_, err = mysql.ParseDSN(url)
	if err != nil {
		t.Error(err)
		return
	}
	redisHost := "127.0.0.1:6379"
	r := redis.NewRedis(redisHost, "node", "")
	ok := r.Ping()
	if !ok {
		fmt.Println("redis connect error")
		return
	}
	fmt.Println("ok")
	conn := sqlx.NewMysql(url)
	m := testmodel.NewTestModel(conn, cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: redisHost,
				Type: "node",
				Pass: "",
			},
			Weight: 100,
		},
	})
	data := "test.data"
	nanoseconds := int64(time.Now().Nanosecond())
	// insert
	ret, err := m.Insert(testmodel.Test{
		Nanosecond: nanoseconds,
		Data:       data,
	})

	assert.Nil(t, err)
	insertId, err := ret.LastInsertId()
	assert.Nil(t, err)

	// findOne
	info, err := m.FindOne(insertId)
	assert.Nil(t, err)
	assert.Equal(t, nanoseconds, info.Nanosecond)
	assert.Equal(t, data, info.Data)

	// fineByMobile
	info, err = m.FindOneByNanosecond(nanoseconds)
	assert.Nil(t, err)
	assert.Equal(t, nanoseconds, info.Nanosecond)
	assert.Equal(t, data, info.Data)

	// update
	updateData := "update.data"
	info.Data = updateData
	_, err = m.Update(*info)
	assert.Nil(t, err)

	// update result check
	updateInfo, err := m.FindOne(info.Id)
	assert.Nil(t, err)
	assert.Equal(t, info.Data, updateInfo.Data)

	// delete
	err = m.Delete(info.Id)
	assert.Nil(t, err)

	// delete result check
	_, err = m.FindOne(info.Id)
	assert.Equal(t, testmodel.ErrNotFound, err)
}

func TestNoCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()
	g := NewDefaultGenerator(source, "./testnocachemodel")
	err := g.Start(false)
	assert.Nil(t, err)

	username := os.Getenv("mysql_user")
	password := os.Getenv("mysql_pass")
	if len(username) == 0 || len(password) == 0 {
		fmt.Println("empty user name or password")
		return
	}
	url := username + ":" + password + "@tcp(127.0.0.1:3306)/gozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	_, err = mysql.ParseDSN(url)
	if err != nil {
		t.Error(err)
		return
	}
	conn := sqlx.NewMysql(url)
	m := testnocachemodel.NewTestModel(conn)
	data := "test.data"
	nanoseconds := int64(time.Now().Nanosecond())
	// insert
	ret, err := m.Insert(testnocachemodel.Test{
		Nanosecond: nanoseconds,
		Data:       data,
	})

	assert.Nil(t, err)
	insertId, err := ret.LastInsertId()
	assert.Nil(t, err)

	// findOne
	info, err := m.FindOne(insertId)
	assert.Nil(t, err)
	assert.Equal(t, nanoseconds, info.Nanosecond)
	assert.Equal(t, data, info.Data)

	// fineByMobile
	info, err = m.FindOneByNanosecond(nanoseconds)
	assert.Nil(t, err)
	assert.Equal(t, nanoseconds, info.Nanosecond)
	assert.Equal(t, data, info.Data)

	// update
	updateData := "update.data"
	info.Data = updateData
	_, err = m.Update(*info)
	assert.Nil(t, err)

	// update result check
	updateInfo, err := m.FindOne(info.Id)
	assert.Nil(t, err)
	assert.Equal(t, info.Data, updateInfo.Data)

	// delete
	err = m.Delete(info.Id)
	assert.Nil(t, err)

	// delete result check
	_, err = m.FindOne(info.Id)
	assert.Equal(t, testnocachemodel.ErrNotFound, err)
}
