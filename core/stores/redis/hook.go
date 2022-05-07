package redis

import (
	"context"
	"strings"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// spanName is the span name of the redis calls.
const spanName = "redis"

var (
	startTimeKey          = contextKey("startTime")
	durationHook          = hook{tracer: otel.GetTracerProvider().Tracer(trace.TraceName)}
	redisCmdsAttributeKey = attribute.Key("redis.cmds")
)

type (
	contextKey string
	hook       struct {
		tracer oteltrace.Tracer
	}
)

func (h hook) BeforeProcess(ctx context.Context, cmd red.Cmder) (context.Context, error) {
	return h.startSpan(context.WithValue(ctx, startTimeKey, timex.Now()), cmd), nil
}

func (h hook) AfterProcess(ctx context.Context, cmd red.Cmder) error {
	err := cmd.Err()
	h.endSpan(ctx, err)

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

func (h hook) BeforeProcessPipeline(ctx context.Context, cmds []red.Cmder) (context.Context, error) {
	if len(cmds) == 0 {
		return ctx, nil
	}

	return h.startSpan(context.WithValue(ctx, startTimeKey, timex.Now()), cmds...), nil
}

func (h hook) AfterProcessPipeline(ctx context.Context, cmds []red.Cmder) error {
	if len(cmds) == 0 {
		return nil
	}

	batchError := errorx.BatchError{}
	for _, cmd := range cmds {
		err := cmd.Err()
		if err == nil {
			continue
		}

		batchError.Add(err)
	}
	h.endSpan(ctx, batchError.Err())

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

func (h hook) startSpan(ctx context.Context, cmds ...red.Cmder) context.Context {
	ctx, span := h.tracer.Start(ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)

	cmdStrs := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.Name())
	}
	span.SetAttributes(redisCmdsAttributeKey.StringSlice(cmdStrs))

	return ctx
}

func (h hook) endSpan(ctx context.Context, err error) {
	span := oteltrace.SpanFromContext(ctx)
	defer span.End()

	if err == nil || err == red.Nil {
		span.SetStatus(codes.Ok, "")
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
