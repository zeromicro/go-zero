package kafka

import "github.com/zeromicro/go-zero/kafka/internal/types"

// message and header types are same as https://github.com/segmentio/kafka-go/blob/main/message.go

type (
	// Message is a data structure representing kafka messages.
	Message = types.Message

	// Header represents a single entry in a list of record headers.
	Header = types.Header
)
