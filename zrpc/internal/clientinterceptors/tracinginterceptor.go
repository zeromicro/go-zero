package clientinterceptors

import (
	"context"
	"io"

	ztrace "github.com/tal-tech/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	receiveEndEvent streamEventType = iota
	errorEvent
)

// UnaryTracingInterceptor returns a grpc.UnaryClientInterceptor for opentelemetry.
func UnaryTracingInterceptor() grpc.UnaryClientInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(ztrace.TraceName)
		name, attr := ztrace.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		defer span.End()

		ztrace.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		ztrace.MessageSent.Event(ctx, 1, req)
		ztrace.MessageReceived.Event(ctx, 1, reply)

		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
			return err
		}

		span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
		return nil
	}
}

// StreamTracingInterceptor returns a grpc.StreamClientInterceptor for opentelemetry.
func StreamTracingInterceptor() grpc.StreamClientInterceptor {
	propagator := otel.GetTextMapPropagator()
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()
		tr := otel.Tracer(ztrace.TraceName)
		name, attr := ztrace.SpanInfo(method, cc.Target())
		ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...))
		ztrace.Inject(ctx, propagator, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			span.SetStatus(codes.Error, grpcStatus.Message())
			span.SetAttributes(ztrace.StatusCodeAttr(grpcStatus.Code()))
			span.End()
			return s, err
		}

		stream := wrapClientStream(ctx, s, desc)

		go func() {
			if err := <-stream.Finished; err != nil {
				s, _ := status.FromError(err)
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
			} else {
				span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
			}

			span.End()
		}()

		return stream, nil
	}
}

type (
	streamEventType int

	streamEvent struct {
		Type streamEventType
		Err  error
	}

	clientStream struct {
		grpc.ClientStream
		Finished          chan error
		desc              *grpc.StreamDesc
		events            chan streamEvent
		eventsDone        chan struct{}
		receivedMessageID int
		sentMessageID     int
	}
)

func (w *clientStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if err == nil && !w.desc.ServerStreams {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err == io.EOF {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.receivedMessageID++
		ztrace.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *clientStream) SendMsg(m interface{}) error {
	err := w.ClientStream.SendMsg(m)
	w.sentMessageID++
	ztrace.MessageSent.Event(w.Context(), w.sentMessageID, m)
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) Header() (metadata.MD, error) {
	md, err := w.ClientStream.Header()
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return md, err
}

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()
	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) sendStreamEvent(eventType streamEventType, err error) {
	select {
	case <-w.eventsDone:
	case w.events <- streamEvent{Type: eventType, Err: err}:
	}
}

// wrapClientStream wraps s with given ctx and desc.
func wrapClientStream(ctx context.Context, s grpc.ClientStream, desc *grpc.StreamDesc) *clientStream {
	events := make(chan streamEvent)
	eventsDone := make(chan struct{})
	finished := make(chan error)

	go func() {
		defer close(eventsDone)

		for {
			select {
			case event := <-events:
				switch event.Type {
				case receiveEndEvent:
					finished <- nil
					return
				case errorEvent:
					finished <- event.Err
					return
				}
			case <-ctx.Done():
				finished <- ctx.Err()
				return
			}
		}
	}()

	return &clientStream{
		ClientStream: s,
		desc:         desc,
		events:       events,
		eventsDone:   eventsDone,
		Finished:     finished,
	}
}
