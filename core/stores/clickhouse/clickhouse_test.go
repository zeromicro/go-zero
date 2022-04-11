package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestClickHouse(t *testing.T) {
	assert.NotNil(t, New("clickhouse", sqlx.DialProvider{}))
}
