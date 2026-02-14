# Multi-Level Cache (sqlc + collection.Cache)

## Problem

Using Redis (`sqlc`) alone for SQL caching incurs a network round-trip on every read. For hot keys
this overhead adds up. `collection.Cache` is an in-process LRU cache with zero network latency, but
on its own it does not survive process restarts and cannot be shared across instances.

## Solution

`cache.NewMultiLevelCache` layers the two together:

```
┌──────────┐  miss   ┌──────────┐  miss   ┌──────┐
│ L1: in-  │ ──────> │ L2: Redis│ ──────> │  DB  │
│ memory   │ <────── │ cache    │ <────── │      │
│ (fast)   │ promote │ (shared) │  query  │      │
└──────────┘         └──────────┘         └──────┘
```

- **Read**: L1 is checked first (sub-microsecond). On an L1 miss, L2 (Redis) is queried. If found
  in L2, the value is promoted into L1 for subsequent requests.
- **Write / Delete**: Both L1 and L2 are updated to keep them consistent.
- **Not-found caching**: When the DB returns not-found, a placeholder is stored in both layers so
  repeated lookups for missing rows do not hit the DB again.

## Quick Start

### Option A — Using sqlc convenience constructors

```go
import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlc"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

db := sqlx.NewMysql(datasource)

// With a single Redis node:
conn, err := sqlc.NewNodeConnWithMultiLevelCache(db, redisClient,
    []cache.MultiLevelCacheOption{
        cache.WithLocalExpire(time.Minute * 5), // L1 TTL (default: 5 min)
        cache.WithLocalLimit(10000),             // L1 max entries (default: 10000)
    },
    cache.WithExpiry(time.Hour*24*7), // L2 (Redis) TTL
)

// With a Redis cluster:
conn, err := sqlc.NewConnWithMultiLevelCache(db, cacheConf,
    []cache.MultiLevelCacheOption{
        cache.WithLocalExpire(time.Minute * 3),
        cache.WithLocalLimit(5000),
    },
    cache.WithExpiry(time.Hour*24*7),
)
```

### Option B — Using the cache layer directly

```go
import (
    "database/sql"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlc"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/core/syncx"
)

// Build the remote (L2) cache as usual
remote := cache.NewNode(redisClient, syncx.NewSingleFlight(),
    cache.NewStat("sqlc"), sql.ErrNoRows,
    cache.WithExpiry(time.Hour*24*7))

// Wrap it with a multi-level cache
mlc, err := cache.NewMultiLevelCache(remote, sql.ErrNoRows,
    cache.WithLocalExpire(time.Minute*5),
    cache.WithLocalLimit(10000))
if err != nil {
    log.Fatal(err)
}

// Use it with sqlc
conn := sqlc.NewConnWithCache(db, mlc)
```

### Using it for queries

Once constructed, use `conn` exactly like a regular `sqlc.CachedConn`:

```go
var user User
err := conn.QueryRowCtx(ctx, &user, userCacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
    return conn.QueryRowCtx(ctx, v, "SELECT id, name FROM user WHERE id = ?", userID)
})
```

The first call hits Redis → DB, subsequent calls hit the in-memory L1 directly.

## Configuration Options

**Multi-level cache options** (`cache.MultiLevelCacheOption`):

- `cache.WithLocalExpire(d time.Duration)` — L1 in-memory TTL (default: 5 minutes)
- `cache.WithLocalLimit(n int)` — L1 max entry count with LRU eviction (default: 10000)

**Standard cache options** (`cache.Option`) are passed through to the L2 (Redis) layer:

- `cache.WithExpiry(d time.Duration)` — L2 TTL (default: 7 days)
- `cache.WithNotFoundExpiry(d time.Duration)` — TTL for not-found placeholders (default: 1 minute)

## When to Use

| Scenario | Recommendation |
|---|---|
| Hot keys read thousands of times/sec | ✅ Use multi-level cache |
| Mostly write-heavy workload | ❌ Stick with Redis-only |
| Single-instance deployment | ✅ Great fit — L1 stays coherent |
| Multi-instance deployment | ⚠️ Use short L1 TTL (seconds) to limit staleness |

## Consistency Notes

- L1 is **per-process** — different instances may briefly serve stale data until L1 expires.
- Set `WithLocalExpire` to a short duration (e.g. 5–30 seconds) in multi-instance deployments.
- `Del`/`Exec` invalidate both L1 and L2 within the same process immediately.
