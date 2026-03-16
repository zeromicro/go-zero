package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/dbtest"
)

func TestUnmarshalRowBool(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value bool
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.True(t, value)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value struct {
			Value bool `db:"value"`
		}
		assert.Error(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select value from users where user=?", "anyone"))
	})
}

func TestUnmarshalRowBoolNotSettable(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value bool
		assert.NotNil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select value from users where user=?", "anyone"))
	})
}

func TestUnmarshalRowInt(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value int
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, 2, value)
	})
}

func TestUnmarshalRowInt8(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value int8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, int8(3), value)
	})
}

func TestUnmarshalRowInt16(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("4")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value int16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.Equal(t, int16(4), value)
	})
}

func TestUnmarshalRowInt32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value int32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.Equal(t, int32(5), value)
	})
}

func TestUnmarshalRowInt64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("6")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value int64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, int64(6), value)
	})
}

func TestUnmarshalRowUint(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value uint
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint(2), value)
	})
}

func TestUnmarshalRowUint8(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value uint8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint8(3), value)
	})
}

func TestUnmarshalRowUint16(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("4")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value uint16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint16(4), value)
	})
}

func TestUnmarshalRowUint32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value uint32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint32(5), value)
	})
}

func TestUnmarshalRowUint64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("6")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value uint64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, uint16(6), value)
	})
}

func TestUnmarshalRowFloat32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("7")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value float32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, float32(7), value)
	})
}

func TestUnmarshalRowFloat64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("8")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value float64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, float64(8), value)
	})
}

func TestUnmarshalRowString(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		const expect = "hello"
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString(expect)
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value string
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowStruct(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Name string
			Age  int
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, "liao", value.Name)
		assert.Equal(t, 5, value.Age)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Name string
			Age  int
		})

		errAny := errors.New("any error")
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, &mockedScanner{
				colErr: errAny,
				next:   1,
			}, true)
		}, "select name, age from users where user=?", "anyone"), errAny)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Name string
			age  *int
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotMatchDestination)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Name string
			age  int
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotMatchDestination)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("8")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		type myString chan int
		var value myString
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select value from users where user=?", "anyone"), ErrUnsupportedValueType)
	})
}

func TestUnmarshalRowStructWithTags(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Age  int    `db:"age"`
			Name string `db:"name"`
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, "liao", value.Name)
		assert.Equal(t, 5, value.Age)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			age  *int   `db:"age"`
			Name string `db:"name"`
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotReadableValue)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			age  int    `db:"age"`
			Name string `db:"name"`
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotReadableValue)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value struct {
			Age  *int    `db:"age"`
			Name *string `db:"name"`
		}

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, "liao", *value.Name)
		assert.Equal(t, 5, *value.Age)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Age  int `db:"age"`
			Name string
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, 5, value.Age)
	})
}

func TestUnmarshalRowStructWithTagsIgnoreFields(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Age    int `db:"age"`
			Name   string
			Ignore bool
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotMatchDestination)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		value := new(struct {
			Age    int `db:"age"`
			Name   string
			Ignore bool `db:"-"`
		})

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("liao,5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(t, 5, value.Age)
	})
}

func TestUnmarshalRowStructWithTagsWrongColumns(t *testing.T) {
	value := new(struct {
		Age  *int   `db:"age"`
		Name string `db:"name"`
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name"}).FromCSVString("liao")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.NotNil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
	})
}

func TestUnmarshalRowsBool(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []bool{true, false}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []bool
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []bool
		assert.Error(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(value, rows, true)
		}, "select value from users where user=?", "anyone"))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value struct {
			value []bool `db:"value"`
		}
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"), ErrUnsupportedValueType)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []bool
		errAny := errors.New("any")
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, &mockedScanner{
				scanErr: errAny,
				next:    1,
			}, true)
		}, "select value from users where user=?", "anyone"), errAny)
	})
}

func TestUnmarshalRowsInt(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []int{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []int
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt8(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []int8{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []int8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt16(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []int16{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []int16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []int32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []int32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []int64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []int64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []uint{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []uint
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint8(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []uint8{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []uint8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint16(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []uint16{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []uint16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []uint32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []uint32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []uint64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []uint64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat32(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []float32{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []float32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat64(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []float64{2, 3}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []float64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsString(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []string{"hello", "world"}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("hello\nworld")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []string
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsBoolPtr(t *testing.T) {
	yes := true
	no := false
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*bool{&yes, &no}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("1\n0")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*bool
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsIntPtr(t *testing.T) {
	two := 2
	three := 3
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*int{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*int
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt8Ptr(t *testing.T) {
	two := int8(2)
	three := int8(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*int8{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*int8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt16Ptr(t *testing.T) {
	two := int16(2)
	three := int16(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*int16{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*int16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt32Ptr(t *testing.T) {
	two := int32(2)
	three := int32(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*int32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*int32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsInt64Ptr(t *testing.T) {
	two := int64(2)
	three := int64(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*int64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*int64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUintPtr(t *testing.T) {
	two := uint(2)
	three := uint(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*uint{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*uint
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint8Ptr(t *testing.T) {
	two := uint8(2)
	three := uint8(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*uint8{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*uint8
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint16Ptr(t *testing.T) {
	two := uint16(2)
	three := uint16(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*uint16{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*uint16
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint32Ptr(t *testing.T) {
	two := uint32(2)
	three := uint32(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*uint32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*uint32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsUint64Ptr(t *testing.T) {
	two := uint64(2)
	three := uint64(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*uint64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*uint64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat32Ptr(t *testing.T) {
	two := float32(2)
	three := float32(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*float32{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*float32
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsFloat64Ptr(t *testing.T) {
	two := float64(2)
	three := float64(3)
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*float64{&two, &three}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("2\n3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*float64
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsStringPtr(t *testing.T) {
	hello := "hello"
	world := "world"
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []*string{&hello, &world}
		rs := sqlmock.NewRows([]string{"value"}).FromCSVString("hello\nworld")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var value []*string
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select value from users where user=?", "anyone"))
		assert.EqualValues(t, expect, value)
	})
}

func TestUnmarshalRowsStruct(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expect := []struct {
			Name string
			Age  int64
		}{
			{
				Name: "first",
				Age:  2,
			},
			{
				Name: "second",
				Age:  3,
			},
		}
		var value []struct {
			Name string
			Age  int64
		}

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value []struct {
			Name string
			Age  int64
		}

		errAny := errors.New("any error")
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, &mockedScanner{
				colErr: errAny,
				next:   1,
			}, true)
		}, "select name, age from users where user=?", "anyone"), errAny)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value []struct {
			Name string
			Age  int64
		}

		errAny := errors.New("any error")
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, &mockedScanner{
				cols:    []string{"name", "age"},
				scanErr: errAny,
				next:    1,
			}, true)
		}, "select name, age from users where user=?", "anyone"), errAny)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value []chan int

		errAny := errors.New("any error")
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, &mockedScanner{
				cols:    []string{"name", "age"},
				scanErr: errAny,
				next:    1,
			}, true)
		}, "select name, age from users where user=?", "anyone"), ErrUnsupportedValueType)
	})
}

func TestUnmarshalRowsStructWithNullStringType(t *testing.T) {
	expect := []struct {
		Name       string
		NullString sql.NullString
	}{
		{
			Name: "first",
			NullString: sql.NullString{
				String: "firstnullstring",
				Valid:  true,
			},
		},
		{
			Name: "second",
			NullString: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	var value []struct {
		Name       string         `db:"name"`
		NullString sql.NullString `db:"value"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "value"}).AddRow(
			"first", "firstnullstring").AddRow("second", nil)
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.NullString.String, value[i].NullString.String)
			assert.Equal(t, each.NullString.Valid, value[i].NullString.Valid)
		}
	})
}

func TestUnmarshalRowsStructWithTags(t *testing.T) {
	expect := []struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []struct {
		Age  int64  `db:"age"`
		Name string `db:"name"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithTagsIgnoreFields(t *testing.T) {
	expect := []struct {
		Name   string
		Age    int64
		Ignore bool
	}{
		{
			Name:   "first",
			Age:    2,
			Ignore: false,
		},
		{
			Name:   "second",
			Age:    3,
			Ignore: false,
		},
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value []struct {
			Age    int64  `db:"age"`
			Name   string `db:"name"`
			Ignore bool
		}

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"), ErrNotMatchDestination)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var value []struct {
			Age    int64  `db:"age"`
			Name   string `db:"name"`
			Ignore bool   `db:"-"`
		}

		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructAndEmbeddedAnonymousStructWithTags(t *testing.T) {
	type Embed struct {
		Value int64 `db:"value"`
	}

	expect := []struct {
		Name  string
		Age   int64
		Value int64
	}{
		{
			Name:  "first",
			Age:   2,
			Value: 3,
		},
		{
			Name:  "second",
			Age:   3,
			Value: 4,
		},
	}
	var value []struct {
		Name string `db:"name"`
		Age  int64  `db:"age"`
		Embed
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age", "value"}).FromCSVString("first,2,3\nsecond,3,4")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, value from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.Value, value[i].Value)
		}
	})
}

func TestUnmarshalRowsStructAndEmbeddedStructPtrAnonymousWithTags(t *testing.T) {
	type Embed struct {
		Value int64 `db:"value"`
	}

	expect := []struct {
		Name  string
		Age   int64
		Value int64
	}{
		{
			Name:  "first",
			Age:   2,
			Value: 3,
		},
		{
			Name:  "second",
			Age:   3,
			Value: 4,
		},
	}
	var value []struct {
		Name string `db:"name"`
		Age  int64  `db:"age"`
		*Embed
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age", "value"}).FromCSVString("first,2,3\nsecond,3,4")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, value from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.Value, value[i].Value)
		}
	})
}

func TestUnmarshalRowsStructPtr(t *testing.T) {
	expect := []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Name string
		Age  int64
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithTagsPtr(t *testing.T) {
	expect := []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Age  int64  `db:"age"`
		Name string `db:"name"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
		}
	})
}

func TestUnmarshalRowsStructWithTagsPtrWithInnerPtr(t *testing.T) {
	expect := []*struct {
		Name string
		Age  int64
	}{
		{
			Name: "first",
			Age:  2,
		},
		{
			Name: "second",
			Age:  3,
		},
	}
	var value []*struct {
		Age  *int64 `db:"age"`
		Name string `db:"name"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "age"}).FromCSVString("first,2\nsecond,3")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age from users where user=?", "anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, *value[i].Age)
		}
	})
}

func TestCommonSqlConn_QueryRowOptional(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"age"}).FromCSVString("5")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		var r struct {
			User string `db:"user"`
			Age  int    `db:"age"`
		}
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(&r, rows, false)
		}, "select age from users where user=?", "anyone"))
		assert.Empty(t, r.User)
		assert.Equal(t, 5, r.Age)
	})
}

func TestUnmarshalRowError(t *testing.T) {
	tests := []struct {
		name     string
		colErr   error
		scanErr  error
		err      error
		next     int
		validate func(err error)
	}{
		{
			name: "with error",
			err:  errors.New("foo"),
			validate: func(err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "without next",
			validate: func(err error) {
				assert.Equal(t, ErrNotFound, err)
			},
		},
		{
			name:    "with error",
			scanErr: errors.New("foo"),
			next:    1,
			validate: func(err error) {
				assert.Equal(t, ErrNotFound, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
				rs := sqlmock.NewRows([]string{"age"}).FromCSVString("5")
				mock.ExpectQuery("select (.+) from users where user=?").
					WithArgs("anyone").WillReturnRows(rs)

				var r struct {
					User string `db:"user"`
					Age  int    `db:"age"`
				}
				test.validate(query(context.Background(), db, func(rows *sql.Rows) error {
					scanner := mockedScanner{
						colErr:  test.colErr,
						scanErr: test.scanErr,
						err:     test.err,
					}
					return unmarshalRow(&r, &scanner, false)
				}, "select age from users where user=?", "anyone"))
			})
		})
	}
}

func TestAnonymousStructPr(t *testing.T) {
	type Score struct {
		Discipline string `db:"discipline"`
		Score      uint   `db:"score"`
	}
	type ClassType struct {
		Grade     sql.NullString `db:"grade"`
		ClassName *string        `db:"class_name"`
	}
	type Class struct {
		*ClassType
		Score
	}
	expect := []*struct {
		Name       string
		Age        int64
		Grade      sql.NullString
		Discipline string
		Score      uint
		ClassName  string
	}{
		{
			Name: "first",
			Age:  2,
			Grade: sql.NullString{
				String: "",
				Valid:  false,
			},
			ClassName:  "experimental class",
			Discipline: "math",
			Score:      100,
		},
		{
			Name: "second",
			Age:  3,
			Grade: sql.NullString{
				String: "grade one",
				Valid:  true,
			},
			ClassName:  "class three grade two",
			Discipline: "chinese",
			Score:      99,
		},
	}
	var value []*struct {
		Age int64 `db:"age"`
		Class
		Name string `db:"name"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{
			"name",
			"age",
			"grade",
			"discipline",
			"class_name",
			"score",
		}).
			AddRow("first", 2, nil, "math", "experimental class", 100).
			AddRow("second", 3, "grade one", "chinese", "class three grade two", 99)
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age,grade,discipline,class_name,score from users where user=?",
			"anyone"))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.ClassName, *value[i].Class.ClassName)
			assert.Equal(t, each.Discipline, value[i].Score.Discipline)
			assert.Equal(t, each.Score, value[i].Score.Score)
			assert.Equal(t, each.Grade, value[i].Class.Grade)
		}
	})
}

func TestAnonymousStructPrError(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		type Score struct {
			Discipline string `db:"discipline"`
			score      uint   `db:"score"`
		}
		type ClassType struct {
			Grade     sql.NullString `db:"grade"`
			ClassName *string        `db:"class_name"`
		}
		type Class struct {
			*ClassType
			Score
		}

		var value []*struct {
			Age int64 `db:"age"`
			Class
			Name string `db:"name"`
		}
		rs := sqlmock.NewRows([]string{
			"name",
			"age",
			"grade",
			"discipline",
			"class_name",
			"score",
		}).
			AddRow("first", 2, nil, "math", "experimental class", 100).
			AddRow("second", 3, "grade one", "chinese", "class three grade two", 99)
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, grade, discipline, class_name, score from users where user=?",
			"anyone"), ErrNotReadableValue)
		if len(value) > 0 {
			assert.Equal(t, value[0].score, 0)
		}
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		type Score struct {
			Discipline string
			score      uint
		}
		type ClassType struct {
			Grade     sql.NullString
			ClassName *string
		}
		type Class struct {
			*ClassType
			Score
		}

		var value []*struct {
			Age int64
			Class
			Name string
		}
		rs := sqlmock.NewRows([]string{
			"name",
			"age",
			"grade",
			"discipline",
			"class_name",
			"score",
		}).
			AddRow("first", 2, nil, "math", "experimental class", 100).
			AddRow("second", 3, "grade one", "chinese", "class three grade two", 99)
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)
		assert.ErrorIs(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, age, grade, discipline, class_name, score from users where user=?",
			"anyone"), ErrNotMatchDestination)
		if len(value) > 0 {
			assert.Equal(t, value[0].score, 0)
		}
	})
}

func TestUnmarshalRowsZeroValueStructPtr(t *testing.T) {
	secondNamePtr := "second_ptr"
	secondAgePtr := int64(30)
	thirdNamePtr := "third_ptr"
	thirdAgePtr := int64(0)

	expect := []struct {
		Name    string
		NamePtr *string
		Age     int64
		AgePtr  *int64
	}{
		{
			Name:    "first",
			NamePtr: nil,
			Age:     2,
			AgePtr:  nil,
		},
		{
			Name:    "second",
			NamePtr: &secondNamePtr,
			Age:     3,
			AgePtr:  &secondAgePtr,
		},
		{
			Name:    "",
			NamePtr: &thirdNamePtr,
			Age:     0,
			AgePtr:  &thirdAgePtr,
		},
	}

	var value []struct {
		Age     int64   `db:"age"`
		AgePtr  *int64  `db:"age_ptr"`
		Name    string  `db:"name"`
		NamePtr *string `db:"name_ptr"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "name_ptr", "age", "age_ptr"}).
			AddRow("first", nil, 2, nil).
			AddRow("second", "second_ptr", 3, 30).
			AddRow("", "third_ptr", 0, 0)

		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, name_ptr, age, age_ptr from users where user=?", "anyone"))

		assert.Equal(t, 3, len(value), "应该返回3行数据")

		for i, each := range expect {

			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)

			if each.NamePtr == nil {
				assert.Nil(t, value[i].NamePtr)
			} else {
				assert.NotNil(t, value[i].NamePtr)
				assert.Equal(t, *each.NamePtr, *value[i].NamePtr)
			}

			if each.AgePtr == nil {
				assert.Nil(t, value[i].AgePtr)
			} else {
				assert.NotNil(t, value[i].AgePtr)
				assert.Equal(t, *each.AgePtr, *value[i].AgePtr)
			}
		}
	})
}

func TestUnmarshalRowsAllNullStructPtrFields(t *testing.T) {
	expect := []struct {
		NamePtr *string
		AgePtr  *int64
	}{
		{
			NamePtr: nil,
			AgePtr:  nil,
		},
		{
			NamePtr: stringPtr("second"),
			AgePtr:  int64Ptr(30),
		},
		{
			NamePtr: nil,
			AgePtr:  nil,
		},
	}

	var value []struct {
		AgePtr  *int64  `db:"age_ptr"`
		NamePtr *string `db:"name_ptr"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name_ptr", "age_ptr"}).
			AddRow(nil, nil).
			AddRow("second", 30).
			AddRow(nil, nil)

		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name_ptr, age_ptr from users where user=?", "anyone"))

		assert.Equal(t, 3, len(value))

		for i, each := range expect {
			if each.NamePtr == nil {
				assert.Nil(t, value[i].NamePtr)
			} else {
				assert.NotNil(t, value[i].NamePtr)
				assert.Equal(t, *each.NamePtr, *value[i].NamePtr)
			}

			if each.AgePtr == nil {
				assert.Nil(t, value[i].AgePtr)
			} else {
				assert.NotNil(t, value[i].AgePtr)
				assert.Equal(t, *each.AgePtr, *value[i].AgePtr)
			}
		}
	})
}

func TestUnmarshalRowsWithSqlNullTypes(t *testing.T) {
	expect := []struct {
		Name       string
		NullName   sql.NullString
		Age        int64
		NullAge    sql.NullInt64
		Score      float64
		NullScore  sql.NullFloat64
		Active     bool
		NullActive sql.NullBool
	}{
		{
			Name: "first",
			NullName: sql.NullString{
				String: "",
				Valid:  false,
			},
			Age: 20,
			NullAge: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Score: 85.5,
			NullScore: sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			},
			Active: true,
			NullActive: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			Name: "second",
			NullName: sql.NullString{
				String: "not_null_name",
				Valid:  true,
			},
			Age: 25,
			NullAge: sql.NullInt64{
				Int64: 30,
				Valid: true,
			},
			Score: 90.0,
			NullScore: sql.NullFloat64{
				Float64: 95.5,
				Valid:   true,
			},
			Active: false,
			NullActive: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			Name: "third",
			NullName: sql.NullString{
				String: "",
				Valid:  false,
			},
			Age: 0,
			NullAge: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			Score: 0,
			NullScore: sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			},
			Active: false,
			NullActive: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}

	var value []struct {
		Name       string          `db:"name"`
		NullName   sql.NullString  `db:"null_name"`
		Age        int64           `db:"age"`
		NullAge    sql.NullInt64   `db:"null_age"`
		Score      float64         `db:"score"`
		NullScore  sql.NullFloat64 `db:"null_score"`
		Active     bool            `db:"active"`
		NullActive sql.NullBool    `db:"null_active"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{
			"name", "null_name", "age", "null_age", "score", "null_score", "active", "null_active",
		}).
			AddRow("first", nil, 20, nil, 85.5, nil, true, nil).
			AddRow("second", "not_null_name", 25, 30, 90.0, 95.5, false, true).
			AddRow("third", nil, 0, nil, 0, nil, false, nil)

		mock.ExpectQuery("select (.+) from users where type=?").
			WithArgs("test").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, null_name, age, null_age, score, null_score, active, null_active from users where type=?", "test"))

		assert.Equal(t, 3, len(value))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.Score, value[i].Score)
			assert.Equal(t, each.Active, value[i].Active)

			assert.Equal(t, each.NullName.Valid, value[i].NullName.Valid)
			if each.NullName.Valid {
				assert.Equal(t, each.NullName.String, value[i].NullName.String)
			}

			assert.Equal(t, each.NullAge.Valid, value[i].NullAge.Valid)
			if each.NullAge.Valid {
				assert.Equal(t, each.NullAge.Int64, value[i].NullAge.Int64)
			}

			assert.Equal(t, each.NullScore.Valid, value[i].NullScore.Valid)
			if each.NullScore.Valid {
				assert.Equal(t, each.NullScore.Float64, value[i].NullScore.Float64)
			}

			assert.Equal(t, each.NullActive.Valid, value[i].NullActive.Valid)
			if each.NullActive.Valid {
				assert.Equal(t, each.NullActive.Bool, value[i].NullActive.Bool)
			}
		}
	})
}

func TestUnmarshalRowsSqlNullWithMixedData(t *testing.T) {
	expect := []struct {
		Name       string
		NullName   sql.NullString
		Age        int64
		NullAge    sql.NullInt64
		IsStudent  bool
		NullActive sql.NullBool
	}{
		{
			Name: "student1",
			NullName: sql.NullString{
				String: "",
				Valid:  false,
			},
			Age: 18,
			NullAge: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			IsStudent: true,
			NullActive: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			Name: "student2",
			NullName: sql.NullString{
				String: "has_nickname",
				Valid:  true,
			},
			Age: 20,
			NullAge: sql.NullInt64{
				Int64: 22,
				Valid: true,
			},
			IsStudent: false,
			NullActive: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
	}

	var value []struct {
		Name       string         `db:"name"`
		NullName   sql.NullString `db:"null_name"`
		Age        int64          `db:"age"`
		NullAge    sql.NullInt64  `db:"null_age"`
		IsStudent  bool           `db:"is_student"`
		NullActive sql.NullBool   `db:"null_active"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "null_name", "age", "null_age", "is_student", "null_active"}).
			AddRow("student1", nil, 18, nil, true, nil).
			AddRow("student2", "has_nickname", 20, 22, false, true)

		mock.ExpectQuery("select (.+) from students where class=?").
			WithArgs("A").WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, null_name, age, null_age, is_student, null_active from students where class=?", "A"))

		assert.Equal(t, 2, len(value))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)
			assert.Equal(t, each.Age, value[i].Age)
			assert.Equal(t, each.IsStudent, value[i].IsStudent)

			assert.Equal(t, each.NullName.Valid, value[i].NullName.Valid)
			if each.NullName.Valid {
				assert.Equal(t, each.NullName.String, value[i].NullName.String)
			}

			assert.Equal(t, each.NullAge.Valid, value[i].NullAge.Valid)
			if each.NullAge.Valid {
				assert.Equal(t, each.NullAge.Int64, value[i].NullAge.Int64)
			}

			assert.Equal(t, each.NullActive.Valid, value[i].NullActive.Valid)
			if each.NullActive.Valid {
				assert.Equal(t, each.NullActive.Bool, value[i].NullActive.Bool)
			}
		}
	})
}

func TestUnmarshalRowsSqlNullTime(t *testing.T) {
	now := time.Now()
	futureTime := now.AddDate(1, 0, 0)

	expect := []struct {
		Name      string
		BirthDate sql.NullTime
		LastLogin sql.NullTime
	}{
		{
			Name: "user1",
			BirthDate: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			LastLogin: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		},
		{
			Name: "user2",
			BirthDate: sql.NullTime{
				Time:  futureTime,
				Valid: true,
			},
			LastLogin: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}

	var value []struct {
		Name      string       `db:"name"`
		BirthDate sql.NullTime `db:"birth_date"`
		LastLogin sql.NullTime `db:"last_login"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "birth_date", "last_login"}).
			AddRow("user1", nil, now).
			AddRow("user2", futureTime, nil)

		mock.ExpectQuery("select (.+) from users").
			WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, birth_date, last_login from users"))

		assert.Equal(t, 2, len(value))

		for i, each := range expect {
			assert.Equal(t, each.Name, value[i].Name)

			assert.Equal(t, each.BirthDate.Valid, value[i].BirthDate.Valid)
			if each.BirthDate.Valid {
				assert.WithinDuration(t, each.BirthDate.Time, value[i].BirthDate.Time, time.Second)
			}

			assert.Equal(t, each.LastLogin.Valid, value[i].LastLogin.Valid)
			if each.LastLogin.Valid {
				assert.WithinDuration(t, each.LastLogin.Time, value[i].LastLogin.Time, time.Second)
			}
		}
	})
}

func TestUnmarshalRowsSqlNullWithEmptyValues(t *testing.T) {
	expect := []struct {
		Name       string
		NullString sql.NullString
		NullInt    sql.NullInt64
		NullFloat  sql.NullFloat64
		NullBool   sql.NullBool
	}{
		{
			Name: "empty_values",
			NullString: sql.NullString{
				String: "",
				Valid:  true,
			},
			NullInt: sql.NullInt64{
				Int64: 0,
				Valid: true,
			},
			NullFloat: sql.NullFloat64{
				Float64: 0.0,
				Valid:   true,
			},
			NullBool: sql.NullBool{
				Bool:  false,
				Valid: true,
			},
		},
		{
			Name: "null_values",
			NullString: sql.NullString{
				String: "",
				Valid:  false,
			},
			NullInt: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			NullFloat: sql.NullFloat64{
				Float64: 0.0,
				Valid:   false,
			},
			NullBool: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
		{
			Name: "mixed_values",
			NullString: sql.NullString{
				String: "actual_value",
				Valid:  true,
			},
			NullInt: sql.NullInt64{
				Int64: 0,
				Valid: true,
			},
			NullFloat: sql.NullFloat64{
				Float64: 0.0,
				Valid:   false,
			},
			NullBool: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
	}

	var value []struct {
		Name       string          `db:"name"`
		NullString sql.NullString  `db:"null_string"`
		NullInt    sql.NullInt64   `db:"null_int"`
		NullFloat  sql.NullFloat64 `db:"null_float"`
		NullBool   sql.NullBool    `db:"null_bool"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "null_string", "null_int", "null_float", "null_bool"}).
			AddRow("empty_values", "", 0, 0.0, false).
			AddRow("null_values", nil, nil, nil, nil).
			AddRow("mixed_values", "actual_value", 0, nil, true)

		mock.ExpectQuery("select (.+) from test_table").
			WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, null_string, null_int, null_float, null_bool from test_table"))

		assert.Equal(t, 3, len(value))

		for i, each := range expect {

			assert.Equal(t, each.Name, value[i].Name)

			assert.Equal(t, each.NullString.Valid, value[i].NullString.Valid)
			if each.NullString.Valid {
				assert.Equal(t, each.NullString.String, value[i].NullString.String)
			} else {
				assert.Equal(t, "", value[i].NullString.String)
			}

			assert.Equal(t, each.NullInt.Valid, value[i].NullInt.Valid)
			if each.NullInt.Valid {
				assert.Equal(t, each.NullInt.Int64, value[i].NullInt.Int64)
			} else {
				assert.Equal(t, int64(0), value[i].NullInt.Int64)
			}

			assert.Equal(t, each.NullFloat.Valid, value[i].NullFloat.Valid)
			if each.NullFloat.Valid {
				assert.Equal(t, each.NullFloat.Float64, value[i].NullFloat.Float64)
			} else {
				assert.Equal(t, 0.0, value[i].NullFloat.Float64)
			}

			assert.Equal(t, each.NullBool.Valid, value[i].NullBool.Valid)
			if each.NullBool.Valid {
				assert.Equal(t, each.NullBool.Bool, value[i].NullBool.Bool)
			} else {
				assert.Equal(t, false, value[i].NullBool.Bool)
			}
		}
	})
}

func TestUnmarshalRowsSqlNullStringEmptyVsNull(t *testing.T) {
	expect := []struct {
		Name         string
		EmptyString  sql.NullString
		NullString   sql.NullString
		NormalString sql.NullString
	}{
		{
			Name: "row1",
			EmptyString: sql.NullString{
				String: "",
				Valid:  true,
			},
			NullString: sql.NullString{
				String: "",
				Valid:  false,
			},
			NormalString: sql.NullString{
				String: "hello",
				Valid:  true,
			},
		},
		{
			Name: "row2",
			EmptyString: sql.NullString{
				String: "   ",
				Valid:  true,
			},
			NullString: sql.NullString{
				String: "",
				Valid:  false,
			},
			NormalString: sql.NullString{
				String: "",
				Valid:  true,
			},
		},
	}

	var value []struct {
		Name         string         `db:"name"`
		EmptyString  sql.NullString `db:"empty_string"`
		NullString   sql.NullString `db:"null_string"`
		NormalString sql.NullString `db:"normal_string"`
	}

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rs := sqlmock.NewRows([]string{"name", "empty_string", "null_string", "normal_string"}).
			AddRow("row1", "", nil, "hello").
			AddRow("row2", "   ", nil, "")

		mock.ExpectQuery("select (.+) from string_test").
			WillReturnRows(rs)

		assert.Nil(t, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRows(&value, rows, true)
		}, "select name, empty_string, null_string, normal_string from string_test"))

		assert.Equal(t, 2, len(value))

		for i, each := range expect {
			assert.True(t, value[i].EmptyString.Valid)
			assert.Equal(t, each.EmptyString.String, value[i].EmptyString.String)

			assert.False(t, value[i].NullString.Valid)
			assert.Equal(t, "", value[i].NullString.String)

			assert.Equal(t, each.NormalString.Valid, value[i].NormalString.Valid)
			if each.NormalString.Valid {
				assert.Equal(t, each.NormalString.String, value[i].NormalString.String)
			}
		}
	})
}

func TestGetValueInterface(t *testing.T) {
	t.Run("non_pointer_field", func(t *testing.T) {
		type testStruct struct {
			Name string
			Age  int
		}
		s := testStruct{}
		v := reflect.ValueOf(&s).Elem()

		nameField := v.Field(0)
		result, err := getValueInterface(nameField)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Should return pointer to the field
		ptr, ok := result.(*string)
		assert.True(t, ok)
		*ptr = "test"
		assert.Equal(t, "test", s.Name)
	})

	t.Run("pointer_field_nil", func(t *testing.T) {
		type testStruct struct {
			NamePtr *string
			AgePtr  *int64
		}
		s := testStruct{}
		v := reflect.ValueOf(&s).Elem()

		// Test with nil pointer field
		namePtrField := v.Field(0)
		assert.True(t, namePtrField.IsNil(), "initial pointer should be nil")

		result, err := getValueInterface(namePtrField)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Should have allocated the pointer
		assert.False(t, namePtrField.IsNil(), "pointer should be allocated after getValueInterface")

		// Should return pointer to pointer field
		ptrPtr, ok := result.(**string)
		assert.True(t, ok)
		testValue := "initialized"
		*ptrPtr = &testValue
		assert.NotNil(t, s.NamePtr)
		assert.Equal(t, "initialized", *s.NamePtr)
	})

	t.Run("pointer_field_already_allocated", func(t *testing.T) {
		type testStruct struct {
			NamePtr *string
		}
		initial := "existing"
		s := testStruct{NamePtr: &initial}
		v := reflect.ValueOf(&s).Elem()

		namePtrField := v.Field(0)
		assert.False(t, namePtrField.IsNil(), "pointer should not be nil initially")

		result, err := getValueInterface(namePtrField)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Should return pointer to pointer field
		ptrPtr, ok := result.(**string)
		assert.True(t, ok)

		// Verify it points to the existing value
		assert.Equal(t, "existing", **ptrPtr)

		// Modify through the returned pointer
		newValue := "modified"
		*ptrPtr = &newValue
		assert.Equal(t, "modified", *s.NamePtr)
	})

	t.Run("pointer_field_zero_value", func(t *testing.T) {
		type testStruct struct {
			IntPtr *int
		}
		s := testStruct{}
		v := reflect.ValueOf(&s).Elem()

		intPtrField := v.Field(0)
		result, err := getValueInterface(intPtrField)
		assert.NoError(t, err)

		// After calling getValueInterface, nil pointer should be allocated
		assert.NotNil(t, s.IntPtr)

		// Set zero value through returned interface
		ptrPtr, ok := result.(**int)
		assert.True(t, ok)
		zero := 0
		*ptrPtr = &zero
		assert.Equal(t, 0, *s.IntPtr)
	})

	t.Run("not_addressable_value", func(t *testing.T) {
		type testStruct struct {
			Name string
		}
		s := testStruct{Name: "test"}
		v := reflect.ValueOf(s) // Non-pointer, not addressable

		nameField := v.Field(0)
		result, err := getValueInterface(nameField)
		assert.Error(t, err)
		assert.Equal(t, ErrNotReadableValue, err)
		assert.Nil(t, result)
	})

	t.Run("multiple_pointer_types", func(t *testing.T) {
		type testStruct struct {
			StringPtr *string
			IntPtr    *int
			Int64Ptr  *int64
			FloatPtr  *float64
			BoolPtr   *bool
		}
		s := testStruct{}
		v := reflect.ValueOf(&s).Elem()

		// Test each pointer type gets properly initialized
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			assert.True(t, field.IsNil(), "field %d should start as nil", i)

			result, err := getValueInterface(field)
			assert.NoError(t, err, "field %d should not error", i)
			assert.NotNil(t, result, "field %d result should not be nil", i)

			// After getValueInterface, pointer should be allocated
			assert.False(t, field.IsNil(), "field %d should be allocated", i)
		}
	})
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func BenchmarkIgnore(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		_ = db.Close()
	}()

	for i := 0; i < b.N; i++ {
		value := new(struct {
			Age    int `db:"age"`
			Name   string
			Ignore bool `db:"-"`
		})

		rs := sqlmock.NewRows([]string{"name", "age", "ignore"}).FromCSVString("liao,5,true")
		mock.ExpectQuery("select (.+) from users where user=?").
			WithArgs("anyone").WillReturnRows(rs)

		assert.Nil(b, query(context.Background(), db, func(rows *sql.Rows) error {
			return unmarshalRow(value, rows, true)
		}, "select name, age from users where user=?", "anyone"))
		assert.Equal(b, 5, value.Age)

	}
}

type mockedScanner struct {
	cols    []string
	colErr  error
	scanErr error
	err     error
	next    int
}

func (m *mockedScanner) Columns() ([]string, error) {
	return m.cols, m.colErr
}

func (m *mockedScanner) Err() error {
	return m.err
}

func (m *mockedScanner) Next() bool {
	if m.next > 0 {
		m.next--
		return true
	}
	return false
}

func (m *mockedScanner) Scan(v ...any) error {
	return m.scanErr
}
