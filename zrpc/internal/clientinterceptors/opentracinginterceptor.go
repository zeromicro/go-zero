package clientinterceptors

import (
	"context"

	opentelemetry2 "github.com/tal-tech/go-zero/core/trace/opentelemetry"
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
		if !opentelemetry2.Enabled() {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(opentelemetry2.TraceName)
		name, attr := opentelemetry2.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		defer span.End()

		opentelemetry2.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		opentelemetry2.MessageSent.Event(ctx, 1, req)
		opentelemetry2.MessageReceived.Event(ctx, 1, reply)

		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(opentelemetry2.StatusCodeAttr(s.Code()))
			return err
		}

		span.SetAttributes(opentelemetry2.StatusCodeAttr(gcodes.OK))
		return nil
	}
}

// StreamOpenTracingInterceptor returns a grpc.StreamClientInterceptor for opentelemetry.
func StreamOpenTracingInterceptor() grpc.StreamClientInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if !opentelemetry2.Enabled() {
			return streamer(ctx, desc, cc, method, opts...)
		}

		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(opentelemetry2.TraceName)
		name, attr := opentelemetry2.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		opentelemetry2.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			span.SetStatus(codes.Error, grpcStatus.Message())
			span.SetAttributes(opentelemetry2.StatusCodeAttr(grpcStatus.Code()))
			span.End()
			return s, err
		}

		stream := opentelemetry2.WrapClientStream(ctx, s, desc)

		go func() {
			if err := <-stream.Finished; err != nil {
				s, _ := status.FromError(err)
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(opentelemetry2.StatusCodeAttr(s.Code()))
			} else {
				span.SetAttributes(opentelemetry2.StatusCodeAttr(gcodes.OK))
			}

			span.End()
		}()

		return stream, nil
	}
}
