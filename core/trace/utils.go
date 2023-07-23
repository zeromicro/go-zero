package trace

import (
	"context"
	"net"
	"strings"

	ztrace "github.com/zeromicro/go-zero/internal/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"
)

const localhost = "127.0.0.1"

var (
	// SpanIDFromContext returns the span id from ctx.
	SpanIDFromContext = ztrace.SpanIDFromContext
	// TraceIDFromContext returns the trace id from ctx.
	TraceIDFromContext = ztrace.TraceIDFromContext
)

// ParseFullMethod returns the method name and attributes.
func ParseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue(nil)
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}

	return name, attrs
}

// PeerAttr returns the peer attributes.
func PeerAttr(addr string) []attribute.KeyValue {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	if len(host) == 0 {
		host = localhost
	}

	return []attribute.KeyValue{
		semconv.NetPeerIPKey.String(host),
		semconv.NetPeerPortKey.String(port),
	}
}

// PeerFromCtx returns the peer from ctx.
func PeerFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return ""
	}

	return p.Addr.String()
}

// SpanInfo returns the span info.
func SpanInfo(fullMethod, peerAddress string) (string, []attribute.KeyValue) {
	attrs := []attribute.KeyValue{RPCSystemGRPC}
	name, mAttrs := ParseFullMethod(fullMethod)
	attrs = append(attrs, mAttrs...)
	attrs = append(attrs, PeerAttr(peerAddress)...)
	return name, attrs
}

// TracerFromContext returns a tracer in ctx, otherwise returns a global tracer.
func TracerFromContext(ctx context.Context) (tracer trace.Tracer) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		tracer = span.TracerProvider().Tracer(TraceName)
	} else {
		tracer = otel.Tracer(TraceName)
	}

	return
}
