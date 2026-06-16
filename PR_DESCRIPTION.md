# Support configurable Redis connection settings

## Description
This PR introduces support for configuring `MaxRetries`, `MinIdleConns`, `DB`, and `PoolSize` in Redis connections via `RedisConf`. 

Previously, these values were hardcoded:
- `MaxRetries`: 3
- `MinIdleConns`: 8
- `DB`: 0
- `PoolSize`: `10 * runtime.GOMAXPROCS`

This change allows users to tune these parameters for better performance and resource management, addressing specific deployment needs such as high-latency networks or limited connection limits.

## Changes
- **Configuration**: Modified `RedisConf` in `core/stores/redis/conf.go` to include:
    - `MaxRetries` (default 3)
    - `MinIdleConns` (default 8)
    - `DB` (default 0)
    - `PoolSize` (optional, defaults to dynamic calculation if 0)
- **Core Update**: Updated `Redis` struct in `core/stores/redis/redis.go` to store these configuration values.
- **Client Logic**: Updated `NewRedis` and `newRedis` to pass these values to the underlying `go-redis` client options.
- **Components**: Updated `redisclientmanager.go`, `redisclustermanager.go`, and `redisblockingnode.go` to utilize the configured values from the `Redis` instance.
- **Tests**: Added unit tests in `core/stores/redis/redis_test.go` (`TestRedisOptions`) to verify that configurations are correctly applied.

## Related Issue
Fixes #4668

## How to Test
1. **Unit Tests**:
   Run the newly added unit tests to verify configuration propagation:
   ```bash
   go test -v core/stores/redis/redis_test.go -run TestRedisOptions
   ```

2. **Manual Verification**:
   You can verify the configuration by initializing a Redis client with custom values:
   ```go
   conf := redis.RedisConf{
       Host:         "localhost:6379",
       Type:         redis.NodeType,
       MaxRetries:   5,
       MinIdleConns: 20,
       DB:           1,
       PoolSize:     50,
   }
   r := redis.MustNewRedis(conf)
   // Verify internally or via behavior (e.g. check connection count in redis-cli)
   ```
