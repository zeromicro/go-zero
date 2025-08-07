package sqlx

import "context"

const (
	// policyRoundRobin round-robin policy for selecting replicas.
	policyRoundRobin = "round-robin"
	// policyRandom random policy for selecting replicas.
	policyRandom = "random"

	// readPrimaryMode indicates that the operation is a read,
	// but should be performed on the primary database instance.
	//
	// This mode is used in scenarios where data freshness and consistency are critical,
	// such as immediately after writes or where replication lag may cause stale reads.
	readPrimaryMode readWriteMode = "read-primary"

	// readReplicaMode indicates that the operation is a read from replicas.
	// This is suitable for scenarios where eventual consistency is acceptable,
	// and the goal is to offload traffic from the primary and improve read scalability.
	readReplicaMode readWriteMode = "read-replica"

	// writeMode indicates that the operation is a write operation (to primary).
	writeMode readWriteMode = "write"

	// notSpecifiedMode indicates that the read/write mode is not specified.
	notSpecifiedMode readWriteMode = ""
)

type readWriteModeKey struct{}

// WithReadPrimary sets the context to read-primary mode.
func WithReadPrimary(ctx context.Context) context.Context {
	return context.WithValue(ctx, readWriteModeKey{}, readPrimaryMode)
}

// WithReadReplica sets the context to read-replica mode.
func WithReadReplica(ctx context.Context) context.Context {
	return context.WithValue(ctx, readWriteModeKey{}, readReplicaMode)
}

// WithWrite sets the context to write mode, indicating that the operation is a write operation.
func WithWrite(ctx context.Context) context.Context {
	return context.WithValue(ctx, readWriteModeKey{}, writeMode)
}

type readWriteMode string

func (m readWriteMode) isValid() bool {
	return m == readPrimaryMode || m == readReplicaMode || m == writeMode
}

func getReadWriteMode(ctx context.Context) readWriteMode {
	if mode := ctx.Value(readWriteModeKey{}); mode != nil {
		if v, ok := mode.(readWriteMode); ok && v.isValid() {
			return v
		}
	}

	return notSpecifiedMode
}

func usePrimary(ctx context.Context) bool {
	return getReadWriteMode(ctx) != readReplicaMode
}
