package clickhouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClickHouse(t *testing.T) {
	assert.NotNil(t, New("clickhouse"))
}
