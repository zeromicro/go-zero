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

const (
	// MetadataHeaderPrefix is the http prefix that represents custom metadata
	// parameters to or from a gRPC call.
	MetadataHeaderPrefix = "Grpc-Metadata-"

	// MetadataTrailerPrefix is prepended to gRPC metadata as it is converted to
	// HTTP headers in a response handled by go-zero gateway
	MetadataTrailerPrefix = "Grpc-Trailer-"
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

func (h *EventHandler) OnReceiveHeaders(md metadata.MD) {
	w, ok := h.writer.(http.ResponseWriter)
	if ok {
		for k, vs := range md {
			header := defaultOutgoingHeaderMatcher(k)
			for _, v := range vs {
				w.Header().Add(header, v)
			}
		}
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

func (h *EventHandler) OnReceiveTrailers(status *status.Status, md metadata.MD) {
	w, ok := h.writer.(http.ResponseWriter)
	if ok {
		for k, vs := range md {
			header := defaultOutgoingTrailerMatcher(k)
			for _, v := range vs {
				w.Header().Add(header, v)
			}
		}
	}

	h.Status = status
}

func (h *EventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *EventHandler) OnSendHeaders(_ metadata.MD) {
}

func defaultOutgoingHeaderMatcher(key string) string {
	return MetadataHeaderPrefix + key
}

func defaultOutgoingTrailerMatcher(key string) string {
	return MetadataTrailerPrefix + key
}
