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
	assert.False(t, logSql.True())
	assert.False(t, logSlowSql.True())

	SetUp(LogConf{
		DisableSqlLog:  true,
		DisableStmtLog: true,
	})
	assert.True(t, logSql.True())
	assert.True(t, logSlowSql.True())
}
