package md

import "context"

// Carrier represents an ability to convert other data into Metadata.
type Carrier interface {
	Extract(ctx context.Context) (context.Context, error)
	Injection(ctx context.Context) error
}
