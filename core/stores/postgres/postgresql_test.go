package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestPostgreSql(t *testing.T) {
	assert.NotNil(t, New("postgre", sqlx.PoolConfig{
		MaxIdleConns: 10,
		MaxOpenConns: 10,
		MaxLifetime:  time.Minute,
	}))
}
