package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	red "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
)

func TestBreakerHook_ProcessHook(t *testing.T) {
	t.Run("breakerHookOpen", func(t *testing.T) {
		s := miniredis.RunT(t)

		rds := MustNewRedis(RedisConf{
			Host: s.Addr(),
			Type: NodeType,
		})

		someError := errors.New("ERR some error")
		s.SetError(someError.Error())

		var err error
		for i := 0; i < 1000; i++ {
			_, err = rds.Get("key")
			if err != nil && err.Error() != someError.Error() {
				break
			}
		}
		assert.Equal(t, breaker.ErrServiceUnavailable, err)
	})

	t.Run("breakerHookClose", func(t *testing.T) {
		s := miniredis.RunT(t)

		rds := MustNewRedis(RedisConf{
			Host: s.Addr(),
			Type: NodeType,
		})

		var err error
		for i := 0; i < 1000; i++ {
			_, err = rds.Get("key")
			if err != nil {
				break
			}
		}
		assert.NotEqual(t, breaker.ErrServiceUnavailable, err)
	})

	t.Run("breakerHook_ignoreCmd", func(t *testing.T) {
		s := miniredis.RunT(t)

		rds := MustNewRedis(RedisConf{
			Host: s.Addr(),
			Type: NodeType,
		})

		someError := errors.New("ERR some error")
		s.SetError(someError.Error())

		var err error

		node, err := getRedis(rds)
		assert.NoError(t, err)

		for i := 0; i < 1000; i++ {
			_, err = rds.Blpop(node, "key")
			if err != nil && err.Error() != someError.Error() {
				break
			}
		}
		assert.Equal(t, someError.Error(), err.Error())
	})

	t.Run("breakerHook_ignoreHello", func(t *testing.T) {
		// hello is issued on connection init and is in ignoreCmds, so repeated
		// failures must never trip the breaker into ErrServiceUnavailable.
		h := breakerHook{brk: breaker.NewBreaker()}
		someError := errors.New("ERR some error")
		process := h.ProcessHook(func(_ context.Context, _ red.Cmder) error {
			return someError
		})

		ctx := context.Background()
		var err error
		for i := 0; i < 1000; i++ {
			err = process(ctx, red.NewCmd(ctx, "hello", 3))
			if err != nil && err.Error() != someError.Error() {
				break
			}
		}
		assert.Equal(t, someError.Error(), err.Error())
	})

	t.Run("breakerHook_notIgnored", func(t *testing.T) {
		// a regular command is not ignored, so repeated failures open the breaker.
		h := breakerHook{brk: breaker.NewBreaker()}
		someError := errors.New("ERR some error")
		process := h.ProcessHook(func(_ context.Context, _ red.Cmder) error {
			return someError
		})

		ctx := context.Background()
		var err error
		for i := 0; i < 1000; i++ {
			err = process(ctx, red.NewCmd(ctx, "get", "key"))
			if err != nil && err.Error() != someError.Error() {
				break
			}
		}
		assert.Equal(t, breaker.ErrServiceUnavailable, err)
	})
}

func TestBreakerHook_ProcessPipelineHook(t *testing.T) {
	t.Run("breakerPipelineHookOpen", func(t *testing.T) {
		s := miniredis.RunT(t)

		rds := MustNewRedis(RedisConf{
			Host: s.Addr(),
			Type: NodeType,
		})

		someError := errors.New("ERR some error")
		s.SetError(someError.Error())

		var err error
		for i := 0; i < 1000; i++ {
			err = rds.Pipelined(
				func(pipe Pipeliner) error {
					pipe.Incr(context.Background(), "pipelined_counter")
					pipe.Expire(context.Background(), "pipelined_counter", time.Hour)
					pipe.ZAdd(context.Background(), "zadd", Z{Score: 12, Member: "zadd"})
					return nil
				},
			)

			if err != nil && err.Error() != someError.Error() {
				break
			}
		}
		assert.Equal(t, breaker.ErrServiceUnavailable, err)
	})

	t.Run("breakerPipelineHookClose", func(t *testing.T) {
		s := miniredis.RunT(t)

		rds := MustNewRedis(RedisConf{
			Host: s.Addr(),
			Type: NodeType,
		})

		var err error
		for i := 0; i < 1000; i++ {
			err = rds.Pipelined(
				func(pipe Pipeliner) error {
					pipe.Incr(context.Background(), "pipelined_counter")
					pipe.Expire(context.Background(), "pipelined_counter", time.Hour)
					pipe.ZAdd(context.Background(), "zadd", Z{Score: 12, Member: "zadd"})
					return nil
				},
			)

			if err != nil {
				break
			}
		}
		assert.NotEqual(t, breaker.ErrServiceUnavailable, err)
	})
}
