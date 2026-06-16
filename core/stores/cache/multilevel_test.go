package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/core/syncx"
)

var errMLTestNotFound = errors.New("not found")

func init() {
	logx.Disable()
}

func TestNewMultiLevelCache(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)

	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)
	assert.NotNil(t, mlc)
}

func TestNewMultiLevelCache_WithOptions(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)

	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound,
		WithLocalExpire(time.Minute),
		WithLocalLimit(500))
	assert.NoError(t, err)
	assert.NotNil(t, mlc)
}

func TestNewMultiLevelCache_InvalidOptions(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)

	// zero/negative values should be ignored and defaults used
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound,
		WithLocalExpire(0),
		WithLocalLimit(0))
	assert.NoError(t, err)
	assert.NotNil(t, mlc)
}

func TestMultiLevelCache_SetAndGet(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	const (
		key   = "test-key"
		value = "test-value"
	)

	// Set a value
	assert.NoError(t, mlc.Set(key, value))

	// Get should return from L1 (in-memory)
	var result string
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, value, result)
}

func TestMultiLevelCache_GetFromRemote(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	const (
		key   = "remote-key"
		value = "remote-value"
	)

	// Set directly in remote (bypassing local)
	assert.NoError(t, remote.Set(key, value))

	// Get should fall through to L2, then promote to L1
	var result string
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, value, result)

	// Second get should come from L1
	var result2 string
	assert.NoError(t, mlc.Get(key, &result2))
	assert.Equal(t, value, result2)
}

func TestMultiLevelCache_GetNotFound(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	var result string
	err = mlc.Get("nonexistent", &result)
	assert.True(t, mlc.IsNotFound(err))
}

func TestMultiLevelCache_Del(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	const key = "del-key"
	assert.NoError(t, mlc.Set(key, "value"))

	// Verify it's there
	var result string
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, "value", result)

	// Delete
	assert.NoError(t, mlc.Del(key))

	// Should be gone from both L1 and L2
	err = mlc.Get(key, &result)
	assert.True(t, mlc.IsNotFound(err))
}

func TestMultiLevelCache_DelMultipleKeys(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	assert.NoError(t, mlc.Set("key1", "val1"))
	assert.NoError(t, mlc.Set("key2", "val2"))

	assert.NoError(t, mlc.Del("key1", "key2"))

	var result string
	assert.True(t, mlc.IsNotFound(mlc.Get("key1", &result)))
	assert.True(t, mlc.IsNotFound(mlc.Get("key2", &result)))
}

func TestMultiLevelCache_DelEmpty(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)
	assert.NoError(t, mlc.Del())
}

func TestMultiLevelCache_SetWithExpire(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	const key = "expire-key"
	assert.NoError(t, mlc.SetWithExpire(key, "value", time.Minute))

	var result string
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, "value", result)
}

func TestMultiLevelCache_Take(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	var queryCount int

	// First Take should trigger the query
	var result string
	err = mlc.Take(&result, "take-key", func(val any) error {
		queryCount++
		*val.(*string) = "db-value"
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "db-value", result)
	assert.Equal(t, 1, queryCount)

	// Second Take should hit L1, no query
	var result2 string
	err = mlc.Take(&result2, "take-key", func(val any) error {
		queryCount++
		*val.(*string) = "should-not-reach"
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "db-value", result2)
	assert.Equal(t, 1, queryCount) // still 1
}

func TestMultiLevelCache_Take_NotFound(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10), WithNotFoundExpiry(time.Second))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	var queryCount int

	// First Take: DB returns not found
	var result string
	err = mlc.Take(&result, "missing-key", func(val any) error {
		queryCount++
		return errMLTestNotFound
	})
	assert.True(t, mlc.IsNotFound(err))
	assert.Equal(t, 1, queryCount)

	// Second Take: should hit L1 not-found placeholder, no query
	err = mlc.Take(&result, "missing-key", func(val any) error {
		queryCount++
		return errMLTestNotFound
	})
	assert.True(t, mlc.IsNotFound(err))
	assert.Equal(t, 1, queryCount) // still 1
}

func TestMultiLevelCache_Take_QueryError(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	errDummy := errors.New("dummy db error")

	var result string
	err = mlc.Take(&result, "err-key", func(val any) error {
		return errDummy
	})
	assert.Equal(t, errDummy, err)
}

func TestMultiLevelCache_TakeWithExpire(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	var queryCount int

	var result string
	err = mlc.TakeWithExpire(&result, "twe-key", func(val any, expire time.Duration) error {
		queryCount++
		*val.(*string) = "twe-value"
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "twe-value", result)
	assert.Equal(t, 1, queryCount)

	// Second call should hit L1
	var result2 string
	err = mlc.TakeWithExpire(&result2, "twe-key", func(val any, expire time.Duration) error {
		queryCount++
		*val.(*string) = "should-not-reach"
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "twe-value", result2)
	assert.Equal(t, 1, queryCount)
}

func TestMultiLevelCache_TakeWithExpire_NotFound(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10), WithNotFoundExpiry(time.Second))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	var queryCount int

	var result string
	err = mlc.TakeWithExpire(&result, "twe-missing", func(val any, expire time.Duration) error {
		queryCount++
		return errMLTestNotFound
	})
	assert.True(t, mlc.IsNotFound(err))
	assert.Equal(t, 1, queryCount)

	// Second call should hit L1 not-found placeholder
	err = mlc.TakeWithExpire(&result, "twe-missing", func(val any, expire time.Duration) error {
		queryCount++
		return errMLTestNotFound
	})
	assert.True(t, mlc.IsNotFound(err))
	assert.Equal(t, 1, queryCount)
}

func TestMultiLevelCache_IsNotFound(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	assert.True(t, mlc.IsNotFound(errMLTestNotFound))
	assert.False(t, mlc.IsNotFound(errors.New("other")))
	assert.False(t, mlc.IsNotFound(nil))
}

func TestMultiLevelCache_SetAndDelInvalidatesLocal(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound)
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	const key = "inv-key"

	// Set and verify
	assert.NoError(t, mlc.Set(key, "v1"))
	var result string
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, "v1", result)

	// Update
	assert.NoError(t, mlc.Set(key, "v2"))
	assert.NoError(t, mlc.Get(key, &result))
	assert.Equal(t, "v2", result)

	// Delete and verify gone
	assert.NoError(t, mlc.Del(key))
	assert.True(t, mlc.IsNotFound(mlc.Get(key, &result)))
}

func TestMultiLevelCache_StructValue(t *testing.T) {
	store := redistest.CreateRedis(t)
	remote := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errMLTestNotFound,
		WithExpiry(time.Second*10))
	mlc, err := NewMultiLevelCache(remote, errMLTestNotFound)
	assert.NoError(t, err)

	type User struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	expected := User{ID: 42, Name: "test-user"}

	// Take with struct value
	var result User
	err = mlc.Take(&result, "user:42", func(val any) error {
		*val.(*User) = expected
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)

	// Get from L1
	var cached User
	assert.NoError(t, mlc.Get("user:42", &cached))
	assert.Equal(t, expected.ID, cached.ID)
	assert.Equal(t, expected.Name, cached.Name)
}
