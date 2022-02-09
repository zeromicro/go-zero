package redis

import (
	"context"
	"strings"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/timex"
)

var (
	startTimeKey = contextKey("startTime")
	durationHook = hook{}
)

type (
	contextKey string
	hook       struct{}
)

func (h hook) BeforeProcess(ctx context.Context, _ red.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, timex.Now()), nil
}

func (h hook) AfterProcess(ctx context.Context, cmd red.Cmder) error {
	val := ctx.Value(startTimeKey)
	if val == nil {
		return nil
	}

	start, ok := val.(time.Duration)
	if !ok {
		return nil
	}

	duration := timex.Since(start)
	if duration > slowThreshold.Load() {
		logDuration(ctx, cmd, duration)
	}

	return nil
}

func (h hook) BeforeProcessPipeline(ctx context.Context, _ []red.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTimeKey, timex.Now()), nil
}

func (h hook) AfterProcessPipeline(ctx context.Context, cmds []red.Cmder) error {
	if len(cmds) == 0 {
		return nil
	}

	val := ctx.Value(startTimeKey)
	if val == nil {
		return nil
	}

	start, ok := val.(time.Duration)
	if !ok {
		return nil
	}

	duration := timex.Since(start)
	if duration > slowThreshold.Load()*time.Duration(len(cmds)) {
		logDuration(ctx, cmds[0], duration)
	}

	return nil
}

func logDuration(ctx context.Context, cmd red.Cmder, duration time.Duration) {
	var buf strings.Builder
	for i, arg := range cmd.Args() {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(mapping.Repr(arg))
	}
	logx.WithContext(ctx).WithDuration(duration).Slowf("[REDIS] slowcall on executing: %s", buf.String())
}
