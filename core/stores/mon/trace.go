package mon

import (
	"context"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/trace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var mongoCmdAttributeKey = attribute.Key("mongo.cmd")

func startSpan(ctx context.Context, cmd string) (context.Context, oteltrace.Span) {
	tracer := trace.TracerFromContext(ctx)
	ctx, span := tracer.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	span.SetAttributes(mongoCmdAttributeKey.String(cmd))

	return ctx, span
}

func endSpan(span oteltrace.Span, err error) {
	defer span.End()

	if err == nil || errorx.In(err, mongo.ErrNoDocuments, mongo.ErrNilValue, mongo.ErrNilDocument) {
		span.SetStatus(codes.Ok, "")
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
