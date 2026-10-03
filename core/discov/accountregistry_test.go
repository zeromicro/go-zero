package discov

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestRegisterAutoSyncInterval(t *testing.T) {
	endpoints := []string{
		"localhost:2379",
	}

	interval := internal.GetAutoSyncInterval(endpoints)
	assert.Equal(t, time.Duration(0), interval)

	RegisterAutoSyncInterval(endpoints, time.Minute)
	assert.Equal(t, time.Minute, internal.GetAutoSyncInterval(endpoints))
}

func TestRegisterAccount(t *testing.T) {
	endpoints := []string{
		"localhost:2379",
	}
	user := "foo" + stringx.Rand()
	RegisterAccount(endpoints, user, "bar")
	account, ok := internal.GetAccount(endpoints)
	assert.True(t, ok)
	assert.Equal(t, user, account.User)
	assert.Equal(t, "bar", account.Pass)
}
