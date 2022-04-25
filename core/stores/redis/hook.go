package redis

import (
	"context"
	"strings"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	tracestd "go.opentelemetry.io/otel/trace"
)

// spanName is the span name of the redis calls.
const spanName = "redis"

var (
	startTimeKey = contextKey("startTime")
	durationHook = hook{tracer: otel.GetTracerProvider().Tracer(trace.TraceName)}
)

type (
	contextKey string
	hook       struct {
		tracer tracestd.Tracer
	}
)

func (h hook) BeforeProcess(ctx context.Context, _ red.Cmder) (context.Context, error) {
	return h.startSpan(context.WithValue(ctx, startTimeKey, timex.Now())), nil
}

func (h hook) AfterProcess(ctx context.Context, cmd red.Cmder) error {
	h.endSpan(ctx)

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
	return h.startSpan(context.WithValue(ctx, startTimeKey, timex.Now())), nil
}

func (h hook) AfterProcessPipeline(ctx context.Context, cmds []red.Cmder) error {
	h.endSpan(ctx)

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

func (h hook) startSpan(ctx context.Context) context.Context {
	ctx, _ = h.tracer.Start(ctx, spanName)
	return ctx
}

func (h hook) endSpan(ctx context.Context) {
	tracestd.SpanFromContext(ctx).End()
}
