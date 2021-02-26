package model

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/redis/redistest"
	mocksql "github.com/tal-tech/go-zero/tools/goctl/model/sql/test"
)

func TestStudentModel(t *testing.T) {
	var (
		testTimeValue          = time.Now()
		testTable              = "`student`"
		testUpdateName         = "gozero1"
		testRowsAffected int64 = 1
		testInsertID     int64 = 1
	)

	var data Student
	data.ID = testInsertID
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

	err := mockStudent(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("insert into %s", testTable)).
			WithArgs(data.Name, data.Age, data.Score).
			WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m StudentModel) {
		r, err := m.Insert(data)
		assert.Nil(t, err)

		lastInsertID, err := r.LastInsertId()
		assert.Nil(t, err)
		assert.Equal(t, testInsertID, lastInsertID)

		rowsAffected, err := r.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, testRowsAffected, rowsAffected)
	})
	assert.Nil(t, err)

	err = mockStudent(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s", testTable)).
			WithArgs(testInsertID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "score", "create_time", "update_time"}).AddRow(testInsertID, data.Name, data.Age, data.Score, testTimeValue, testTimeValue))
	}, func(m StudentModel) {
		result, err := m.FindOne(testInsertID)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockStudent(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("update %s", testTable)).WithArgs(testUpdateName, data.Age, data.Score, testInsertID).WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m StudentModel) {
		data.Name = testUpdateName
		err := m.Update(data)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)

	err = mockStudent(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s ", testTable)).
			WithArgs(testInsertID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "score", "create_time", "update_time"}).AddRow(testInsertID, data.Name, data.Age, data.Score, testTimeValue, testTimeValue))
	}, func(m StudentModel) {
		result, err := m.FindOne(testInsertID)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockStudent(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("delete from %s where `id` = ?", testTable)).WithArgs(testInsertID).WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m StudentModel) {
		err := m.Delete(testInsertID)
		assert.Nil(t, err)
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
		testInsertID     int64 = 1
	)

	var data User
	data.ID = testInsertID
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
			WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m UserModel) {
		r, err := m.Insert(data)
		assert.Nil(t, err)

		lastInsertID, err := r.LastInsertId()
		assert.Nil(t, err)
		assert.Equal(t, testInsertID, lastInsertID)

		rowsAffected, err := r.RowsAffected()
		assert.Nil(t, err)
		assert.Equal(t, testRowsAffected, rowsAffected)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s", testTable)).
			WithArgs(testInsertID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user", "name", "password", "mobile", "gender", "nickname", "create_time", "update_time"}).AddRow(testInsertID, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, testTimeValue, testTimeValue))
	}, func(m UserModel) {
		result, err := m.FindOne(testInsertID)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("update %s", testTable)).WithArgs(data.User, testUpdateName, data.Password, data.Mobile, data.Gender, data.Nickname, testInsertID).WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m UserModel) {
		data.Name = testUpdateName
		err := m.Update(data)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(fmt.Sprintf("select (.+) from %s ", testTable)).
			WithArgs(testInsertID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user", "name", "password", "mobile", "gender", "nickname", "create_time", "update_time"}).AddRow(testInsertID, data.User, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, testTimeValue, testTimeValue))
	}, func(m UserModel) {
		result, err := m.FindOne(testInsertID)
		assert.Nil(t, err)
		assert.Equal(t, *result, data)
	})
	assert.Nil(t, err)

	err = mockUser(func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(fmt.Sprintf("delete from %s where `id` = ?", testTable)).WithArgs(testInsertID).WillReturnResult(sqlmock.NewResult(testInsertID, testRowsAffected))
	}, func(m UserModel) {
		err := m.Delete(testInsertID)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

// with cache
func mockStudent(mockFn func(mock sqlmock.Sqlmock), fn func(m StudentModel)) error {
	db, mock, err := sqlmock.New()
	if err != nil {
		return err
	}

	defer db.Close()

	mock.ExpectBegin()
	mockFn(mock)
	mock.ExpectCommit()

	conn := mocksql.NewMockConn(db)
	r, clean, err := redistest.CreateRedis()
	if err != nil {
		return err
	}

	defer clean()

	m := NewStudentModel(conn, cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: "node",
			},
			Weight: 100,
		},
	})
	fn(m)
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
