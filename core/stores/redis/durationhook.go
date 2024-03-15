package redis

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// spanName is the span name of the redis calls.
const spanName = "redis"

var (
	defaultDurationHook   = durationHook{}
	redisCmdsAttributeKey = attribute.Key("redis.cmds")
)

type durationHook struct {
}

func (h durationHook) DialHook(next red.DialHook) red.DialHook {
	return next
}

func (h durationHook) ProcessHook(next red.ProcessHook) red.ProcessHook {
	return func(ctx context.Context, cmd red.Cmder) error {
		start := timex.Now()
		ctx, endSpan := h.startSpan(ctx, cmd)

		err := next(ctx, cmd)

		endSpan(err)
		duration := timex.Since(start)

		if duration > slowThreshold.Load() {
			logDuration(ctx, []red.Cmder{cmd}, duration)
			metricSlowCount.Inc(cmd.Name())
		}

		metricReqDur.Observe(duration.Milliseconds(), cmd.Name())
		if msg := formatError(err); len(msg) > 0 {
			metricReqErr.Inc(cmd.Name(), msg)
		}

		return err
	}
}

func (h durationHook) ProcessPipelineHook(next red.ProcessPipelineHook) red.ProcessPipelineHook {
	return func(ctx context.Context, cmds []red.Cmder) error {
		if len(cmds) == 0 {
			return next(ctx, cmds)
		}

		start := timex.Now()
		ctx, endSpan := h.startSpan(ctx, cmds...)

		err := next(ctx, cmds)

		endSpan(err)
		duration := timex.Since(start)
		if duration > slowThreshold.Load()*time.Duration(len(cmds)) {
			logDuration(ctx, cmds, duration)
		}

		metricReqDur.Observe(duration.Milliseconds(), "Pipeline")
		if msg := formatError(err); len(msg) > 0 {
			metricReqErr.Inc("Pipeline", msg)
		}

		return err
	}
}

func (h durationHook) startSpan(ctx context.Context, cmds ...red.Cmder) (context.Context, func(err error)) {
	tracer := trace.TracerFromContext(ctx)

	ctx, span := tracer.Start(ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)

	cmdStrs := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.Name())
	}
	span.SetAttributes(redisCmdsAttributeKey.StringSlice(cmdStrs))

	return ctx, func(err error) {
		defer span.End()

		if err == nil || errors.Is(err, red.Nil) {
			span.SetStatus(codes.Ok, "")
			return
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

func formatError(err error) string {
	if err == nil || errors.Is(err, red.Nil) {
		return ""
	}

	var opErr *net.OpError
	ok := errors.As(err, &opErr)
	if ok && opErr.Timeout() {
		return "timeout"
	}

	switch {
	case errors.Is(err, io.EOF):
		return "eof"
	case errors.Is(err, context.DeadlineExceeded):
		return "context deadline"
	case errors.Is(err, breaker.ErrServiceUnavailable):
		return "breaker open"
	default:
		return "unexpected error"
	}
}

func logDuration(ctx context.Context, cmds []red.Cmder, duration time.Duration) {
	var buf strings.Builder
	for k, cmd := range cmds {
		if k > 0 {
			buf.WriteByte('\n')
		}
		var build strings.Builder
		for i, arg := range cmd.Args() {
			if i > 0 {
				build.WriteByte(' ')
			}
			build.WriteString(mapping.Repr(arg))
		}
		buf.WriteString(build.String())
	}
	logx.WithContext(ctx).WithDuration(duration).Slowf("[REDIS] slowcall on executing: %s", buf.String())
}
