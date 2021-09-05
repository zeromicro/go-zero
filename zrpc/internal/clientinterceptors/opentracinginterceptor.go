package clientinterceptors

import (
	"context"

	"github.com/tal-tech/go-zero/core/trace/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// OpenTracingInterceptor returns a grpc.UnaryClientInterceptor for opentelemetry.
func OpenTracingInterceptor() grpc.UnaryClientInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !opentelemetry.Enabled() {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(opentelemetry.TraceName)
		name, attr := opentelemetry.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		defer span.End()

		opentelemetry.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		opentelemetry.MessageSent.Event(ctx, 1, req)
		opentelemetry.MessageReceived.Event(ctx, 1, reply)

		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(opentelemetry.StatusCodeAttr(s.Code()))
			return err
		}

		span.SetAttributes(opentelemetry.StatusCodeAttr(gcodes.OK))
		return nil
	}
}

// StreamOpenTracingInterceptor returns a grpc.StreamClientInterceptor for opentelemetry.
func StreamOpenTracingInterceptor() grpc.StreamClientInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if !opentelemetry.Enabled() {
			return streamer(ctx, desc, cc, method, opts...)
		}

		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(opentelemetry.TraceName)
		name, attr := opentelemetry.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		opentelemetry.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			span.SetStatus(codes.Error, grpcStatus.Message())
			span.SetAttributes(opentelemetry.StatusCodeAttr(grpcStatus.Code()))
			span.End()
			return s, err
		}

		stream := opentelemetry.WrapClientStream(ctx, s, desc)

		go func() {
			if err := <-stream.Finished; err != nil {
				s, _ := status.FromError(err)
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(opentelemetry.StatusCodeAttr(s.Code()))
			} else {
				span.SetAttributes(opentelemetry.StatusCodeAttr(gcodes.OK))
			}

			span.End()
		}()

		return stream, nil
	}
}
