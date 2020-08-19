package example

// todo: In order to pass the go test on github,
// todo: the code is commented. If you need go test,
// todo: modify the configuration file and delete the comment and execute it.

//import (
//	"flag"
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/tal-tech/go-zero/core/conf"
//	"github.com/tal-tech/go-zero/core/jsonx"
//	"github.com/tal-tech/go-zero/core/logx"
//	"github.com/tal-tech/go-zero/core/stores/redis"
//	"github.com/tal-tech/go-zero/core/stores/sqlx"
//	"github.com/tal-tech/go-zero/tools/goctl/model/sql/example/config"
//	"github.com/tal-tech/go-zero/tools/goctl/model/sql/example/model"
//)
//
//var (
//	userModel   *model.UserModel
//	courseModel *model.UserCourseModel
//	// verify redis cache
//	r *redis.Redis
//)
//var configFile = flag.String("f", "etc/config.json", "the config file")
//
//func TestMain(m *testing.M) {
//	flag.Parse()
//	var c config.Config
//	conf.MustLoad(*configFile, &c)
//
//	logx.MustSetup(c.LogConf)
//	conn := sqlx.NewMysql(c.Mysql.DataSource)
//	if conn == nil {
//		return
//	}
//	userModel = model.NewUserModel(conn, c.CacheRedis, c.Mysql.Table.User)
//	courseModel = model.NewUserCourseModel(conn, c.Mysql.Table.Course)
//	r = redis.NewRedis(c.Redis.Host, c.Redis.Type, c.Redis.Pass)
//	if userModel == nil || courseModel == nil || r == nil {
//		return
//	}
//	m.Run()
//}
//
////  cache model
//func TestUser(t *testing.T) {
//	var user model.User
//	user.Name = "test"
//	user.Password = "123456"
//	user.Mobile = "136****0001"
//	user.Gender = "男"
//	user.Nickname = "Keson"
//	insert, err := userModel.Insert(user)
//	assert.Nil(t, err)
//	id, err := insert.LastInsertId()
//	assert.Nil(t, err)
//
//	// select
//	ret, err := userModel.FindOneByName("test")
//	assert.Nil(t, err)
//	assert.Equal(t, user.Mobile, ret.Mobile)
//
//	// should cache
//	// expected primary key
//	var redisId int64
//	err = get(&redisId, "cache#User#name#test")
//	assert.Nil(t, err)
//	assert.Equal(t, ret.Id, redisId)
//
//	//expected user from cache
//	var redisData model.User
//	err = get(&redisData, fmt.Sprintf("cache#User#id#%v", redisId))
//	assert.Nil(t, err)
//	assert.Equal(t, redisData.Nickname, user.Nickname)
//
//	// update
//	ret.Nickname = "Keson after"
//	err = userModel.Update(*ret)
//	assert.Nil(t, err)
//	// expected cache delete
//	exist, err := r.Exists(fmt.Sprintf("cache#User#id#%v", ret.Id))
//	assert.Nil(t, err)
//	assert.Equal(t, false, exist)
//
//	// select
//	ret, err = userModel.FindOne(id)
//	assert.Nil(t, err)
//	assert.Equal(t, "Keson after", ret.Nickname)
//
//	// delete
//	err = userModel.Delete(ret.Id)
//	assert.Nil(t, err)
//	exist, err = r.Exists(fmt.Sprintf("cache#User#id#%v", ret.Id))
//	assert.Nil(t, err)
//	assert.Equal(t, false, exist)
//
//	// verify
//	_, err = userModel.FindOne(id)
//	assert.Equal(t, model.ErrNotFound, err)
//
//}
//
//// no cache model
//func TestCourse(t *testing.T) {
//	// new user
//	insert, err := userModel.Insert(model.User{
//		Name:     "courseUser",
//		Password: "123456",
//		Mobile:   "136****1111",
//		Gender:   "男",
//		Nickname: "Keson",
//	})
//	assert.Nil(t, err)
//	id, err := insert.LastInsertId()
//	assert.Nil(t, err)
//
//	var course model.UserCourse
//	course.Id = id
//	course.CourseName = "Java"
//
//	courseInsert, err := courseModel.Insert(course)
//	assert.Nil(t, err)
//	courseId, err := courseInsert.LastInsertId()
//	assert.Nil(t, err)
//
//	// select
//	ret, err := courseModel.FindOne(courseId)
//	assert.Nil(t, err)
//	assert.Equal(t, course.CourseName, ret.CourseName)
//
//	// update
//	ret.CourseName = "Golang"
//	err = courseModel.Update(*ret)
//	assert.Nil(t, err)
//
//	// select
//	ret, err = courseModel.FindOne(courseId)
//	assert.Nil(t, err)
//	assert.Equal(t, "Golang", ret.CourseName)
//
//	// delete
//	err = courseModel.Delete(ret.Id)
//	assert.Nil(t, err)
//
//	// verify
//	_, err = courseModel.FindOne(courseId)
//	assert.Equal(t, model.ErrNotFound, err)
//
//	// clean user
//	err = userModel.Delete(id)
//	assert.Nil(t, err)
//}
//
//func get(data interface{}, key string) error {
//	v, err := r.Get(key)
//	if err != nil {
//		return err
//	}
//	return jsonx.Unmarshal([]byte(v), data)
//}
