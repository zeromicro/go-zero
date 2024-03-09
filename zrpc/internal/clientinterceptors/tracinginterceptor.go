package clientinterceptors

import (
	"context"
	"errors"
	"io"

	ztrace "github.com/zeromicro/go-zero/core/trace"
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
func UnaryTracingInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx, span := startSpan(ctx, method, cc.Target())
	defer span.End()

	ztrace.MessageSent.Event(ctx, 1, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	ztrace.MessageReceived.Event(ctx, 1, reply)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	}

	span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
	return nil
}

// StreamTracingInterceptor returns a grpc.StreamClientInterceptor for opentelemetry.
func StreamTracingInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx, span := startSpan(ctx, method, cc.Target())
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(ztrace.StatusCodeAttr(st.Code()))
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		return s, err
	}

	stream := wrapClientStream(ctx, s, desc)

	go func() {
		if err := <-stream.Finished; err != nil {
			s, ok := status.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.Message())
				span.SetAttributes(ztrace.StatusCodeAttr(s.Code()))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
		} else {
			span.SetAttributes(ztrace.StatusCodeAttr(gcodes.OK))
		}

		span.End()
	}()

	return stream, nil
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

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()
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

func (w *clientStream) RecvMsg(m any) error {
	err := w.ClientStream.RecvMsg(m)
	if err == nil && !w.desc.ServerStreams {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if errors.Is(err, io.EOF) {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.receivedMessageID++
		ztrace.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *clientStream) SendMsg(m any) error {
	err := w.ClientStream.SendMsg(m)
	w.sentMessageID++
	ztrace.MessageSent.Event(w.Context(), w.sentMessageID, m)
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

func startSpan(ctx context.Context, method, target string) (context.Context, trace.Span) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	tr := otel.Tracer(ztrace.TraceName)
	name, attr := ztrace.SpanInfo(method, target)
	ctx, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...))
	ztrace.Inject(ctx, otel.GetTextMapPropagator(), &md)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx, span
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
