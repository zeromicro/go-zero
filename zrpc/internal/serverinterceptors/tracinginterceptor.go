package serverinterceptors

import (
	"context"

	ztrace "github.com/tal-tech/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryTracingInterceptor returns a grpc.UnaryServerInterceptor for opentelemetry.
func UnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		var md metadata.MD
		requestMetadata, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = requestMetadata.Copy()
		}
		bags, spanCtx := ztrace.Extract(ctx, propagator, &md)
		ctx = baggage.ContextWithBaggage(ctx, bags)
		tr := otel.Tracer(ztrace.TraceName)
		name, attr := ztrace.SpanInfo(info.FullMethod, ztrace.PeerFromCtx(ctx))
		ctx, span := tr.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), name,
			trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attr...))
		defer span.End()

		ztrace.MessageReceived.Event(ctx, 1, req)
		resp, err := handler(ctx, req)
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.Message())
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
			span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
			ztrace.MessageSent.Event(ctx, 1, s.Proto())
			return nil, err
		}

		span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
		ztrace.MessageSent.Event(ctx, 1, resp)

		return resp, nil
	}
}

// StreamTracingInterceptor returns a grpc.StreamServerInterceptor for opentelemetry.
func StreamTracingInterceptor() grpc.StreamServerInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var md metadata.MD
		ctx := ss.Context()
		requestMetadata, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = requestMetadata.Copy()
		}
		bags, spanCtx := ztrace.Extract(ctx, propagator, &md)
		ctx = baggage.ContextWithBaggage(ctx, bags)
		tr := otel.Tracer(ztrace.TraceName)
		name, attr := ztrace.SpanInfo(info.FullMethod, ztrace.PeerFromCtx(ctx))
		ctx, span := tr.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), name,
			trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attr...))
		defer span.End()

		if err := handler(srv, wrapServerStream(ctx, ss)); err != nil {
			s, ok := status.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.Message())
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
			span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
			return err
		}

		span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
		return nil
	}
}

// serverStream wraps around the embedded grpc.ServerStream,
// and intercepts the RecvMsg and SendMsg method call.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context

	receivedMessageID int
	sentMessageID     int
}

func (w *serverStream) Context() context.Context {
	return w.ctx
}

func (w *serverStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.receivedMessageID++
		ztrace.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *serverStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	w.sentMessageID++
	ztrace.MessageSent.Event(w.Context(), w.sentMessageID, m)

	return err
}

// wrapServerStream wraps the given grpc.ServerStream with the given context.
func wrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}
