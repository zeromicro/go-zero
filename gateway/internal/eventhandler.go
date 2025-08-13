package internal

import (
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type EventHandler struct {
	Status    *status.Status
	writer    io.Writer
	marshaler jsonpb.Marshaler
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

func (h *EventHandler) OnReceiveResponse(message proto.Message) {
	if err := h.marshaler.Marshal(h.writer, message); err != nil {
		logx.Error(err)
	}
}

func (h *EventHandler) OnReceiveTrailers(status *status.Status, md metadata.MD) {
	w, ok := h.writer.(http.ResponseWriter)
	if ok {
		for k, v := range md {
			for _, v2 := range v {
				w.Header().Add(k, v2)
			}
		}
	}
	h.Status = status
}
func (h *EventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *EventHandler) OnSendHeaders(_ metadata.MD) {
}

func (h *EventHandler) OnReceiveHeaders(md metadata.MD) {
	w, ok := h.writer.(http.ResponseWriter)
	if ok {
		for k, v := range md {
			for _, v2 := range v {
				w.Header().Add(k, v2)
			}
		}
	}
}
