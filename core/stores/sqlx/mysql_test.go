package sqlx

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
)

func init() {
	stat.SetReporter(nil)
}

func TestBreakerOnDuplicateEntry(t *testing.T) {
	logx.Disable()

	err := tryOnDuplicateEntryError(t, mysqlAcceptable)
	assert.Equal(t, duplicateEntryCode, err.(*mysql.MySQLError).Number)
}

func TestBreakerOnNotHandlingDuplicateEntry(t *testing.T) {
	logx.Disable()

	var found bool
	for i := 0; i < 100; i++ {
		if errors.Is(tryOnDuplicateEntryError(t, nil), breaker.ErrServiceUnavailable) {
			found = true
		}
	}
	assert.True(t, found)
}

func TestMysqlAcceptable(t *testing.T) {
	conn := NewMysql("nomysql").(*commonSqlConn)
	withMysqlAcceptable()(conn)
	assert.True(t, mysqlAcceptable(nil))
	assert.False(t, mysqlAcceptable(errors.New("any")))
	assert.False(t, mysqlAcceptable(new(mysql.MySQLError)))
}

func tryOnDuplicateEntryError(t *testing.T, accept func(error) bool) error {
	logx.Disable()

	conn := commonSqlConn{
		brk:    breaker.NewBreaker(),
		accept: accept,
	}
	for i := 0; i < 1000; i++ {
		assert.NotNil(t, conn.brk.DoWithAcceptable(func() error {
			return &mysql.MySQLError{
				Number: duplicateEntryCode,
			}
		}, conn.acceptable))
	}
	return conn.brk.DoWithAcceptable(func() error {
		return &mysql.MySQLError{
			Number: duplicateEntryCode,
		}
	}, conn.acceptable)
}
