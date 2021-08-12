package serverinterceptors

import (
	"context"

	"github.com/tal-tech/go-zero/core/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryOpenTracingInterceptor() grpc.UnaryServerInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !opentelemetry.Enabled() {
			return handler(ctx, req)
		}

		requestMetadata, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := requestMetadata.Copy()

		bags, spanCtx := opentelemetry.Extract(ctx, propagator, &metadataCopy)
		ctx = baggage.ContextWithBaggage(ctx, bags)

		tr := otel.Tracer(opentelemetry.TraceName)
		name, attr := opentelemetry.SpanInfo(info.FullMethod, opentelemetry.PeerFromCtx(ctx))

		var span trace.Span
		ctx, span = tr.Start(
			trace.ContextWithRemoteSpanContext(ctx, spanCtx),
			name,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attr...),
		)
		defer span.End()

		opentelemetry.MessageReceived.Event(ctx, 1, req)

		resp, err := handler(ctx, req)

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(opentelemetry.StatusCodeAttr(s.Code()))
			opentelemetry.MessageSent.Event(ctx, 1, s.Proto())
		} else {
			span.SetAttributes(opentelemetry.StatusCodeAttr(grpc_codes.OK))
			opentelemetry.MessageSent.Event(ctx, 1, resp)
		}

		return resp, err
	}
}
