package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type EventHandler struct {
	Status       *status.Status
	writer       io.Writer
	marshaler    jsonpb.Marshaler
	ctx          context.Context
	useOkHandler bool
}

func NewEventHandler(writer io.Writer, resolver jsonpb.AnyResolver) *EventHandler {
	return &EventHandler{
		writer: writer,
		marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			AnyResolver:  resolver,
		},
	}
}

// NewEventHandlerWithContext creates an EventHandler that supports httpx.OkHandler callbacks
func NewEventHandlerWithContext(ctx context.Context, w http.ResponseWriter, resolver jsonpb.AnyResolver, useOkHandler bool) *EventHandler {
	return &EventHandler{
		writer: w,
		marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			AnyResolver:  resolver,
		},
		ctx:          ctx,
		useOkHandler: useOkHandler,
	}
}

func (h *EventHandler) OnReceiveResponse(message proto.Message) {
	if h.useOkHandler {
		// Use httpx.OkJsonCtx to trigger the OkHandler callback
		var buf bytes.Buffer
		if err := h.marshaler.Marshal(&buf, message); err != nil {
			logx.Error(err)
			return
		}

		result := buf.Bytes()
		httpx.OkJsonCtx(h.ctx, h.writer.(http.ResponseWriter), result)
	} else {
		// Fallback to original behavior
		if err := h.marshaler.Marshal(h.writer, message); err != nil {
			logx.Error(err)
		}
	}
}

func (h *EventHandler) OnReceiveTrailers(status *status.Status, _ metadata.MD) {
	h.Status = status
}

func (h *EventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *EventHandler) OnSendHeaders(_ metadata.MD) {
}

func (h *EventHandler) OnReceiveHeaders(_ metadata.MD) {
}
