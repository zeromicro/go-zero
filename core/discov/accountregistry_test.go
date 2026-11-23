package discov

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/stringx"
)

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

func TestRegisterAutoSyncInterval(t *testing.T) {
	endpoints := []string{
		stringx.Rand() + ":2379",
	}

	_, ok := internal.GetAutoSyncInterval(endpoints)
	assert.False(t, ok)

	RegisterAutoSyncInterval(endpoints, time.Minute)
	interval, ok := internal.GetAutoSyncInterval(endpoints)
	assert.True(t, ok)
	assert.Equal(t, time.Minute, interval)

	// a zero interval can be registered explicitly to disable the auto sync.
	RegisterAutoSyncInterval(endpoints, 0)
	interval, ok = internal.GetAutoSyncInterval(endpoints)
	assert.True(t, ok)
	assert.Equal(t, time.Duration(0), interval)
}
