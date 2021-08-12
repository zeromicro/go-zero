package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	grpc_codes "google.golang.org/grpc/codes"
)

const (
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")

	// Name of message transmitted or received.
	RPCNameKey = attribute.Key("name")

	// Type of message transmitted or received.
	RPCMessageTypeKey = attribute.Key("message.type")

	// Identifier of message transmitted or received.
	RPCMessageIDKey = attribute.Key("message.id")

	// The compressed size of the message transmitted or received in bytes.
	RPCMessageCompressedSizeKey = attribute.Key("message.compressed_size")

	// The uncompressed size of the message transmitted or received in
	// bytes.
	RPCMessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")
)

// Semantic conventions for common RPC attributes.
var (
	// Semantic convention for gRPC as the remoting system.
	RPCSystemGRPC = semconv.RPCSystemKey.String("grpc")

	// Semantic convention for a message named message.
	RPCNameMessage = RPCNameKey.String("message")

	// Semantic conventions for RPC message types.
	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

func StatusCodeAttr(c grpc_codes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(c))
}
