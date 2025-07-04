package sqlx

import "context"

type (
	readWriteMode    string
	readWriteModeKey struct{}
)

const (
	// PolicyRoundRobin round-robin policy for selecting replicas.
	PolicyRoundRobin = "round-robin"
	// PolicyRandom random policy for selecting replicas.
	PolicyRandom = "random"

	// readMode indicates that the operation is a read operation.
	readMode readWriteMode = "read"
	// writeMode indicates that the operation is a write operation.
	writeMode readWriteMode = "write"
	// notSpecifiedMode indicates that the read/write mode is not specified.
	notSpecifiedMode readWriteMode = ""
)

func (m readWriteMode) isValid() bool {
	return m == readMode || m == writeMode
}

// WithReadMode sets the context to read mode, indicating that the operation is a read operation.
func WithReadMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, readWriteModeKey{}, readMode)
}

// WithWriteMode sets the context to write mode, indicating that the operation is a write operation.
func WithWriteMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, readWriteModeKey{}, writeMode)
}

func getReadWriteMode(ctx context.Context) readWriteMode {
	if mode := ctx.Value(readWriteModeKey{}); mode != nil {
		if v, ok := mode.(readWriteMode); ok && v.isValid() {
			return v
		}
	}

	return notSpecifiedMode
}

func isReadonly(ctx context.Context) bool {
	return getReadWriteMode(ctx) == readMode
}
