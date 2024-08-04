package internal

import (
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	EventHandler struct {
		Status    *status.Status
		writer    io.Writer
		marshaler jsonpb.Marshaler

		Message     proto.Message
		RespHandler func(writer io.Writer, status *status.Status, message proto.Message)
	}

	HandlerOption func(handler *EventHandler)
)

func NewEventHandler(writer io.Writer, resolver jsonpb.AnyResolver, opts ...HandlerOption) *EventHandler {
	handler := &EventHandler{
		writer: writer,
		marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			AnyResolver:  resolver,
		},
	}
	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

func (h *EventHandler) OnReceiveResponse(message proto.Message) {

	if h.RespHandler != nil {
		h.Message = message
		return
	}

	if err := h.marshaler.Marshal(h.writer, message); err != nil {
		logx.Error(err)
	}
}

func (h *EventHandler) OnReceiveTrailers(status *status.Status, _ metadata.MD) {
	h.Status = status

	if h.RespHandler != nil {
		h.RespHandler(h.writer, h.Status, h.Message)
	}
}

func (h *EventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *EventHandler) OnSendHeaders(_ metadata.MD) {
}

func (h *EventHandler) OnReceiveHeaders(_ metadata.MD) {
}
