package internal

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MetadataHeaderPrefix is the http prefix that represents custom metadata
// parameters to or from a gRPC call.
const MetadataHeaderPrefix = "Grpc-Metadata-"

// MetadataTrailerPrefix is prepended to gRPC metadata as it is converted to
// HTTP headers in a response handled by go-zero gateway
const MetadataTrailerPrefix = "Grpc-Trailer-"

func defaultOutgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", MetadataHeaderPrefix, key), true
}

func defaultOutgoingTrailerMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", MetadataTrailerPrefix, key), true
}

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

func (h *EventHandler) OnReceiveHeaders(md metadata.MD) {
	w, ok := h.writer.(http.ResponseWriter)
	if ok {
		for k, vs := range md {
			if h, ok := defaultOutgoingHeaderMatcher(k); ok {
				for _, v := range vs {
					w.Header().Add(h, v)
				}
			}
		}
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
		for k, vs := range md {
			if h, ok := defaultOutgoingTrailerMatcher(k); ok {
				for _, v := range vs {
					w.Header().Add(h, v)
				}
			}
		}
	}

	h.Status = status
}

func (h *EventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *EventHandler) OnSendHeaders(_ metadata.MD) {
}
