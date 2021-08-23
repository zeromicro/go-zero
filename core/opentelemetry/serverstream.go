package opentelemetry

import (
	"context"

	"google.golang.org/grpc"
)

// serverStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
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
		MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *serverStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)

	w.sentMessageID++
	MessageSent.Event(w.Context(), w.sentMessageID, m)

	return err
}

func WrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}
