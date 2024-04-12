package redis

import (
	"context"
	"errors"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/timex"
)

const minTimeout = time.Millisecond * 100

var ignoreCmds = map[string]lang.PlaceholderType{
	"blpop": {},
}

type breakerHook struct {
	brk breaker.Breaker
}

func (h breakerHook) DialHook(next red.DialHook) red.DialHook {
	return next
}

func (h breakerHook) ProcessHook(next red.ProcessHook) red.ProcessHook {
	return func(ctx context.Context, cmd red.Cmder) error {
		if _, ok := ignoreCmds[cmd.Name()]; ok {
			return next(ctx, cmd)
		}

		start := timex.Now()
		return h.brk.DoWithAcceptable(func() error {
			return next(ctx, cmd)
		}, protectedAcceptable(start))
	}
}

func (h breakerHook) ProcessPipelineHook(next red.ProcessPipelineHook) red.ProcessPipelineHook {
	return func(ctx context.Context, cmds []red.Cmder) error {
		start := timex.Now()
		return h.brk.DoWithAcceptable(func() error {
			return next(ctx, cmds)
		}, protectedAcceptable(start))
	}
}

func acceptable(err error) bool {
	return err == nil || errors.Is(err, red.Nil) || errors.Is(err, context.Canceled)
}

func protectedAcceptable(start time.Duration) breaker.Acceptable {
	return func(err error) bool {
		if acceptable(err) {
			return true
		}

		return errors.Is(err, context.DeadlineExceeded) && timex.Since(start) < minTimeout
	}
}
