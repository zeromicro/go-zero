package redis

import (
	"context"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/lang"
)

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

		return h.brk.DoWithAcceptableCtx(ctx, func() error {
			return next(ctx, cmd)
		}, acceptable)
	}
}

func (h breakerHook) ProcessPipelineHook(next red.ProcessPipelineHook) red.ProcessPipelineHook {
	return func(ctx context.Context, cmds []red.Cmder) error {
		return h.brk.DoWithAcceptableCtx(ctx, func() error {
			return next(ctx, cmds)
		}, acceptable)
	}
}
