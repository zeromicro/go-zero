package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	mocksql "github.com/zeromicro/go-zero/tools/goctl/model/sql/test"
)

func TestStudentModel(t *testing.T) {
	var (
		testTimeValue          = time.Now()
		testTable              = "`student`"
		testUpdateName         = "gozero1"
		testRowsAffected int64 = 1
		testInsertId     int64 = 1
		class                  = "一年级1班"
	)

	var data Student
	data.Id = testInsertId
	data.Class = class
	data.Name = "gozero"
	data.Age = sql.NullInt64{
		Int64: 1,
		Valid: true,
	}
	data.Score = sql.NullFloat64{
		Float64: 100,
		Valid:   true,
	}
	data.CreateTime = testTimeValue
	data.UpdateTime = sql.NullTime{
		Time:  testTimeValue,
		Valid: true,
	}

	err := mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("insert into %s", testTable)).
			WithArgs(data.Class, data.Name, data.Age, data.Score).
			WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m StudentModel, redis *redis.Redis) {
		r, err := m.Insert(data)
		assert.Nil(t, err)

		lastInsertId, err := r.LastInsertId()
		assert.Nil(t, err)
		assert.Equal(t, testInsertId, lastInsertId)

		rowsAffected, err := r.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, testRowsAffected, rowsAffected)
	})
	assert.Nil(t, err)

	err = mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s", testTable)).
			WithArgs(testInsertId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "class", "name", "age", "score", "create_time", "update_time"}).AddRow(testInsertId, data.Class, data.Name, data.Age, data.Score, testTimeValue, testTimeValue))
	}, func(m StudentModel, redis *redis.Redis) {
		result, err := m.FindOne(testInsertId)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)

		var resp Student
		val, err := redis.Get(fmt.Sprintf("%s%v", cacheStudentIdPrefix, testInsertId))
		assert.Nil(t, err)
		err = json.Unmarshal([]byte(val), &resp)
		assert.Nil(t, err)
		assert.Equal(t, resp.Name, data.Name)
	})
	assert.Nil(t, err)

	err = mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("update %s", testTable)).WithArgs(data.Class, testUpdateName, data.Age, data.Score, testInsertId).WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m StudentModel, redis *redis.Redis) {
		data.Name = testUpdateName
		err := m.Update(data)
		assert.Nil(t, err)

		val, err := redis.Get(fmt.Sprintf("%s%v", cacheStudentIdPrefix, testInsertId))
		assert.Nil(t, err)
		assert.Equal(t, "", val)
	})
	assert.Nil(t, err)

	data.Name = testUpdateName
	err = mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s ", testTable)).
			WithArgs(testInsertId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "class", "name", "age", "score", "create_time", "update_time"}).AddRow(testInsertId, data.Class, data.Name, data.Age, data.Score, testTimeValue, testTimeValue))
	}, func(m StudentModel, redis *redis.Redis) {
		result, err := m.FindOne(testInsertId)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)

		var resp Student
		val, err := redis.Get(fmt.Sprintf("%s%v", cacheStudentIdPrefix, testInsertId))
		assert.Nil(t, err)
		err = json.Unmarshal([]byte(val), &resp)
		assert.Nil(t, err)
		assert.Equal(t, testUpdateName, data.Name)
	})
	assert.Nil(t, err)

	err = mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s ", testTable)).
			WithArgs(class, testUpdateName).
			WillReturnRows(sqlmock.NewRows([]string{"id", "class", "name", "age", "score", "create_time", "update_time"}).AddRow(testInsertId, data.Class, data.Name, data.Age, data.Score, testTimeValue, testTimeValue))
	}, func(m StudentModel, redis *redis.Redis) {
		result, err := m.FindOneByClassName(class, testUpdateName)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)

		val, err := redis.Get(fmt.Sprintf("%s%v%v", cacheStudentClassNamePrefix, class, testUpdateName))
		assert.Nil(t, err)
		assert.Equal(t, "1", val)
	})
	assert.Nil(t, err)

	err = mockStudent(t, func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("delete from %s where `id` = ?", testTable)).WithArgs(testInsertId).WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m StudentModel, redis *redis.Redis) {
		err = m.Delete(testInsertId, class, testUpdateName)
		assert.Nil(t, err)

		val, err := redis.Get(fmt.Sprintf("%s%v", cacheStudentIdPrefix, testInsertId))
		assert.Nil(t, err)
		assert.Equal(t, "", val)

		val, err = redis.Get(fmt.Sprintf("%s%v%v", cacheStudentClassNamePrefix, class, testUpdateName))
		assert.Nil(t, err)
		assert.Equal(t, "", val)
	})

	assert.Nil(t, err)
}

func TestUserModel(t *testing.T) {
	var (
		testTimeValue          = time.Now()
		testTable              = "`user`"
		testUpdateName         = "gozero1"
		testUser               = "gozero"
		testPassword           = "test"
		testMobile             = "test_mobile"
		testGender             = "男"
		testNickname           = "test_nickname"
		testRowsAffected int64 = 1
		testInsertId     int64 = 1
	)

	var data User
	data.ID = testInsertId
	data.User = testUser
	data.Name = "gozero"
	data.Password = testPassword
	data.Mobile = testMobile
	data.Gender = testGender
	data.Nickname = testNickname
	data.CreateTime = testTimeValue
	data.UpdateTime = testTimeValue

	err := mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("insert into %s", testTable)).
			WithArgs(data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname).
			WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m UserModel) {
		r, err := m.Insert(data)
		assert.Nil(t, err)

		lastInsertId, err := r.LastInsertId()
		assert.Nil(t, err)
		assert.Equal(t, testInsertId, lastInsertId)

		rowsAffected, err := r.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, testRowsAffected, rowsAffected)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s", testTable)).
			WithArgs(testInsertId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user", "name", "password", "mobile", "gender", "nickname", "create_time", "update_time"}).AddRow(testInsertId, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, testTimeValue, testTimeValue))
	}, func(m UserModel) {
		result, err := m.FindOne(testInsertId)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("update %s", testTable)).WithArgs(data.User, testUpdateName, data.Password, data.Mobile, data.Gender, data.Nickname, testInsertId).WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m UserModel) {
		data.Name = testUpdateName
		err := m.Update(data)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s ", testTable)).
			WithArgs(testInsertId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user", "name", "password", "mobile", "gender", "nickname", "create_time", "update_time"}).AddRow(testInsertId, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, testTimeValue, testTimeValue))
	}, func(m UserModel) {
		result, err := m.FindOne(testInsertId)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("delete from %s where `id` = ?", testTable)).WithArgs(testInsertId).WillReturnResult(sqlmock.NewResult(testInsertId, testRowsAffected))
	}, func(m UserModel) {
		err := m.Delete(testInsertId)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

// with cache
func mockStudent(t *testing.T, mockFn func(mock sqlmock.Sqlmock), fn func(m StudentModel, r *redis.Redis)) error {
	db, mock, err := sqlmock.New()
	if err != nil {
		return err
	}

	defer db.Close()

	mock.ExpectBegin()
	mockFn(mock)
	mock.ExpectCommit()

	conn := mocksql.NewMockConn(db)
	r := redistest.CreateRedis(t)
	m := NewStudentModel(conn, cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: "node",
			},
			Weight: 100,
		},
	})
	mock.ExpectBegin()
	fn(m, r)
	mock.ExpectCommit()
	return nil
}

// without cache
func mockUser(mockFn func(mock sqlmock.Sqlmock), fn func(m UserModel)) error {
	db, mock, err := sqlmock.New()
	if err != nil {
		return err
	}

	defer db.Close()

	mock.ExpectBegin()
	mockFn(mock)
	mock.ExpectCommit()

	conn := mocksql.NewMockConn(db)
	m := NewUserModel(conn)
	fn(m)
	return nil
}
