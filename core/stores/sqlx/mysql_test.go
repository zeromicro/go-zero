package sqlx

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/breaker"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stat"
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
		if tryOnDuplicateEntryError(t, nil) == breaker.ErrServiceUnavailable {
			found = true
		}
	}
	assert.True(t, found)
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
