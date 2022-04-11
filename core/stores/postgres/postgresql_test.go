package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestPostgreSql(t *testing.T) {
	assert.NotNil(t, New("postgre", sqlx.DialProvider{}))
}
