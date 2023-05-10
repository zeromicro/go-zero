package sqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUp(t *testing.T) {
	SetUp(LogConf{
		DisableSqlLog:  false,
		DisableStmtLog: false,
	})
	defer func() {
		logSql.Set(true)
		logSlowSql.Set(true)
	}()
	assert.True(t, logSql.True())
	assert.True(t, logSlowSql.True())

	SetUp(LogConf{
		DisableSqlLog:  true,
		DisableStmtLog: true,
	})
	assert.False(t, logSql.True())
	assert.False(t, logSlowSql.True())
}
